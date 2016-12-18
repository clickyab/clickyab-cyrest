// Package cat is the models for category module
package chn

import (
	"common/assert"
	"common/models/common"
	"time"
)

//'pending', 'rejected','accepted','archive'
const (
	ChannelStatusPending  ChannelStatus = "pending"
	ChannelStatusRejected ChannelStatus = "rejected"
	ChannelStatusAccepted ChannelStatus = "accepted"
	ChannelStatusArchive  ChannelStatus = "archive"
)

type (
	// ChannelStatus is the channel status
	// @Enum{
	// }
	ChannelStatus string
)

// Category model
// @Model {
//		table = channels
//		primary = true, id
//		find_by = id,user_id
//		list = yes
// }
type Channel struct {
	ID        int64             `db:"id" json:"id"`
	UserID    int64             `json:"user_id" db:"user_id"`
	Name      string            `json:"name" db:"name"`
	Link      common.NullString `json:"link" db:"link"`
	Admin     common.NullString `json:"admin" db:"admin"`
	Status    ChannelStatus     `json:"status" db:"status"`
	CreatedAt time.Time         `db:"created_at" json:"created_at"`
	UpdatedAt time.Time         `db:"updated_at" json:"updated_at"`
}

func (c *Channel) Initialize() {

}

func (m *Manager) Create(admin, link, name string, status ChannelStatus, userID int64) *Channel {

	ch := &Channel{
		Admin:  common.NullString{Valid: admin != "", String: admin},
		Link:   common.NullString{Valid: link != "", String: link},
		Name:   name,
		Status: status,
		UserID: userID,
	}
	assert.Nil(m.CreateChannel(ch))
	return ch
}
