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
func (m *Manager) FindBundleChannelAd(channelID, bundleID, adID int64) *BundleChannelAd {
	var res BundleChannelAd
	q := fmt.Sprintf(
		"SELECT * FROM %s WHERE bundle_id=? AND channel_id=? AND ad_id=?", BundleChannelAdTableFull)
	err := m.GetDbMap().SelectOne(
		&res,
		q,
		bundleID,
		channelID,
		adID,
	)
	assert.Nil(err)
	return &res
}

//FindBundleChannelAdActiveType type bundle channel active ad
type FindBundleChannelAdActiveType struct {
	BundleChannelAd
	TargetView   int64             `db:"target_view"`
	Position     int64             `db:"position"`
	CliMessageID common.NullString `db:"cli_message_id"`
	PromoteData  common.NullString `db:"promote_data"`
	Src          common.NullString `db:"src"`
	Code         string            `db:"code"`
}

// FindBundleChannelAdActive return the adID base on its ad_id,ActiveStatus
func (m *Manager) FindBundleChannelAdActive() ([]FindBundleChannelAdActiveType, error) {
	res := []FindBundleChannelAdActiveType{}
	_, err := m.GetDbMap().Select(
		&res,
		fmt.Sprintf("SELECT %[1]s.*,%[2]s.cli_message_id,%[2]s.promote_data,%[2]s.src, %[3]s.code"+
			"(%[3]s.view -((%[3]s.view * %[3]s.percent_finish)/100)) AS target_view,%[3]s.position"+
			" FROM %[1]s "+
			" INNER JOIN %[2]s ON %[2]s.id=%[1]s.ad_id AND %[1]s.ad_id = %[3]s.target_ad "+
			" INNER JOIN %[3]s ON %[3]s.id=%[1]s.bundle_id "+
			" AND %[1]s.active=? ",
			BundleChannelAdTableFull,
			AdTableFull,
			BundlesTableFull,
		),
		ActiveStatusYes,
	)

	if err != nil {
		return nil, err
	}

	return res, nil
}

// UpdateBundleChannelAds update bundle channel ads
func (m *Manager) UpdateBundleChannelAds(ca []BundleChannelAd) error {
	var q = fmt.Sprintf("UPDATE %s SET  updated_at=?,warning=? , view=? WHERE channel_id=? AND ad_id=? AND bundle_id = ?", BundleChannelAdTableFull)
	for i := range ca {
		_, err := m.GetDbMap().Exec(
			q,
			time.Now(),
			ca[i].Warning,
			ca[i].View,
			ca[i].ChannelID,
			ca[i].AdID,
			ca[i].BundleID,
		)
		if err != nil {
			return err
		}
	}

	return nil
}
