package balance

import (
	"common/assert"
	"common/controllers/base"
	"modules/balance/acc"
	"modules/balance/middlewares"
	user "modules/user/controllers"

	"github.com/gin-gonic/gin"
)

// Controller is the base controller for balance module
// @Route {
//		group = /balance/accounts/:account
//		middleware = authz.Authenticate , accm.AccountCheck
//		_account_ = integer,the account id
// }
type Controller struct {
	user.BaseController
}

// MustGetUnit return the unit or panic if there is none
func (c *Controller) MustGetAccount(ctx *gin.Context) *acc.Account {
	u, ok := accm.GetAccount(ctx)
	assert.True(ok, "[BUG] account is not in the context")

	return u
}

func init() {
	base.Register(&Controller{})
}
