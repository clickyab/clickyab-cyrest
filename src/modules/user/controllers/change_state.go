package user

import (
	"fmt"
	"modules/user/aaa"
	"strconv"

	"github.com/labstack/echo"
)

type changeStatePayload struct {
	Status aaa.UserStatus `json:"status"`
}

func (csp *changeStatePayload) Validate(ctx echo.Context) (bool, map[string]string) {
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
func (u *Controller) changeState(ctx echo.Context) error {
	status := u.MustGetPayload(ctx).(*changeStatePayload)
	uID, err := strconv.ParseInt(ctx.Param("user_id"), 10, 0)
	if err != nil {
		return u.NotFoundResponse(ctx, err)

	}
	m := aaa.NewAaaManager()
	usr, err := m.FindUserByID(uID)
	if err != nil {
		return u.NotFoundResponse(ctx, err)

	}

	usr.Status = status.Status
	if m.UpdateUser(usr) != nil {
		return u.BadResponse(ctx, err)

	}
	return u.OKResponse(ctx, nil)
}
