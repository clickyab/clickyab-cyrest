package main

import (
	"common/assert"
	"common/config"
	"common/initializer"
	"common/rabbit"
	"common/utils"
	"common/version"
	"modules/telegram/ad/bot/worker"
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
	go func() {
		err := rabbit.RunWorker(
			&bot.SendWarn{}, bot.SendWarnAction, 10,
		)
		assert.Nil(err)

	}()
	utils.WaitExitSignal()

	version.LogVersion().Info("Goodbye")
}
