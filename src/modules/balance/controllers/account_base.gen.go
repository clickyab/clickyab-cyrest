package balance

// AUTO GENERATED CODE. DO NOT EDIT!

import (
	"common/middlewares"
	"common/utils"
	"modules/user/middlewares"

	"github.com/gin-gonic/gin"
)

// Routes return the route registered with this
func (abc *AccountBaseController) Routes(r *gin.Engine, mountPoint string) {

	groupMiddleware := gin.HandlersChain{
		authz.Authenticate,
	}

	group := r.Group(mountPoint+"/balance", groupMiddleware...)

	// Route {/add/account POST AccountBaseController.addAccount balance []  AccountBaseController abc addAccountPayload } with key 0
	m0 := gin.HandlersChain{}

	// Make sure payload is the last middleware
	m0 = append(m0, middlewares.PayloadUnMarshallerGenerator(addAccountPayload{}))
	m0 = append(m0, abc.addAccount)
	group.POST("/add/account", m0...)
	// End route {/add/account POST AccountBaseController.addAccount balance []  AccountBaseController abc addAccountPayload } with key 0

	// Route {/list/accounts GET AccountBaseController.listAccounts balance []  AccountBaseController abc  } with key 1
	m1 := gin.HandlersChain{}

	m1 = append(m1, abc.listAccounts)
	group.GET("/list/accounts", m1...)
	// End route {/list/accounts GET AccountBaseController.listAccounts balance []  AccountBaseController abc  } with key 1

	// Route {/list/transactions GET AccountBaseController.listTransactions balance []  AccountBaseController abc  } with key 2
	m2 := gin.HandlersChain{}

	m2 = append(m2, abc.listTransactions)
	group.GET("/list/transactions", m2...)
	// End route {/list/transactions GET AccountBaseController.listTransactions balance []  AccountBaseController abc  } with key 2

	utils.DoInitialize(abc)
}
