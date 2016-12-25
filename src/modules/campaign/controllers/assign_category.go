package campaign

import (
	"common/assert"
	"fmt"
	"modules/campaign/cmp"
	"modules/category/cat"
	"modules/user/aaa"

	"modules/user/middlewares"

	echo "gopkg.in/labstack/echo.v3"
)

// @Validate {
// }
type assignCategoryPayload struct {
	CampaignID  int64   `json:"campaign_id" validate:"required"`
	CategoryIDs []int64 `json:"category_id" validate:"required"`
}

// assignCategory
// @Route {
//		url	=	/assign/category
//		method	=	post
//		payload	=	assignCategoryPayload
//		resource=	assign_category:self
//		middleware = authz.Authenticate
//		200	=	cmp.CampaignCategory
//		400	=	base.ErrorResponseSimple
// }
func (u *Controller) assignCategory(ctx echo.Context) error {
	pl := u.MustGetPayload(ctx).(*assignCategoryPayload)
	usrManager := aaa.NewAaaManager()
	cmpManager := cmp.NewCmpManager()
	catManager := cat.NewCatManager()
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
	err = cmpManager.DeleteCat(campaign.ID)
	fmt.Println(err)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}

	for _, v := range pl.CategoryIDs {
		category, err := catManager.FindCategoryByID(v)
		if err != nil {
			return u.NotFoundResponse(ctx, nil)
		}
		black := &cmp.CampaignCategory{CampaignID: campaign.ID, CategoryID: category.ID}
		assert.Nil(cmpManager.CreateCampaignCategory(black))
	}

	return u.OKResponse(ctx, nil)

}
