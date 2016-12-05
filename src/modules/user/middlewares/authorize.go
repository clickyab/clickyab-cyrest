package middlewares

import (
	"net/http"

	"github.com/labstack/echo"
)

// AuthorizeGenerator generate middleware for specified action
func AuthorizeGenerator(resource string, scope string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			st := struct {
				Error string `json:"error"`
			}{
				Error: http.StatusText(http.StatusForbidden),
			}
			// get user
			u := MustGetUserData(c)

			//check if the user has the specified perm
			if !u.HasPerm(resource, scope) {
				c.Request().Header().Set("error", st.Error)
				c.JSON(
					http.StatusForbidden,
					st,
				)
			}

			return next(c)
		}
	}
}
