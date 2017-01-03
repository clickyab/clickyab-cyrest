package main

import (
	"common/assert"
	"common/config"
	"common/initializer"
	"common/utils"
	"common/version"
	"modules/cyborg/worker"
	"net"
)

func main() {
	config.Initialize()
	config.InitApplication()
	config.Config.AMQP.Publisher = 1

	defer initializer.Initialize().Finalize()

	version.LogVersion().Infof("Application started")

	_, err := worker.NewMultiWorker(net.IPv4(127, 0, 0, 1), 9999)
	assert.Nil(err)

	utils.WaitExitSignal()

	version.LogVersion().Info("Goodbye")
}
