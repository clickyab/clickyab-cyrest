package user

import "common/controllers/base"

type BaseController struct {
	base.Controller
}

func init() {
	base.Register(&Controller{})
}
