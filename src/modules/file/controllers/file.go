package file

import "modules/misc/base"

// Controller is the controller for the file package
// @Route {
//		group = /file
// }
type Controller struct {
	base.Controller
}

func init() {
	base.Register(&Controller{})
}
