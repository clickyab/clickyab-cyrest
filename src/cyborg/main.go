package main

import (
	"common/assert"
	"common/config"
	"common/initializer"
	"common/utils"
	"common/version"
	"modules/telegram/cyborg/worker"
	"net"
)

func main() {
	config.Initialize()
	config.InitApplication()
	config.Config.AMQP.Publisher = 1

	defer initializer.Initialize().Finalize()

	version.LogVersion().Infof("Application started")
	ip, err := net.LookupIP(config.Config.Telegram.CLIAddress)
	assert.Nil(err)
	assert.True(len(ip) > 0, "no ip found")
	_, err = worker.NewMultiWorker(ip[0], config.Config.Telegram.CLIPort)
	assert.Nil(err)

	utils.WaitExitSignal()

	version.LogVersion().Info("Goodbye")
}
