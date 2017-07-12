package bundle

import (
	"modules/misc/base"

	"gopkg.in/labstack/echo.v3"
)

// Controller is the controller for the bundle package
// @Route {
//		group = /bundle
// }
type Controller struct {
	base.Controller
}

func (*Controller) Routes(r *echo.Echo, mountPoint string) {
	panic("implement me")
}

func init() {
	base.Register(&Controller{})
}
