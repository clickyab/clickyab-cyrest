package fila

import (
	"common/config"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"modules/file/config"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"time"

	"common/assert"
	"common/utils"
	"net/url"
	"sort"
	"strings"

	"gopkg.in/labstack/echo.v3"
)

const (
	// FileTypeImage is the image type
	FileTypeImage FileType = "image"
	// FileTypeVideo is the video type
	FileTypeVideo FileType = "video"
	// FileTypeDocument is the document type
	FileTypeDocument FileType = "document"
)

type (
	// FileType is the file type
	// @Enum{
	// }
	FileType string
)

// File model
// @Model {
//		table = files
//		primary = true, id
//		find_by = id,user_id
//		list = yes
// }
type File struct {
	ID        int64      `db:"id" json:"id" sort:"true" title:"ID"`
	UserID    int64      `json:"user_id" db:"user_id" title:"UserID"`
	Src       string     `json:"src" db:"src" title:"Src"`
	RealName  string     `json:"real_name" db:"real_name" title:"RealName"`
	DBName    string     `json:"db_name" db:"db_name" title:"DBName"`
	Type      FileType   `json:"type" db:"type" title:"Type"`
	Size      int64      `json:"size" db:"size" title:"Size"`
	CreatedAt *time.Time `db:"created_at" json:"created_at" sort:"true" title:"Created at"`
	UpdatedAt *time.Time `db:"updated_at" json:"updated_at" sort:"true" title:"Updated at"`
}

var (
	// DefaultDirPermissions is the default permissions for directories created by gongflow
	DefaultDirPermissions os.FileMode = 0777
	// DefaultFilePermissions is the default permissions for directories created by gongflow
	DefaultFilePermissions os.FileMode = 0600
	// ErrNoTempDir is returned when the temp directory doesn't exist
	ErrNoTempDir = errors.New("gongflow: the temporary directory doesn't exist")
	// ErrCantCreateDir is returned wwhen the temporary directory doesn't exist
	ErrCantCreateDir = errors.New("gongflow: can't create a directory under the temp directory")
	// ErrCantWriteFile is returned when it can't create a directory under the temp directory
	ErrCantWriteFile = errors.New("gongflow: can't write to a file under the temp directory")
	// ErrCantReadFile is returned when it can't read a file under the temp directory (or got back bad data)
	ErrCantReadFile = errors.New("gongflow: can't read a file under the temp directory (or got back bad data)")
	// ErrCantDelete is return when it can't delete a file/directory under the temp directory
	ErrCantDelete                   = errors.New("gongflow: can't delete a file/directory under the temp directory")
	alreadyCheckedDirectory         = false
	lastCheckedDirectoryError error // = nil
	// ValidImgExtension is the image valid ext
	ValidImgExtension = []string{".jpg", ".jpeg", ".png"}
	// ValidVideoExtension is the valid video ext
	ValidVideoExtension = []string{".mov", ".mp4"}
	// ValidDocumentExtension is the valid document ext
	ValidDocumentExtension = []string{".pdf"}
)

// NgFlowData is all the data listed in the "How do I set it up with my server?" section of the ng-flow
// README.md https://github.com/flowjs/flow.js/blob/master/README.md
type NgFlowData struct {
	// ChunkNumber is the index of the chunk in the current upload. First chunk is 1 (no base-0 counting here).
	ChunkNumber int
	// TotalChunks is the total number of chunks.
	TotalChunks int
	// ChunkSize is the general chunk size. Using this value and TotalSize you can calculate the total number of chunks. The "final chunk" can be anything less than 2x chunk size.
	ChunkSize int
	// TotalSize is the total file size.
	TotalSize int
	// TotalSize is a unique identifier for the file contained in the request.
	Identifier string
	// Filename is the original file name (since a bug in Firefox results in the file name not being transmitted in chunk multichunk posts).
	Filename string
	// RelativePath is the file's relative path when selecting a directory (defaults to file name in all browsers except Chrome)
	RelativePath string
}

// ByChunk is the sortable mode
type ByChunk []os.FileInfo

func (a ByChunk) Len() int      { return len(a) }
func (a ByChunk) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByChunk) Less(i, j int) bool {
	ai, _ := strconv.Atoi(a[i].Name())
	aj, _ := strconv.Atoi(a[j].Name())
	return ai < aj
}

// ChunkFlowData does exactly what it says on the tin, it extracts all the flow data from a request object and puts
// it into a nice little struct for you
func ChunkFlowData(r echo.Context) (NgFlowData, error) {
	var err error
	ngfd := NgFlowData{}
	ngfd.ChunkNumber, err = strconv.Atoi(r.FormValue("flowChunkNumber"))
	if err != nil {
		return ngfd, errors.New("Bad ChunkNumber")
	}
	ngfd.TotalChunks, err = strconv.Atoi(r.FormValue("flowTotalChunks"))
	if err != nil {
		return ngfd, errors.New("Bad TotalChunks")
	}
	ngfd.ChunkSize, err = strconv.Atoi(r.FormValue("flowChunkSize"))
	if err != nil {
		return ngfd, errors.New("Bad ChunkSize")
	}
	ngfd.TotalSize, err = strconv.Atoi(r.FormValue("flowTotalSize"))
	if err != nil {
		return ngfd, errors.New("Bad TotalSize")
	}
	ngfd.Identifier = r.FormValue("flowIdentifier")
	if ngfd.Identifier == "" {
		return ngfd, errors.New("Bad Identifier")
	}
	ngfd.Filename = r.FormValue("flowFilename")
	if ngfd.Filename == "" {
		return ngfd, errors.New("Bad Filename")
	}
	ngfd.RelativePath = r.FormValue("flowRelativePath")
	if ngfd.RelativePath == "" {
		return ngfd, errors.New("Bad RelativePath")
	}
	return ngfd, nil
}

// ChunkUpload is used to handle a POST from ng-flow, it will return an empty string for chunk upload (incomplete) and when
// all the chunks have been uploaded, it will return the path to the reconstituted file.  So, you can just keep calling it
// until you get back the path to a file.
func ChunkUpload(tempDir string, ngfd NgFlowData, r echo.Context) (string, error) {
	err := checkDirectory(tempDir)
	if err != nil {
		return "", err
	}
	fileDir, chunkFile := buildPathChunks(tempDir, ngfd)
	err = storeChunk(fileDir, chunkFile, ngfd, r)
	if err != nil {
		return "", errors.New("Unable to store chunk" + err.Error())
	}
	if allChunksUploaded(tempDir, ngfd) {
		file, err := combineChunks(fileDir, ngfd)
		if err != nil {
			return "", err
		}
		return file, nil
	}
	return "", nil
}

// ChunkStatus is used to handle a GET from ng-flow, it will return a (message, 200) for when it already has a chunk, and it
// will return a (message, 404 | 500) when a chunk is incomplete or not started.
func ChunkStatus(tempDir string, ngfd NgFlowData) (string, int) {
	err := checkDirectory(tempDir)
	if err != nil {
		return "Directory is broken: " + err.Error(), http.StatusInternalServerError
	}
	_, chunkFile := buildPathChunks(tempDir, ngfd)
	ChunkNumberString := strconv.Itoa(ngfd.ChunkNumber)
	dat, err := ioutil.ReadFile(chunkFile)
	if err != nil {
		// every thing except for 200, 201, 202, 404, 415. 500, 501
		return "The chunk " + ngfd.Identifier + ":" + ChunkNumberString + " isn't started yet!", http.StatusNotAcceptable
	}
	// An exception for large last chunks, according to ng-flow the last chunk can be anywhere less
	// than 2x the chunk size unless you haave forceChunkSize on... seems like idiocy to me, but alright.
	if ngfd.ChunkNumber != ngfd.TotalChunks && ngfd.ChunkSize != len(dat) {
		return "The chunk " + ngfd.Identifier + ":" + ChunkNumberString + " is the wrong size!", http.StatusInternalServerError
	}

	return "The chunk " + ngfd.Identifier + ":" + ChunkNumberString + " looks great!", http.StatusOK
}

// ChunksCleanup is used to go through the tempDir and remove any chunks and directories older than
// than the timeoutDur, best to set this VERY conservatively.
func ChunksCleanup(tempDir string, timeoutDur time.Duration) error {
	files, err := ioutil.ReadDir(tempDir)
	if err != nil {
		return err
	}
	for _, f := range files {
		fl := path.Join(tempDir, f.Name())
		finfo, err := os.Stat(fl)
		if err != nil {
			return err
		}

		log.Println(f.Name())
		log.Println(time.Since(finfo.ModTime()))
		if time.Since(finfo.ModTime()) > timeoutDur {
			err = os.RemoveAll(fl)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// buildPathChunks simply builds the paths to the ID of the upload, and to the specific Chunk
func buildPathChunks(tempDir string, ngfd NgFlowData) (string, string) {
	filePath := path.Join(tempDir, ngfd.Identifier)
	chunkFile := path.Join(filePath, strconv.Itoa(ngfd.ChunkNumber))
	return filePath, chunkFile
}

// combineChunks will take the chunks uploaded, and combined them into a single file with the
// name as uploaded from the NgFlowData, and it will clean up the chunks as it goes.
func combineChunks(fileDir string, ngfd NgFlowData) (string, error) {
	combinedName := path.Join(fileDir, ngfd.Filename)
	cn, err := os.Create(combinedName)
	if err != nil {
		return "", err
	}

	files, err := ioutil.ReadDir(fileDir)
	sort.Sort(ByChunk(files))
	if err != nil {
		return "", err
	}
	for _, f := range files {
		fl := path.Join(fileDir, f.Name())
		// make sure, we not copy the same file in the final file.
		// the files array contain the full uploaded file name too.
		if fl != combinedName {
			dat, err := ioutil.ReadFile(fl)
			if err != nil {
				return "", err
			}
			_, err = cn.Write(dat)
			if err != nil {
				return "", err
			}
			err = os.Remove(fl)
			if err != nil {
				return "", err
			}
		}
	}

	err = cn.Close()
	if err != nil {
		return "", err
	}
	return combinedName, nil
}

// allChunksUploaded checks if the file is completely uploaded (based on total size)
func allChunksUploaded(tempDir string, ngfd NgFlowData) bool {
	chunksPath := path.Join(tempDir, ngfd.Identifier)
	files, err := ioutil.ReadDir(chunksPath)
	if err != nil {
		log.Println(err)
	}
	totalSize := int64(0)
	for _, f := range files {
		fi, err := os.Stat(path.Join(chunksPath, f.Name()))
		if err != nil {
			log.Println(err)
		}
		totalSize += fi.Size()
	}

	return totalSize == int64(ngfd.TotalSize)
}

// storeChunk puts the chunk in the request into the right place on disk
func storeChunk(tempDir string, tempFile string, ngfd NgFlowData, r echo.Context) error {
	err := os.MkdirAll(tempDir, DefaultDirPermissions)
	if err != nil {
		return errors.New("Bad directory")
	}
	file, _, err := r.Request().FormFile("file")
	if err != nil {
		return errors.New("Can't access file field")
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return errors.New("Can't read file")
	}
	err = ioutil.WriteFile(tempFile, data, DefaultDirPermissions)
	if err != nil {
		return errors.New("Can't write file")
	}
	return nil
}

// checkDirectory makes sure that we have all the needed permissions to the temp directory to
// read/write/delete.  Expensive operation, so it only does it once.
func checkDirectory(d string) error {
	if alreadyCheckedDirectory {
		return lastCheckedDirectoryError
	}

	alreadyCheckedDirectory = true

	if !directoryExists(d) {
		lastCheckedDirectoryError = ErrNoTempDir
		return lastCheckedDirectoryError
	}

	testName := "5d58061677944334bb616ba19cec5cc4"
	testChunk := "42"
	contentName := "foobie"
	testContent := `For instance, on the planet Earth, man had always assumed that he was more intelligent than
	dolphins because he had achieved so much—the wheel, New York, wars and so on—whilst all the dolphins had
	ever done was muck about in the water having a good time. But conversely, the dolphins had always believed
	that they were far more intelligent than man—for precisely the same reasons.`

	p := path.Join(d, testName, testChunk)
	err := os.MkdirAll(p, DefaultDirPermissions)
	if err != nil {
		lastCheckedDirectoryError = ErrCantCreateDir
		return lastCheckedDirectoryError
	}

	f := path.Join(p, contentName)
	err = ioutil.WriteFile(f, []byte(testContent), DefaultFilePermissions)
	if err != nil {
		lastCheckedDirectoryError = ErrCantWriteFile
		return lastCheckedDirectoryError
	}

	b, err := ioutil.ReadFile(f)
	if err != nil {
		lastCheckedDirectoryError = ErrCantReadFile
		return lastCheckedDirectoryError
	}
	if string(b) != testContent {
		lastCheckedDirectoryError = ErrCantReadFile // TODO: This should probably be a different error
		return lastCheckedDirectoryError
	}

	err = os.RemoveAll(path.Join(d, testName))
	if err != nil {
		lastCheckedDirectoryError = ErrCantDelete
		return lastCheckedDirectoryError
	}

	if os.TempDir() == d {
		log.Println("You should really have a directory just for upload temp (different from system temp).  It is OK, but consider making a subdirectory for it.")
	}

	return nil
}

// directoryExists checks if the directory exists of course!
func directoryExists(d string) bool {
	finfo, err := os.Stat(d)

	if err == nil && finfo.IsDir() {
		return true
	}
	return false
}

// UploadFromURL is the upload main route
func UploadFromURL(link string, uID int64) (string, error) {
	year := fmt.Sprintf("%d", time.Now().Year())
	month := time.Now().Month().String()
	basePath := filepath.Join(config.Config.StaticRoot, year, month)
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		err = os.MkdirAll(basePath, DefaultDirPermissions)
		assert.Nil(err)
	}

	extension := strings.ToLower(filepath.Ext(link))

	tokens := strings.Split(link, "/")
	realFileName := tokens[len(tokens)-1]

	//check extension
	typ := FileTypeImage
	if utils.StringInArray(extension, ValidImgExtension...) {
		typ = FileTypeImage
	} else if utils.StringInArray(extension, ValidVideoExtension...) {
		typ = FileTypeVideo
	} else if utils.StringInArray(extension, ValidDocumentExtension...) {
		typ = FileTypeDocument
	} else {
		return "", errors.New("error file type")
	}

	//check file type

	hash := utils.Sha1(fmt.Sprintf("%d", time.Now().UnixNano()))
	newFileName := fmt.Sprintf("%s%s", hash, extension)
	filePath := filepath.Join(basePath, newFileName)
	out, err := os.Create(filePath)
	defer func() {
		err = out.Close()
		assert.Nil(err)
	}()
	if err != nil {
		err = os.Remove(filePath)
		assert.Nil(err)
		return "", errors.New("error while uploading file")
	}
	resp, err := http.Get(link)
	if err != nil {
		err = os.Remove(filePath)
		assert.Nil(err)
		return "", errors.New("error while uploading file")
	}
	if err != nil {
		err = os.Remove(filePath)
		assert.Nil(err)
		return "", errors.New("error while uploading file")
	}
	defer func() {
		err = resp.Body.Close()
		assert.Nil(err)
	}()
	if err != nil {
		err = os.Remove(filePath)
		assert.Nil(err)
		return "", errors.New("error while uploading file")
	}
	downSize, err := io.Copy(out, resp.Body)
	if err != nil {
		err = os.Remove(filePath)
		assert.Nil(err)
		return "", errors.New("error while uploading file")
	}
	if downSize > fcfg.Fcfg.Size.MaxDownload {
		err = os.Remove(filePath)
		assert.Nil(err)
		return "", errors.New("size not valid")
	}
	fpath := filepath.Join(year, month, newFileName)

	newFile := &File{
		DBName:   newFileName,
		RealName: realFileName,
		Src:      fpath,
		Type:     typ,
		Size:     downSize,
		UserID:   uID,
	}
	assert.Nil(NewFilaManager().CreateFile(newFile))
	return "/" + fpath, nil
}

// CheckUpload is the check upload func
func CheckUpload(link string, uID int64) (string, error) {
	urlObj, err := url.Parse(link)
	if err != nil {
		return "", errors.New("url not valid")
	}
	host := urlObj.Host
	if host == fcfg.Fcfg.File.SameUploadPath {
		return strings.Replace(urlObj.Path, fcfg.Fcfg.File.UploadURLReplace, "", 1), nil
	}
	return UploadFromURL(link, uID)
}
