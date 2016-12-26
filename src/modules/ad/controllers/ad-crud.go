package ad

import (
	"modules/ad/ads"
	"modules/misc/trans"
	"modules/user/aaa"
	"net/http"
	"strconv"

	"modules/user/middlewares"

	"common/models/common"

	echo "gopkg.in/labstack/echo.v3"
)

// @Validate {
// }
type AdPayload struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Link        string `json:"link"`
}

//	editAd
//	@Route	{
//	url	=	/:id
//	method	= put
//	payload	= AdPayload
//	resource = edit_ad:self
//	middleware = authz.Authenticate
//	200 = ads.Ad
//	400 = base.ErrorResponseSimple
//	}
func (u *Controller) editAd(ctx echo.Context) error {
	pl := u.MustGetPayload(ctx).(*AdPayload)
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 0)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	m := ads.NewAdsManager()
	ad, err := m.FindAdByID(id)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	owner, err := aaa.NewAaaManager().FindUserByID(ad.UserID)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	currentUser, ok := authz.GetUser(ctx)
	if !ok {
		return u.NotFoundResponse(ctx, nil)
	}
	_, b := currentUser.HasPermOn("edit_ad", owner.ID, owner.DBParentID.Int64)
	if !b {
		return ctx.JSON(http.StatusForbidden, trans.E("user can't access"))
	}
	ad.ID = id
	ad.Link = common.NullString{Valid: pl.Link != "", String: pl.Link}
	ad.Description = common.NullString{Valid: pl.Description != "", String: pl.Description}
	ad.Name = pl.Name
	ch := m.EditAd(owner, ad)
	return u.OKResponse(ctx, ch)
}
