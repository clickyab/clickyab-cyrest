package ads

import (
	"common/models/common"
	"fmt"
	"time"
)

// ChannelAdDetail is the list of ad in channel for cyborg
// @Model {
//		table = channel_ad_detail
//		primary = true, id
//		find_by = channel_id, ad_id
// }
type ChannelAdDetail struct {
	ID        int64            `db:"id" json:"id"`
	ChannelID int64            `db:"channel_id" json:"channel_id"`
	AdID      int64            `db:"ad_id" json:"ad_id"`
	View      int64            `db:"view" json:"view"`
	Position  common.NullInt64 `db:"position" json:"position"`
	Warning   int64            `db:"warning" json:"warning"`
	CreatedAt time.Time        `db:"created_at" json:"created_at" sort:"true"`
}

// FindChannelIDAdIDDetail return the ChannelAdDetail base on its ad_id , channel_id
func (m *Manager) FindChannelIDAdIDDetail(c int64, a int64) (*ChannelAdDetail, error) {
	var res ChannelAdDetail
	err := m.GetDbMap().SelectOne(
		&res,
		fmt.Sprintf("SELECT * FROM %s WHERE ad_id=? AND channel_id=?", ChannelAdDetailTableFull),
		a,
		c,
	)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// CreateChannelAdDetails create ad channel details
func (m *Manager) CreateChannelAdDetails(cad []*ChannelAdDetail) error {
	h := make([]interface{}, len(cad))
	for i := range cad {
		h[i] = cad[i]
	}
	return m.GetDbMap().Insert(h...)

}
