package accm

import (
	"common/controllers/base"
	"modules/misc/trans"
	"modules/user/middlewares"
	"net/http"

	"github.com/gin-gonic/gin"
)

// IsOwner is a check if the user is owner of this unit or not
// Route {
//		401 = base.ErrorResponseSimple
// }
func IsOwner(c *gin.Context) {
	if u, ok := authz.GetUser(c); ok {
		if a, ok := GetAccount(c); ok && a.OwnerID == u.ID {
			return
		}
	}

	c.JSON(
		http.StatusUnauthorized, // User is logged in but have no access to this
		base.ErrorResponseSimple{Error: trans.T("must_be_owner")},
	)
	c.Abort()
}
