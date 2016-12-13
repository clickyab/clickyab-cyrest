package user

import (
	"modules/misc/trans"

	"common/assert"
	"modules/user/aaa"

	"gopkg.in/labstack/echo.v3"

	"errors"
)

type assignParentPayload struct {
	UserID   int64 `json:"user_id"`
	ParentID int64 `json:"parent_id"`
}

func (pl *assignParentPayload) Validate(ctx echo.Context) error {
	if pl.UserID == pl.ParentID || pl.UserID <= 0 || pl.ParentID <= 0 {
		return errors.New("invalid payload")
	}
	return nil
}

// changePassword
// @Route {
//		url	=	/assign
//		method	=	post
//		payload	=	assignParentPayload
//		resource=	assign_parent:global
//		200	=	base.NormalResponse
//		400	=	base.ErrorResponseSimple
// }
func (u *Controller) assignParent(ctx echo.Context) error {
	pl := u.MustGetPayload(ctx).(*assignParentPayload)
	m := aaa.NewAaaManager()
	//var usr *aaa.User
	usr, err := m.FindUserByID(pl.UserID)
	if err != nil {
		return u.BadResponse(ctx, errors.New("user not found"))
	}
	parent, err := m.FindUserByID(pl.ParentID)
	if err != nil {
		return u.BadResponse(ctx, errors.New("parent not found"))
	}

	//check user has not parent
	if parent.ParentID.Valid {
		return u.BadResponse(ctx, trans.E("parent must not be child of another user"))
	}

	usr.ParentID.Int64 = pl.ParentID
	usr.ParentID.Valid = true

	assert.Nil(m.UpdateUser(usr))
	return u.OKResponse(ctx, nil)

}
