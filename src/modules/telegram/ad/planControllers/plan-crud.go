package plan

import (
	"common/assert"

	"modules/telegram/ad/ads"

	"strconv"

	"gopkg.in/labstack/echo.v3"
)

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
