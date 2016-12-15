package main

import (
	_ "modules/category" // category module
	"modules/category/cat"
	_ "modules/misc" // misc controller
	_ "modules/user" // user module
)

func init() {
	cat.RegisterScopes("channel")
}
