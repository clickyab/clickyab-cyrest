package bil

import (
	"common/models/common"
	"time"
)

// Billing model
// @Model {
//		table = billings
//		primary = true, id
//		find_by = id
//		list = yes
// }
type Billing struct {
	ID        int64             `db:"id" json:"id" sort:"true" title:"ID"`
	UserID    int64             `json:"user_id" db:"user_id" title:"UserID"`
	Amount    int64             `json:"amount" db:"amount" title:"Amount"`
	Reason    common.NullString `json:"reason" db:"reason" title:"Reason"`
	CreatedAt time.Time         `db:"created_at" json:"created_at" sort:"true" title:"Created at"`
	UpdatedAt time.Time         `db:"updated_at" json:"updated_at" sort:"true" title:"Updated at"`
}


