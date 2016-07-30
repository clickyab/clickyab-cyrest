package user

import (
	"modules/user/aaa"

	"github.com/labstack/echo"
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
func (u *Controller) listRoles(ctx echo.Context) error {
	m := aaa.NewAaaManager()

	return u.OKResponse(ctx, m.ListRoles())
}
