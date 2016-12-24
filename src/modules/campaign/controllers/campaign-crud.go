package campaign

import (
	"common/assert"
	"common/models/common"
	"modules/campaign/cmp"
	"modules/misc/trans"
	"modules/user/aaa"
	"modules/user/middlewares"
	"net/http"
	"strconv"
	"time"

	"gopkg.in/labstack/echo.v3"
)

// @Validate {
// }
type campaignPayload struct {
	UserID int64     `json:"user_id"`
	Name   string    `json:"name" validate:"required"`
	Start  time.Time `json:"start"`
	End    time.Time `json:"stop" `
}

//	createCampaign
//	@Route	{
//	url	=	/create
//	method	= post
//	payload	= campaignPayload
//	resource = create_campaign:self
//	middleware = authz.Authenticate
//	200 = cmp.Campaign
//	400 = base.ErrorResponseSimple
//	}
func (u *Controller) createCampaign(ctx echo.Context) error {
	pl := u.MustGetPayload(ctx).(*campaignPayload)
	m := cmp.NewCmpManager()
	currentUser, ok := authz.GetUser(ctx)
	if !ok {
		return u.NotFoundResponse(ctx, nil)
	}
	if pl.UserID == 0 {
		pl.UserID = currentUser.ID
	}
	owner, err := aaa.NewAaaManager().FindUserByID(pl.UserID)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	_, b := currentUser.HasPermOn("create_campaign", owner.ID, owner.DBParentID.Int64)
	if !b {
		return ctx.JSON(http.StatusForbidden, trans.E("user can't access"))
	}
	c := m.Create(owner, pl.Name, pl.Start, pl.End)
	return u.OKResponse(ctx, c)
}

//	active
//	@Route	{
//	url	=	/active/:id
//	method	= put
//	resource = active_campaign:self
//	middleware = authz.Authenticate
//	200 = cmp.Campaign
//	400 = base.ErrorResponseSimple
//	}
func (u *Controller) active(ctx echo.Context) error {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 0)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	m := cmp.NewCmpManager()
	campaign, err := m.FindCampaignByID(id)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	currentUser, ok := authz.GetUser(ctx)
	if !ok {
		return u.NotFoundResponse(ctx, nil)
	}
	owner, err := aaa.NewAaaManager().FindUserByID(campaign.UserID)
	_, b := currentUser.HasPermOn("active_channel", owner.ID, owner.DBParentID.Int64)
	if !b {
		return ctx.JSON(http.StatusForbidden, trans.E("user can't access"))
	}
	ch := m.ChangeActive(owner, campaign)
	return u.OKResponse(ctx, ch)
}

//	editCampaign
//	@Route	{
//	url	=	/:id
//	method	= put
//	payload	= campaignPayload
//	resource = edit_campaign:self
//	middleware = authz.Authenticate
//	200 = cmp.Campaign
//	400 = base.ErrorResponseSimple
//	}
func (u *Controller) editCampaign(ctx echo.Context) error {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 0)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	pl := u.MustGetPayload(ctx).(*campaignPayload)
	cmpManager := cmp.NewCmpManager()
	usrManager := aaa.NewAaaManager()
	currentUser, ok := authz.GetUser(ctx)
	if !ok {
		return u.NotFoundResponse(ctx, nil)
	}
	campaign, err := cmpManager.FindCampaignByID(id)
	if err != nil {
		return u.ForbiddenResponse(ctx, nil)
	}
	owner, err := usrManager.FindUserByID(campaign.UserID)
	if err != nil {
		return u.ForbiddenResponse(ctx, nil)
	}
	_, b := currentUser.HasPermOn("edit_campaign", owner.ID, owner.DBParentID.Int64)
	if !b {
		return ctx.JSON(http.StatusForbidden, trans.E("user can't access"))
	}
	campaign.Name = pl.Name
	campaign.Start = common.NullTime{Valid: !pl.Start.IsZero(), Time: pl.Start}
	campaign.End = common.NullTime{Valid: !pl.End.IsZero(), Time: pl.End}
	assert.Nil(cmpManager.UpdateCampaign(campaign))
	return u.OKResponse(ctx, campaign)
}
