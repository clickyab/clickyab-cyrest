// Package cmp is the models for campaign module
package cmp

import "time"

// Campaign model
// @Model {
//		table = campaigns
//		primary = true, id
//		find_by = id,user_id
//		list = yes
// }
type Campaign struct {
	ID        int64          `db:"id" json:"id"`
	UserID    int64          `json:"user_id" db:"user_id"`
	Name      string         `json:"name" db:"name"`
	Status    CampaignStatus `json:"status" db:"status"`
	Active    CampaignActive `json:"active" db:"active"`
	Start     time.Time      `json:"start" db:"start"`
	End       time.Time      `json:"end" db:"end"`
	CreatedAt time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt time.Time      `db:"updated_at" json:"updated_at"`
}

const (
	CampaignStatusPending  CampaignStatus = "pending"
	CampaignStatusRejected CampaignStatus = "rejected"
	CampaignStatusAccepted CampaignStatus = "accepted"
	CampaignStatusArchive  CampaignStatus = "archive"

	CampaignActiveStart CampaignActive = "start"
	CampaignActiveStop  CampaignActive = "stop"
)

type (
	// CampaignStatus is the campaign status
	// @Enum{
	// }
	CampaignStatus string

	// CampaignActive is the campaign active
	// @Enum{
	// }
	CampaignActive string
)

func (c *Campaign) Initialize() {

}
