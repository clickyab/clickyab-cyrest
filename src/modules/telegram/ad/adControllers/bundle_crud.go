package ad

import (
	"common/models/common"
	"modules/telegram/ad/ads"

	"gopkg.in/labstack/echo.v3"
)

// @Validate {
// }
type checkAdInBundlePayload struct {
	Data common.CommaArray `json:"data" validate:"required" error:"status is required"`
}

/*
// Validate custom validation for user scope
func (lp *checkAdInBundlePayload) ValidateExtra(ctx echo.Context) error {
	if !lp.Data.IsValid() {
		return middlewares.GroupError{
			"status": trans.E("status is invalid"),
		}
	}
	return nil
}*/

//SliceInt64 structure export slice int64
type SliceInt64 []int64

//	checkAdInBundle check exist ad in bundle
//	@Route	{
//		url	=	/check_ad/
//		method	= post
//		payload	= checkAdInBundlePayload
//		resource = bundle:global
//		middleware = authz.Authenticate
//		200 = SliceInt64
//		400 = base.ErrorResponseSimple
//	}
func (u *Controller) checkAdInBundle(ctx echo.Context) error {
	pl := u.MustGetPayload(ctx).(*checkAdInBundlePayload)
	m := ads.NewAdsManager()
	bundleChannelAd, err := m.CheckAdActiveAdBundle(pl.Data)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	var IDs []int64
	for i := range bundleChannelAd {
		IDs = append(IDs, bundleChannelAd[i].AdID)
	}

	return u.OKResponse(ctx, IDs)
}
