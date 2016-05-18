package user

import (
	"fmt"
	"modules/user/aaa"
	"strconv"

	"github.com/gin-gonic/gin"
)

type changeStatePayload struct {
	Status aaa.UserStatus `json:"status"`
}

func (csp *changeStatePayload) Validate(ctx *gin.Context) (bool, map[string]string) {
	if !csp.Status.IsValid() {
		return false, map[string]string{"status": fmt.Sprintf("invaid status, valids are : %s, %s, %s", aaa.UserStatusBanned, aaa.UserStatusRegistered, aaa.UserStatusVerified)}
	}

	return true, nil
}

// changeState change user state, ban or verify. admin actions
// @Route {
// 		url = /state/:user_id
//		method = post
//      payload = changeStatePayload
//		resource = user_admin
//      200 = base.NormalResponse
//      400 = base.ErrorResponseSimple
// }
func (u *Controller) changeState(ctx *gin.Context) {
	status := u.MustGetPayload(ctx).(*changeStatePayload)
	uID, err := strconv.ParseInt(ctx.Param("user_id"), 10, 0)
	if err != nil {
		u.NotFoundResponse(ctx, err)
		return
	}
	m := aaa.NewAaaManager()
	usr, err := m.FindUserByID(uID)
	if err != nil {
		u.NotFoundResponse(ctx, err)
		return
	}

	usr.Status = status.Status
	if m.UpdateUser(usr) != nil {
		u.BadResponse(ctx, err)
		return
	}
	u.OKResponse(ctx, nil)
}
