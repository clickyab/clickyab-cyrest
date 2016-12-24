package cmp

import "fmt"

// CampaignBlack model
// @Model {
//		table = campaign_black
//		primary = false, channel_id, campaign_id
// }
type CampaignBlack struct {
	CampaignID int64 `db:"campaign_id" json:"campaign_id"`
	ChannelID  int64 `db:"channel_id" json:"channel_id"`
}

// DeleteBlack is try to delete black list
func (m *Manager) DeleteBlack(ID int64) (err error) {
	query := fmt.Sprintf("DELETE FROM %s WHERE campaign_id=?", CampaignBlackTableFull)
	_, err = m.Manager.GetDbMap().Exec(
		query,
		ID,
	)
	return err
}
