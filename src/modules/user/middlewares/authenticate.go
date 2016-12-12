package authz

import (
	"common/assert"
	"common/redis"
	"modules/user/aaa"
	"net/http"

	"modules/user/config"

	"gopkg.in/labstack/echo.v3"
)

const userData = "__user_data__"
const tokenData = "__token__"

// Authenticate is the middleware for authenticating user
func Authenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("token")
		st := struct {
			Error string `json:"error"`
		}{
			Error: http.StatusText(http.StatusUnauthorized),
		}
		if token != "" {
			//get the token on redis
			accessToken, err := aredis.GetHashKey(token, "token", true, ucfg.Cfg.TokenTimeout)
			if err != nil { //user not authenticated
				return c.JSON(http.StatusUnauthorized, st)
			}
			//check if the accessToken exists in users table
			user, err := aaa.NewAaaManager().FetchByToken(accessToken)
			if err != nil { //user not found
				return c.JSON(http.StatusUnauthorized, st)
			}
			//all good put user in context
			c.Set(userData, user)
			c.Set(tokenData, token)
		}
		return next(c)
	}
}

// GetUser is the helper function to extract user data from context
func GetUser(ctx echo.Context) (*aaa.User, bool) {
	rd, ok := ctx.Get(userData).(*aaa.User)
	if !ok {
		return nil, false
	}

	return rd, true
}

// MustGetUser try to get user data, or panic if there is no user data
func MustGetUser(ctx echo.Context) *aaa.User {
	rd, ok := GetUser(ctx)
	assert.True(ok, "[BUG] no user in context")
	return rd
}

// GetToken return the token in context
func GetToken(ctx echo.Context) (string, bool) {
	t, ok := ctx.Get(tokenData).(string)
	if !ok {
		return "", false
	}

	return t, true
}

// MustGetToken return the token in context
func MustGetToken(ctx echo.Context) string {
	t, ok := GetToken(ctx)
	assert.True(ok, "[BUG] no token in context")
	return t
}
