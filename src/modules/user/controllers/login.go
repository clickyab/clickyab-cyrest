package user

import (
	"errors"
	"modules/misc/trans"
	"modules/user/aaa"
	"strings"

	"github.com/gin-gonic/gin"
)

type payloadLoginData struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type responseLoginOK struct {
	UserID    int64    `json:"user_id"`
	Username  string   `json:"username"`
	Contact   string   `json:"contact"`
	Token     string   `json:"token"`
	Resources []string `json:"resources"`
}

// Login is the login route for REST requests
// @Route {
// 		url = /login
//		method = post
//      payload = payloadLoginData
//      200 = responseLoginOK
//      400 = base.ErrorResponseSimple
// }
func (u *Controller) login(ctx *gin.Context) {
	payload := u.MustGetPayload(ctx).(*payloadLoginData)
	m := aaa.NewAaaManager()
	token, user, err := m.LoginUserByPassword(strings.ToLower(payload.UserName), payload.Password)
	if err != nil {
		audit(payload.UserName, "LoginFail", "error", err)
		u.BadResponse(ctx, errors.New(trans.T("invalid username/password")))
		return
	}
	// Ignore the result, not a big deal
	_ = m.UpdateLastLogin(user)
	audit(user.Username, "LoginOK", "success", err)
	u.OKResponse(
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
