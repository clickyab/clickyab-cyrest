package ad

import (
	"common/assert"
	"modules/ad/ads"
	"modules/plan/pln"

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
	ads, err := adManager.FindPlanByID(pl.AdID)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}

	assert.Nil(plnManager.UpdatePlan(plan))

	return u.OKResponse(ctx, ads)

}
