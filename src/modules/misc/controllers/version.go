package misc

import (
	"common/version"

	"github.com/gin-gonic/gin"
)

// getVersion get the version information
// @Route {
// 		url = /version
//		method = get
//      200 = version.Version
// }
func (u *Controller) getVersion(ctx *gin.Context) {
	u.OKResponse(ctx, version.GetVersion())
}
