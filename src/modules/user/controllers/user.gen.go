package user

// AUTO GENERATED CODE. DO NOT EDIT!

import (
	"common/middlewares"
	"common/utils"
	"modules/user/middlewares"

	"github.com/gin-gonic/gin"
)

// Routes return the route registered with this
func (u *Controller) Routes(r *gin.Engine, mountPoint string) {

	groupMiddleware := gin.HandlersChain{}

	group := r.Group(mountPoint+"/user", groupMiddleware...)

	// Route {/avatar/:user_id/:size/avatar.png GET Controller.getAvatar user []  Controller u  } with key 0
	m0 := gin.HandlersChain{}

	m0 = append(m0, u.getAvatar)
	group.GET("/avatar/:user_id/:size/avatar.png", m0...)
	// End route {/avatar/:user_id/:size/avatar.png GET Controller.getAvatar user []  Controller u  } with key 0

	// Route {/challenge POST Controller.challengeCreate user []  Controller u reserveUserPayload } with key 1
	m1 := gin.HandlersChain{}

	// Make sure payload is the last middleware
	m1 = append(m1, middlewares.PayloadUnMarshallerGenerator(reserveUserPayload{}))
	m1 = append(m1, u.challengeCreate)
	group.POST("/challenge", m1...)
	// End route {/challenge POST Controller.challengeCreate user []  Controller u reserveUserPayload } with key 1

	// Route {/state/:user_id POST Controller.changeState user [authz.Authenticate]  Controller u changeStatePayload user_admin} with key 2
	m2 := gin.HandlersChain{
		authz.Authenticate,
	}

	m2 = append(m2, authz.AuthorizeGenerator("user_admin"))

	// Make sure payload is the last middleware
	m2 = append(m2, middlewares.PayloadUnMarshallerGenerator(changeStatePayload{}))
	m2 = append(m2, u.changeState)
	group.POST("/state/:user_id", m2...)
	// End route {/state/:user_id POST Controller.changeState user [authz.Authenticate]  Controller u changeStatePayload user_admin} with key 2

	// Route {/role POST Controller.createRole user [authz.Authenticate]  Controller u createRolePayload user_admin} with key 3
	m3 := gin.HandlersChain{
		authz.Authenticate,
	}

	m3 = append(m3, authz.AuthorizeGenerator("user_admin"))

	// Make sure payload is the last middleware
	m3 = append(m3, middlewares.PayloadUnMarshallerGenerator(createRolePayload{}))
	m3 = append(m3, u.createRole)
	group.POST("/role", m3...)
	// End route {/role POST Controller.createRole user [authz.Authenticate]  Controller u createRolePayload user_admin} with key 3

	// Route {/role/:id PUT Controller.updateRole user [authz.Authenticate]  Controller u createRolePayload create_role} with key 4
	m4 := gin.HandlersChain{
		authz.Authenticate,
	}

	m4 = append(m4, authz.AuthorizeGenerator("create_role"))

	// Make sure payload is the last middleware
	m4 = append(m4, middlewares.PayloadUnMarshallerGenerator(createRolePayload{}))
	m4 = append(m4, u.updateRole)
	group.PUT("/role/:id", m4...)
	// End route {/role/:id PUT Controller.updateRole user [authz.Authenticate]  Controller u createRolePayload create_role} with key 4

	// Route {/roles GET Controller.listRoles user [authz.Authenticate]  Controller u  } with key 5
	m5 := gin.HandlersChain{
		authz.Authenticate,
	}

	m5 = append(m5, u.listRoles)
	group.GET("/roles", m5...)
	// End route {/roles GET Controller.listRoles user [authz.Authenticate]  Controller u  } with key 5

	// Route {/users GET Controller.listUsers user [authz.Authenticate]  Controller u  user_admin} with key 6
	m6 := gin.HandlersChain{
		authz.Authenticate,
	}

	m6 = append(m6, authz.AuthorizeGenerator("user_admin"))

	m6 = append(m6, u.listUsers)
	group.GET("/users", m6...)
	// End route {/users GET Controller.listUsers user [authz.Authenticate]  Controller u  user_admin} with key 6

	// Route {/login POST Controller.login user []  Controller u payloadLoginData } with key 7
	m7 := gin.HandlersChain{}

	// Make sure payload is the last middleware
	m7 = append(m7, middlewares.PayloadUnMarshallerGenerator(payloadLoginData{}))
	m7 = append(m7, u.login)
	group.POST("/login", m7...)
	// End route {/login POST Controller.login user []  Controller u payloadLoginData } with key 7

	// Route {/logout GET Controller.logout user [authz.Authenticate]  Controller u  } with key 8
	m8 := gin.HandlersChain{
		authz.Authenticate,
	}

	m8 = append(m8, u.logout)
	group.GET("/logout", m8...)
	// End route {/logout GET Controller.logout user [authz.Authenticate]  Controller u  } with key 8

	// Route {/register POST Controller.registerUser user []  Controller u registrationPayload } with key 9
	m9 := gin.HandlersChain{}

	// Make sure payload is the last middleware
	m9 = append(m9, middlewares.PayloadUnMarshallerGenerator(registrationPayload{}))
	m9 = append(m9, u.registerUser)
	group.POST("/register", m9...)
	// End route {/register POST Controller.registerUser user []  Controller u registrationPayload } with key 9

	// Route {/role/:id DELETE Controller.removeRole user [authz.Authenticate]  Controller u  create_role} with key 10
	m10 := gin.HandlersChain{
		authz.Authenticate,
	}

	m10 = append(m10, authz.AuthorizeGenerator("create_role"))

	m10 = append(m10, u.removeRole)
	group.DELETE("/role/:id", m10...)
	// End route {/role/:id DELETE Controller.removeRole user [authz.Authenticate]  Controller u  create_role} with key 10

	// Route {/test POST Controller.testFunction user []  Controller u tmp } with key 11
	m11 := gin.HandlersChain{}

	// Make sure payload is the last middleware
	m11 = append(m11, middlewares.PayloadUnMarshallerGenerator(tmp{}))
	m11 = append(m11, u.testFunction)
	group.POST("/test", m11...)
	// End route {/test POST Controller.testFunction user []  Controller u tmp } with key 11

	utils.DoInitialize(u)
}
