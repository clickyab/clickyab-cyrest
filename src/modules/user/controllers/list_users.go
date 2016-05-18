package user

import (
	"common/utils"
	"modules/user/aaa"
	"strings"

	"github.com/gin-gonic/gin"
)

type userListResponse struct {
	Total int64      `json:"total"`
	List  []aaa.User `json:"list"`
}

// listUsers list all users in site
// @Route {
// 		url = /users
//		method = get
//      #payload = payload class
//		resource = user_admin
//		_username_ = string, the username filter
//		_status_ = string, the status filter
//		_p_ = integer, the page number
//		_c_ = integer, per page item
//      200 = userListResponse
// }
func (u *Controller) listUsers(ctx *gin.Context) {
	o, c := utils.GetPageAndCount(ctx.Request, true)

	username := strings.Trim(ctx.Request.URL.Query().Get("username"), " ")
	status := aaa.UserStatus(ctx.Request.URL.Query().Get("status"))

	res := &userListResponse{}
	res.List, res.Total = aaa.NewAaaManager().ListUserFilterByUsername(o, c, username, status)

	u.OKResponse(ctx, res)
}
