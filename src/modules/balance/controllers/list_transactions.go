package balance

import (
	"common/utils"
	"modules/balance/acc"

	"strings"

	"strconv"

	"github.com/gin-gonic/gin"
)

type transactionsResponse struct {
	Total int64             `json:"total"`
	List  []acc.Transaction `json:"list"`
}

// listTransactions list all transaction belong to accounts in the list
// @Route {
// 		url = /list/transactions
//		method = get
//      #payload = payload class
//		#resource = resource_name
//		_accounts_ = string, comma separated list of id of accounts to get. empty means all accounts
//		_p_ = integer, the page number
//		_c_ = integer, per page item
//      200 = transactionsResponse
//      400 = base.ErrorResponseSimple
// }
func (abc AccountBaseController) listTransactions(ctx *gin.Context) {
	usr := abc.MustGetUser(ctx)
	accStr := strings.Split(ctx.Query("accounts"), ",")
	var accArray []int64
	for i := range accStr {
		tmp := strings.Trim(accStr[i], " \n\t")
		if tmp == "" {
			continue
		}
		id, err := strconv.ParseInt(tmp, 10, 0)
		if err != nil {
			continue
		}
		if id > 0 {
			accArray = append(accArray, id)
		}
	}
	offset, perPage := utils.GetPageAndCount(ctx.Request, true)
	m := acc.NewAccManager()
	accObj := m.ListUserAccountsWithID(usr, accArray...)
	res := transactionsResponse{
		Total: m.CountTransactionsByUserAndAccounts(usr, accObj...),
		List:  m.ListTransactionsByUserAndAccounts(usr, offset, perPage, accObj...),
	}

	abc.OKResponse(ctx, res)
}
