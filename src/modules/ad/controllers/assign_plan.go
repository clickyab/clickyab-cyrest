package ad

import (
	"common/assert"
	"modules/ad/ads"
	"modules/plan/pln"

	"common/models/common"
	"modules/user/aaa"
	"modules/user/middlewares"

	"gopkg.in/labstack/echo.v3"
)

// @Validate {
// }
type assignPlanPayload struct {
	AdID   int64 `json:"Ad_id" validate:"required"`
	PlanID int64 `json:"plan_id" validate:"required"`
}

// assignPlan
// @Route {
//		url	=	/assign/plan
//		method	=	post
//		payload	=	assignPlanPayload
//		resource=	assign_plan:self
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
	assert.Nil(err)

	//find owner of ads
	owner, err := aaa.NewAaaManager().FindUserByID(ads.UserID)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	//current user
	currentUser, _ := authz.GetUser(ctx)

	//check current user has permisiion
	_, b := currentUser.HasPermOn("assign_plan", owner.ID, owner.DBParentID.Int64)
	if !b {
		return u.ForbiddenResponse(ctx, nil)

	}
	ads.PlanID = common.NullInt64{Valid: true, Int64: plan.ID}
	assert.Nil(adManager.UpdateAd(ads))

	return u.OKResponse(ctx, ads)

}
