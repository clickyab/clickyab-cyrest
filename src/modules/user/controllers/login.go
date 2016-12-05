package user

import (
	"modules/user/aaa"

	"github.com/labstack/echo"
)

type response2LoginOK struct {
	UserID      int64  `json:"user_id"`
	AccessToken string `json:"Accesstoken"`
}

type loginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *loginPayload) Validate(ctx echo.Context) (bool, map[string]string) {
	var res = make(map[string]string)
	var fail bool
	if len(r.Password) < 6 {
		res["password"] = "password is invalid"
		fail = true
	}
	valid, _:= ValidateEmail(r.Email)
	if len(r.Email) < 6  || valid{
		res["Email"] = "Email is invalid"
		fail = true
	}

	if fail {
		return false, res
	}
	return true, nil
}

// loginUser login user in system
// @Route {
// 		url = /login
//		method = post
//      payload = loginPayload
//      200 = responseLoginOK
//      400 = base.ErrorResponseSimple
// }
func (u *Controller) loginUser(ctx echo.Context) error {
	pl := u.MustGetPayload(ctx).(*loginPayload)
	m := aaa.NewAaaManager()

	user, err := m.LoginUserByPassword(pl.Email, pl.Password)
	if err != nil {
		return u.BadResponse(ctx, err)

	}

	token := m.GetNewToken(user.AccessToken)
	return u.OKResponse(
		ctx,
		responseLoginOK{
			UserID:      user.ID,
			Email:       user.Email,
			AccessToken: token,
		},
	)
}
