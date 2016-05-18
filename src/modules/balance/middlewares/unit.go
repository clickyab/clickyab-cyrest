package accm

import (
	"common/controllers/base"
	"modules/balance/acc"
	"modules/user/middlewares"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	// ContextUnit is for user unmarshal object
	ContextUnit string = "_acc_account"
)

// AccountCheck is a check for unit in the process, and check if the current user has access to
// this must be added with authz.Authenticate middleware
// Route {
//		403 = base.ErrorResponseSimple
// }
func AccountCheck(c *gin.Context) {
	u, ok := authz.GetUser(c)
	if !ok {
		c.JSON(
			http.StatusForbidden,
			base.ErrorResponseSimple{Error: http.StatusText(http.StatusForbidden)},
		)
		c.Abort()
		return
	}
	accountID, err := strconv.ParseInt(c.Param("account"), 10, 0)
	if err != nil {
		c.JSON(
			http.StatusNotFound,
			base.ErrorResponseSimple{Error: err.Error()},
		)
		c.Abort()
		return
	}

	m := acc.NewAccManager()
	account, err := m.FindAccountByID(accountID)
	if err != nil {
		c.JSON(
			http.StatusNotFound,
			base.ErrorResponseSimple{Error: err.Error()},
		)
		c.Abort()
		return
	}

	// TODO : add support for access check for other user
	if account.OwnerID != u.ID {
		c.JSON(
			http.StatusForbidden, // User is logged in but have no access to this
			base.ErrorResponseSimple{Error: http.StatusText(http.StatusForbidden)},
		)
		c.Abort()
		return
	}

	c.Set(ContextUnit, account)
}

// GetAccount from the request
func GetAccount(c *gin.Context) (*acc.Account, bool) {
	if u, ok := c.Get(ContextUnit); ok {
		user, ok := u.(*acc.Account)
		return user, ok
	}

	return nil, false
}
