package ad

import "modules/misc/base"

// Controller is the controller for the ad package
// @Route {
//		group = /campaign
// }
type Controller struct {
	base.Controller
}

func init() {
	base.Register(&Controller{})
}
