package cmp


// CampaignBlack model
// @Model {
//		table = campaign_black
//		primary = false, channel_id, campaign_id
// }
type CampaignBlack struct {
	CampaignID int64     `db:"campaign_id" json:"campaign_id"`
	ChannelID  int64     `db:"channel_id" json:"channel_id"`
}
