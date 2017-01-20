package main

import (
	"common/assert"
	"common/config"
	"common/initializer"
	"common/utils"
	"common/version"
	"modules/telegram/common/tgbot"
)

func main() {
	config.Initialize()
	config.InitApplication()

	defer initializer.Initialize().Finalize()
	version.LogVersion().Infof("Application started")

	go func() {
		assert.Nil(tgbot.Start())
	}()
	utils.WaitExitSignal()

	version.LogVersion().Info("Goodbye")
}
