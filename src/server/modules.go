package main

import (
	_ "modules/category" // category module
	"modules/category/cat"
	_ "modules/location"          // location module
	_ "modules/misc"              // misc controller
	_ "modules/telegram/channel"  // channel module
	_ "modules/telegram/plan"     // plan module
	_ "modules/telegram/teleuser" // teleuser module

	_ "modules/file"          // file module
	_ "modules/telegram/ad"   // ad module
	_ "modules/telegram/plan" // plan module
	_ "modules/user"          // user module
)

func init() {
	cat.RegisterScopes("channel")
}
