package user

import (
	"modules/user/aaa"

	"github.com/gin-gonic/gin"
)

// listRoles list all roles in system, no pagination
// @Route {
// 		url = /roles
//		method = get
//      #payload = payload class
//		#resource = resource_name
//		middleware = authz.Authenticate
//      200 = roles
//      400 = base.ErrorResponseSimple
// }
func (u *Controller) listRoles(ctx *gin.Context) {
	m := aaa.NewAaaManager()

	u.OKResponse(ctx, m.ListRoles())
}
