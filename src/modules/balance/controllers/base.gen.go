package balance

// AUTO GENERATED CODE. DO NOT EDIT!

import (
	"common/middlewares"
	"common/utils"
	"modules/balance/middlewares"
	"modules/user/middlewares"

	"github.com/gin-gonic/gin"
)

// Routes return the route registered with this
func (c *Controller) Routes(r *gin.Engine, mountPoint string) {

	groupMiddleware := gin.HandlersChain{
		authz.Authenticate,
		accm.AccountCheck,
	}

	group := r.Group(mountPoint+"/balance/accounts/:account", groupMiddleware...)

	// Route {/transactions POST Controller.addTransaction balance []  Controller u transactionMultiPayload } with key 0
	m0 := gin.HandlersChain{}

	// Make sure payload is the last middleware
	m0 = append(m0, middlewares.PayloadUnMarshallerGenerator(transactionMultiPayload{}))
	m0 = append(m0, c.addTransaction)
	group.POST("/transactions", m0...)
	// End route {/transactions POST Controller.addTransaction balance []  Controller u transactionMultiPayload } with key 0

	// Route {/ DELETE Controller.deleteAccount balance [accm.IsOwner]  Controller c  } with key 1
	m1 := gin.HandlersChain{
		accm.IsOwner,
	}

	m1 = append(m1, c.deleteAccount)
	group.DELETE("/", m1...)
	// End route {/ DELETE Controller.deleteAccount balance [accm.IsOwner]  Controller c  } with key 1

	// Route {/ PUT Controller.editAccount balance [accm.IsOwner]  Controller c theAccountPayload } with key 2
	m2 := gin.HandlersChain{
		accm.IsOwner,
	}

	// Make sure payload is the last middleware
	m2 = append(m2, middlewares.PayloadUnMarshallerGenerator(theAccountPayload{}))
	m2 = append(m2, c.editAccount)
	group.PUT("/", m2...)
	// End route {/ PUT Controller.editAccount balance [accm.IsOwner]  Controller c theAccountPayload } with key 2

	// Route {/ GET Controller.getAccount balance [accm.AccountCheck]  Controller c  } with key 3
	m3 := gin.HandlersChain{
		accm.AccountCheck,
	}

	m3 = append(m3, c.getAccount)
	group.GET("/", m3...)
	// End route {/ GET Controller.getAccount balance [accm.AccountCheck]  Controller c  } with key 3

	// Route {/transaction/:transaction_id GET Controller.getTransaction balance []  Controller c  } with key 4
	m4 := gin.HandlersChain{}

	m4 = append(m4, c.getTransaction)
	group.GET("/transaction/:transaction_id", m4...)
	// End route {/transaction/:transaction_id GET Controller.getTransaction balance []  Controller c  } with key 4

	utils.DoInitialize(c)
}
