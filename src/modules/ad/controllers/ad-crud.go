package ad

import (
	"modules/ad/ads"

	"modules/user/middlewares"

	"common/assert"

	echo "gopkg.in/labstack/echo.v3"
)

// @Validate {
// }
type AdPayload struct {
	Name string `json:"name" validate:"required" error:"name is required"`
}

//	create create ad
//	@Route	{
//		url	=	/
//		method	= post
//		payload	= AdPayload
//		resource = create_ad:self
//		middleware = authz.Authenticate
//		200 = ads.Ad
//		400 = base.ErrorResponseSimple
//	}
func (u *Controller) create(ctx echo.Context) error {
	pl := u.MustGetPayload(ctx).(*AdPayload)
	m := ads.NewAdsManager()
	currentUser, ok := authz.GetUser(ctx)

	if !ok {
		return u.NotFoundResponse(ctx, nil)
	}

	newAd := &ads.Ad{
		Name:            pl.Name,
		AdArchiveStatus: ads.AdArchiveStatusNo,
		AdPayStatus:     ads.AdPayStatusNo,
		AdAdminStatus:   ads.AdAdminStatusPending,
		UserID:          currentUser.ID,
	}
	assert.Nil(m.CreateAd(newAd))
	return u.OKResponse(ctx, newAd)
}
