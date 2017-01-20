package channel

import "modules/misc/base"

// Controller is the controller for the channel package
// @Route {
//		group = /channel
// }
type Controller struct {
	base.Controller
}

func init() {
	base.Register(&Controller{})
}
