package main

import (
	"common/config"
	"common/initializer"

	"common/rabbit"
	"common/utils"
	"modules/cyborg/commands"
)

func main() {
	config.Initialize()
	config.InitApplication()

	defer initializer.Initialize().Finalize()
	j := commands.GetLastCommand{
		HashKey: <-utils.ID,
		Channel: "tst1234567",
	}
	rabbit.MustPublish(j)
}
