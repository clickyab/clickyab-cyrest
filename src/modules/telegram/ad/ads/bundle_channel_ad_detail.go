package ads

import (
	"common/models/common"
	"time"
)

// BundleChannelAdDetail is the list of ad in channel for cyborg
// @Model {
//		table = bundle_channel_ad_detail
//		primary = true, id
//		find_by = channel_id, ad_id,bundle_id
// }
type BundleChannelAdDetail struct {
	ID        int64            `db:"id" json:"id"`
	ChannelID int64            `db:"channel_id" json:"channel_id"`
	BundleID  int64            `db:"bundle_id" json:"bundle_id"`
	AdID      int64            `db:"ad_id" json:"ad_id"`
	View      int64            `db:"view" json:"view"`
	Position  common.NullInt64 `db:"position" json:"position"`
	Warning   int64            `db:"warning" json:"warning"`
	CreatedAt *time.Time       `db:"created_at" json:"created_at" sort:"true"`
}

// CreateBundleChannelAdDetails create bundle channel ad details
func (m *Manager) CreateBundleChannelAdDetails(cad []*BundleChannelAdDetail) error {
	h := make([]interface{}, len(cad))
	for i := range cad {
		h[i] = cad[i]
	}
	return m.GetDbMap().Insert(h...)

}
