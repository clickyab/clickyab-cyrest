package user

import (
	"modules/misc/trans"
	"modules/user/aaa"
	"modules/user/middlewares"

	"modules/misc/base"

	"modules/telegram/teleuser/tlu"

	"gopkg.in/labstack/echo.v3"
)

type responseLoginOK struct {
	UserID      int64                                `json:"user_id"`
	Email       string                               `json:"email"`
	AccessToken string                               `json:"token"`
	Resolve     string                               `json:"resolve,omitempty"`
	Permissions map[base.UserScope][]base.Permission `json:"perm"`
	Profile     *aaa.UserProfile                     `json:"profile,omitempty"`
}

// @Validate {
// }
type loginPayload struct {
	Email    string `json:"email" validate:"email" error:"email is invalid"`
	Password string `json:"password" validate:"gt=5" error:"password is too short"`
}

func createLoginResponse(u *aaa.User, t string) responseLoginOK {
	x := u.GetPermission()
	res := make(map[base.UserScope][]base.Permission)

	for i := range x {
		for j := range x[i] {
			res[i] = append(res[i], j)
		}
	}
	profile := u.GetProfile()
	result := responseLoginOK{
		UserID:      u.ID,
		Email:       u.Email,
		AccessToken: t,
		Permissions: res,
	}
	result.Profile = profile
	return result
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
	if err != nil {
		return u.BadResponse(ctx, trans.E("user/password is invalid"))
	}

	if !usr.VerifyPassword(pl.Password) || usr.Status == aaa.UserStatusBlocked {
		return u.BadResponse(ctx, trans.E("user/password is invalid"))
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
	result := createLoginResponse(usr, token)
	//find if user has been resolved
	count := tlu.NewTluManager().GetActiveCountByUserID(usr.ID)
	if count >= 1 {
		result.Resolve = "yes"
	} else {
		result.Resolve = "no"
	}
	return u.OKResponse(
		ctx,
		result,
	)
}
