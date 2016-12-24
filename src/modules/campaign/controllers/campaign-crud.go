package campaign

import (
	"modules/campaign/cmp"
	"modules/misc/trans"
	"modules/user/aaa"
	"net/http"
	"strconv"
	"time"

	"modules/user/middlewares"

	echo "gopkg.in/labstack/echo.v3"
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
	ch := m.ChangeActive(campaign.ID, campaign.UserID, campaign.Active, campaign.CreatedAt)
	return u.OKResponse(ctx, ch)
}

