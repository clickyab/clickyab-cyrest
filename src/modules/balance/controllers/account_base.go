package balance

import (
	"common/controllers/base"
	user "modules/user/controllers"
)

// Controller is the base controller for balance module
// @Route {
//		group = /balance
//		middleware = authz.Authenticate
// }
type AccountBaseController struct {
	user.BaseController
}

func init() {
	base.Register(&AccountBaseController{})
}
