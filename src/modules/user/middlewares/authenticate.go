package authz

import (
	"errors"
	"modules/user/aaa"
	"net/http"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
)

const (
	// ContextUser is for user unmarshal object
	ContextUser string = "_user"
	// ContextToken is login token
	ContextToken string = "_token"
)

// Authenticate is a middleware to handle authentication for a user. all route with this middleware
// need to has a token in their headed
// Route {
//		403 = base.ErrorResponseSimple
// }
func Authenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := strings.Trim(c.Request().Header().Get("token"), "\n\t\" ")
		logrus.Infof("token '%s' is recieved", token)
		if token == "" {
			st := struct {
				Error string `json:"error"`
			}{
				Error: http.StatusText(http.StatusForbidden),
			}
			c.Request().Header().Set("error", st.Error)
			c.JSON(
				http.StatusForbidden,
				st,
			)
			return errors.New(st.Error)
		}

		m := aaa.NewAaaManager()
		u, err := m.FindUserByIndirectToken(token)
		if err != nil {
			st := struct {
				Error string `json:"error"`
			}{
				Error: http.StatusText(http.StatusForbidden),
			}
			c.Request().Header().Set("error", st.Error)
			c.JSON(
				http.StatusForbidden,
				st,
			)
			logrus.Infof("forbiden due to error %s", err)
			return err
		}
		c.Set(ContextUser, u)
		c.Set(ContextToken, token)
		return next(c)
	}
}

// GetUser from the request
func GetUser(c echo.Context) (*aaa.User, bool) {
	u := c.Get(ContextUser)
	user, ok := u.(*aaa.User)
	return user, ok
}
