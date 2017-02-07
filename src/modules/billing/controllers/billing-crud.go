package billing

import (
	"common/assert"
	"modules/billing/bil"
	"modules/misc/trans"
	"modules/telegram/ad/ads"
	"modules/user/aaa"
	"modules/user/middlewares"
	"net/http"
	"strconv"
	"time"

	"github.com/Sirupsen/logrus"
	"gopkg.in/labstack/echo.v3"
)

type weeklyReport struct {
	ChannelName string           `json:"name"`
	Detail      map[string]int64 `json:"detail"`
}

type weeklyReportArr struct {
	Weeklyreport []weeklyReport `json:"report"`
}

//	dashboard shows views per channel
//	@Route	{
//		url	=	/dashboard/:id
//		method	= get
//		resource = create_ad:self
//		middleware = authz.Authenticate
//		200 = weeklyReportArr
//		400 = base.ErrorResponseSimple
//	}
func (u *Controller) dashboard(ctx echo.Context) error {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 0)
	if err != nil {
		return u.NotFoundResponse(ctx, err)
	}
	var channels []ads.Channel
	var res weeklyReportArr

	m := ads.NewAdsManager()

	channels, err = m.FindChannelsByUserID(id)
	if err != nil {
		return u.NotFoundResponse(ctx, err)
	}

	for i := range channels {
		reports, err := m.GetChanDailyViewByID(channels[i].ID)
		assert.Nil(err)

		temp := weeklyReport{
			ChannelName: channels[i].Name,
			Detail:      map[string]int64{},
		}
		for j := range reports {
			temp.Detail[reports[j].Date.Format(time.RFC3339)] = reports[j].View
		}
		res.Weeklyreport = append(res.Weeklyreport, temp)
		logrus.Warn(res)
	}

	return u.OKResponse(ctx, res)
}

//	billing billing for ad
//	@Route	{
//		url	=	/billing/:id
//		method	= get
//		resource = get_billing:self
//		middleware = authz.Authenticate
//		200 = bil.Billing
//		400 = base.ErrorResponseSimple
//	}
func (u *Controller) billing(ctx echo.Context) error {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 0)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	m := bil.NewBilManager()
	currentBilling, err := m.FindBillingByID(id)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	currentUser := authz.MustGetUser(ctx)
	owner, err := aaa.NewAaaManager().FindUserByID(currentBilling.UserID)
	assert.Nil(err)
	_, b := currentUser.HasPermOn("get_billing", owner.ID, owner.DBParentID.Int64)
	if !b {
		return ctx.JSON(http.StatusForbidden, trans.E("user can't access"))
	}
	return u.OKResponse(ctx, currentBilling)
}

//	payment payment for payment
//	@Route	{
//		url	=	/payment/:id
//		method	= get
//		resource = get_payment:self
//		middleware = authz.Authenticate
//		200 = bil.Payment
//		400 = base.ErrorResponseSimple
//	}
func (u *Controller) payment(ctx echo.Context) error {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 0)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	m := bil.NewBilManager()
	currentPayment, err := m.FindPaymentByID(id)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	currentUser := authz.MustGetUser(ctx)
	owner, err := aaa.NewAaaManager().FindUserByID(currentPayment.UserID)
	assert.Nil(err)
	_, b := currentUser.HasPermOn("get_payment", owner.ID, owner.DBParentID.Int64)
	if !b {
		return ctx.JSON(http.StatusForbidden, trans.E("user can't access"))
	}
	return u.OKResponse(ctx, currentPayment)
}
