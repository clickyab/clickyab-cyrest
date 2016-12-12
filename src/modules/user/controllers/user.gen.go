package user

// AUTO GENERATED CODE. DO NOT EDIT!

import (
	"common/controllers/base"
	"common/middlewares"
	"common/utils"
	"modules/user/middlewares"

	"gopkg.in/labstack/echo.v3"
)

// Routes return the route registered with this
func (u *Controller) Routes(r *echo.Echo, mountPoint string) {

	groupMiddleware := []echo.MiddlewareFunc{}

	group := r.Group(mountPoint+"/user", groupMiddleware...)

	// Route {/change-password POST Controller.changePassword user [authz.Authenticate]  Controller u changePasswordPayload  } with key 0
	m0 := []echo.MiddlewareFunc{
		authz.Authenticate,
	}

	// Make sure payload is the last middleware
	m0 = append(m0, middlewares.PayloadUnMarshallerGenerator(changePasswordPayload{}))
	group.POST("/change-password", u.changePassword, m0...)
	// End route {/change-password POST Controller.changePassword user [authz.Authenticate]  Controller u changePasswordPayload  } with key 0

	// Route {/forgot/call POST Controller.forgotPassword user []  Controller u forgotPayload  } with key 1
	m1 := []echo.MiddlewareFunc{}

	// Make sure payload is the last middleware
	m1 = append(m1, middlewares.PayloadUnMarshallerGenerator(forgotPayload{}))
	group.POST("/forgot/call", u.forgotPassword, m1...)
	// End route {/forgot/call POST Controller.forgotPassword user []  Controller u forgotPayload  } with key 1

	// Route {/forgot/callback/:code GET Controller.forgotGeneratePassword user []  Controller u   } with key 2
	m2 := []echo.MiddlewareFunc{}

	group.GET("/forgot/callback/:code", u.forgotGeneratePassword, m2...)
	// End route {/forgot/callback/:code GET Controller.forgotGeneratePassword user []  Controller u   } with key 2

	// Route {/login POST Controller.loginUser user []  Controller u loginPayload  } with key 3
	m3 := []echo.MiddlewareFunc{}

	// Make sure payload is the last middleware
	m3 = append(m3, middlewares.PayloadUnMarshallerGenerator(loginPayload{}))
	group.POST("/login", u.loginUser, m3...)
	// End route {/login POST Controller.loginUser user []  Controller u loginPayload  } with key 3

	// Route {/register POST Controller.registerUser user []  Controller u registrationPayload  } with key 4
	m4 := []echo.MiddlewareFunc{}

	// Make sure payload is the last middleware
	m4 = append(m4, middlewares.PayloadUnMarshallerGenerator(registrationPayload{}))
	group.POST("/register", u.registerUser, m4...)
	// End route {/register POST Controller.registerUser user []  Controller u registrationPayload  } with key 4

	// Route {/logout GET Controller.logout user [authz.Authenticate]  Controller u   } with key 5
	m5 := []echo.MiddlewareFunc{
		authz.Authenticate,
	}

	group.GET("/logout", u.logout, m5...)
	// End route {/logout GET Controller.logout user [authz.Authenticate]  Controller u   } with key 5

	// Route {/sessions GET Controller.activeSessions user [authz.Authenticate]  Controller u   } with key 6
	m6 := []echo.MiddlewareFunc{
		authz.Authenticate,
	}

	group.GET("/sessions", u.activeSessions, m6...)
	// End route {/sessions GET Controller.activeSessions user [authz.Authenticate]  Controller u   } with key 6

	// Route {/session/terminate/:id GET Controller.terminateSession user [authz.Authenticate]  Controller u   } with key 7
	m7 := []echo.MiddlewareFunc{
		authz.Authenticate,
	}

	group.GET("/session/terminate/:id", u.terminateSession, m7...)
	// End route {/session/terminate/:id GET Controller.terminateSession user [authz.Authenticate]  Controller u   } with key 7

	// Route {/sessions/terminate GET Controller.terminateAllSession user [authz.Authenticate]  Controller u   } with key 8
	m8 := []echo.MiddlewareFunc{
		authz.Authenticate,
	}

	group.GET("/sessions/terminate", u.terminateAllSession, m8...)
	// End route {/sessions/terminate GET Controller.terminateAllSession user [authz.Authenticate]  Controller u   } with key 8

	// Route { GET Controller.listUser user [authz.Authenticate]  Controller u  user_list parent} with key 9
	m9 := []echo.MiddlewareFunc{
		authz.Authenticate,
	}

	m9 = append(m9, authz.AuthorizeGenerator("user_list", base.UserScope("parent")))

	group.GET("", u.listUser, m9...)
	// End route { GET Controller.listUser user [authz.Authenticate]  Controller u  user_list parent} with key 9

	utils.DoInitialize(u)
}
