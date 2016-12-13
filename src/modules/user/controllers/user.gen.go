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

	// Route {/assign/parent POST Controller.assignParent user [authz.Authenticate]  Controller u assignParentPayload assign_parent global} with key 0
	m0 := []echo.MiddlewareFunc{
		authz.Authenticate,
	}

	m0 = append(m0, authz.AuthorizeGenerator("assign_parent", base.UserScope("global")))

	// Make sure payload is the last middleware
	m0 = append(m0, middlewares.PayloadUnMarshallerGenerator(assignParentPayload{}))
	group.POST("/assign/parent", u.assignParent, m0...)
	// End route {/assign/parent POST Controller.assignParent user [authz.Authenticate]  Controller u assignParentPayload assign_parent global} with key 0

	// Route {/change-password POST Controller.changePassword user [authz.Authenticate]  Controller u changePasswordPayload  } with key 1
	m1 := []echo.MiddlewareFunc{
		authz.Authenticate,
	}

	// Make sure payload is the last middleware
	m1 = append(m1, middlewares.PayloadUnMarshallerGenerator(changePasswordPayload{}))
	group.POST("/change-password", u.changePassword, m1...)
	// End route {/change-password POST Controller.changePassword user [authz.Authenticate]  Controller u changePasswordPayload  } with key 1

	// Route {/forgot/call POST Controller.forgotPassword user []  Controller u forgotPayload  } with key 2
	m2 := []echo.MiddlewareFunc{}

	// Make sure payload is the last middleware
	m2 = append(m2, middlewares.PayloadUnMarshallerGenerator(forgotPayload{}))
	group.POST("/forgot/call", u.forgotPassword, m2...)
	// End route {/forgot/call POST Controller.forgotPassword user []  Controller u forgotPayload  } with key 2

	// Route {/forgot/callback/:code GET Controller.forgotGeneratePassword user []  Controller u   } with key 3
	m3 := []echo.MiddlewareFunc{}

	group.GET("/forgot/callback/:code", u.forgotGeneratePassword, m3...)
	// End route {/forgot/callback/:code GET Controller.forgotGeneratePassword user []  Controller u   } with key 3

	// Route {/login POST Controller.loginUser user []  Controller u loginPayload  } with key 4
	m4 := []echo.MiddlewareFunc{}

	// Make sure payload is the last middleware
	m4 = append(m4, middlewares.PayloadUnMarshallerGenerator(loginPayload{}))
	group.POST("/login", u.loginUser, m4...)
	// End route {/login POST Controller.loginUser user []  Controller u loginPayload  } with key 4

	// Route {/authenticate/:action GET Controller.oauthInit user []  Controller u   } with key 5
	m5 := []echo.MiddlewareFunc{}

	group.GET("/authenticate/:action", u.oauthInit, m5...)
	// End route {/authenticate/:action GET Controller.oauthInit user []  Controller u   } with key 5

	// Route {/oauth/callback GET Controller.oauthCallback user []  Controller u   } with key 6
	m6 := []echo.MiddlewareFunc{}

	group.GET("/oauth/callback", u.oauthCallback, m6...)
	// End route {/oauth/callback GET Controller.oauthCallback user []  Controller u   } with key 6

	// Route {/register POST Controller.registerUser user []  Controller u registrationPayload  } with key 7
	m7 := []echo.MiddlewareFunc{}

	// Make sure payload is the last middleware
	m7 = append(m7, middlewares.PayloadUnMarshallerGenerator(registrationPayload{}))
	group.POST("/register", u.registerUser, m7...)
	// End route {/register POST Controller.registerUser user []  Controller u registrationPayload  } with key 7

	// Route {/logout GET Controller.logout user [authz.Authenticate]  Controller u   } with key 8
	m8 := []echo.MiddlewareFunc{
		authz.Authenticate,
	}

	group.GET("/logout", u.logout, m8...)
	// End route {/logout GET Controller.logout user [authz.Authenticate]  Controller u   } with key 8

	// Route {/sessions GET Controller.activeSessions user [authz.Authenticate]  Controller u   } with key 9
	m9 := []echo.MiddlewareFunc{
		authz.Authenticate,
	}

	group.GET("/sessions", u.activeSessions, m9...)
	// End route {/sessions GET Controller.activeSessions user [authz.Authenticate]  Controller u   } with key 9

	// Route {/session/terminate/:id GET Controller.terminateSession user [authz.Authenticate]  Controller u   } with key 10
	m10 := []echo.MiddlewareFunc{
		authz.Authenticate,
	}

	group.GET("/session/terminate/:id", u.terminateSession, m10...)
	// End route {/session/terminate/:id GET Controller.terminateSession user [authz.Authenticate]  Controller u   } with key 10

	// Route {/sessions/terminate GET Controller.terminateAllSession user [authz.Authenticate]  Controller u   } with key 11
	m11 := []echo.MiddlewareFunc{
		authz.Authenticate,
	}

	group.GET("/sessions/terminate", u.terminateAllSession, m11...)
	// End route {/sessions/terminate GET Controller.terminateAllSession user [authz.Authenticate]  Controller u   } with key 11

	// Route {/users GET Controller.listUser user [authz.Authenticate]  Controller u  user_list parent} with key 12
	m12 := []echo.MiddlewareFunc{
		authz.Authenticate,
	}

	m12 = append(m12, authz.AuthorizeGenerator("user_list", base.UserScope("parent")))

	group.GET("/users", u.listUser, m12...)
	// End route {/users GET Controller.listUser user [authz.Authenticate]  Controller u  user_list parent} with key 12

	utils.DoInitialize(u)
}
