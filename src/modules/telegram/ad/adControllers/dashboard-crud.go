package ad

import (
	"modules/telegram/ad/ads"
	"modules/user/middlewares"

	"gopkg.in/labstack/echo.v3"
)

//	getSpecificAd shows ad with specific details
//	@Route	{
//		url	=	/chart
//		method	= get
//		resource = get_ad_chart:self
//		middleware = authz.Authenticate
//		200 = ads.AdDashboard
//		400 = base.ErrorResponseSimple
//	}
func (u *Controller) adDashboardChart(ctx echo.Context) error {
	m := ads.NewAdsManager()
	currentUser := authz.MustGetUser(ctx)
	adChartData := m.PieChartAd(currentUser.ID)
	return u.OKResponse(ctx, adChartData)
}
