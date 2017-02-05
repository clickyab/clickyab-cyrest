package billing

import (
	"common/assert"
	"modules/billing/bil"
	"modules/misc/trans"
	"modules/user/aaa"
	"modules/user/middlewares"
	"net/http"
	"strconv"

	"gopkg.in/labstack/echo.v3"
)

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
