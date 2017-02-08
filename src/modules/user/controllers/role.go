package user

import (
	"modules/misc/base"
	"modules/misc/middlewares"
	"modules/misc/trans"
	"modules/user/aaa"
	"strconv"

	"gopkg.in/labstack/echo.v3"
)

// @Validate {
// }
type rolePayLoad struct {
	Name        string                               `json:"name" validate:"gt=3" error:"name must be valid"`
	Description string                               `json:"description" validate:"gt=3" error:"description must be valid"`
	Perm        map[base.UserScope][]base.Permission `json:"perm"`
}

type roleResponse struct {
	Role aaa.Role                                    `json:"role"`
	Perm map[base.UserScope]map[base.Permission]bool `json:"perm"`
}

// Validate custom validation for user scope
func (lp *rolePayLoad) ValidateExtra(ctx echo.Context) error {
	for i := range lp.Perm {
		if !i.IsValid() {
			return middlewares.GroupError{
				string(i): trans.E("scope is invalid"),
			}
		}
	}
	return nil
}

// createRole register user in system
// @Route {
// 		url = /role/create/
// 		resource = create_role:global
//		method = post
//      payload = rolePayLoad
//		200 = aaa.Role
//		400 = base.ErrorResponseSimple
// }
func (u *Controller) createRole(ctx echo.Context) error {
	pl := u.MustGetPayload(ctx).(*rolePayLoad)
	m := aaa.NewAaaManager()

	//insert new role to database
	role, err := m.RegisterRole(pl.Name, pl.Description, pl.Perm)
	if err != nil {
		return u.BadResponse(ctx, trans.E("can not create role"))
	}

	return u.OKResponse(
		ctx,
		role,
	)
}

// deleteRole delete specified role in system
// @Route {
// 		url = /role/delete/:id
// 		resource = delete_role:global
//		method = delete
//		:id = true, int, id of role to be deleted
//		200 = aaa.Role
//		400 = base.ErrorResponseSimple
// }
func (u *Controller) deleteRole(ctx echo.Context) (err error) {
	//var role *aaa.Role
	var m = aaa.NewAaaManager()
	ID, err := strconv.ParseInt(ctx.Param("id"), 10, 0)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}

	role, err := m.DeleteRole(ID)
	if err != nil {
		return u.BadResponse(ctx, trans.E("can not delete role"))
	}
	return u.OKResponse(
		ctx,
		role,
	)
}

// getRole getRole by id
// @Route {
// 		url = /role/:id
// 		resource = get_role:global
//		method = get
//		:id = true, int, id of role to be deleted
//		200 = aaa.Role
//		400 = base.ErrorResponseSimple
// }
func (u *Controller) getRole(ctx echo.Context) (err error) {
	var m = aaa.NewAaaManager()
	ID, err := strconv.ParseInt(ctx.Param("id"), 10, 0)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	role, err := m.FindRoleByID(ID)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	perms := m.GetPermissionMap(*role)
	return u.OKResponse(ctx, roleResponse{
		Role: *role,
		Perm: perms,
	})
}

// updateRole register user in system
// @Route {
// 		url = /role/update/:id
// 		resource = update_role:global
// 		:id = true, int, id of role to be updated
//		method = put
//      	payload = rolePayLoad
//		200 = aaa.Role
//		400 = base.ErrorResponseSimple
// }
func (u *Controller) updateRole(ctx echo.Context) error {
	ID, err := strconv.ParseInt(ctx.Param("id"), 10, 0)
	if err != nil {
		return u.BadResponse(ctx, trans.E("can not update role"))
	}
	pl := u.MustGetPayload(ctx).(*rolePayLoad)
	m := aaa.NewAaaManager()

	//update role in db
	role, err := m.UpdateRoleWithPerm(ID, pl.Name, pl.Description, pl.Perm)
	if err != nil {
		return u.BadResponse(ctx, trans.E("can not update role"))
	}

	return u.OKResponse(
		ctx,
		role,
	)
}

// allRoles allRoles
// @Route {
// 		url = /allroles
// 		resource = get_all_roles:parent
//		method = get
//		200 = aaa.Role
//		400 = base.ErrorResponseSimple
// }
func (u *Controller) allRoles(ctx echo.Context) (err error) {
	m := aaa.NewAaaManager()
	roles, err := m.GetAllRole()
	if err != nil {
		return u.BadResponse(ctx, trans.E("error while getting roles"))
	}
	return u.OKResponse(ctx, roles)
}
