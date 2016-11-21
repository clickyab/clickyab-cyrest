package misc

import (
	"common/version"

	"github.com/labstack/echo"
)

// getVersion get the version information
// @Route {
// 		url = /version
//		method = get
//      200 = version.Version
// }
func (u *Controller) getVersion(ctx echo.Context) error {
	return u.OKResponse(ctx, version.GetVersion())
}
