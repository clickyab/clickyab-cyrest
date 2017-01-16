package telegram

import "modules/misc/base"

// Controller is the controller for the teleuser package
// @Route {
//		group = /telegram
// }
type Controller struct {
	base.Controller
}

func init() {
	base.Register(&Controller{})
}
