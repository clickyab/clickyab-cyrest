package main

import (
	"common/assert"
	"common/config"
	"common/initializer"
	"common/tgbot"
	"common/version"
)

func main() {
	config.Initialize()
	config.InitApplication()

	defer initializer.Initialize().Finalize()
	version.LogVersion().Infof("Application started")

	go func() {
		assert.Nil(tgbot.Start())
	}()

	version.LogVersion().Info("Goodbye")
}
