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
	ips, err := net.LookupIP(config.Config.Telegram.CLIAddress)
	assert.Nil(err)
	var (
		ip    net.IP
		found bool
	)
	for i := range ips {
		if ips[i].To4() != nil {
			ip = ips[i]
			found = true
			break
		}
	}
	assert.True(found, "no ip found")
	_, err = worker.NewMultiWorker(ip, config.Config.Telegram.CLIPort)
	assert.Nil(err)

	utils.WaitExitSignal()

	version.LogVersion().Info("Goodbye")
}
