package ads

import (
	"common/models/common"
	"time"
)

const (
	AdStatusPending  AdStatus = "pending"
	AdStatusRejected AdStatus = "rejected"
	AdStatusAccepted AdStatus = "accepted"
	AdStatusArchive  AdStatus = "archive"

	AdTypeImg      AdType = "img"
	AdTypeDocument AdType = "document"
	AdTypeVideo    AdType = "video"
)

type (
	// AdStatus is the ad status
	// @Enum{
	// }
	AdStatus string

	// AdType is the ad status
	// @Enum{
	// }
	AdType string
)

// Ad model
// @Model {
//		table = ads
//		primary = true, id
//		find_by = id,user_id
//		list = yes
// }
type Ad struct {
	ID          int64             `db:"id" json:"id" sort:"true" title:"ID"`
	UserID      int64             `json:"user_id" db:"user_id" title:"UserID"`
	Name        string            `json:"name" db:"name" search:"true" title:"Name"`
	Media       common.NullString `json:"media" db:"media" title:"Media"`
	Description string            `json:"description" db:"description" title:"Description"`
	Link        common.NullString `json:"link" db:"link" search:"true" title:"Link"`
	Type        AdType            `json:"type" db:"type" filter:"true" title:"Type"`
	Status      AdStatus          `json:"status" db:"status" filter:"true" title:"Status"`
	CreatedAt   time.Time         `db:"created_at" json:"created_at" sort:"true" title:"Created at"`
	UpdatedAt   time.Time         `db:"updated_at" json:"updated_at" sort:"true" title:"Updated at"`
}
