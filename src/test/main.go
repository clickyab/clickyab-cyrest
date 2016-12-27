package main

import (
	"common/assert"
	"common/tgo"
	"net"
)

func main() {
	_, err := tgo.NewTelegramCli(net.IPv4(172, 17, 0, 1), 9010)
assert.Nil(err)

}
