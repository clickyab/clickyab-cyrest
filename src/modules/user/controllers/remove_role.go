package user

import (
	"fmt"
	"modules/user/aaa"
	"strconv"

	"github.com/gin-gonic/gin"
)

// removeRole remove a role from database
// @Route {
// 		url = /role/:id
//		method = delete
//      #payload = payload class
//		resource = create_role
//		:id = true, int, the role id to remove
//      200 = base.NormalResponse
//      400 = base.ErrorResponseSimple
// }
func (u *Controller) removeRole(ctx *gin.Context) {
	m := aaa.NewAaaManager()
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		u.NotFoundResponse(ctx, err)
		return
	}
	role, err := m.FindRoleByID(id)
	if err != nil {
		u.NotFoundResponse(ctx, nil)
		return
	}

	cnt := m.CountRoleUsers(role)
	if cnt > 0 {
		u.BadResponse(ctx, fmt.Errorf("there is %d user in this role, can not remove it", cnt))
		return
	}

	_, err = m.GetDbMap().Delete(role)
	if err != nil {
		u.BadResponse(ctx, err)
		return
	}

	u.OKResponse(ctx, nil)
}
