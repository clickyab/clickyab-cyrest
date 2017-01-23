package main

import (
	"common/config"
	"common/initializer"
	"common/rabbit"
	"modules/telegram/cyborg/commands"
)

func main() {
	config.Initialize()
	config.InitApplication()

	defer initializer.Initialize().Finalize()
	rabbit.MustPublish(commands.SelectAd{
		ChannelID: 2,
		ChatID:    1,
	})
	rabbit.MustPublish(commands.UpdateMessage{
		CLiChannelName: "daratest",
	})

}
