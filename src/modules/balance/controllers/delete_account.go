package balance

import (
	"modules/balance/acc"

	"github.com/gin-gonic/gin"
)

type deleteAccountResponse struct {
	Delete bool `json:"delete"`
}

// deleteAccount delete the account in system
// @Route {
// 		url = /
//		method = delete
//      #payload = addAccountPayload
//		middleware = accm.IsOwner
//      200 = deleteAccountResponse
//      400 = base.ErrorResponseSimple
// }
func (c *Controller) deleteAccount(ctx *gin.Context) {
	// TODO : do not delete the main account for a user (aka Wallet)
	account := c.MustGetAccount(ctx)
	usr := c.MustGetUser(ctx)

	m := acc.NewAccManager()
	f, err := m.DeleteAccount(usr, account)
	if err != nil {
		c.BadResponse(ctx, err)
		return
	}

	c.OKResponse(ctx, deleteAccountResponse{Delete: f})
}
