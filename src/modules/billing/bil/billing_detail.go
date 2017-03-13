package bil

import (
	"common/models/common"
	"time"
)

// BillingDetail model
// @Model {
//		table = billing_detail
//		primary = true, id
//		find_by = id,user_id,billing_id
//		list = yes
// }
type BillingDetail struct {
	ID        int64             `db:"id" json:"id" sort:"true" title:"ID" visible:"false"`
	BillingID int64             `json:"billing_id" db:"billing_id" title:"BillingID"`
	UserID    int64             `json:"user_id" db:"user_id" title:"UserID" `
	Reason    common.NullString `json:"reason" db:"reason" title:"Reason"`
	CreatedAt *time.Time        `json:"created_at" db:"created_at" sort:"true" title:"Created at"`
}
