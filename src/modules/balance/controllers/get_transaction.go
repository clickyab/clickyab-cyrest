package balance

import (
	"strconv"

	"modules/balance/acc"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

// getTransaction get transaction
// @Route {
// 		url = /transaction/:transaction_id
//		method = get
//		_transaction_id_ = integer, the transaction id
//      #payload = transactionMultiPayload
//		#resource = resource_name
//      200 = acc.Transaction
//		404 = base.ErrorResponseSimple
// }
func (c *Controller) getTransaction(ctx *gin.Context) {
	account := c.MustGetAccount(ctx)
	tID, err := strconv.ParseInt(ctx.Param("transaction_id"), 10, 0)
	if err != nil {
		logrus.Debug(err)
		c.NotFoundResponse(ctx, nil)
		return
	}

	m := acc.NewAccManager()
	tr, err := m.FindTransactionByAccountAndID(account, tID)
	if err != nil {
		logrus.Debug(err)
		c.NotFoundResponse(ctx, nil)
		return
	}

	c.OKResponse(ctx, tr)
}
