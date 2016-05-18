package aaa

import "time"

// MessageLog model
// @Model {
//		table = message_logs
//		schema = aaa
//		primary = true, id
//		find_by = id
//		belong_to = User:user_id
//		list = yes
// }
type MessageLog struct {
	ID        int64     `db:"id" json:"id"`
	UserID    int64     `db:"user_id" json:"user_id"`
	Body      string    `db:"body" json:"body"`
	Contact   string    `db:"contact" json:"contact"`
	MediaType string    `db:"media_type" json:"media_type"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
