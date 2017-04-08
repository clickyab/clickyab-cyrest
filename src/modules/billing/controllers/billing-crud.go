package billing

import (
	"common/assert"
	"encoding/json"
	"modules/billing/bil"
	"modules/misc/trans"
	"modules/telegram/ad/ads"
	"modules/user/aaa"
	"modules/user/middlewares"
	"net/http"
	"strconv"

	"common/models/common"
	"time"

	"modules/misc/middlewares"

	"math"

	"common/config"
	"common/mail"

	"gopkg.in/labstack/echo.v3"
)

// WeeklyReport show
type WeeklyReport struct {
	ChannelName string   `json:"name"`
	Report      []Report `json:"report"`
}

// Report shows view if its ended
type Report struct {
	View int64           `json:"view"`
	End  common.NullTime `json:"end"`
}

type weeklyReportArr struct {
	Weeklyreport []WeeklyReport `json:"chandetail"`
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
		temp, err := m.GetChanViewByID(channels[i].ID)
		assert.Nil(err)

		wk := WeeklyReport{}
		wk.ChannelName = channels[i].Name
		for k := range temp {
			rep := Report{}
			rep.View = temp[k].View
			if time.Now().Before(temp[k].End) {
				rep.End = common.NullTime{Time: temp[k].End}
				rep.End.Valid = true
			}
			wk.Report = append(wk.Report, rep)
		}
		res.Weeklyreport = append(res.Weeklyreport, wk)
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

// @Validate {
// }
type bilChangeStatusPayload struct {
	Status   bil.BillingStatus `json:"status" validate:"required" error:"Status is required"`
	Describe string            `json:"describe"`
}

// Validate custom validation for billing Status
func (lp *bilChangeStatusPayload) ValidateExtra(ctx echo.Context) error {
	if !lp.Status.IsValid() {
		return middlewares.GroupError{
			"Status": trans.E("Status is invalid"),
		}
	}
	return nil
}

//	changeStatus change Status for withdrawal
//	@Route	{
//		url = /list/status/:id
//		method = put
//		resource = change_status_billing:global
//		payload	= bilChangeStatusPayload
//		middleware = authz.Authenticate
//		200 = bil.Billing
//		400 = base.ErrorResponseSimple
//	}
func (u *Controller) changeStatus(ctx echo.Context) error {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 0)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	pl := u.MustGetPayload(ctx).(*bilChangeStatusPayload)
	m := bil.NewBilManager()
	currentBil, err := m.FindBillingByID(id)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	currentUser := authz.MustGetUser(ctx)
	owner, err := aaa.NewAaaManager().FindUserByID(currentBil.UserID)
	assert.Nil(err)
	_, b := currentUser.HasPermOn("change_status_billing", owner.ID, owner.DBParentID.Int64)
	if !b {
		return ctx.JSON(http.StatusForbidden, trans.E("user can't access"))
	}
	currentBil.Status = pl.Status
	currentBil.Reason = common.MakeNullString(pl.Describe)
	assert.Nil(m.UpdateBilling(currentBil))
	if currentBil.Status == bil.BilStatusAccepted && currentBil.Type == bil.BilTypeWithdrawal {
		mail.SendByTemplateName(trans.T("withdrawal accepted").Translate("fa_IR"), "withdrawal", struct {
			Name  string
			Price int64
		}{
			Name:  owner.Email,
			Price: currentBil.Amount,
		}, config.Config.Mail.From, owner.Email)
	}
	if currentBil.Status == bil.BilStatusRejected && currentBil.Type == bil.BilTypeWithdrawal {
		mail.SendByTemplateName(trans.T("withdrawal rejected").Translate("fa_IR"), "rejectWitdrawal", struct {
			Name string
		}{
			Name: owner.Email,
		}, config.Config.Mail.From, owner.Email)
	}
	return u.OKResponse(ctx, currentBil)
}

// @Validate {
// }
type bilChangeDepositPayload struct {
	Deposit  bil.BillingDeposit `json:"deposit" validate:"required" error:"Deposit is required"`
	Describe string             `json:"describe"`
}

// Validate custom validation for billing Status
func (lp *bilChangeDepositPayload) ValidateExtra(ctx echo.Context) error {
	if !lp.Deposit.IsValid() {
		return middlewares.GroupError{
			"Status": trans.E("Deposit is invalid"),
		}
	}
	return nil
}

//	changeDeposit change Deposit for withdrawal
//	@Route	{
//		url = /list/deposit/:id
//		method = put
//		resource = change_deposit_billing:global
//		payload	= bilChangeDepositPayload
//		middleware = authz.Authenticate
//		200 = bil.Billing
//		400 = base.ErrorResponseSimple
//	}
func (u *Controller) changeDeposit(ctx echo.Context) error {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 0)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	pl := u.MustGetPayload(ctx).(*bilChangeDepositPayload)
	m := bil.NewBilManager()
	currentBil, err := m.FindBillingByID(id)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	if currentBil.Type != bil.BilTypeWithdrawal {
		return ctx.JSON(http.StatusBadRequest, trans.E("this route just for withdrawal"))
	}
	currentUser := authz.MustGetUser(ctx)
	owner, err := aaa.NewAaaManager().FindUserByID(currentBil.UserID)
	assert.Nil(err)
	_, b := currentUser.HasPermOn("change_deposit_billing", owner.ID, owner.DBParentID.Int64)
	if !b {
		return ctx.JSON(http.StatusForbidden, trans.E("user can't access"))
	}
	currentBil.Deposit = pl.Deposit
	currentBil.Reason = common.MakeNullString(pl.Describe)
	assert.Nil(m.UpdateBilling(currentBil))
	return u.OKResponse(ctx, currentBil)
}

// @Validate {
// }
type bilCreatePayload struct {
	Amount int64 `json:"amount" validate:"required" error:"Amount is required"`
	UserID int64 `json:"user_id"`
}

//	createWithdrawal create withdrawal
//	@Route	{
//		url = /withdrawal
//		method = post
//		resource = request_withdrawal:self
//		payload	= bilCreatePayload
//		middleware = authz.Authenticate
//		200 = bil.Billing
//		400 = base.ErrorResponseSimple
//	}
func (u *Controller) createWithdrawal(ctx echo.Context) error {

	m := bil.NewBilManager()
	err := m.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			assert.Nil(m.Rollback())
		} else {
			err = m.Commit()
		}

	}()
	pl := u.MustGetPayload(ctx).(*bilCreatePayload)
	/*if pl.Amount < bcfg.Bcfg.Withdrawal.MinWithdrawal {
		return ctx.JSON(http.StatusBadRequest, trans.E("you can not withdrawal under %d", bcfg.Bcfg.Withdrawal.MinWithdrawal))
	}*/
	currentUser := authz.MustGetUser(ctx)
	sumBil := m.SumBilling(currentUser.ID)

	if sumBil.Balance.Int64 < pl.Amount {
		return ctx.JSON(http.StatusBadRequest, trans.E("your withdrawal under your billing"))
	}
	if pl.UserID == 0 {
		pl.UserID = currentUser.ID
	}
	owner, err := aaa.NewAaaManager().FindUserByID(pl.UserID)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	_, b := currentUser.HasPermOn("request_withdrawal", owner.ID, owner.DBParentID.Int64)
	if !b {
		return ctx.JSON(http.StatusForbidden, trans.E("user can't access"))
	}

	withdrawal := bil.Billing{}
	f := float64(pl.Amount)
	f = math.Abs(f) * -1
	pl.Amount = int64(f)
	withdrawal.Amount = pl.Amount
	withdrawal.Type = bil.BilTypeWithdrawal
	withdrawal.Status = bil.BilStatusPending
	withdrawal.Deposit = bil.BilDepositNo
	withdrawal.UserID = pl.UserID

	err = m.CreateBilling(&withdrawal)
	if err != nil {
		return nil
	}
	jsonBilling, err := json.Marshal(withdrawal)
	if err != nil {
		return err
	}
	bilDetail := bil.BillingDetail{
		UserID:    currentUser.ID,
		BillingID: withdrawal.ID,
		Reason:    common.MakeNullString(string(jsonBilling)),
	}
	err = m.CreateBillingDetail(&bilDetail)
	if err != nil {
		return nil
	}

	return u.OKResponse(ctx, withdrawal)
}

type showBilling struct {
	Billing int64
}

//	showBilling show billing
//	@Route	{
//		url = /billing/show
//		method = get
//		resource = get_billing:self
//		middleware = authz.Authenticate
//		200 = showBilling
//		400 = base.ErrorResponseSimple
//	}
func (u *Controller) showBilling(ctx echo.Context) error {
	currentUser := authz.MustGetUser(ctx)
	m := bil.NewBilManager()
	sumBil := m.SumBilling(currentUser.ID)
	return u.OKResponse(ctx, showBilling{Billing: sumBil.Balance.Int64})
}
