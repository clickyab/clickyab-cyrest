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

	// Route {/login POST Controller.loginUser user []  Controller u loginPayload } with key 0
	m0 := []echo.MiddlewareFunc{}

	// Make sure payload is the last middleware
	m0 = append(m0, middlewares.PayloadUnMarshallerGenerator(loginPayload{}))
	group.POST("/login", u.loginUser, m0...)
	// End route {/login POST Controller.loginUser user []  Controller u loginPayload } with key 0

	// Route {/register POST Controller.registerUser user []  Controller u registrationPayload } with key 1
	m1 := []echo.MiddlewareFunc{}

	// Make sure payload is the last middleware
	m1 = append(m1, middlewares.PayloadUnMarshallerGenerator(registrationPayload{}))
	group.POST("/register", u.registerUser, m1...)
	// End route {/register POST Controller.registerUser user []  Controller u registrationPayload } with key 1

	utils.DoInitialize(u)
}
