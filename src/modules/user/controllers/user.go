package user

import "common/controllers/base"

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

func init() {
	base.Register(&Controller{})
}
