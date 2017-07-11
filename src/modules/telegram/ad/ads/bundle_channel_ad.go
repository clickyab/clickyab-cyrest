package ads

import (
	"common/assert"
	"common/models/common"
	"fmt"
	"modules/user/aaa"
	"time"
)

// BundleChannelAd list bundle in ad channel
// @Model {
//		table = bundle_channel_ad
//		primary = false, channel_id,ad_id,bundle_id
// }
type BundleChannelAd struct {
	ChannelID    int64             `db:"channel_id" json:"channel_id"`
	AdID         int64             `db:"ad_id" json:"ad_id"`
	BundleID     int64             `db:"bundle_id" json:"bundle_id"`
	View         int64             `db:"view" json:"view"`
	Shot         common.NullString `db:"shot" json:"shot"`
	Warning      int64             `db:"warning" json:"warning"`
	BotChatID    int64             `db:"bot_chat_id" json:"bot_chat_id"`
	BotMessageID common.NullInt64  `db:"bot_message_id" json:"bot_message_id"`
	Active       ActiveStatus      `db:"active" json:"active"`
	Start        common.NullTime   `db:"start" json:"start"`
	End          common.NullTime   `db:"end" json:"end"`
	CreatedAt    *time.Time        `db:"created_at" json:"created_at" sort:"true"`
	UpdatedAt    *time.Time        `db:"updated_at" json:"updated_at" sort:"true"`
}

// FindActiveChannelBundleByUserID try to find by user id
func (m *Manager) FindActiveChannelBundleByUserID(userID int64, channelID int64, bundleID int64) []ActiveAdUser {
	var res []ActiveAdUser
	q := fmt.Sprintf(
		"SELECT a.*,c.id AS channel_id "+
			"FROM %s AS ba "+
			"INNER JOIN %s AS c "+
			"ON ba.channel_id = c.id "+
			"INNER JOIN %s AS b "+
			"ON ba.bundle_id = b.id "+
			"INNER JOIN %s AS u "+
			"ON u.id = c.user_id"+
			"INNER JOIN %s as a "+
			"ON ba.ad_id = a.id "+
			"WHERE c.user_id = u.id "+
			"AND ba.active = ? "+
			"AND u.id = ?"+
			"AND b.id = ?"+
			"AND c.id = ?", BundleChannelAdTableFull, ChannelTableFull, BundlesTableFull, aaa.UserTableFull, AdTableFull)
	_, err := m.GetDbMap().Select(
		&res,
		q,
		ActiveStatusYes,
		userID,
		bundleID,
		channelID,
	)
	assert.Nil(err)
	return res
}

// FindBundleChannelAd try to find by user id
func (m *Manager) FindBundleChannelAd(channelID, bundleID int64) *BundleChannelAd {
	var res BundleChannelAd
	q := fmt.Sprintf(
		"SELECT * FROM %s WHERE bundle_id=? AND channel_id=?", BundleChannelAdTableFull)
	err := m.GetDbMap().SelectOne(
		&res,
		q,
		bundleID,
		channelID,
	)
	assert.Nil(err)
	return &res
}
