package main

import (
	"common/config"
	"common/initializer"
	"common/rabbit"
	"modules/cyborg/commands"
)

func main() {
	config.Initialize()
	config.InitApplication()

	defer initializer.Initialize().Finalize()
	rabbit.MustPublish(commands.GetChanCommand{
		ChannelID: 1,
		Count:     20,
	})

}
