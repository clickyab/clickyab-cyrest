// Code generated by go-bindata.
// sources:
// db/migrations/.gitignore
// db/migrations/20160117193701_user_base.sql
// DO NOT EDIT!

package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data, name string) ([]byte, error) {
	gz, err := gzip.NewReader(strings.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _dbMigrationsGitignore = "\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xd2\xe2\x52\xd4\x4b\xcf\x2c\xc9\x4c\xcf\xcb\x2f\x4a\xe5\x02\x04\x00\x00\xff\xff\x17\xeb\x9a\xa9\x0e\x00\x00\x00"

func dbMigrationsGitignoreBytes() ([]byte, error) {
	return bindataRead(
		_dbMigrationsGitignore,
		"db/migrations/.gitignore",
	)
}

func dbMigrationsGitignore() (*asset, error) {
	bytes, err := dbMigrationsGitignoreBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "db/migrations/.gitignore", size: 14, mode: os.FileMode(436), modTime: time.Unix(1472965876, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _dbMigrations20160117193701_user_baseSql = "\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xcc\x59\x4d\x73\xe2\x38\x13\xbe\xf3\x2b\x74\x03\xea\x4d\xaa\x92\xcc\x3b\x73\xd8\x9c\xbc\xe0\xec\x52\xcb\x47\xc6\x31\x5b\x33\x27\x47\xb6\x05\x68\x31\x96\xcb\x36\xc9\xb0\xbf\x7e\x25\xdb\xb2\x24\x4b\x36\x04\xd8\xda\x99\xca\x01\xd4\xea\x47\xdd\xad\xd6\xd3\xdd\x43\xef\xf6\x16\xfc\x6f\x87\xd7\x29\xcc\x11\x58\x26\xec\xeb\xcb\xd7\x29\xc0\x31\xc8\x50\x90\x63\x12\x83\xfe\x32\xe9\x03\x9c\x01\xf4\x03\x05\xfb\x1c\x85\xe0\x7d\x83\x62\x90\x6f\xe8\x52\xa9\xc7\x36\xd1\x2f\x30\x49\x22\x8c\xc2\x5e\x6f\xe4\xd8\x96\x6b\x03\xd7\xfa\x75\x6a\x83\x7d\x86\xd2\xac\x37\xe8\x01\xfa\x0f\x87\x60\x32\x77\x07\xf7\xf7\x43\xf0\xec\x4c\x66\x96\xf3\x1d\xfc\x61\x7f\x07\xf3\x85\x0b\xe6\xcb\xe9\x14\x58\x4b\x77\xe1\x4d\xe6\x54\x7f\x66\xcf\xdd\x9b\x42\x07\xed\x20\x8e\xc0\x9f\x96\x33\xfa\xdd\x72\x06\x0f\x9f\x3f\x0f\xcb\xf5\x04\x66\xd9\x3b\x49\x43\x83\x88\x44\xa1\xd7\x21\x86\x41\x80\xb2\xcc\xcb\xc9\x96\xba\xa1\x8b\x33\xb2\x4f\x03\x04\xec\xf9\x72\x36\xe8\x07\xe9\xae\x7f\x03\xfa\x41\x84\x83\xed\x01\xfa\xfd\x6a\x0f\x73\xca\xcb\x0f\x09\xdf\x96\x50\x1f\x49\x0c\xa3\x62\x2f\x49\x13\x52\x46\xa5\x5f\xdb\x9a\xa2\x38\xf7\x84\xfb\x95\x21\x6f\x30\x87\xa9\xc9\x84\x1c\xe6\xfb\xac\xc2\x4e\xd1\x1a\x67\x39\x4a\x51\xc8\xd0\xdf\x50\x8a\x57\xb8\xfc\xec\x47\x24\xd8\xd2\x8f\x95\x56\x90\x22\x7a\x87\xa1\x07\x73\xe0\x4e\x66\xf6\x8b\x6b\xcd\x9e\xc1\xd8\x7e\xb2\x96\x53\x17\x8c\x96\x8e\x43\x63\xea\x09\x09\x8f\x7a\xe5\x51\x12\xb6\x2b\xf7\xef\xe8\xbf\xdb\xe2\x0f\xdc\xdd\xfd\x52\xfc\xf5\x6b\x80\xde\xf0\x51\xbf\x72\x2f\x25\x11\xaa\xae\xbd\xf8\x2e\xdd\xbd\x80\xed\x37\xac\x60\x4a\xa7\xed\x34\x39\x5b\x4a\x46\x8b\xf9\x8b\xeb\x58\x14\x02\xbc\x56\x59\xf6\xaa\xa4\xdb\xa0\xb2\xe7\x86\x1f\x37\x34\x7b\x90\xa4\x64\x85\xe9\x06\x7e\xb9\x2d\xde\x98\x32\xb9\xb4\x64\x85\xd3\x2c\xf7\x62\xb8\x43\x86\x3b\x8e\x60\xbb\xcc\xc7\x69\xbe\x09\xe1\x01\x8c\xa9\x49\xcc\xbd\x72\x79\x8d\xe2\x10\xa5\x55\x5a\xcc\x60\x84\x58\x12\x3c\xd1\x07\x42\x3f\xf1\x1c\x40\x51\x94\x6c\x48\x6c\x42\x6d\x5b\x87\x61\x98\xd2\x07\x01\x5c\xfb\x5b\xf5\xe6\xfe\xc6\x89\x17\x90\xd0\xb4\x39\x2e\x12\x1b\x46\xaa\xfc\xfe\x8e\x9f\x4f\xf6\x71\x9e\x1e\xb4\x54\xa7\xb1\x7c\xc3\x71\x80\x34\x41\x80\x73\x7d\xb7\x74\xb9\x3c\x02\xed\xe9\xaa\xed\xe8\xbe\x4d\xe9\x79\x7e\xf8\x42\x73\x9c\x47\xa6\xa0\xa0\x80\xc4\x64\x87\x83\x32\x28\x8a\x2b\xfc\xf5\x1a\x44\xe5\x7d\xa8\x8c\x50\x5d\x85\x7e\xc4\xbf\x16\x58\x39\x9e\xf2\xba\x39\x8a\x2b\x1c\xc3\x38\xc0\xf5\x63\x38\x87\xd1\x9b\xf1\x56\x63\xec\xc3\x78\xdb\xf6\x2e\x28\x73\xb3\x30\x78\x1b\xca\xf0\xc8\x44\x9c\x01\x4c\x43\x2f\xde\xef\x7c\x55\xda\xd4\xd7\x77\x70\xe2\xdd\x20\x1f\xb6\x8b\xdb\x49\xc7\x44\x9f\xe6\x08\xd2\x8a\xf2\xe1\xbc\x23\x29\x5e\x63\xf6\xf2\xe2\xb5\x17\x21\x68\x2a\x6a\xc1\x3e\xcb\xc9\x4e\xc9\xb3\x87\x4a\xb4\xc6\x42\xe1\xcb\xff\x39\xfd\x50\x18\xaf\xaa\x33\xc5\xee\x3a\x5f\xe9\xba\x7f\x60\x66\x16\xeb\x0f\x2d\xe4\x08\xf3\x3c\xc5\x3e\x6d\x07\x2e\xa9\xed\x8d\x08\x94\x8b\xaf\x5b\x74\x78\x95\x98\x85\x53\xcb\x1b\x8c\xf6\x48\x59\xd7\x2c\x63\x7c\x7e\x89\x3d\x2d\x69\xf7\xd3\x14\xd6\xa2\x5e\xd1\x72\xb4\xc3\x59\x26\xf8\xeb\x1c\x4f\x9b\x85\x56\xb5\x59\x1c\x61\x7a\x23\x01\xa9\xfb\x9e\x75\x44\xfc\xb2\xeb\x29\x5b\x1c\xf6\x89\xbc\xc7\x3f\x5d\x47\xc2\x59\xf2\x92\xe4\xd0\x19\xb8\x61\x30\xcb\x1e\xcf\xd8\x52\x16\x92\x15\x34\x48\x58\x79\x20\xb1\x78\xb4\x0f\xd7\x61\x9a\xd2\x56\x7c\xaa\xbf\x25\x32\xce\x88\xb0\xd0\xe4\x9c\xf2\xf8\x74\xe5\x4f\xb5\xfc\x53\xa3\x68\x15\xee\xbd\xcc\xac\xe9\x94\x59\xf1\x45\xae\x7f\xc2\xf5\x2b\x91\x2c\xad\x7b\x27\xfb\x6d\xbc\x67\x43\x45\xbd\xc2\x45\x17\xe5\x58\xef\x01\x2e\x74\x76\x47\x9b\x05\xb8\x46\x5e\x44\xd6\xd7\xa6\xe1\x6c\xef\xff\x45\x07\x40\x63\x2b\x12\xe7\xf4\xad\x4b\xfd\xe2\x7f\xf5\xce\x8d\xd3\x12\x2c\x66\xa3\x84\xf6\xca\xb4\x5c\xb2\x8f\x19\x23\xa6\xa1\x36\x1c\x48\xb1\xf3\x8a\x21\x95\x86\xc0\x5b\x6d\xc1\xd3\xc2\xb1\x27\xbf\xcd\x95\x51\x61\x08\x1c\xfb\xc9\xa6\x6e\x8c\xec\x97\x72\xa2\x05\x83\x6a\x6e\xb0\xa6\xae\xed\xc8\xb3\x2e\xb0\xc6\x63\x15\xa4\x1e\xff\x5a\x60\xea\x4b\x9d\xcc\xc7\xf6\xb7\x52\xa0\x98\xb4\x98\xf3\xdd\x02\xaa\x56\x5a\xce\x27\x5f\x97\xb2\x2e\x53\xda\x63\x3a\x2a\xfc\x90\x14\x0b\x8d\xa6\xad\xc5\x90\xa6\xdb\xdb\xed\xf4\xc9\x30\x7c\xba\x92\x61\x8a\x0a\xdd\xe6\x74\x81\xe3\x55\x6a\x92\xdf\x25\x7e\x8d\x67\x30\xa0\x39\xaa\x5d\xc1\xa7\xe3\x90\xa2\x20\x28\xa8\x35\xf7\x9e\x8f\x2c\x51\x90\x02\x5d\x97\xb1\x0b\x8c\x2e\xc7\x02\xd5\xe2\x82\x32\x4d\xb7\xc2\x71\xc4\xf4\x54\xaa\xcb\x97\xa3\x9d\x59\x9f\x71\x14\xac\x0e\xe0\x11\x3c\x11\xe8\x63\x90\x52\xe4\xba\x31\xe5\x10\x9b\x9f\x92\x86\x6d\x7e\x5b\x06\x6c\x9e\x6b\x1d\x17\x24\x4d\xa1\x57\xcc\xd5\x4e\xd4\x4b\xd3\xb5\x13\xfc\xe2\x8c\xed\x36\xfd\x43\x49\x2b\x41\x9d\x92\xb7\xf2\xc9\x6d\xa9\x6b\x84\xec\xcc\x5e\x15\xb5\x35\x81\x4d\xc0\x47\x72\x58\x41\x3e\x9e\xc6\xa6\x13\x8e\x64\xb2\x72\x42\x47\x32\xd7\xff\x19\xf0\xf1\x14\x6e\xad\x5d\x35\xa6\xc1\x3e\x71\x5e\x5b\x01\x11\xda\xdc\x45\x29\x80\x92\x7a\x87\x53\x6c\xf0\xbd\x9e\x3b\x14\xad\x2d\xda\xec\xa0\x2e\x43\xc4\x84\x7d\x45\x7b\x04\xa8\xc1\x22\xe9\xc4\xb6\x00\x4b\xfa\x86\x06\x45\x41\x10\xae\x99\xac\xe1\x55\x5e\x98\x20\xb7\x06\x72\x30\x1a\x43\xef\xd9\x7d\x86\x7e\xbc\x00\xd5\x2d\x91\x4f\x34\x04\xa3\xa9\xaf\x36\x2d\x1a\x80\xb9\x75\x11\xe4\x78\x2e\x59\xab\xc5\x90\x73\x81\xc6\x4c\x12\x0b\x9b\xa8\x48\x2d\x7b\x12\xfb\x88\x90\x34\x79\xdc\xa4\x39\xe2\xd0\x19\x91\x34\x65\x93\x33\x72\x44\xb5\x98\x9f\x8c\xba\x4c\xd2\x08\x5f\x55\x01\xce\xae\x45\x46\xb6\xac\xea\x84\x64\x43\x7b\x9d\x61\x7b\x75\xb2\xe6\x0a\x46\x6a\x2e\x35\xdb\x46\x0f\xaa\x2d\x89\xe4\x27\xa4\xfc\x4a\x37\x26\xef\x31\xff\x9d\xae\xfe\x91\x8e\x2d\x9e\xf4\x33\x1d\xcd\xc5\x88\x4a\x7d\x18\x6c\x7b\xbd\xb1\xb3\x78\x6e\x76\xf2\x8f\xf2\xaa\x64\xcf\xa3\xb6\xbb\xd9\xfa\xb4\xef\x90\x4a\x8a\xbe\xa9\x26\x66\x5d\x44\x49\x52\x5f\x14\x2c\xa3\xc8\x1a\xaf\x4e\x93\xa9\xbb\xcb\x6b\x52\x96\xea\x04\x51\x37\xf2\x24\xd4\x0c\xa1\x2b\xff\x04\x00\x00\xff\xff\x42\xf8\xa1\x41\x3e\x1d\x00\x00"

func dbMigrations20160117193701_user_baseSqlBytes() ([]byte, error) {
	return bindataRead(
		_dbMigrations20160117193701_user_baseSql,
		"db/migrations/20160117193701_user_base.sql",
	)
}

func dbMigrations20160117193701_user_baseSql() (*asset, error) {
	bytes, err := dbMigrations20160117193701_user_baseSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "db/migrations/20160117193701_user_base.sql", size: 7486, mode: os.FileMode(436), modTime: time.Unix(1480422304, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"db/migrations/.gitignore": dbMigrationsGitignore,
	"db/migrations/20160117193701_user_base.sql": dbMigrations20160117193701_user_baseSql,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}
var _bintree = &bintree{nil, map[string]*bintree{
	"db": &bintree{nil, map[string]*bintree{
		"migrations": &bintree{nil, map[string]*bintree{
			".gitignore": &bintree{dbMigrationsGitignore, map[string]*bintree{}},
			"20160117193701_user_base.sql": &bintree{dbMigrations20160117193701_user_baseSql, map[string]*bintree{}},
		}},
	}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}

