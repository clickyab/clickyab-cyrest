package main

import (
	"common/config"
	"common/initializer"
	"common/utils"
	"common/version"
	"modules/misc/base"
)

func main() {
	config.Initialize()
	config.InitApplication()

	defer initializer.Initialize().Finalize()

	version.LogVersion().Infof("Application started")

	go base.Initialize(config.Config.MountPoint).Start(config.Config.Port)

	utils.WaitExitSignal()

	version.LogVersion().Info("Goodbye")
}
