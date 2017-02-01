package fcfg

import (
	"common/config"

	"os"
	"path/filepath"

	"gopkg.in/fzerorubigd/onion.v2"
)

// Fcfg is the file module settings
var Fcfg Config

// Config is the file module config
type Config struct {
	Size struct {
		MaxUpload   int64 `onion:"max_upload"`
		MaxDownload int64 `onion:"max_download"`
	} `onion:"size"`
	File struct {
		ValidExtension    []string `onion:"valid_extension"`
		ServerPath        string   `onion:"server_path"`
		TempDirectoryPath string   `onion:"temp_directory_path"`
		UploadPath        string   `onion:"upload_path"`
	} `onion:"file"`
}

type configLoader struct {
	o *onion.Onion
}

// Initialize is called when the module is going to add its layer
func (c *configLoader) Initialize(o *onion.Onion) []onion.Layer {
	c.o = o
	def := onion.NewDefaultLayer()
	_ = def.SetDefault("size.max_upload", 104857600)
	_ = def.SetDefault("size.max_download", 104857600)
	_ = def.SetDefault("file.valid_extension", []string{".jpg", ".jpeg", ".mp4", ".png", ".pdf"})
	_ = def.SetDefault("file.server_path", "statics.clickgram.com")
	_ = def.SetDefault("file.upload_path", "http://rubik.clickyab.ae/statics")
	_ = def.SetDefault("file.temp_directory_path", filepath.Join(os.TempDir(), "upload"))
	return []onion.Layer{def}
}

// Loaded inform the modules that all layer are ready
func (c *configLoader) Loaded() {
	c.o.GetStruct("size", &Fcfg.Size)
	c.o.GetStruct("file", &Fcfg.File)
}

func init() {
	config.Register(&configLoader{})
}
