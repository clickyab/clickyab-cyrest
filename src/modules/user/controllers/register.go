package user

import (
	"modules/user/aaa"

	"github.com/labstack/echo"
)

type responseLoginOK struct {
	UserID      int64  `json:"user_id"`
	Email       string `json:"email"`
	AccessToken string `json:"Accesstoken"`
}

type registrationPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Personal bool   `json:"personal"`
}

func (r *registrationPayload) Validate(ctx echo.Context) (bool, map[string]string) {
	var res = make(map[string]string)
	var fail bool
	if len(r.Password) < 6 {
		res["password"] = "password is invalid"
		fail = true
	}

	if fail {
		return false, res
	}
	return true, nil
}

// registerUser register user in system
// @Route {
// 		url = /register
//		method = post
//      payload = registrationPayload
//      200 = responseLoginOK
//      400 = base.ErrorResponseSimple
// }
func (u *Controller) registerUser(ctx echo.Context) error {
	pl := u.MustGetPayload(ctx).(*registrationPayload)
	m := aaa.NewAaaManager()

	user, err := m.RegisterUser(pl.Email, pl.Password, pl.Personal)
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
