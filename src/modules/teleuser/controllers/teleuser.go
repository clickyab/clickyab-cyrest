package teleuser

import "modules/misc/base"

// Controller is the controller for the teleuser package
// @Route {
//		group = /teleuser
// }
type Controller struct {
	base.Controller
}

func init() {
	base.Register(&Controller{})
}
