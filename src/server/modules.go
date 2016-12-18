package main

import (
	_ "modules/category" // category module
	"modules/category/cat"
	_ "modules/channel" // channel module
	_ "modules/misc"    // misc controller
	_ "modules/user"    // user module
	_ "modules/location" // location module
)

func init() {
	cat.RegisterScopes("channel")
}
