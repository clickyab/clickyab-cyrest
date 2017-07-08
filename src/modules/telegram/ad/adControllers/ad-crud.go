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
	"path/filepath"
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

	"common/redis"

	"common/config"
	"common/mail"

	"modules/category/cat"

	"gopkg.in/labstack/echo.v3"
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
	AdminStatus ads.AdminStatus `json:"admin_status" validate:"required" error:"status is required"`
}

// @Validate {
// }
type adDescriptionPayLoad struct {
	Body common.MB4String `json:"body" validate:"required" error:"body is required"`
}

// @Validate {
// }
type adPlanPayLoad struct {
	ID int64 `json:"plan_id" validate:"required" error:"plan is required"`
}

// @Validate {
// }
type adArchiveStatusPayload struct {
	ActiveStatus ads.ActiveStatus `json:"archive_status" validate:"required" error:"status is required"`
}

// @Validate {
// }
type adActiveStatusPayload struct {
	ActiveStatus ads.ActiveStatus `json:"active_status" validate:"required" error:"status is required"`
}

// Validate custom validation for user scope
func (lp *adActiveStatusPayload) ValidateExtra(ctx echo.Context) error {
	if !lp.ActiveStatus.IsValid() {
		return middlewares.GroupError{
			"status": trans.E("status is invalid"),
		}
	}
	return nil
}

// Validate custom validation for user scope
func (lp *adArchiveStatusPayload) ValidateExtra(ctx echo.Context) error {
	if !lp.ActiveStatus.IsValid() {
		return middlewares.GroupError{
			"status": trans.E("status is invalid"),
		}
	}
	return nil
}

// Validate custom validation for user scope
func (lp *adAdminStatusPayload) ValidateExtra(ctx echo.Context) error {
	if !lp.AdminStatus.IsValid() {
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
		AdArchiveStatus: ads.ActiveStatusNo,
		AdPayStatus:     ads.ActiveStatusNo,
		AdAdminStatus:   ads.AdminStatusPending,
		AdActiveStatus:  ads.ActiveStatusNo,
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
	currentAd.AdAdminStatus = pl.AdminStatus
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
	_, b := currentUser.HasPermOn("change_archive_ad", owner.ID, owner.DBParentID.Int64)
	if !b {
		return ctx.JSON(http.StatusForbidden, trans.E("user can't access"))
	}
	currentAd.AdArchiveStatus = pl.ActiveStatus
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
	currentAd.Description = pl.Body
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
		currentAd.Extension = common.MakeNullString(filepath.Ext(currentAd.Src.String))
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
	//get message from redis by cli ID
	messageHistory, err := aredis.GetKey(pl.CliMessageID, false, 0)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	currentAd.CliMessageID = common.MakeNullString(pl.CliMessageID)
	currentAd.PromoteData = common.MakeNullString(messageHistory)
	currentAd.Src = common.MakeNullString("")
	currentAd.BotChatID = common.NullInt64{}
	currentAd.BotMessageID = common.NullInt64{}
	currentAd.Description = nil
	assert.Nil(m.UpdateAd(currentAd))
	defer func() {
		rabbit.MustPublish(&commands.IdentifyAD{
			AdID: currentAd.ID,
		})
	}()
	return u.OKResponse(ctx, currentAd)
}

//	changeActiveStatus change active status for ad
//	@Route	{
//		url = /list/active_status/:id
//		method = put
//		payload	= adActiveStatusPayload
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
	//currentUser := authz.MustGetUser(ctx)
	assert.Nil(err)
	owner, err := aaa.NewAaaManager().FindUserByID(currentAd.UserID)
	assert.Nil(err)
	//check everything is good
	if currentAd.AdPayStatus == ads.ActiveStatusYes && currentAd.AdAdminStatus == ads.AdminStatusPending {
		currentAd.AdActiveStatus = pl.ActiveStatus
		// send mail
		go func() {
			mail.SendByTemplateName(trans.T("AD activated").Translate("fa_IR"), "active-ad", struct {
				Ad   string
				Name string
			}{
				Ad:   currentAd.Name,
				Name: owner.Email,
			}, config.Config.Mail.From, owner.Email)
		}()
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
	if currentAd.AdActiveStatus == ads.ActiveStatusYes || currentAd.AdAdminStatus == ads.AdminStatusAccepted {
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
		currentAd.Extension = common.MakeNullString(filepath.Ext(currentAd.Src.String))
	}
	if currentAd.PromoSrc.Valid {
		g, err := url.Parse(fcfg.Fcfg.File.UploadPath)
		if err != nil {
			return u.NotFoundResponse(ctx, nil)
		}
		g.Path = path.Join(g.Path, currentAd.PromoSrc.String)
		s := g.String()
		currentAd.PromoSrc = common.MakeNullString(s)
		currentAd.Extension = common.MakeNullString(filepath.Ext(currentAd.PromoSrc.String))
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
	if currentAd.AdPayStatus != ads.ActiveStatusNo {
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
		Description: plan.Description,
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
	println("//////////")
	fmt.Println(ctx.Param("id"))
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 0)
	if err != nil {
		return ctx.Redirect(http.StatusFound, frontNOk)
	}
	currentAd, err := adManager.FindAdByID(id)
	if err != nil {
		return ctx.Redirect(http.StatusFound, frontNOk)
	}
	plan, err := adManager.FindPlanByID(currentAd.PlanID.Int64)
	if err != nil {
		return ctx.Redirect(http.StatusFound, frontNOk)
	}
	if strings.ToLower(ctx.FormValue("Status")) != "ok" {
		return ctx.Redirect(http.StatusFound, frontNOk)
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
		return ctx.Redirect(http.StatusFound, frontNOk)
	}
	if resp.Status == bcfg.Bcfg.Gate.MerchantOkStatus {
		currentUser, err := aaa.NewAaaManager().FindUserByID(currentAd.UserID)
		assert.Nil(err)
		billing, err := bil.NewBilManager().RegisterBilling(currentUser, ctx.FormValue("Authority"), resp.RefID, plan.Price, resp.Status, id)
		if err != nil {
			return ctx.Redirect(http.StatusFound, frontNOk)
		}
		//update ad pay status
		currentAd.AdPayStatus = ads.ActiveStatusYes
		currentAd.AdActiveStatus = ads.ActiveStatusYes
		assert.Nil(adManager.UpdateAd(currentAd))
		//call worker
		// send mail
		go func() {
			mail.SendByTemplateName(trans.T("new plan bought").Translate("fa_IR"), "charge", struct {
				Name     string
				Price    int64
				Campaign string
			}{
				Name:     currentUser.Email,
				Price:    plan.Price,
				Campaign: currentAd.Name,
			}, config.Config.Mail.From, currentUser.Email)
		}()
		return ctx.Redirect(http.StatusMovedPermanently, fmt.Sprintf("%s%d", frontOk, billing.PaymentID.Int64))
	}
	return ctx.Redirect(http.StatusMovedPermanently, frontNOk)
}

//	pieChartAdvertiser show per campaign view
//	@Route	{
//		url	=	/dashboard/pie-chart
//		method	= get
//		resource = pie_chart_advertiser:self
//		middleware = authz.Authenticate
//		200 = ads.PieChart
//		400 = base.ErrorResponseSimple
//	}
func (u *Controller) pieChartAdvertiser(ctx echo.Context) error {
	currentUser := authz.MustGetUser(ctx)
	m := ads.NewAdsManager()
	pieChart, err := m.PieChartAdvertiser(currentUser.ID)
	if err != nil {
		return u.BadResponse(ctx, nil)
	}
	return u.OKResponse(ctx, pieChart)

}

//	callIdentifyAd call identify
//	@Route	{
//		url	=	/identify/:id
//		method	= get
//		resource = get_ad:global
//		middleware = authz.Authenticate
//		200 = base.NormalResponse
//		400 = base.ErrorResponseSimple
//	}
func (u *Controller) callIdentifyAd(ctx echo.Context) error {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 0)
	assert.Nil(err)
	rabbit.MustPublish(&commands.IdentifyAD{
		AdID: id,
	})
	return u.OKResponse(ctx, nil)
}

type assignCat struct {
	Categories []int64 `json:"categories"`
}

//	assignCategory assignCategory for ad
//	@Route	{
//		url	=	/category/:id
//		method	= put
//		payload	= assignCat
//		resource = assign_category:self
//		middleware = authz.Authenticate
//		200 = ads.Ad
//		400 = base.ErrorResponseSimple
//	}
func (u *Controller) assignCategory(ctx echo.Context) error {
	currentUser := authz.MustGetUser(ctx)
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 0)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	pl := u.MustGetPayload(ctx).(*assignCat)
	adManager := ads.NewAdsManager()
	currentAd, err := adManager.FindAdByID(id)
	if err != nil {
		return u.NotFoundResponse(ctx, trans.E("ad not found"))
	}
	owner, err := aaa.NewAaaManager().FindUserByID(currentAd.UserID)
	assert.Nil(err)
	_, b := currentUser.HasPermOn("assign_category", currentAd.UserID, owner.DBParentID.Int64)
	if !b {
		return ctx.JSON(http.StatusForbidden, trans.E("user can't access"))
	}
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	catManager := cat.NewCatManager()
	err = catManager.AssignCat(currentAd.ID, pl.Categories)
	if err != nil {
		return u.BadResponse(ctx, trans.E("error while assigning role"))
	}
	return u.OKResponse(ctx, trans.T("category assigned successfully"))
}
