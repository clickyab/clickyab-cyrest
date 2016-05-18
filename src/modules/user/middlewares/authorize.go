package authz

import (
	"common/utils"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

const (
	godResource string = "god"
)

// AuthorizeGenerator create a middleware for check authorize request, must use it after authenticate
// or else this block the request
func AuthorizeGenerator(resource string) gin.HandlerFunc {
	return func(c *gin.Context) {
		u, ok := GetUser(c)
		if !ok {
			st := struct {
				Error string `json:"error"`
			}{
				Error: http.StatusText(http.StatusForbidden),
			}
			c.Header("error", st.Error)
			c.JSON(
				http.StatusForbidden,
				st,
			)
			c.Abort()
			return
		}
		m := u.GetResources()
		if !utils.StringInArray(resource, m...) && !utils.StringInArray(godResource, m...) {
			st := struct {
				Error string `json:"error"`
			}{
				Error: http.StatusText(http.StatusForbidden),
			}
			c.Header("error", st.Error)
			c.JSON(
				http.StatusForbidden,
				st,
			)
			logrus.Infof("forbidden since the %s is not available", resource)
			c.Abort()
			return
		}

		c.Next()
	}
}
