package plan

import (
	"common/assert"

	"modules/telegram/ad/ads"

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
	plns, err := m.GetAllActivePlans()
	assert.Nil(err)
	return u.OKResponse(ctx, plns)
}
