package ad

import (
	"modules/telegram/ad/ads"

	"modules/misc/base"
	"modules/user/middlewares"

	"gopkg.in/labstack/echo.v3"
)

//	publisherTotalViewChart publisher total view chart plan
//	@Route	{
//		url	=	/chart/pubtotalview
//		method	= get
//		resource = get_ad_chart:self
//		middleware = authz.Authenticate
//		200 = ads.PubDashboardTotalView
//		400 = base.ErrorResponseSimple
//	}
func (u *Controller) publisherTotalViewChart(ctx echo.Context) error {
	m := ads.NewAdsManager()
	currentUser := authz.MustGetUser(ctx)
	scope, _ := currentUser.HasPerm(base.ScopeGlobal, "get_ad_chart")
	pubChartData := m.PubDashboardTotalView(currentUser.ID, scope)
	return u.OKResponse(ctx, pubChartData)
}
