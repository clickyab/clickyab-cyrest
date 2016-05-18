package try

import "github.com/fzerorubigd/block/catch"

var try = catch.New()

// Try check error in multiple function and return the final transformed erro
func Try(err error) error {
	return try.Try(err)
}

// Catch register a global catcher
func CatchHook(f interface{}) {
	try.Catch(f)
}

// FinallyHook add a finally call
func FinallyHook(f interface{}) {
	try.Finally(f)
}
