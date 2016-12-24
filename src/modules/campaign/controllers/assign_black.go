package campaign

import (
	"common/assert"
	"modules/user/aaa"

	"modules/campaign/cmp"
	"modules/channel/chn"
	"modules/user/middlewares"

	"gopkg.in/labstack/echo.v3"
)

// @Validate {
// }
type assignBlackPayload struct {
	CampaignID int64   `json:"campaign_id" validate:"required"`
	ChannelIDs []int64 `json:"channel_id" validate:"required"`
}

// assignBlack
// @Route {
//		url	=	/assign/black
//		method	=	post
//		payload	=	assignBlackPayload
//		resource=	assign_black:self
//		middleware = authz.Authenticate
//		200	=	cmp.CampaignBlack
//		400	=	base.ErrorResponseSimple
// }
func (u *Controller) assignBlack(ctx echo.Context) error {
	pl := u.MustGetPayload(ctx).(*assignBlackPayload)
	usrManager := aaa.NewAaaManager()
	cmpManager := cmp.NewCmpManager()
	chnManager := chn.NewChnManager()
	currentUser, ok := authz.GetUser(ctx)
	if !ok {
		return u.NotFoundResponse(ctx, nil)
	}
	// find campaign
	campaign, err := cmpManager.FindCampaignByID(pl.CampaignID)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	owner, err := usrManager.FindUserByID(campaign.UserID)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}

	//check permission
	_, b := currentUser.HasPermOn("assign_black", owner.ID, owner.DBParentID.Int64)
	if !b {
		return u.NotFoundResponse(ctx, nil)
	}

	//delete previous black-list
	err = cmpManager.DeleteBlack(campaign.ID)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}

	for _, v := range pl.ChannelIDs {
		channel, err := chnManager.FindChannelByID(v)
		if err != nil {
			return u.NotFoundResponse(ctx, nil)
		}
		black := &cmp.CampaignBlack{CampaignID: campaign.ID, ChannelID: channel.ID}
		assert.Nil(cmpManager.CreateCampaignBlack(black))
	}

	return u.OKResponse(ctx, nil)

}
