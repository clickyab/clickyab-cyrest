package authz

import (
	"common/assert"
	"common/controllers/base"
	"net/http"

	"gopkg.in/labstack/echo.v3"
)

const (
	godResource string = "god"

	scopeGranted = "__granted_scope"
	permGranted  = "__granted_perm"
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
			granted := resource
			grantedScope, ok := u.HasPerm(scope, granted)
			if !ok {
				granted = godResource
				grantedScope, ok = u.HasPerm(base.ScopeGlobal, granted)
			}
			if !ok {
				c.Request().Header.Set("error", st.Error)
				return c.JSON(
					http.StatusForbidden,
					st,
				)
			}
			c.Set(scopeGranted, grantedScope)
			c.Set(permGranted, granted)

			return next(c)
		}
	}
}
