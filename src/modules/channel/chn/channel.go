// Package cat is the models for category module
package chn

import "time"

// Category model
// @Model {
//		table = channels
//		primary = true, id
//		find_by = id
//		list = yes
// }
type Channel struct {
	ID        int64     `db:"id" json:"id"`
	UserID    int64     `json:"user_id" db:"user_id"`
	Name      string    `json:"name" db:"name"`
	Link      string    `json:"link" db:"link"`
	AdminID   int64     `json:"admin_id" db:"admin_id"`
	Status    int       `json:"status" db:"status"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

func (c *Channel) Initialize() {

}
