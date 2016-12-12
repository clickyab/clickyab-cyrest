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

	info := bindataFileInfo{name: "db/migrations/.gitignore", size: 14, mode: os.FileMode(420), modTime: time.Unix(1470570111, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _dbMigrations20160117193701_user_baseSql = "\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xcc\x58\x5d\x6f\xdb\x36\x14\x7d\xf7\xaf\xe0\x9b\x6c\xac\x05\x94\x15\xe9\xc3\xf2\xa4\x39\xcc\x66\xcc\x96\x5b\x45\x1a\xda\x27\x95\x96\xae\x63\xc2\x34\x29\x90\x74\xd2\xec\xd7\x0f\xfa\xa6\x28\xc9\x36\xe6\x14\x1d\x60\x04\x8a\xee\xe5\x21\x79\xee\x39\x57\x94\x26\xef\xdf\xa3\x5f\x0e\xf4\x49\x12\x0d\x28\xca\xf2\x7f\x1f\x3f\x2f\x11\xe5\x48\x41\xa2\xa9\xe0\xc8\x89\x32\x07\x51\x85\xe0\x3b\x24\x47\x0d\x29\x7a\xd9\x01\x47\x7a\x47\x15\x2a\xc7\xe5\x49\x54\x21\x92\x65\x8c\x42\x3a\x99\xcc\x03\xec\x85\x18\x85\xde\xef\x4b\x8c\x8e\x0a\xa4\x9a\x4c\x27\x08\x21\x44\x53\xb4\xf0\xc3\xe9\xcd\xcd\x0c\x7d\x0a\x16\x2b\x2f\xf8\x8a\xfe\xc2\x5f\x91\xbf\x0e\x91\x1f\x2d\x97\xc8\x8b\xc2\x75\xbc\xf0\xe7\x01\x5e\x61\x3f\x7c\x57\x8c\x81\x03\xa1\x0c\xfd\xed\x05\xf3\x3f\xbd\x60\x7a\xeb\xce\xda\xf4\xc8\x5f\x7c\x8e\x30\x2a\xf3\x32\xa2\xd4\x8b\x90\x69\x93\xfa\xd1\x9d\x35\x99\x65\x8a\x60\x69\xdc\x4b\xfb\xe0\xce\xca\x28\x49\x12\x50\x2a\xd6\x62\x0f\xdc\x8c\x22\x0b\x25\xdf\x50\xac\x5f\x33\x40\xd8\x8f\x56\x53\x27\x03\xa9\x04\x27\xcc\x79\x87\x9c\x44\xc8\x4c\x94\x8c\x38\xb3\x7a\x5d\x12\xb8\x8e\xdb\xad\x57\xb3\x3d\x13\x4d\x64\x33\xcf\xcd\xc7\x7a\x19\x4a\x13\x7d\x54\x15\xb6\x84\x27\xaa\x34\x48\x48\x73\xf4\x67\x90\x74\x4b\xcb\xeb\x0d\x13\xc9\x1e\x52\xa7\xb7\xbc\x44\x02\xd1\x90\xc6\x44\xa3\x70\xb1\xc2\x8f\xa1\xb7\xfa\x84\xee\xf1\x83\x17\x2d\x43\x34\x8f\x82\x00\xfb\x61\xdc\x46\xea\xd1\xd5\xde\xb2\x74\x7c\xb0\xe3\xba\xae\xfb\xbe\xf8\x21\xd7\xfd\xad\xf8\x39\x16\xc0\x7c\xed\x3f\x86\x81\xb7\xf0\xc3\x92\x27\x9a\x33\x5e\x11\x10\x6f\xf7\xe8\x61\x1d\xe0\xc5\x1f\x7e\x51\xf6\x69\x13\x99\xa1\x00\x3f\xe0\x00\xfb\x73\xfc\x58\x0a\x06\x4d\x69\x3a\x9b\xcc\xee\x1a\x31\x2d\xfc\x7b\xfc\xa5\x8c\xc5\xe5\xdf\x12\x70\xed\xd7\x03\x5a\xb4\xbb\xc9\xc4\x12\xa1\x14\x0c\xae\x11\x21\x27\x07\x68\x6a\xf5\xeb\xed\xed\xec\xe7\x72\x9d\x13\x53\xed\xaf\xf2\x40\x49\x4f\xbe\xcd\x9c\x97\x23\xe5\x29\x7c\xcf\xb9\x29\x36\x5e\x90\x79\x37\xc4\x49\x9c\x81\x3c\x50\xa5\xa8\xe0\x57\xb0\x53\x4d\xdb\x0c\xec\xee\xb3\x9d\x62\x80\x41\x95\x88\xc6\x48\x4f\x4c\x6c\x4a\x1b\x95\xa5\xcc\xaf\x14\xb0\xad\xf3\x93\xd9\xee\x29\xdb\x62\x2e\x2e\x48\x1e\x12\x78\x45\x4c\x47\xde\x6d\x45\x4e\x56\xd1\x80\xef\x15\xd4\x08\x56\xa5\xed\x98\x64\x68\x79\xad\x5b\x7a\x00\xf5\x22\xef\x06\x1a\x77\x4c\xb4\x96\x74\x73\xd4\x57\xb9\xa7\xea\x04\xdd\x06\xf8\x6d\x0f\xaf\xdf\xda\xfe\xe7\xd6\xfd\xef\x99\xb0\x23\x0c\xdc\xb7\x3b\x4b\xbb\xb2\x4e\x3f\xe8\xf0\x5f\x4d\x7c\xa2\xbd\x0c\xf1\x6f\xc3\x77\xf8\xb7\x82\x43\xfc\x9f\x5a\xde\x10\x42\xbd\xca\xc1\x02\x6c\x29\x27\x3c\xa1\x84\xbd\x1d\xff\x96\xb0\x37\x84\xef\xe3\x91\x06\x47\x92\x44\x1c\xb9\x8e\x77\x82\xa5\x20\x87\x3a\x20\x91\x69\xcc\x8f\x87\x4d\x37\x6a\x8f\xef\x67\xd4\xfe\xdf\xc1\x86\x8c\x87\xff\x37\x96\xef\xd6\xe2\x87\x28\xae\x45\xef\x0b\xae\x89\x8d\xea\xad\xbb\x36\x4b\x6d\xc6\xf0\x93\x62\xcb\xa4\xd8\x52\x06\xb1\x71\x8c\xa9\x64\x77\x5a\x42\x9a\x6a\xd6\x95\x8f\x95\x00\x89\xe0\xe2\x40\x93\x38\x11\xa9\x91\x58\x7b\xbb\x3e\xe6\xf4\xc2\xf6\xa3\x64\x27\xb8\xd1\x1b\x3e\xd4\x2a\x4b\x53\x09\x4a\x0d\xe9\x27\x57\x9f\x7c\xed\x35\x9f\x4c\x8a\x67\xca\x13\xe8\x05\x12\xaa\xfb\xd9\x86\x0a\xef\xbd\x10\xe7\x62\x32\x15\xd6\xbf\x6f\x0b\x67\x80\xd7\x37\x90\x50\x59\x7d\x13\xb3\x99\xa7\xdc\x85\x29\x81\x81\x25\xa0\x69\x95\x77\x11\x64\xc3\xe4\x79\xd4\x26\xf5\x12\x60\xa3\x14\x67\x91\x8d\xdc\x11\x27\x0d\xcd\x50\x1b\xc2\xf2\xd4\xe0\x0c\x17\xd9\xa3\x3e\xf1\x5f\xe4\x8d\x2d\x95\x4a\x8f\xf5\x57\x46\xc6\x63\x1b\x2a\xf5\x2e\x25\xaf\x96\xb6\x9e\x80\xe7\xad\xb8\x3c\x32\x1d\x08\x83\xfc\x98\xb4\x85\xe2\xaa\xd6\x2b\x30\xd6\xb5\x4a\x8b\x3a\x76\xbf\xf6\x50\x88\xbf\x54\x4f\x8d\x7f\x68\x66\xf9\xb1\x49\xe6\x05\x5b\x84\x75\xe3\x37\xee\x8f\x76\xdd\x78\x83\x1f\xc9\x18\x33\x62\x5d\xc1\x37\x73\x61\x03\x78\x81\x05\xeb\xdc\x51\xff\xf5\xc1\x4e\x9a\xcf\xc0\x1b\x75\x5e\x0f\xf2\x8c\xed\x5a\xcc\xf3\x9e\xeb\x61\x9f\x31\x5c\x8b\x6d\xb8\x6d\xc0\x6e\xf9\xc1\x74\xc4\x62\xed\x53\xdc\x7e\x66\xdb\xef\x22\xe3\x99\x57\x9d\x2d\x6c\x61\x15\xd3\xfe\x47\x31\x9d\x40\xbc\xfe\xcd\xc2\x38\x22\x34\x88\x56\xcd\xf3\x5b\xe6\x6b\xc0\xc8\xc8\xa1\x63\x6c\x39\xd4\x2c\xa3\xf9\x35\xe9\x5e\xbc\xf0\xfa\x7b\x52\xf3\x31\x29\xbf\x79\xd1\xe7\x24\x29\x18\x83\x14\x6d\x48\xb2\x9f\xfc\x1b\x00\x00\xff\xff\x89\x12\x9e\x19\xa5\x12\x00\x00"

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

	info := bindataFileInfo{name: "db/migrations/20160117193701_user_base.sql", size: 4773, mode: os.FileMode(420), modTime: time.Unix(1481612293, 0)}
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

