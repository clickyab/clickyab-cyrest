package plan

import "modules/misc/base"

// Controller is the controller for the plan package
// @Route {
//		group = /plan
// }
type Controller struct {
	base.Controller
}

func init() {
	base.Register(&Controller{})
}
