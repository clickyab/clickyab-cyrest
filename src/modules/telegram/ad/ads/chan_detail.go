package ads

import "time"

// ChanDetail is the list of all known channels for cyborg
// @Model {
//		table = channel_details
//		primary = true, id
//		find_by = id, name,channel_id
// }
type ChanDetail struct {
	ID         int64     `db:"id" json:"id"`
	Name       string    `db:"name" json:"name"`
	ChannelID  int64     `db:"channel_id" json:"channel_id"`
	Title      string    `db:"title" json:"title"`
	Info       string    `db:"info" json:"info"`
	TelegramID string    `db:"cli_telegram_id" json:"cli_telegram_id"`
	UserCount  int64     `db:"user_count" json:"user_count"`
	AdminCount int       `db:"admin_count" json:"admin_count"`
	PostCount  int       `db:"post_count" json:"post_count"`
	TotalView  int       `db:"total_view" json:"total_view"`
	CreatedAt  time.Time `db:"created_at" json:"created_at" sort:"true"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at" sort:"true"`
}
