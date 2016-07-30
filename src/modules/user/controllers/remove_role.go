package user

import (
	"fmt"
	"modules/user/aaa"
	"strconv"

	"github.com/labstack/echo"
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
func (u *Controller) removeRole(ctx echo.Context) error {
	m := aaa.NewAaaManager()
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return u.NotFoundResponse(ctx, err)

	}
	role, err := m.FindRoleByID(id)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}

	cnt := m.CountRoleUsers(role)
	if cnt > 0 {
		return u.BadResponse(ctx, fmt.Errorf("there is %d user in this role, can not remove it", cnt))

	}

	_, err = m.GetDbMap().Delete(role)
	if err != nil {
		return u.BadResponse(ctx, err)

	}

	return u.OKResponse(ctx, nil)
}
