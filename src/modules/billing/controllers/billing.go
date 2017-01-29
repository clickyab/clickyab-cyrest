package billing

import "modules/misc/base"

// Controller is the controller for the billing package
// @Route {
//		group = /billing
// }
type Controller struct {
	base.Controller
}

func init() {
	base.Register(&Controller{})
}
