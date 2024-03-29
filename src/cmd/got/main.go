package main

import (
	"common/assert"
	"common/config"
	"common/initializer"
	"common/rabbit"
	"common/utils"
	"common/version"

	"modules/telegram/ad/worker"
	bot "modules/telegram/bot/worker"
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
		rabbit.RunWorker(
			&bot.SendWarn{}, bot.SendWarnAction, 10,
		)
	}()
	go func() {
		rabbit.RunWorker(
			&bot.ClearTrans{}, bot.ClearTransAction, 10,
		)
	}()
	go func() {
		rabbit.RunWorker(
			&worker.AdDelivery{},
			worker.AdDeliveryAction,
			10,
		)
	}()
	utils.WaitExitSignal()

	version.LogVersion().Info("Goodbye")
}
