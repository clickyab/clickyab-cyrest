package misc

// AUTO GENERATED CODE. DO NOT EDIT!

import (
	"common/utils"

	"gopkg.in/labstack/echo.v3"
)

// Routes return the route registered with this
func (u *Controller) Routes(r *echo.Echo, mountPoint string) {

	groupMiddleware := []echo.MiddlewareFunc{}

	group := r.Group(mountPoint+"/misc", groupMiddleware...)

	// Route {/version GET Controller.getVersion misc []  Controller u   } with key 0
	m0 := []echo.MiddlewareFunc{}

	group.GET("/version", u.getVersion, m0...)
	// End route {/version GET Controller.getVersion misc []  Controller u   } with key 0

	utils.DoInitialize(u)
}
