package user

import (
	"common/assert"
	"modules/user/aaa"

	"gopkg.in/labstack/echo.v3"

	"errors"
)

// @Validate {
// }
type assignRolesPayload struct {
	UserID  int64   `json:"user_id" validate:"required"`
	RoleIDs []int64 `json:"role_id" validate:"required"`
}

// assignRoles
// @Route {
//		url	=	/assign/roles
//		method	=	post
//		payload	=	assignRolesPayload
//		resource=	assign_roles:global
//		middleware = authz.Authenticate
//		200	=	aaa.UserRole
//		400	=	base.ErrorResponseSimple
// }
func (u *Controller) assignRoles2User(ctx echo.Context) error {
	pl := u.MustGetPayload(ctx).(*assignRolesPayload)
	m := aaa.NewAaaManager()
	//var usr *aaa.User
	usr, err := m.FindUserByID(pl.UserID)
	if err != nil {
		return u.BadResponse(ctx, errors.New("user not found"))
	}

	for _, v := range pl.RoleIDs {
		role, err := m.FindRoleByID(v)
		if err != nil {
			return u.BadResponse(ctx, errors.New("role not found"))
		}
		ur := &aaa.UserRole{RoleID: role.ID, UserID: usr.ID}
		assert.Nil(m.CreateUserRole(ur))
	}

	return u.OKResponse(ctx, nil)

}
