package user

import (
	"common/controllers/base"
	"fmt"
	"modules/misc/trans"
	"modules/user/aaa"
	"modules/user/middlewares"

	"gopkg.in/labstack/echo.v3"
)

var (
	userPasswordError = trans.E("user/password is invalid")
)

type responseLoginOK struct {
	UserID      int64                              `json:"user_id"`
	Email       string                             `json:"email"`
	AccessToken string                             `json:"token"`
	Permissions map[base.UserScope]map[string]bool `json:"perm"`
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
		Permissions: u.GetPermission(),
	}
}

// loginUser login user in system
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

// pingUser is the ping user and get data again
// @Route {
// 		url = /ping
//		method = get
//		middleware = authz.Authenticate
//		200 = responseLoginOK
// }
func (u *Controller) pingUser(ctx echo.Context) error {
	usr := authz.MustGetUser(ctx)
	token := authz.MustGetToken(ctx)

	return u.OKResponse(
		ctx,
		createLoginResponse(usr, token),
	)
}
