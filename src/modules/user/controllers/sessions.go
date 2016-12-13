package user

import (
	"common/assert"
	"common/redis"
	"errors"
	"modules/user/aaa"
	"modules/user/config"
	"modules/user/middlewares"

	"gopkg.in/labstack/echo.v3"
)

type activeSessions struct {
	More     bool         `json:"more"`
	Sessions aaa.Sessions `json:"sessions"`
}

// logout is for the logout from the system
// @Route {
// 		url = /logout
//		method = get
//		middleware = authz.Authenticate
//      200 = base.NormalResponse
// }
func (u *Controller) logout(ctx echo.Context) error {
	token := authz.MustGetToken(ctx)
	m := aaa.NewAaaManager()
	assert.Nil(m.EraseToken(token))
	return u.OKResponse(ctx, nil)
}

// @Route {
// 		url = /sessions
//		method = get
//		middleware = authz.Authenticate
//      200 = activeSessions
// }
func (u *Controller) activeSessions(ctx echo.Context) error {
	token := authz.MustGetToken(ctx)
	m := aaa.NewAaaManager()
	res := activeSessions{}
	res.Sessions, res.More = m.GetSessions(authz.MustGetUser(ctx), token, 10)
	return u.OKResponse(ctx, res)
}

// @Route {
// 		url = /session/terminate/:id
//		method = get
//		middleware = authz.Authenticate
//		:id = true, string , session id to terminate, must be your session
//      200 = activeSessions
//      400 = base.ErrorResponseSimple
// }
func (u *Controller) terminateSession(ctx echo.Context) error {
	token := authz.MustGetToken(ctx)
	inner, err := aredis.GetHashKey(token, "token", false, 0)
	assert.Nil(err)

	other := ctx.Param("id")
	inner2, err := aredis.GetHashKey(other, "token", false, 0)
	assert.Nil(err)

	if inner != inner2 {
		return u.BadResponse(ctx, errors.New("the token is invalid"))
	}
	if other == token {
		return u.BadResponse(ctx, errors.New("cant kill current session"))
	}

	m := aaa.NewAaaManager()
	assert.Nil(m.EraseToken(other))

	res := activeSessions{}
	res.Sessions, res.More = m.GetSessions(authz.MustGetUser(ctx), token, 10)
	return u.OKResponse(ctx, res)
}

// @Route {
// 		url = /sessions/terminate
//		method = get
//		middleware = authz.Authenticate
//      200 = activeSessions
//      400 = base.ErrorResponseSimple
// }
func (u *Controller) terminateAllSession(ctx echo.Context) error {
	token := authz.MustGetToken(ctx)
	usr := authz.MustGetUser(ctx)

	m := aaa.NewAaaManager()
	assert.Nil(m.LogoutAllSession(usr))
	assert.Nil(aredis.StoreHashKey(token, "token", usr.AccessToken, ucfg.Cfg.TokenTimeout))

	res := activeSessions{}
	res.Sessions, res.More = m.GetSessions(usr, token, 10)
	return u.OKResponse(ctx, res)
}
