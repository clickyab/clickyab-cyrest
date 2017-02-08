package ad

import (
	"common/assert"
	"common/models/common"
	"modules/file/fila"
	"modules/misc/middlewares"
	"modules/misc/trans"
	"modules/telegram/ad/ads"
	"modules/user/aaa"
	"modules/user/middlewares"
	"net/http"
	"net/url"
	"strconv"

	"common/rabbit"
	"modules/telegram/cyborg/commands"

	"fmt"

	"modules/billing/config"

	"modules/billing/bil"

	"strings"

	"common/payment"

	"path"

	"modules/file/config"

	echo "gopkg.in/labstack/echo.v3"
)

// @Validate {
// }
type adPayload struct {
	Name string `json:"name" validate:"required" error:"name is required"`
}

// @Validate {
// }
type adPromotePayload struct {
	CliMessageID string `json:"cli_message_id" validate:"required" error:"cli_message_id is required"`
}

// @Validate {
// }
type adUploadPayload struct {
	URL string `json:"url" validate:"required" error:"url is required"`
}

// @Validate {
// }
type adAdminStatusPayload struct {
	AdAdminStatus ads.AdAdminStatus `json:"admin_status" validate:"required" error:"status is required"`
}

// @Validate {
// }
type adDescriptionPayLoad struct {
	Body string `json:"body" validate:"required" error:"body is required"`
}

// @Validate {
// }
type adPlanPayLoad struct {
	ID int64 `json:"plan_id" validate:"required" error:"plan is required"`
}

// @Validate {
// }
type adArchiveStatusPayload struct {
	AdArchiveStatus ads.AdArchiveStatus `json:"archive_status" validate:"required" error:"status is required"`
}

// @Validate {
// }
type adActiveStatusPayload struct {
	AdActiveStatus ads.AdActiveStatus `json:"active_status" validate:"required" error:"status is required"`
}

// Validate custom validation for user scope
func (lp *adActiveStatusPayload) ValidateExtra(ctx echo.Context) error {
	if !lp.AdActiveStatus.IsValid() {
		return middlewares.GroupError{
			"status": trans.E("status is invalid"),
		}
	}
	return nil
}

// Validate custom validation for user scope
func (lp *adArchiveStatusPayload) ValidateExtra(ctx echo.Context) error {
	if !lp.AdArchiveStatus.IsValid() {
		return middlewares.GroupError{
			"status": trans.E("status is invalid"),
		}
	}
	return nil
}

// Validate custom validation for user scope
func (lp *adAdminStatusPayload) ValidateExtra(ctx echo.Context) error {
	if !lp.AdAdminStatus.IsValid() {
		return middlewares.GroupError{
			"status": trans.E("status is invalid"),
		}
	}
	return nil
}

//	create create ad
//	@Route	{
//		url	=	/
//		method	= post
//		payload	= adPayload
//		resource = create_ad:self
//		middleware = authz.Authenticate
//		200 = ads.Ad
//		400 = base.ErrorResponseSimple
//	}
func (u *Controller) create(ctx echo.Context) error {
	pl := u.MustGetPayload(ctx).(*adPayload)
	m := ads.NewAdsManager()
	currentUser := authz.MustGetUser(ctx)
	newAd := &ads.Ad{
		Name:            pl.Name,
		AdArchiveStatus: ads.AdArchiveStatusNo,
		AdPayStatus:     ads.AdPayStatusNo,
		AdAdminStatus:   ads.AdAdminStatusPending,
		AdActiveStatus:  ads.AdActiveStatusNo,
		UserID:          currentUser.ID,
	}
	assert.Nil(m.CreateAd(newAd))
	return u.OKResponse(ctx, newAd)
}

//	changeAdminStatus change admin status for ad
//	@Route	{
//		url	=	/list/admin_status/:id
//		method	= put
//		payload	= adAdminStatusPayload
//		resource = change_admin_ad:parent
//		middleware = authz.Authenticate
//		200 = ads.Ad
//		400 = base.ErrorResponseSimple
//	}
func (u *Controller) changeAdminStatus(ctx echo.Context) error {
	pl := u.MustGetPayload(ctx).(*adAdminStatusPayload)
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 0)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	m := ads.NewAdsManager()
	currentAd, err := m.FindAdByID(id)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	currentUser := authz.MustGetUser(ctx)
	owner, err := aaa.NewAaaManager().FindUserByID(currentAd.UserID)
	assert.Nil(err)

	_, b := currentUser.HasPermOn("change_admin_ad", owner.ID, owner.DBParentID.Int64)
	if !b {
		return ctx.JSON(http.StatusForbidden, trans.E("user can't access"))
	}
	currentAd.AdAdminStatus = pl.AdAdminStatus
	assert.Nil(m.UpdateAd(currentAd))

	return u.OKResponse(ctx, currentAd)
}

//	changeArchiveStatus change archive status for ad
//	@Route	{
//		url = /list/archive_status/:id
//		method = put
//		resource = change_archive_ad:self
//		payload	= adArchiveStatusPayload
//		middleware = authz.Authenticate
//		200 = ads.Ad
//		400 = base.ErrorResponseSimple
//	}
func (u *Controller) changeArchiveStatus(ctx echo.Context) error {
	pl := u.MustGetPayload(ctx).(*adArchiveStatusPayload)
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 0)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	m := ads.NewAdsManager()
	currentAd, err := m.FindAdByID(id)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	currentUser := authz.MustGetUser(ctx)
	owner, err := aaa.NewAaaManager().FindUserByID(currentAd.UserID)
	assert.Nil(err)
	_, b := currentUser.HasPermOn("change_admin_ad", owner.ID, owner.DBParentID.Int64)
	if !b {
		return ctx.JSON(http.StatusForbidden, trans.E("user can't access"))
	}
	if pl.AdArchiveStatus == ads.AdArchiveStatusNo && currentAd.AdArchiveStatus == ads.AdArchiveStatusYes {
		currentAd.AdArchiveStatus = ads.AdArchiveStatusNo
	} else if pl.AdArchiveStatus == ads.AdArchiveStatusYes && currentAd.AdArchiveStatus == ads.AdArchiveStatusNo {
		currentAd.AdArchiveStatus = ads.AdArchiveStatusYes
	}
	assert.Nil(m.UpdateAd(currentAd))
	return u.OKResponse(ctx, currentAd)
}

//	addDescription add description to ad
//	@Route	{
//		url	=	/desc/:id
//		method	= put
//		payload	= adDescriptionPayLoad
//		resource = add_description_ad:self
//		middleware = authz.Authenticate
//		200 = ads.Ad
//		400 = base.ErrorResponseSimple
//	}
func (u *Controller) addDescription(ctx echo.Context) error {
	pl := u.MustGetPayload(ctx).(*adDescriptionPayLoad)
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 0)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	m := ads.NewAdsManager()
	currentAd, err := m.FindAdByID(id)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	currentUser := authz.MustGetUser(ctx)
	owner, err := aaa.NewAaaManager().FindUserByID(currentAd.UserID)
	assert.Nil(err)
	_, b := currentUser.HasPermOn("add_description_ad", owner.ID, owner.DBParentID.Int64)
	if !b {
		return ctx.JSON(http.StatusForbidden, trans.E("user can't access"))
	}
	currentAd.Description = common.MakeNullString(pl.Body)
	currentAd.CliMessageID = common.MakeNullString("")
	assert.Nil(m.UpdateAd(currentAd))
	return u.OKResponse(ctx, currentAd)
}

//	uploadBanner uploadBanner for ad
//	@Route	{
//		url	=	/upload/:id
//		method	= put
//		payload	= adUploadPayload
//		resource = upload_ad:self
//		middleware = authz.Authenticate
//		200 = ads.Ad
//		400 = base.ErrorResponseSimple
//	}
func (u *Controller) uploadBanner(ctx echo.Context) error {
	pl := u.MustGetPayload(ctx).(*adUploadPayload)
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 0)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	dURL := pl.URL
	_, err = url.Parse(dURL)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	m := ads.NewAdsManager()
	currentAd, err := m.FindAdByID(id)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	currentUser := authz.MustGetUser(ctx)
	owner, err := aaa.NewAaaManager().FindUserByID(currentAd.UserID)
	assert.Nil(err)
	_, b := currentUser.HasPermOn("upload_ad", owner.ID, owner.DBParentID.Int64)
	if !b {
		return ctx.JSON(http.StatusForbidden, trans.E("user can't access"))
	}

	//upload
	file, err := fila.CheckUpload(dURL, currentUser.ID)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	currentAd.Src = common.MakeNullString(file)
	currentAd.CliMessageID = common.MakeNullString("")
	assert.Nil(m.UpdateAd(currentAd))
	if currentAd.Src.Valid {
		g, err := url.Parse(fcfg.Fcfg.File.UploadPath)
		if err != nil {
			return u.NotFoundResponse(ctx, nil)
		}
		g.Path = path.Join(g.Path, currentAd.Src.String)
		s := g.String()
		currentAd.Src = common.MakeNullString(s)
	}
	return u.OKResponse(ctx, currentAd)
}

//	promoteAd promoteAd for ad
//	@Route	{
//		url	=	/promote/:id
//		method	= put
//		payload	= adPromotePayload
//		resource = promote_ad:self
//		middleware = authz.Authenticate
//		200 = ads.Ad
//		400 = base.ErrorResponseSimple
//	}
func (u *Controller) promoteAd(ctx echo.Context) error {
	pl := u.MustGetPayload(ctx).(*adPromotePayload)
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 0)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	m := ads.NewAdsManager()
	currentAd, err := m.FindAdByID(id)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	currentUser := authz.MustGetUser(ctx)
	owner, err := aaa.NewAaaManager().FindUserByID(currentAd.UserID)
	assert.Nil(err)
	_, b := currentUser.HasPermOn("promote_ad", owner.ID, owner.DBParentID.Int64)
	if !b {
		return ctx.JSON(http.StatusForbidden, trans.E("user can't access"))
	}
	currentAd.CliMessageID = common.MakeNullString(pl.CliMessageID)
	currentAd.Src = common.MakeNullString("")
	currentAd.BotChatID = common.NullInt64{}
	currentAd.BotMessageID = common.NullInt64{}
	currentAd.Description = common.MakeNullString("")
	assert.Nil(m.UpdateAd(currentAd))
	return u.OKResponse(ctx, currentAd)
}

//	changeActiveStatus change active status for ad
//	@Route	{
//		url = /list/active_status/:id
//		method = put
//		payload	= adActiveStatusPayload
//		resource = change_active_ad:self
//		middleware = authz.Authenticate
//		200 = ads.Ad
//		400 = base.ErrorResponseSimple
//	}
func (u *Controller) changeActiveStatus(ctx echo.Context) error {
	pl := u.MustGetPayload(ctx).(*adActiveStatusPayload)
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 0)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	m := ads.NewAdsManager()
	currentAd, err := m.FindAdByID(id)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	currentUser := authz.MustGetUser(ctx)
	owner, err := aaa.NewAaaManager().FindUserByID(currentAd.UserID)
	assert.Nil(err)
	_, b := currentUser.HasPermOn("change_active_ad", owner.ID, owner.DBParentID.Int64)
	if !b {
		return ctx.JSON(http.StatusForbidden, trans.E("user can't access"))
	}
	//check everything is good TODO: check pay status later
	if currentAd.AdAdminStatus == "accepted" && currentAd.Name != "" && pl.AdActiveStatus == ads.AdActiveStatusYes && currentAd.AdActiveStatus == ads.AdActiveStatusNo {
		currentAd.AdActiveStatus = ads.AdActiveStatusYes
		//check to add job
		if currentAd.CliMessageID.Valid {
			rabbit.MustPublish(commands.IdentifyAD{AdID: currentAd.ID})
		}
	} else if pl.AdActiveStatus == ads.AdActiveStatusNo && currentAd.AdActiveStatus == ads.AdActiveStatusYes {
		currentAd.AdActiveStatus = ads.AdActiveStatusNo
	}
	assert.Nil(m.UpdateAd(currentAd))
	return u.OKResponse(ctx, currentAd)
}

//	edit edit ad
//	@Route	{
//		url	=	/:id
//		method	= put
//		payload	= adPayload
//		resource = edit_ad:self
//		middleware = authz.Authenticate
//		200 = ads.Ad
//		400 = base.ErrorResponseSimple
//	}
func (u *Controller) edit(ctx echo.Context) error {
	pl := u.MustGetPayload(ctx).(*adPayload)
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 0)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	m := ads.NewAdsManager()
	currentAd, err := m.FindAdByID(id)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	currentUser := authz.MustGetUser(ctx)
	owner, err := aaa.NewAaaManager().FindUserByID(currentAd.UserID)
	assert.Nil(err)
	_, b := currentUser.HasPermOn("edit_ad", owner.ID, owner.DBParentID.Int64)
	if !b {
		return ctx.JSON(http.StatusForbidden, trans.E("user can't access"))
	}
	//check if it can be edited
	if currentAd.AdActiveStatus == ads.AdActiveStatusYes || currentAd.AdAdminStatus == ads.AdAdminStatusAccepted {
		return ctx.JSON(http.StatusForbidden, trans.E("user can't access"))
	}
	currentAd.Name = pl.Name
	assert.Nil(m.UpdateAd(currentAd))
	return u.OKResponse(ctx, currentAd)
}

//	getAd getAd for ad
//	@Route	{
//		url	=	/get/:id
//		method	= get
//		resource = get_ad:self
//		middleware = authz.Authenticate
//		200 = ads.Ad
//		400 = base.ErrorResponseSimple
//	}
func (u *Controller) getAd(ctx echo.Context) error {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 0)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	m := ads.NewAdsManager()
	currentAd, err := m.FindAdByID(id)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	currentUser := authz.MustGetUser(ctx)
	owner, err := aaa.NewAaaManager().FindUserByID(currentAd.UserID)
	assert.Nil(err)
	_, b := currentUser.HasPermOn("promote_ad", owner.ID, owner.DBParentID.Int64)
	if !b {
		return ctx.JSON(http.StatusForbidden, trans.E("user can't access"))
	}
	if currentAd.Src.Valid {
		g, err := url.Parse(fcfg.Fcfg.File.UploadPath)
		if err != nil {
			return u.NotFoundResponse(ctx, nil)
		}
		g.Path = path.Join(g.Path, currentAd.Src.String)
		s := g.String()
		currentAd.Src = common.MakeNullString(s)
	}
	return u.OKResponse(ctx, currentAd)
}

//	assignPlan assignPlan for ad
//	@Route	{
//		url	=	/plan/:id
//		method	= put
//		payload	= adPlanPayLoad
//		resource = assign_plan:self
//		middleware = authz.Authenticate
//		200 = ads.Ad
//		400 = base.ErrorResponseSimple
//	}
func (u *Controller) assignPlan(ctx echo.Context) error {
	m := ads.NewAdsManager()
	pl := u.MustGetPayload(ctx).(*adPlanPayLoad)
	//find plan
	plan, err := m.FindPlanByID(pl.ID)
	if err != nil || plan.Active != ads.ActiveStatusYes {
		return u.NotFoundResponse(ctx, nil)
	}
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 0)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	currentAd, err := m.FindAdByID(id)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	currentUser := authz.MustGetUser(ctx)
	owner, err := aaa.NewAaaManager().FindUserByID(currentAd.UserID)
	assert.Nil(err)
	_, b := currentUser.HasPermOn("assign_plan", owner.ID, owner.DBParentID.Int64)
	if !b {
		return ctx.JSON(http.StatusForbidden, trans.E("user can't access"))
	}
	//check if he/she can assign plan (flow)
	if currentAd.AdPayStatus != ads.AdPayStatusNo {
		return ctx.JSON(http.StatusForbidden, trans.E("user can't access"))
	}
	currentAd.PlanID = common.NullInt64{Valid: true, Int64: plan.ID}
	assert.Nil(m.UpdateAd(currentAd))
	return u.OKResponse(ctx, currentAd)
}

//	charge charge user plan
//	@Route	{
//		url	=	/pay/:ad_id
//		method	= get
//		resource = pay_ad:self
//		middleware = authz.Authenticate
//		200 = ads.Ad
//		400 = base.ErrorResponseSimple
//	}
func (u *Controller) charge(ctx echo.Context) error {
	//find ad
	currentUser := authz.MustGetUser(ctx)
	id, err := strconv.ParseInt(ctx.Param("ad_id"), 10, 0)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	adManager := ads.NewAdsManager()
	bilManager := bil.NewBilManager()
	currentAd, err := adManager.FindAdByID(id)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	//set callback url
	callbackURL := fmt.Sprintf("%s%s", ctx.Scheme()+"://", path.Join(ctx.Request().Host, bcfg.Bcfg.Gate.CallbackURL, fmt.Sprintf("%d", currentAd.ID)))
	//find plan
	plan, err := adManager.FindPlanByID(currentAd.PlanID.Int64)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}

	price := plan.Price
	client := payment.NewPaymentGatewayImplementationServicePortType("", false, nil)
	// Create a new payment request to Zarinpal
	resp, err := client.Request(&payment.Request{
		MerchantID:  bcfg.Bcfg.Gate.MerchantID,
		Amount:      price,
		Description: bcfg.Bcfg.Gate.Description,
		Email:       bcfg.Bcfg.Gate.Email,
		Mobile:      bcfg.Bcfg.Gate.Mobile,
		CallbackURL: callbackURL,
	})
	if err != nil {
		return u.BadResponse(ctx, nil)
	}

	if resp.Status == bcfg.Bcfg.Gate.MerchantOkStatus {
		payment := &bil.Payment{
			UserID:    currentUser.ID,
			Amount:    price,
			Status:    bil.PayStatusPending,
			Authority: common.MakeNullString(resp.Authority),
		}
		assert.Nil(bilManager.CreatePayment(payment))
		return u.OKResponse(ctx, fmt.Sprintf("%s%s", bcfg.Bcfg.Gate.ZarinURL, resp.Authority))
	}
	return u.BadResponse(ctx, nil)
}

//	charge charge user plan
//	@Route	{
//		url	=	/verify/:id
//		method	= get
//		200 = ads.Ad
//		400 = base.ErrorResponseSimple
//	}
func (u *Controller) verify(ctx echo.Context) error {
	adManager := ads.NewAdsManager()
	frontURL := fmt.Sprintf("%s%s", ctx.Scheme()+"://", path.Join(ctx.Request().Host, bcfg.Bcfg.Gate.FrontCallbackURL))
	frontOk := fmt.Sprintf("%s%s", frontURL, "?success=yes&payment=")
	frontNOk := fmt.Sprintf("%s%s", frontURL, "?success=no")
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 0)
	if err != nil {
		return ctx.Redirect(http.StatusMovedPermanently, frontNOk)
	}
	currentAd, err := adManager.FindAdByID(id)
	if err != nil {
		return ctx.Redirect(http.StatusMovedPermanently, frontNOk)
	}
	plan, err := adManager.FindPlanByID(currentAd.PlanID.Int64)
	if err != nil {
		return ctx.Redirect(http.StatusMovedPermanently, frontNOk)
	}
	if strings.ToLower(ctx.FormValue("Status")) != "ok" {
		return ctx.Redirect(http.StatusMovedPermanently, frontNOk)
	}
	client := payment.NewPaymentGatewayImplementationServicePortType("", false, nil)
	// Create a new payment request to Zarinpal
	resp, err := client.Verification(&payment.Verification{
		MerchantID: bcfg.Bcfg.Gate.MerchantID,
		Amount:     plan.Price,
		Authority:  ctx.FormValue("Authority"),
	})
	// Check if response is error free
	if err != nil {
		return ctx.Redirect(http.StatusMovedPermanently, frontNOk)
	}
	if resp.Status == bcfg.Bcfg.Gate.MerchantOkStatus {
		billing, err := bil.NewBilManager().RegisterBilling(ctx.FormValue("Authority"), resp.RefID, plan.Price, resp.Status)
		if err != nil {
			return ctx.Redirect(http.StatusMovedPermanently, frontNOk)
		}
		//update ad pay status
		currentAd.AdPayStatus = ads.AdPayStatusYes
		assert.Nil(adManager.UpdateAd(currentAd))
		return ctx.Redirect(http.StatusMovedPermanently, fmt.Sprintf("%s%d", frontOk, billing.PaymentID.Int64))
		//call worker
		defer func() {
			rabbit.MustPublish(&commands.IdentifyAD{
				AdID: currentAd.ID,
			})
		}()
	}
	return ctx.Redirect(http.StatusMovedPermanently, frontNOk)

}
