package cmp

import "fmt"

// CampaignCategory model
// @Model {
//		table = campaign_category
//		primary = false, category_id, campaign_id
// }
type CampaignCategory struct {
	CampaignID int64 `db:"campaign_id" json:"campaign_id"`
	CategoryID int64 `db:"category_id" json:"category_id"`
}

// DeleteBlack is try to delete black list
func (m *Manager) DeleteCat(ID int64) (err error) {
	query := fmt.Sprintf("DELETE FROM %s WHERE campaign_id=?", CampaignCategoryTableFull)
	_, err = m.Manager.GetDbMap().Exec(
		query,
		ID,
	)
	return err
}
