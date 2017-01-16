package main

import (
	"common/config"
	"common/initializer"
)

func main() {
	config.Initialize()
	config.InitApplication()

	defer initializer.Initialize().Finalize()

}
