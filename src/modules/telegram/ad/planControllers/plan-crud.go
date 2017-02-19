package plan

import (
	"common/assert"

	"modules/telegram/ad/ads"

	"strconv"

	"modules/misc/trans"
	"modules/user/aaa"
	"modules/user/middlewares"
	"net/http"

	"gopkg.in/labstack/echo.v3"
)

// @Validate {
// }
type planPayload struct {
	AdID int64 `json:"ad_id" validate:"required" error:"ad id is required"`
}

//	allPlan get all plans
//	@Route	{
//	url	=	/
//	method	= get
//	resource = get_plan:self
//	middleware = authz.Authenticate
//	200 = plans
//	400 = base.ErrorResponseSimple
//	}
func (u *Controller) allPlan(ctx echo.Context) error {
	m := ads.NewAdsManager()
	plans, err := m.GetAllActivePlans()
	assert.Nil(err)
	return u.OKResponse(ctx, plans)
}

//	allIndividualPlan get all individual plans
//	@Route	{
//	url	=	/individual
//	method	= get
//	resource = get_plan:self
//	middleware = authz.Authenticate
//	200 = plans
//	400 = base.ErrorResponseSimple
//	}
func (u *Controller) allIndividualPlan(ctx echo.Context) error {
	m := ads.NewAdsManager()
	plans, err := m.GetAllIndividualActivePlans()
	assert.Nil(err)
	return u.OKResponse(ctx, plans)
}

//	getAllAppropriatePlan get all appropriate plans
//	@Route	{
//	url	=	/appropriate
//	method	= get
//	resource = get_plan:self
//	middleware = authz.Authenticate
//	200 = plans
//	400 = base.ErrorResponseSimple
//	}
func (u *Controller) getAllAppropriatePlan(ctx echo.Context) error {
	var plans []ads.Plan
	pl := u.MustGetPayload(ctx).(*planPayload)
	m := ads.NewAdsManager()
	currentAd, err := m.FindAdByID(pl.AdID)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	currentUser := authz.MustGetUser(ctx)
	owner, err := aaa.NewAaaManager().FindUserByID(currentAd.UserID)
	assert.Nil(err)
	_, b := currentUser.HasPermOn("get_plan", owner.ID, owner.DBParentID.Int64)
	if !b {
		return ctx.JSON(http.StatusForbidden, trans.E("user can't access"))
	}
	if currentAd.CliMessageID.Valid { //get forwarded plan
		plans, err = m.GetAllPromotionActivePlans()
		assert.Nil(err)
	} else {
		plans, err = m.GetAllIndividualActivePlans()
		assert.Nil(err)
	}
	return u.OKResponse(ctx, plans)
}

//	allPromotionPlan get all promotion plans
//	@Route	{
//	url	=	/promotion
//	method	= get
//	resource = get_plan:self
//	middleware = authz.Authenticate
//	200 = plans
//	400 = base.ErrorResponseSimple
//	}
func (u *Controller) allPromotionPlan(ctx echo.Context) error {
	m := ads.NewAdsManager()
	plans, err := m.GetAllPromotionActivePlans()
	assert.Nil(err)
	return u.OKResponse(ctx, plans)
}

//	allPlan get all plans
//	@Route	{
//	url	=	/:id
//	method	= get
//	resource = get_plan:self
//	middleware = authz.Authenticate
//	200 = plans
//	400 = base.ErrorResponseSimple
//	}
func (u *Controller) getPlan(ctx echo.Context) error {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 0)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	m := ads.NewAdsManager()
	plan, err := m.FindPlanByID(id)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	return u.OKResponse(ctx, plan)
}
