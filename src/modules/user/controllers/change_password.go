package user

import (
	"modules/user/aaa"

	"modules/misc/trans"

	"github.com/gin-gonic/gin"
)

type changePasswordPayload struct {
	OldPass string `json:"old_pass,omitempty"`
	NewPass string `json:"new_pass"`
}

func (p *changePasswordPayload) Validate(ctx *gin.Context) (fail bool, res map[string]string) {
	res = make(map[string]string)

	if len(p.NewPass) < 6 {
		res["password"] = trans.T("password is invalid")
		fail = true
	}

	return
}

// changePassword change user password
// @Route {
// 		url = /password
//		method = post
//      payload = changePasswordPayload
//		middleware = authz.Authenticate
//      200 = base.NormalResponse
//      400 = base.ErrorResponseSimple
// }
func (u *Controller) changePassword(ctx *gin.Context) {
	password := u.MustGetPayload(ctx).(*changePasswordPayload)
	usr := u.MustGetUser(ctx)

	if usr.HasPassword() {
		if !usr.VerifyPassword(password.OldPass) {
			u.BadResponse(ctx, trans.E("old password is invalid"))
			return
		}
	}
	usr.Password = password.NewPass
	m := aaa.NewAaaManager()
	err := m.UpdateUser(usr)
	if err != nil {
		u.BadResponse(ctx, err)
		return
	}

	u.OKResponse(ctx, nil)
}
