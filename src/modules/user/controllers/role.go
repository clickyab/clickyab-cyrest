package user

import (
	"gopkg.in/labstack/echo.v3"
	"modules/user/aaa"
	"common/controllers/base"
	"gopkg.in/go-playground/validator.v9"
	"common/middlewares"
	"modules/misc/trans"
	"strconv"
	"fmt"
)

type rolePayLoad struct {
	Name string `json:"name" validate:"gt=3" error:"name must be valid"`
	Description string `json:"description" validator:"gt=3" error:"description must be valid"`
	Perm map[base.UserScope][]string `json:"perm"`
}

// Validate custom validation for user scope
func (lp *rolePayLoad)Validate(ctx echo.Context) error {
	for i:= range lp.Perm{
		if !i.IsValid(){
			return middlewares.GroupError{
				string(i) : "scope is invalid",
			}
		}
	}
	return validator.New().Struct(lp)
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
	m:=aaa.NewAaaManager()

	//insert new role to database
	role,err:=m.RegisterRole(pl.Name,pl.Description,pl.Perm)
	if err != nil {
		return u.BadResponse(ctx, trans.E("can not create role"))
	}

	return u.OKResponse(
		ctx,
		role,
	)
	return nil
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
	fmt.Println("here")
	//var role *aaa.Role
	var m= aaa.NewAaaManager()
	ID,err:=strconv.ParseInt(ctx.Param("id"),10,0)
	role,err:=m.DeleteRole(ID)
	if err!=nil{
		return u.BadResponse(ctx,trans.E("can not delete role"))
	}
	return u.OKResponse(
		ctx,
		role,
	)

}