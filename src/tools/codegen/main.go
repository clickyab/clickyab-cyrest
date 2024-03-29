package main

import (
	_ "tools/codegen/datatable" // ENUM plugin
	_ "tools/codegen/enum"      // ENUM plugin
	_ "tools/codegen/gin"       // Gin plugin
	_ "tools/codegen/models"    // Models plugin
	"tools/codegen/plugins"
	_ "tools/codegen/swagger"  // Raml plugin
	_ "tools/codegen/validate" // ENUM plugin

	"github.com/Sirupsen/logrus"
	"github.com/goraz/humanize"
	"github.com/ogier/pflag"
)

var (
	pkg = pflag.StringP("package", "p", "", "the package to scan for gin controller")
)

func main() {
	pflag.Parse()

	p, err := humanize.ParsePackage(*pkg)
	if err != nil {
		logrus.Fatal(err)
	}

	err = plugins.ProcessPackage(*p)
	if err != nil {
		logrus.Fatal(err)
	}

	err = plugins.Finalize(*p)
	if err != nil {
		logrus.Fatal(err)
	}
}
