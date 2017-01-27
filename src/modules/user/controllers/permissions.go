package user

import (
	"modules/misc/base"

	"gopkg.in/labstack/echo.v3"
)

// listPermissions list of all permissions
// @Route {
// 		url = /permissions
//		method = get
//		middleware = authz.Authenticate
//		200 = perms
// }
func (u *Controller) listPermissions(ctx echo.Context) error {
	return u.OKResponse(ctx, base.GetAllPermission())
}
