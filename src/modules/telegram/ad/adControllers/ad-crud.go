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
//		url	=	/change-admin/:id
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
//		url = /change-archive/:id
//		method = put
//		resource = change_archive_ad:self
//		middleware = authz.Authenticate
//		200 = ads.Ad
//		400 = base.ErrorResponseSimple
//	}
func (u *Controller) changeArchiveStatus(ctx echo.Context) error {
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
	if currentAd.AdArchiveStatus == ads.AdArchiveStatusYes {
		currentAd.AdArchiveStatus = ads.AdArchiveStatusNo
	} else {
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
//		url = /change-active/:id
//		method = put
//		resource = change_active_ad:self
//		middleware = authz.Authenticate
//		200 = ads.Ad
//		400 = base.ErrorResponseSimple
//	}
func (u *Controller) changeActiveStatus(ctx echo.Context) error {
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
	if currentAd.AdAdminStatus == "accepted" && currentAd.AdArchiveStatus == "no" && currentAd.Name != "" && currentAd.AdActiveStatus == ads.AdActiveStatusNo {
		currentAd.AdActiveStatus = ads.AdActiveStatusYes
		//check to add job
		if currentAd.CliMessageID.Valid {
			rabbit.MustPublish(commands.IdentifyAD{AdID: currentAd.ID})
		}
	} else {
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
//		url	=	/:id
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
	return u.OKResponse(ctx, currentAd)
}
