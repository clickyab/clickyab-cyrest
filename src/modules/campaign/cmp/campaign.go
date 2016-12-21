// Package cmp is the models for campaign module
package cmp

import (
	"common/assert"
	"common/models/common"
	"modules/user/aaa"
	"time"
)

// Campaign model
// @Model {
//		table = campaigns
//		primary = true, id
//		find_by = id,user_id
//		list = yes
// }
type Campaign struct {
	ID     int64  `db:"id" json:"id"`
	UserID int64  `json:"user_id" db:"user_id"`
	Name   string `json:"name" db:"name"`
	//Status    CampaignStatus `json:"status" db:"status"`
	Active    CampaignActive  `json:"active" db:"active"`
	Start     common.NullTime `json:"start" db:"start"`
	End       common.NullTime `json:"end" db:"end"`
	CreatedAt time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt time.Time       `db:"updated_at" json:"updated_at"`
}

const (
	/*	CampaignStatusPending  CampaignStatus = "pending"
		CampaignStatusRejected CampaignStatus = "rejected"
		CampaignStatusAccepted CampaignStatus = "accepted"
		CampaignStatusArchive  CampaignStatus = "archive"*/

	CampaignActiveStart CampaignActive = "yes"
	CampaignActiveStop  CampaignActive = "no"
)

type (

	// CampaignActive is the campaign active
	// @Enum{
	// }
	CampaignActive string
)

func (c *Campaign) Initialize() {

}

// Create campaign
func (m *Manager) Create(user *aaa.User, name string, start, end time.Time) *Campaign {

	c := &Campaign{
		UserID: user.ID,
		Name:   name,
		Active: CampaignActiveStop,
		Start:  common.NullTime{Valid: !start.IsZero(), Time: start},
		End:    common.NullTime{Valid: !end.IsZero(), Time: end},
	}
	assert.Nil(m.CreateCampaign(c))
	return c
}
