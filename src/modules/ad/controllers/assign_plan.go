package ad

import (
	"common/assert"
	"modules/ad/ads"
	"modules/plan/pln"

	"common/models/common"
	"modules/misc/trans"
	"modules/user/aaa"
	"modules/user/middlewares"
	"net/http"

	"gopkg.in/labstack/echo.v3"
)

// @Validate {
// }
type assignPlanPayload struct {
	AdID   int64 `json:"Ad_id" validate:"required"`
	PlanID int64 `json:"plan_id" vaSlidate:"required"`
}

// assignPlan
// @Route {
//		url	=	/assign/plan
//		method	=	post
//		payload	=	assignPlanPayload
//		resource=	assign_plan:self
//		middleware = authz.Authenticate
//		200	=	ads.Ad
//		400	=	base.ErrorResponseSimple
//// }
func (u *Controller) assignPlan(ctx echo.Context) error {
	pl := u.MustGetPayload(ctx).(*assignPlanPayload)
	plnManager := pln.NewPlnManager()
	adManager := ads.NewAdsManager()

	//find plan
	plan, err := plnManager.FindPlanByID(pl.PlanID)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	//find ads
	ads, err := adManager.FindAdByID(pl.AdID)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}

	//find owner of ads
	owner, err := aaa.NewAaaManager().FindUserByID(ads.UserID)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	//current user
	currentUser, ok := authz.GetUser(ctx)
	if !ok {
		return u.NotFoundResponse(ctx, nil)
	}

	//check current user has permisiion
	_, b := currentUser.HasPermOn("assign_plan", owner.ID, owner.DBParentID.Int64)
	if !b {
		return ctx.JSON(http.StatusForbidden, trans.E("user can't access"))
	}
	ads.PlanID = common.NullInt64{Valid: true, Int64: plan.ID}
	assert.Nil(adManager.UpdateAd(ads))

	return u.OKResponse(ctx, ads)

}
