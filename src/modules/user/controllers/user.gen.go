package user

// AUTO GENERATED CODE. DO NOT EDIT!

import (
	"common/middlewares"
	"common/utils"
	"modules/user/middlewares"

	"github.com/labstack/echo"
)

// Routes return the route registered with this
func (u *Controller) Routes(r *echo.Echo, mountPoint string) {

	groupMiddleware := []echo.MiddlewareFunc{}

	group := r.Group(mountPoint+"/user", groupMiddleware...)

	// Route {/avatar/:user_id/:size/avatar.png GET Controller.getAvatar user []  Controller u  } with key 0
	m0 := []echo.MiddlewareFunc{}

	group.GET("/avatar/:user_id/:size/avatar.png", u.getAvatar, m0...)
	// End route {/avatar/:user_id/:size/avatar.png GET Controller.getAvatar user []  Controller u  } with key 0

	// Route {/challenge POST Controller.challengeCreate user []  Controller u reserveUserPayload } with key 1
	m1 := []echo.MiddlewareFunc{}

	// Make sure payload is the last middleware
	m1 = append(m1, middlewares.PayloadUnMarshallerGenerator(reserveUserPayload{}))
	group.POST("/challenge", u.challengeCreate, m1...)
	// End route {/challenge POST Controller.challengeCreate user []  Controller u reserveUserPayload } with key 1

	// Route {/password POST Controller.changePassword user [authz.Authenticate]  Controller u changePasswordPayload } with key 2
	m2 := []echo.MiddlewareFunc{
		authz.Authenticate,
	}

	// Make sure payload is the last middleware
	m2 = append(m2, middlewares.PayloadUnMarshallerGenerator(changePasswordPayload{}))
	group.POST("/password", u.changePassword, m2...)
	// End route {/password POST Controller.changePassword user [authz.Authenticate]  Controller u changePasswordPayload } with key 2

	// Route {/state/:user_id POST Controller.changeState user [authz.Authenticate]  Controller u changeStatePayload user_admin} with key 3
	m3 := []echo.MiddlewareFunc{
		authz.Authenticate,
	}

	m3 = append(m3, authz.AuthorizeGenerator("user_admin"))

	// Make sure payload is the last middleware
	m3 = append(m3, middlewares.PayloadUnMarshallerGenerator(changeStatePayload{}))
	group.POST("/state/:user_id", u.changeState, m3...)
	// End route {/state/:user_id POST Controller.changeState user [authz.Authenticate]  Controller u changeStatePayload user_admin} with key 3

	// Route {/role POST Controller.createRole user [authz.Authenticate]  Controller u createRolePayload user_admin} with key 4
	m4 := []echo.MiddlewareFunc{
		authz.Authenticate,
	}

	m4 = append(m4, authz.AuthorizeGenerator("user_admin"))

	// Make sure payload is the last middleware
	m4 = append(m4, middlewares.PayloadUnMarshallerGenerator(createRolePayload{}))
	group.POST("/role", u.createRole, m4...)
	// End route {/role POST Controller.createRole user [authz.Authenticate]  Controller u createRolePayload user_admin} with key 4

	// Route {/role/:id PUT Controller.updateRole user [authz.Authenticate]  Controller u createRolePayload create_role} with key 5
	m5 := []echo.MiddlewareFunc{
		authz.Authenticate,
	}

	m5 = append(m5, authz.AuthorizeGenerator("create_role"))

	// Make sure payload is the last middleware
	m5 = append(m5, middlewares.PayloadUnMarshallerGenerator(createRolePayload{}))
	group.PUT("/role/:id", u.updateRole, m5...)
	// End route {/role/:id PUT Controller.updateRole user [authz.Authenticate]  Controller u createRolePayload create_role} with key 5

	// Route {/roles GET Controller.listRoles user [authz.Authenticate]  Controller u  } with key 6
	m6 := []echo.MiddlewareFunc{
		authz.Authenticate,
	}

	group.GET("/roles", u.listRoles, m6...)
	// End route {/roles GET Controller.listRoles user [authz.Authenticate]  Controller u  } with key 6

	// Route {/users GET Controller.listUsers user [authz.Authenticate]  Controller u  user_admin} with key 7
	m7 := []echo.MiddlewareFunc{
		authz.Authenticate,
	}

	m7 = append(m7, authz.AuthorizeGenerator("user_admin"))

	group.GET("/users", u.listUsers, m7...)
	// End route {/users GET Controller.listUsers user [authz.Authenticate]  Controller u  user_admin} with key 7

	// Route {/login POST Controller.login user []  Controller u payloadLoginData } with key 8
	m8 := []echo.MiddlewareFunc{}

	// Make sure payload is the last middleware
	m8 = append(m8, middlewares.PayloadUnMarshallerGenerator(payloadLoginData{}))
	group.POST("/login", u.login, m8...)
	// End route {/login POST Controller.login user []  Controller u payloadLoginData } with key 8

	// Route {/logout GET Controller.logout user [authz.Authenticate]  Controller u  } with key 9
	m9 := []echo.MiddlewareFunc{
		authz.Authenticate,
	}

	group.GET("/logout", u.logout, m9...)
	// End route {/logout GET Controller.logout user [authz.Authenticate]  Controller u  } with key 9

	// Route {/authenticate/:action GET Controller.oauthInit user []  Controller u  } with key 10
	m10 := []echo.MiddlewareFunc{}

	group.GET("/authenticate/:action", u.oauthInit, m10...)
	// End route {/authenticate/:action GET Controller.oauthInit user []  Controller u  } with key 10

	// Route {/oauth/callback GET Controller.oauthCallback user []  Controller u  } with key 11
	m11 := []echo.MiddlewareFunc{}

	group.GET("/oauth/callback", u.oauthCallback, m11...)
	// End route {/oauth/callback GET Controller.oauthCallback user []  Controller u  } with key 11

	// Route {/register POST Controller.registerUser user []  Controller u registrationPayload } with key 12
	m12 := []echo.MiddlewareFunc{}

	// Make sure payload is the last middleware
	m12 = append(m12, middlewares.PayloadUnMarshallerGenerator(registrationPayload{}))
	group.POST("/register", u.registerUser, m12...)
	// End route {/register POST Controller.registerUser user []  Controller u registrationPayload } with key 12

	// Route {/role/:id DELETE Controller.removeRole user [authz.Authenticate]  Controller u  create_role} with key 13
	m13 := []echo.MiddlewareFunc{
		authz.Authenticate,
	}

	m13 = append(m13, authz.AuthorizeGenerator("create_role"))

	group.DELETE("/role/:id", u.removeRole, m13...)
	// End route {/role/:id DELETE Controller.removeRole user [authz.Authenticate]  Controller u  create_role} with key 13

	// Route {/test POST Controller.testFunction user []  Controller u tmp } with key 14
	m14 := []echo.MiddlewareFunc{}

	// Make sure payload is the last middleware
	m14 = append(m14, middlewares.PayloadUnMarshallerGenerator(tmp{}))
	group.POST("/test", u.testFunction, m14...)
	// End route {/test POST Controller.testFunction user []  Controller u tmp } with key 14

	utils.DoInitialize(u)
}
