package user

import (
	"common/assert"
	"modules/user/aaa"
	"modules/user/middlewares"

	"github.com/labstack/echo"
)

// logout is for the logout from the system
// @Route {
// 		url = /logout
//		method = get
//		middleware = authz.Authenticate
//      200 = base.NormalResponse
// }
func (u *Controller) logout(ctx echo.Context) error {
	token, ok := ctx.Get(authz.ContextToken).(string)
	assert.True(ok, "[BUG] no token on logout route")
	m := aaa.NewAaaManager()
	assert.Nil(m.EraseToken(token))
	return u.OKResponse(ctx, nil)
}
