package authz

import (
	"common/utils"
	"net/http"

	"errors"

	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
)

const (
	godResource string = "god"
)

// AuthorizeGenerator create a middleware for check authorize request, must use it after authenticate
// or else this block the request
func AuthorizeGenerator(resource string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			u, ok := GetUser(c)
			if !ok {
				st := struct {
					Error string `json:"error"`
				}{
					Error: http.StatusText(http.StatusForbidden),
				}
				c.Response().Header().Set("error", st.Error)
				c.JSON(
					http.StatusForbidden,
					st,
				)
				return errors.New(st.Error)
			}
			m := u.GetResources()
			if !utils.StringInArray(resource, m...) && !utils.StringInArray(godResource, m...) {
				st := struct {
					Error string `json:"error"`
				}{
					Error: http.StatusText(http.StatusForbidden),
				}
				c.Response().Header().Set("error", st.Error)
				c.JSON(
					http.StatusForbidden,
					st,
				)
				logrus.Infof("forbidden since the %s is not available", resource)
				return errors.New(st.Error)
			}

			return next(c)
		}
	}
}
