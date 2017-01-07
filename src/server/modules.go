package main

import (
	_ "modules/campaign" // campaign module
	_ "modules/category" // category module
	"modules/category/cat"
	_ "modules/channel"  // channel module
	_ "modules/location" // location module
	_ "modules/misc"     // misc controller

	_ "modules/user" // user module
	_ "modules/ad" // ad module
	_ "modules/plan" // plan module
)

func init() {
	cat.RegisterScopes("channel")
}
