package user

import (
	"errors"
	"modules/misc/trans"
	"modules/user/aaa"
	"strconv"

	"github.com/gin-gonic/gin"
)

type (
	createRolePayload struct {
		Name        string   `json:"name"`
		Description string   `json:"description"`
		Resources   []string `json:"resources"`
	}

	createRoleResponse struct {
		ID int64 `json:"id"`
	}
)

// createRole is for creating a role in system
// @Route {
// 		url = /role
//		method = post
//      payload = createRolePayload
//		resource = user_admin
//      200 = createRoleResponse
//      400 = base.ErrorResponseSimple
// }
func (u *Controller) createRole(ctx *gin.Context) {
	payload := u.MustGetPayload(ctx).(*createRolePayload)
	role := aaa.Role{
		Name:        payload.Name,
		Description: payload.Description,
		Resources:   append([]string{}, payload.Resources...),
	}
	m := aaa.NewAaaManager()

	err := m.CreateRole(&role)
	if err != nil {
		u.BadResponse(ctx, errors.New(trans.T("can not create role")))
		return
	}
	u.OKResponse(
		ctx,
		createRoleResponse{
			ID: role.ID,
		},
	)
}

// updateRole is for changing a role in system
// @Route {
// 		url = /role/:id
//		method = put
//      payload = createRolePayload
//		resource = create_role
//		:id = true, int, the role id to edit
//      200 = base.NormalResponse
//      400 = base.ErrorResponseSimple
// }
func (u *Controller) updateRole(ctx *gin.Context) {
	m := aaa.NewAaaManager()
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		u.NotFoundResponse(ctx, err)
		return
	}
	role, err := m.FindRoleByID(id)
	if err != nil {
		u.NotFoundResponse(ctx, nil)
		return
	}
	payload := u.MustGetPayload(ctx).(*createRolePayload)

	role.Name = payload.Name
	role.Description = payload.Description
	role.Resources = append([]string{}, payload.Resources...)
	err = m.UpdateRole(role)
	if err != nil {
		u.BadResponse(ctx, errors.New(trans.T("can not update role")))
		return
	}
	u.OKResponse(
		ctx,
		nil,
	)
}
