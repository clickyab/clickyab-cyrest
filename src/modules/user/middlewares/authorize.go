package authz

import (
	"common/assert"
	"common/controllers/base"
	"net/http"

	"gopkg.in/labstack/echo.v3"
)

// AuthorizeGenerator generate middleware for specified action
func AuthorizeGenerator(resource string, scope base.UserScope) echo.MiddlewareFunc {
	assert.True(scope.IsValid(), "[BUG] invalid scope")
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			st := struct {
				Error string `json:"error"`
			}{
				Error: http.StatusText(http.StatusForbidden),
			}
			// get user
			u := MustGetUser(c)

			//check if the user has the specified perm
			if _, ok := u.HasPerm(scope, resource); !ok {
				c.Request().Header.Set("error", st.Error)
				return c.JSON(
					http.StatusForbidden,
					st,
				)
			}

			return next(c)
		}
	}
}
