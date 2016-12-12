package user

import (
	"modules/user/aaa"

	"modules/misc/trans"

	"gopkg.in/go-playground/validator.v9"
	"gopkg.in/labstack/echo.v3"
)

var (
	userPasswordError = trans.E("user/password is invalid")
)

type responseLoginOK struct {
	UserID      int64  `json:"user_id"`
	Email       string `json:"email"`
	AccessToken string `json:"token"`
}

type loginPayload struct {
	Email    string `json:"email" validate:"email"`
	Password string `json:"password" validate:"gt=6"`
}

func (lp *loginPayload) Validate(ctx echo.Context) error {
	return validator.New().Struct(lp)
}

func createLoginResponse(u *aaa.User, t string) responseLoginOK {
	return responseLoginOK{
		UserID:      u.ID,
		Email:       u.Email,
		AccessToken: t,
	}
}

// registerUser register user in system
// @Route {
// 		url = /login
//		method = post
//      payload = loginPayload
//		200 = responseLoginOK
//		400 = base.ErrorResponseSimple
// }
func (u *Controller) loginUser(ctx echo.Context) error {
	pl := u.MustGetPayload(ctx).(*loginPayload)
	m := aaa.NewAaaManager()

	usr, err := m.FindUserByEmail(pl.Email)
	if err != nil {
		return u.BadResponse(ctx, userPasswordError)
	}

	if !usr.VerifyPassword(pl.Password) {
		return u.BadResponse(ctx, userPasswordError)
	}

	token := m.GetNewToken(usr, ctx.Request().UserAgent(), ctx.RealIP())
	return u.OKResponse(
		ctx,
		createLoginResponse(usr, token),
	)
}
