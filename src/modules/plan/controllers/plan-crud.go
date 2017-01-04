package plan

import (
	"modules/channel/chn"
	"modules/misc/middlewares"
	"modules/misc/trans"

	echo "gopkg.in/labstack/echo.v3"
)

// @Validate {
// }
type planPayload struct {
	UserID int64  `json:"user_id"`
	Name   string `json:"name" validate:"required"`
	Link   string `json:"link" `
	Admin  string `json:"admin"`
}

//
////
////	getPlan
////	@Route	{
////	url	=	/:id
////	method	= get
////	resource = list_Plan:global
////	middleware = authz.Authenticate
////	200 = pln.Plan
////	400 = base.ErrorResponseSimple
////	}
//func (u *Controller) getPlan(ctx echo.Context) error {
//	id, err := strconv.ParseInt(ctx.Param("id"), 10, 0)
//	if err != nil {
//		return u.NotFoundResponse(ctx, nil)
//	}
//	m := chn.NewChnManager()
//	plan, err := m.FindChannelByID(id)
//	if err != nil {
//		return u.NotFoundResponse(ctx, nil)
//	}
//	owner, err := aaa.NewAaaManager().FindUserByID(plan.UserID)
//	if err != nil {
//		return u.NotFoundResponse(ctx, nil)
//	}
//	currentUser, ok := authz.GetUser(ctx)
//	if !ok {
//		return u.NotFoundResponse(ctx, nil)
//	}
//	_, b := currentUser.HasPermOn("list_plan", owner.ID, owner.DBParentID.Int64)
//	if !b {
//		return ctx.JSON(http.StatusForbidden, trans.E("user can't access"))
//	}
//	return u.OKResponse(ctx, plan)
//}

// @Validate {
// }
type statusPayload struct {
	Status chn.ChannelStatus `json:"status" validate:"required"`
}

// Validate custom validation for user scope
func (lp *statusPayload) ValidateExtra(ctx echo.Context) error {
	if !lp.Status.IsValid() {
		return middlewares.GroupError{
			"status": trans.E("is invalid"),
		}
	}
	return nil
}
