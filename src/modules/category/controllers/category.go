package controllers

import "common/controllers/base"

// Controller is the controller for the category package
// @Route {
//		group = /category
// }
type Controller struct {
	base.Controller
}

func init() {
	base.Register(&Controller{})
}
