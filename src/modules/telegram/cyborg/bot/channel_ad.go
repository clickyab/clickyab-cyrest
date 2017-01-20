package bot

import (
	"common/models/common"
	"fmt"
	"time"
)

// ActiveStatus is the channel_ad active
// @Enum {
// }
type ActiveStatus string

const (
	//ActiveStatusYes is the yes status
	ActiveStatusYes ActiveStatus = "yes"
	// ActiveStatusNo is the no status
	ActiveStatusNo ActiveStatus = "no"
)

// ChannelAd is the list of ad in channel for cyborg
// @Model {
//		table = channel_ad
//		primary = false, channel_id,ad_id
//		find_by = channel_id, ad_id
// }
type ChannelAd struct {
	ChannelID    int64           `db:"channel_id" json:"channel_id"`
	AdID         int64           `db:"ad_id" json:"ad_id"`
	View         int64           `db:"view" json:"view"`
	CliMessageID string          `db:"cli_message_id" json:"cli_message_id"`
	Active       ActiveStatus    `db:"active" json:"active"`
	Start        common.NullTime `db:"start" json:"start"`
	End          common.NullTime `db:"end" json:"end"`
	Warning      int64           `db:"warning" json:"warning"`
	PossibleView int64           `db:"possible_view" json:"possible_view"`
	CreatedAt    time.Time       `db:"created_at" json:"created_at" sort:"true"`
	UpdatedAt    time.Time       `db:"updated_at" json:"updated_at" sort:"true"`
}

// FindChannelIDAdByAdID return the ChannelAd base on its ad_id
func (m *Manager) FindChannelIDAdByAdID(c int64, a int64) (*ChannelAd, error) {
	var res ChannelAd
	err := m.GetDbMap().SelectOne(
		&res,
		fmt.Sprintf("SELECT * FROM %s WHERE ad_id=? AND channel_id=?", ChannelAdTableFull),
		a,
		c,
	)

	if err != nil {
		return nil, err
	}

	return &res, nil
}
