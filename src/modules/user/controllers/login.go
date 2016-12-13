package user

import (
	"modules/user/aaa"

	"modules/misc/trans"

	"fmt"

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

// @Validate {
// }
type loginPayload struct {
	Email    string `json:"email" validate:"email" error:"email is invalid"`
	Password string `json:"password" validate:"gt=5" error:"password is too short"`
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
	fmt.Printf("%+v", pl)
	m := aaa.NewAaaManager()

	usr, err := m.FindUserByEmail(pl.Email)
	fmt.Println(err)
	if err != nil {
		return u.BadResponse(ctx, userPasswordError)
	}

	if !usr.VerifyPassword(pl.Password) || usr.Status == aaa.UserStatusBlocked {
		return u.BadResponse(ctx, userPasswordError)
	}

	token := m.GetNewToken(usr, ctx.Request().UserAgent(), ctx.RealIP())
	return u.OKResponse(
		ctx,
		createLoginResponse(usr, token),
	)
}
