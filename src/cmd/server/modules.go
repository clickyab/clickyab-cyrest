package main

import (
	_ "modules/billing"           // billing module
	_ "modules/category"          // category module
	_ "modules/file"              // file module
	_ "modules/location"          // location module
	_ "modules/misc"              // misc controller
	_ "modules/telegram/ad"       // ad module
	_ "modules/telegram/teleuser" // teleuser module
	_ "modules/user"              // user module
)

//func init() {
//	cat.RegisterScopes("channel")
//}
