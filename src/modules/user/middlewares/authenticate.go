package authz

import (
	"modules/user/aaa"
	"net/http"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
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
func Authenticate(c *gin.Context) {
	token := strings.Trim(c.Request.Header.Get("token"), "\n\t\" ")
	logrus.Infof("token '%s' is recieved", token)
	if token == "" {
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

	m := aaa.NewAaaManager()
	u, err := m.FindUserByIndirectToken(token)
	if err != nil {
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
		logrus.Infof("forbiden due to error %s", err)
		return
	}
	c.Set(ContextUser, u)
	c.Set(ContextToken, token)
}

// GetUser from the request
func GetUser(c *gin.Context) (*aaa.User, bool) {
	if u, ok := c.Get(ContextUser); ok {
		user, ok := u.(*aaa.User)
		return user, ok
	}

	return nil, false
}
