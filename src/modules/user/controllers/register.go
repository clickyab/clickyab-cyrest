package user

import (
	"modules/user/aaa"

	"modules/misc/trans"

	"gopkg.in/labstack/echo.v3"
)

//@Validate{
// }
type registrationPayload struct {
	Email    string `json:"email" validate:"email"`
	Password string `json:"password" validate:"gt=5"`
}

// registerUser register user in system
// @Route {
// 		url = /register
//		method = post
//      payload = registrationPayload
//		200 = responseLoginOK
//		400 = base.ErrorResponseSimple
// }
func (u *Controller) registerUser(ctx echo.Context) error {
	pl := u.MustGetPayload(ctx).(*registrationPayload)
	m := aaa.NewAaaManager()

	user, err := m.RegisterUser(pl.Email, pl.Password)
	if err != nil {
		return u.BadResponse(ctx, trans.E("duplicate user"))
	}

	token := m.GetNewToken(user, ctx.Request().UserAgent(), ctx.RealIP())
	return u.OKResponse(
		ctx,
		responseLoginOK{
			UserID:      user.ID,
			Email:       user.Email,
			AccessToken: token,
		},
	)
}
