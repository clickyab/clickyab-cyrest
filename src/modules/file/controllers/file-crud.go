package file

import (
	"modules/user/middlewares"

	"common/utils"

	"os"
	"path/filepath"

	"modules/file/config"

	"common/assert"
	"common/config"
	"fmt"
	"modules/file/fila"
	"path"
	"time"

	"errors"

	"strings"

	"io"

	"gopkg.in/labstack/echo.v3"
)

//	upload upload file
//	@Route	{
//		url	=	/upload
//		method	= post
//		resource = upload_file:self
//		middleware = authz.Authenticate
//		400 = base.ErrorResponseSimple
//	}
func (u *Controller) upload(ctx echo.Context) error {
	year := fmt.Sprintf("%d", time.Now().Year())
	month := time.Now().Month().String()
	currentUser, ok := authz.GetUser(ctx)
	if !ok {
		return u.NotFoundResponse(ctx, nil)
	}
	flowData, err := fila.ChunkFlowData(ctx)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	// TODO change in deploy
	/*tmpPath := filepath.Join(os.TempDir(), "upload")
	if _, err := os.Stat(tmpPath); os.IsNotExist(err) {
		os.MkdirAll(tmpPath, fila.DefaultDirPermissions)
	}*/
	if _, err := os.Stat(fcfg.Fcfg.File.TempDirectoryPath); os.IsNotExist(err) {
		os.MkdirAll(fcfg.Fcfg.File.TempDirectoryPath, fila.DefaultDirPermissions)
	}
	file, err := fila.ChunkUpload(fcfg.Fcfg.File.TempDirectoryPath, flowData, ctx)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	if file == "" {
		return u.OKResponse(ctx, nil)
	}
	fileObj, err := os.Open(file)
	defer fileObj.Close()
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	fileInfo, err := fileObj.Stat()
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	extension := strings.ToLower(filepath.Ext(fileObj.Name()))

	//check extension to be valid
	// TODO :// check extension from config
	//check extension
	typ := fila.FileTypeImage
	if utils.StringInArray(extension, fila.ValidImgExtension...) {
		typ = fila.FileTypeImage
	} else if utils.StringInArray(extension, fila.ValidVideoExtension...) {
		typ = fila.FileTypeVideo
	} else if utils.StringInArray(extension, fila.ValidDocumentExtension...) {
		typ = fila.FileTypeDocument
	} else {
		return u.NotFoundResponse(ctx, nil)
	}
	size := fileInfo.Size()
	//check size to be valid
	if size > fcfg.Fcfg.Size.MaxUpload {
		return u.NotFoundResponse(ctx, errors.New("too big"))
	}

	basePath := filepath.Join(config.Config.StaticRoot, fmt.Sprintf("%d", time.Now().Year()), time.Now().Month().String())
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		os.MkdirAll(basePath, fila.DefaultDirPermissions)
	}

	hash := utils.Sha1(fmt.Sprintf("%d", time.Now().UnixNano()))
	newFilename := fmt.Sprintf("%s%s", hash, extension)
	out, err := os.Create(path.Join(basePath, newFilename))
	if err != nil {
		os.Remove(path.Join(basePath, newFilename))
		return u.NotFoundResponse(ctx, nil)
	}
	defer out.Close()
	_, err = io.Copy(out, fileObj)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	err = os.RemoveAll(filepath.Dir(fileObj.Name()))
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	srcPath := filepath.Join(year, month, newFilename)
	m := fila.NewFilaManager()
	newFile := &fila.File{
		RealName: fileInfo.Name(),
		DBName:   newFilename,
		Src:      srcPath,
		Type:     typ,
		UserID:   currentUser.ID,
		Size:     size,
	}
	assert.Nil(m.CreateFile(newFile))
	// TODO : Dirty hack to handle upload
	newFile.Src = fmt.Sprintf("%s%s", fcfg.Fcfg.File.UploadPath, filepath.Join(year, month, newFilename))
	return u.OKResponse(ctx, newFile)
}
