package bil

import (
	"common/models/common"
	"time"
)

const (
	// PayStatusPending is the payment pending status
	PayStatusPending PayStatus = "pending"
	// PayStatusRejected is the payment rejected status
	PayStatusRejected PayStatus = "rejected"
	// PayStatusPaid is the payment paid status
	PayStatusPaid PayStatus = "paid"
)

type (
	// PayStatus is the payment status
	// @Enum{
	// }
	PayStatus string
)

// Payment model
// @Model {
//		table = payments
//		primary = true, id
//		find_by = id
//		list = yes
// }
type Payment struct {
	ID         int64             `db:"id" json:"id" sort:"true" title:"ID"`
	UserID     int64             `json:"user_id" db:"user_id" title:"UserID"`
	Amount     int64             `json:"amount" db:"amount" title:"Amount"`
	Status     PayStatus         `json:"status" db:"status" title:"Status"`
	Authority  common.NullString `json:"authority" db:"authority" title:"Authority"`
	RefID      common.NullInt64  `json:"ref_id" db:"ref_id" title:"RefID"`
	StatusCode common.NullInt64  `json:"status_code" db:"status_code" title:"StatusCode"`
	CreatedAt  time.Time         `db:"created_at" json:"created_at" sort:"true" title:"Created at"`
	UpdatedAt  time.Time         `db:"updated_at" json:"updated_at" sort:"true" title:"Updated at"`
}
