package balance

import (
	"common/utils"
	"modules/balance/acc"

	"github.com/gin-gonic/gin"
)

type listAccountResponse struct {
	List  []acc.Account `json:"list"`
	Total int64         `json:"total"`
}

// listAccounts list accounts in the units
// @Route {
// 		url = /list/accounts
//		method = get
//		_all_ = integer,show all data with no pagination
//		_dis_ = integer,show disabled account too
//		_p_ = integer, the page number
//		_c_ = integer, per page item
//      200 = listAccountResponse
//      400 = base.ErrorResponseSimple
// }
func (abc *AccountBaseController) listAccounts(ctx *gin.Context) {
	usr := abc.MustGetUser(ctx)
	offset, perPage := utils.GetPageAndCount(ctx.Request, true)
	all := ctx.Query("all") != ""
	dis := ctx.Query("dis") != ""

	list, cnt := acc.NewAccManager().ListUserAccounts(usr, dis, all, offset, perPage)
	abc.OKResponse(ctx, listAccountResponse{
		List:  list,
		Total: cnt,
	})
}
