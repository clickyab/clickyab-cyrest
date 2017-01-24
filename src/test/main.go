package main

import (
	"common/config"
	"common/initializer"
	"common/rabbit"
	"modules/telegram/ad/bot/worker"
)

func main() {
	config.Initialize()
	config.InitApplication()

	defer initializer.Initialize().Finalize()
	rabbit.MustPublish(bot.SendWarn{
		Msg:       "erhabi",
		ChatID:    70018667,
		ChannelID: 1,
		AdID:      1,
	})

}
