package balance

import (
	"modules/balance/acc"

	"github.com/gin-gonic/gin"
)

type addAccountPayload struct {
	theAccountPayload
	Initial acc.Money `json:"initial"`
}

// addAccount add new account to the system
// @Route {
// 		url = /add/account
//		method = post
//      payload = addAccountPayload
//      200 = addAccountResponse
//      400 = base.ErrorResponseSimple
// }
func (abc *AccountBaseController) addAccount(ctx *gin.Context) {
	// TODO : check account limit on creating new one
	payload := abc.MustGetPayload(ctx).(*addAccountPayload)
	usr := abc.MustGetUser(ctx)

	m := acc.NewAccManager()
	acc, err := m.AddAccount(
		usr,
		payload.Initial,
		payload.Title,
		payload.Description,
		payload.Disabled,
	)
	if err != nil {
		abc.BadResponse(ctx, err)
		return
	}

	abc.OKResponse(ctx, addAccountResponse{AccountID: acc.ID})
}
