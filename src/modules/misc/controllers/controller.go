package misc

import "common/controllers/base"

// Controller is the misc controller
// @Route {
//		group = /misc
// }
type Controller struct {
	base.Controller
}

// Initialize is the initialization method for ths controller
// TODO : the trans module is somehow dangled in the system find its place!
func (u *Controller) Initialize() {
	//try.CatchHook(func(e error) error {
	//	return fmt.Errorf(trans.T(e.Error()))
	//})
}

func init() {
	base.Register(&Controller{})
}
