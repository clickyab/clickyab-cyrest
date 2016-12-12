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

	// Route {/logout GET Controller.logout user [authz.Authenticate]  Controller u   } with key 4
	m4 := []echo.MiddlewareFunc{
		authz.Authenticate,
	}

	group.GET("/logout", u.logout, m4...)
	// End route {/logout GET Controller.logout user [authz.Authenticate]  Controller u   } with key 4

	// Route {/sessions GET Controller.activeSessions user [authz.Authenticate]  Controller u   } with key 5
	m5 := []echo.MiddlewareFunc{
		authz.Authenticate,
	}

	group.GET("/sessions", u.activeSessions, m5...)
	// End route {/sessions GET Controller.activeSessions user [authz.Authenticate]  Controller u   } with key 5

	// Route {/session/terminate/:id GET Controller.terminateSession user [authz.Authenticate]  Controller u   } with key 6
	m6 := []echo.MiddlewareFunc{
		authz.Authenticate,
	}

	group.GET("/session/terminate/:id", u.terminateSession, m6...)
	// End route {/session/terminate/:id GET Controller.terminateSession user [authz.Authenticate]  Controller u   } with key 6

	// Route {/sessions/terminate GET Controller.terminateAllSession user [authz.Authenticate]  Controller u   } with key 7
	m7 := []echo.MiddlewareFunc{
		authz.Authenticate,
	}

	group.GET("/sessions/terminate", u.terminateAllSession, m7...)
	// End route {/sessions/terminate GET Controller.terminateAllSession user [authz.Authenticate]  Controller u   } with key 7

	// Route {/users GET Controller.listUser user [authz.Authenticate]  Controller u  user_list parent} with key 8
	m8 := []echo.MiddlewareFunc{
		authz.Authenticate,
	}

	m8 = append(m8, authz.AuthorizeGenerator("user_list", base.UserScope("parent")))

	group.GET("/users", u.listUser, m8...)
	// End route {/users GET Controller.listUser user [authz.Authenticate]  Controller u  user_list parent} with key 8

	utils.DoInitialize(u)
}
