package aaa

import (
	"common/models/common"
	"time"
)

// UserFinancial model
// @Model {
//		table = user_financial
//		primary = true, id
//		find_by = id, user_id
//		list = yes
// }
type UserFinancial struct {
	ID            int64             `db:"id" json:"id"`
	UserID        int64             `db:"user_id" json:"user_id"`
	BankName      common.NullString `db:"bank_name" json:"bank_name"`
	AccountHolder common.NullString `db:"account_holder" json:"account_holder"`
	CardNumber    common.NullString `db:"card_number" json:"card_number"`
	AccountNumber common.NullString `db:"account_number" json:"account_number"`
	ShebaNumber   common.NullString `db:"sheba_number" json:"sheba_number"`
	CreatedAt     time.Time         `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time         `db:"updated_at" json:"updated_at"`
}
