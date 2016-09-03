package user

// AUTO GENERATED CODE. DO NOT EDIT!

import (
	"common/middlewares"
	"common/utils"

	"github.com/labstack/echo"
)

// Routes return the route registered with this
func (u *Controller) Routes(r *echo.Echo, mountPoint string) {

	groupMiddleware := []echo.MiddlewareFunc{}

	group := r.Group(mountPoint+"/user", groupMiddleware...)

	// Route {/register POST Controller.registerUser user []  Controller u registrationPayload } with key 0
	m0 := []echo.MiddlewareFunc{}

	// Make sure payload is the last middleware
	m0 = append(m0, middlewares.PayloadUnMarshallerGenerator(registrationPayload{}))
	group.POST("/register", u.registerUser, m0...)
	// End route {/register POST Controller.registerUser user []  Controller u registrationPayload } with key 0

	utils.DoInitialize(u)
}
