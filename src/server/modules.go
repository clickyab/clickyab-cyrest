package main

import (
	_ "modules/category" // category module
	"modules/category/cat"
	_ "modules/channel"  // channel module
	_ "modules/location" // location module
	_ "modules/misc"     // misc controller
	_ "modules/plan"     // plan module
	_ "modules/teleuser" // teleuser module

	_ "modules/ad"   // ad module
	_ "modules/file" // file module
	_ "modules/plan" // plan module
	_ "modules/user" // user module
	_ "modules/plan" // plan module

)

func init() {
	cat.RegisterScopes("channel")
}
