package middlewares

import (
	"common/redis"
	"time"

	"github.com/labstack/echo"

	"common/assert"
	"errors"
	"modules/user/aaa"
	"net/http"
)

const userData = "__user_data__"
const tokenData = "__token__"

// Auth is the middleware for authenticating user
func Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header().Get("token")
		st := struct {
			Error string `json:"error"`
		}{
			Error: http.StatusText(http.StatusUnauthorized),
		}
		if token != "" {
			//get the token on redis
			accessToken, err := aredis.GetKey(token, true, 24*time.Hour)
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

// GetUserData is the helper function to extract user data from context
func GetUserData(ctx echo.Context) (*aaa.User, error) {
	rd, ok := ctx.Get(userData).(*aaa.User)
	if !ok {
		return nil, errors.New("not valid data in context")
	}

	return rd, nil
}

// MustGetUserData try to get user data, or panic if there is no user data
func MustGetUserData(ctx echo.Context) *aaa.User {
	rd, err := GetUserData(ctx)
	assert.Nil(err)
	return rd
}
