package balance

import (
	"modules/balance/acc"
	"modules/misc/trans"

	"github.com/gin-gonic/gin"
)

type theAccountPayload struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Disabled    bool   `json:"disabled"`
}

type addAccountResponse struct {
	AccountID int64 `json:"account_id"`
}

func (r *theAccountPayload) Validate(ctx *gin.Context) (bool, map[string]string) {
	var res = make(map[string]string)
	var fail bool
	if len(r.Title) < 3 {
		res["title"] = trans.T("title_is_too_short")
		fail = true
	}

	if fail {
		return false, res
	}
	return true, nil
}

// editAccount edit the account in system
// @Route {
// 		url = /
//		method = put
//      payload = theAccountPayload
//		middleware = accm.IsOwner
//      200 = base.NormalResponse
//      400 = base.ErrorResponseSimple
// }
func (c *Controller) editAccount(ctx *gin.Context) {
	account := c.MustGetAccount(ctx)
	usr := c.MustGetUser(ctx)
	payload := c.MustGetPayload(ctx).(*theAccountPayload)

	m := acc.NewAccManager()
	err := m.EditAccount(
		usr,
		account,
		payload.Title,
		payload.Description,
		payload.Disabled,
	)
	if err != nil {
		c.BadResponse(ctx, err)
		return
	}

	c.OKResponse(ctx, nil)
}
