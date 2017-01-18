package bot

import "time"

// ChanDetails is the list of all known channels for cyborg
// @Model {
//		table = channel_details
//		primary = true, id
//		find_by = id, name
// }
type ChanDetail struct {
	ID         int64     `db:"id" json:"id"`
	Name       string    `db:"name" json:"name"`
	ChannelID  int64     `db:"channel_id" json:"channel_id"`
	Title      string    `db:"title" json:"title"`
	Info       string    `db:"info" json:"info"`
	TelegramID string    `db:"telegram_id" json:"telegram_id"`
	UserCount  int64     `db:"user_count" json:"user_count"`
	AdminCount int       `db:"admin_count" json:"admin_count"`
	Num        int       `db:"num" json:"num"`
	TotalView  int       `db:"total_view" json:"total_view"`
	CreatedAt  time.Time `db:"created_at" json:"created_at" sort:"true"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at" sort:"true"`
}
