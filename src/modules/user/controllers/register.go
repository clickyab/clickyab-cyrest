package user

import (
	"modules/misc/trans"
	"modules/user/aaa"
	"modules/user/utils"
	"regexp"

	"github.com/labstack/echo"
)

type registrationPayload struct {
	Token    string `json:"token"`
	Contact  string `json:"contact"`
	Username string `json:"username"`
	Password string `json:"password"`
}

var (
	// UsernameValidator is a simple regexp to validate user names
	usernameValidator = regexp.MustCompile("^[a-zA-Z][A-Za-z0-9._%@-]+$")
)

func (r *registrationPayload) Validate(ctx echo.Context) (bool, map[string]string) {
	_, err := utils.DetectContactType(r.Contact)
	var res = make(map[string]string)
	var fail bool
	if err != nil {
		res["contact"] = trans.T("only accept email and phone number")
		fail = true
	}

	if len(r.Password) < 6 {
		res["password"] = trans.T("password is invalid")
		fail = true
	}

	if !usernameValidator.MatchString(r.Username) {
		res["username"] = trans.T("only a-z and 0-9 is allowed")
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
	user, err := m.RegisterUserByToken(pl.Token, pl.Contact, pl.Username, pl.Password)
	if err != nil {
		return u.BadResponse(ctx, err)

	}

	token := m.GetNewToken(user.Token)
	return u.OKResponse(
		ctx,
		responseLoginOK{
			UserID:    user.ID,
			Username:  user.Username,
			Contact:   user.Contact,
			Token:     token,
			Resources: user.GetResources(),
		},
	)
}
