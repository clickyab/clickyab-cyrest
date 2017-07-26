package ads

import (
	"common/models/common"
	"time"
)

const (
	// BTypeBanner type banner
	BTypeBanner BType = "banner"
	// BTypeBannerRep type banner rep
	BTypeBannerRep BType = "banner+rep"
	// BTypeRepBanner type rep banner
	BTypeRepBanner BType = "rep+banner"
	// BTypeRepBannerRep type rep banner rep
	BTypeRepBannerRep BType = "rep+banner+rep"
)

// BType bundle type
//	@Enum{
//	}
type BType string

// Bundles model
// @Model {
//		table = bundles
//		primary = true, id
//		find_by = id,user_id,code
//		list = yes
//	}
type Bundles struct {
	ID            int64             `db:"id" json:"id"`
	UserID        int64             `db:"user_id" json:"user_id"`
	Place         int64             `db:"place" json:"place"`
	View          int64             `db:"view" json:"view"`
	Price         int64             `db:"price" json:"price"`
	PercentFinish int64             `db:"percent_finish" json:"percent_finish"`
	BundleType    BType             `db:"bundle_type" json:"bundle_type"`
	Rules         common.NullString `json:"rules" db:"rules"`
	Code          string            `json:"code" db:"code"`
	AdminStatus   ActiveStatus      `json:"admin_status" db:"admin_status"`
	ActiveStatus  ActiveStatus      `json:"active_status" db:"active_status"`
	Ads           common.CommaArray `json:"ads" db:"ads"`
	TargetAd      int64             `json:"target_ad" db:"target_ad"`
	Start         common.NullTime   `json:"start" db:"start"`
	End           common.NullTime   `json:"end" db:"end"`
	CreatedAt     *time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt     *time.Time        `json:"updated_at" db:"updated_at"`
}
