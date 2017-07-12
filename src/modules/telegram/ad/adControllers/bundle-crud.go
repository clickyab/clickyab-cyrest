package ad

import (
	"common/models/common"
	"modules/telegram/ad/ads"
	"modules/user/middlewares"

	"modules/misc/middlewares"
	"modules/misc/trans"

	"common/assert"
	"modules/user/aaa"
	"net/http"
	"strconv"

	"fmt"
	"strings"

	"github.com/Sirupsen/logrus"
	"gopkg.in/labstack/echo.v3"
)

// @Validate {
// }
type bundleCreatePayload struct {
	Position      int64            `json:"position" validate:"required" error:"position is required"`
	Price         common.NullInt64 `json:"price"`
	PercentFinish int64            `json:"percent_finish" validate:"required" error:"percent_finish is required"`
	BundleType    ads.BType        `json:"bundle_type" validate:"required" error:"bundle_type is required"`
	Rules         string           `json:"rules"`
	TargetAd      int64            `json:"target_ad"`
}

// Validate custom validation for user scope
func (lp *bundleCreatePayload) ValidateExtra(ctx echo.Context) error {
	if !lp.BundleType.IsValid() {
		return middlewares.GroupError{
			"status": trans.E("bundle type is invalid"),
		}
	}
	return nil
}

//	createBundle create bundle for admin
//	@Route	{
//		url	=	/bundle
//		method	= post
//		payload	= bundleCreatePayload
//		resource = create_bundle:global
//		middleware = authz.Authenticate
//		200 = ads.Bundles
//		400 = base.ErrorResponseSimple
//	}
func (u *Controller) createBundle(ctx echo.Context) error {
	pl := u.MustGetPayload(ctx).(*bundleCreatePayload)
	currentUser := authz.MustGetUser(ctx)
	bundle := ads.Bundles{
		UserID:        currentUser.ID,
		Price:         pl.Price,
		AdminStatus:   ads.ActiveStatusNo,
		ActiveStatus:  ads.ActiveStatusNo,
		BundleType:    pl.BundleType,
		Place:         pl.Position,
		Rules:         common.MakeNullString(pl.Rules),
		TargetAd:      common.MakeNullInt64(pl.TargetAd),
		PercentFinish: pl.PercentFinish,
	}
	m := ads.NewAdsManager()
	err := m.CreateBundles(&bundle)
	logrus.Warn(err)
	if err != nil {
		return u.BadResponse(ctx, trans.E("error while creating bundle"))
	}
	return u.OKResponse(ctx, bundle)
}

// @Validate {
// }
type bundleAssignPayload struct {
	Ads    []int64 `json:"ads" validate:"required" error:"ads is required"`
	Target int64   `json:"target" validate:"required" error:"target is required"`
}

//	assignBundle assign bundle for admin
//	@Route	{
//		url	=	/bundle/assign/:id
//		method	= put
//		payload	= bundleAssignPayload
//		resource = assign_bundle:global
// 		middleware = authz.Authenticate
//		200 = ads.Bundles
//		400 = base.ErrorResponseSimple
//	}
func (u *Controller) assignBundle(ctx echo.Context) error {
	currentUser := authz.MustGetUser(ctx)
	pl := u.MustGetPayload(ctx).(*bundleAssignPayload)
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 0)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	m := ads.NewAdsManager()
	currentBundle, err := m.FindBundlesByID(id)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	owner, err := aaa.NewAaaManager().FindUserByID(currentBundle.UserID)
	assert.Nil(err)

	_, b := currentUser.HasPermOn("assign_bundle", owner.ID, owner.DBParentID.Int64)
	if !b {
		return ctx.JSON(http.StatusForbidden, trans.E("user can't access"))
	}
	// TODO  check all ads and target exists and are active one
	var adIDs []string
	for i := range pl.Ads {
		adIDs = append(adIDs, fmt.Sprintf("%d", pl.Ads[i]))
	}
	currentBundle.Ads = common.CommaArray(strings.Join(adIDs, ","))
	currentBundle.TargetAd = common.MakeNullInt64(pl.Target)
	err = m.UpdateBundles(currentBundle)
	if err != nil {
		return u.BadResponse(ctx, trans.E("error while updating bundle"))
	}
	return u.OKResponse(ctx, currentBundle)
}
