package user

// AUTO GENERATED CODE. DO NOT EDIT!

import (
	"common/middlewares"
	"common/utils"
	"modules/user/aaa"
	"modules/user/middlewares"

	"gopkg.in/labstack/echo.v3"
)

// Routes return the route registered with this
func (u *Controller) Routes(r *echo.Echo, mountPoint string) {

	groupMiddleware := []echo.MiddlewareFunc{}

	group := r.Group(mountPoint+"/user", groupMiddleware...)

	// Route {/forgot/call POST Controller.forgotPassword user []  Controller u forgotPayload  } with key 0
	m0 := []echo.MiddlewareFunc{}

	// Make sure payload is the last middleware
	m0 = append(m0, middlewares.PayloadUnMarshallerGenerator(forgotPayload{}))
	group.POST("/forgot/call", u.forgotPassword, m0...)
	// End route {/forgot/call POST Controller.forgotPassword user []  Controller u forgotPayload  } with key 0

	// Route {/forgot/callback/:code GET Controller.forgotGeneratePassword user []  Controller u   } with key 1
	m1 := []echo.MiddlewareFunc{}

	group.GET("/forgot/callback/:code", u.forgotGeneratePassword, m1...)
	// End route {/forgot/callback/:code GET Controller.forgotGeneratePassword user []  Controller u   } with key 1

	// Route {/login POST Controller.loginUser user []  Controller u loginPayload  } with key 2
	m2 := []echo.MiddlewareFunc{}

	// Make sure payload is the last middleware
	m2 = append(m2, middlewares.PayloadUnMarshallerGenerator(loginPayload{}))
	group.POST("/login", u.loginUser, m2...)
	// End route {/login POST Controller.loginUser user []  Controller u loginPayload  } with key 2

	// Route {/register POST Controller.registerUser user []  Controller u registrationPayload  } with key 3
	m3 := []echo.MiddlewareFunc{}

	// Make sure payload is the last middleware
	m3 = append(m3, middlewares.PayloadUnMarshallerGenerator(registrationPayload{}))
	group.POST("/register", u.registerUser, m3...)
	// End route {/register POST Controller.registerUser user []  Controller u registrationPayload  } with key 3

	// Route { GET Controller.listUser user [authz.Authenticate]  Controller u  user_list parent} with key 4
	m4 := []echo.MiddlewareFunc{
		authz.Authenticate,
	}

	m4 = append(m4, authz.AuthorizeGenerator("user_list", aaa.ScopePerm("parent")))

	group.GET("", u.listUser, m4...)
	// End route { GET Controller.listUser user [authz.Authenticate]  Controller u  user_list parent} with key 4

	utils.DoInitialize(u)
}
