package location

import "modules/misc/base"

// Controller is the controller for the country package
// @Route {
//		group = /location
// }
type Controller struct {
	base.Controller
}

func init() {
	base.Register(&Controller{})
}
