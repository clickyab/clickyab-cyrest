package balance

import "github.com/gin-gonic/gin"

// getAccount get a single account
// @Route {
// 		url = /
//		method = get
//      #payload = payload class
//		#resource = resource_name
//		middleware = accm.AccountCheck
//      200 = acc.Account
//      400 = base.ErrorResponseSimple
// }
func (c *Controller) getAccount(ctx *gin.Context) {
	account := c.MustGetAccount(ctx)

	c.OKResponse(ctx, account)
}
