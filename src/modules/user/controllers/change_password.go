package user

import (
	"modules/user/aaa"

	"common/assert"

	"modules/user/middlewares"

	"modules/misc/trans"

	"gopkg.in/labstack/echo.v3"
)

// @Validate {
// }
type changePasswordPayload struct {
	OldPassword string `json:"old_password" validate:"gt=5"`
	NewPassword string `json:"new_password" validate:"gt=5"`
}

// changePassword
// @Route {
//		url	=	/change-password
//		method	=	post
//		payload	= changePasswordPayload
//		middleware = authz.Authenticate
//		200	=	base.NormalResponse
//		400	=	base.ErrorResponseSimple
// }
func (u *Controller) changePassword(ctx echo.Context) error {
	pl := u.MustGetPayload(ctx).(changePasswordPayload)

	//var usr *aaa.User
	usr := authz.MustGetUser(ctx)
	if !usr.VerifyPassword(pl.OldPassword) {
		return u.BadResponse(ctx, trans.E("old password is not true"))
	}
	usr.Password.String = pl.NewPassword

	m := aaa.NewAaaManager()
	assert.Nil(m.UpdateUser(usr))
	return u.OKResponse(ctx, trans.T("change password successful"))
}