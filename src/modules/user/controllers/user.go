package user

import (
	"common/controllers/base"
	"common/redis"
	"modules/user/config"

	"gopkg.in/labstack/echo.v3"
)

// Controller is the controller for the user package
// @Route {
//		group = /user
// }
type Controller struct {
	base.Controller
}

//type userAudit struct {
//	Username string      `json:"username"`
//	Action   string      `json:"action"`
//	Class    string      `json:"class"`
//	Data     interface{} `json:"data"`
//}

var (
	_ = base.ErrorResponseMap{}
	_ = base.ErrorResponseSimple{}
)

//// String make this one a stringer
//func (u userAudit) String() string {
//	r, _ := json.Marshal(u)
//
//	return string(r)
//}

//func audit(username, action, class string, data interface{}) {
//	hub.Publish(
//		"audit",
//		userAudit{
//			Username: username,
//			Action:   action,
//			Class:    class,
//			Data:     data,
//		},
//	)
//}

func (c Controller) storeData(ctx echo.Context, token string) error {
	err := aredis.StoreHashKey(token, "ua", ctx.Request().UserAgent(), ucfg.Cfg.TokenTimeout)
	if err != nil {
		return err
	}
	return aredis.StoreHashKey(token, "ip", ctx.RealIP(), ucfg.Cfg.TokenTimeout)
}

func init() {
	base.Register(&Controller{})
}
