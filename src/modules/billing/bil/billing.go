package bil

import (
	"common/assert"
	"common/models/common"
	"fmt"
	"time"
)

// Billing model
// @Model {
//		table = billings
//		primary = true, id
//		find_by = id,payment_id
//		list = yes
// }
type Billing struct {
	ID        int64             `db:"id" json:"id" sort:"true" title:"ID"`
	UserID    int64             `json:"user_id" db:"user_id" title:"UserID"`
	PaymentID common.NullInt64  `json:"payment_id" db:"payment_id" title:"PaymentID"`
	Amount    int64             `json:"amount" db:"amount" title:"Amount"`
	Reason    common.NullString `json:"reason" db:"reason" title:"Reason"`
	CreatedAt time.Time         `db:"created_at" json:"created_at" sort:"true" title:"Created at"`
	UpdatedAt time.Time         `db:"updated_at" json:"updated_at" sort:"true" title:"Updated at"`
}

// FindPaymentByAuthority return the Payment base on its authority
func (m *Manager) FindPaymentByAuthority(a common.NullString) (*Payment, error) {
	var res Payment
	err := m.GetDbMap().SelectOne(
		&res,
		fmt.Sprintf("SELECT * FROM %s WHERE authority=?", PaymentTableFull),
		a.String,
	)

	if err != nil {
		return nil, err
	}

	return &res, nil
}

// RegisterBilling is try to register billing
func (m *Manager) RegisterBilling(authority string, refID int64, price int64, statusCode int64) (*Billing, error) {
	payment, err := m.FindPaymentByAuthority(common.MakeNullString(authority))
	if err != nil {
		return nil, err
	}
	payment.Status = PayStatusPaid
	payment.RefID = common.NullInt64{Valid: true, Int64: refID}
	payment.StatusCode = common.NullInt64{Valid: true, Int64: statusCode}
	err = m.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			assert.Nil(m.Rollback())
		} else {
			err = m.Commit()
		}

	}()
	err = m.UpdatePayment(payment)
	if err != nil {
		return nil, err
	}
	billing := &Billing{
		PaymentID: common.NullInt64{Valid: true, Int64: payment.ID},
		Amount:    price,
		Reason:    common.MakeNullString("for buying our plan"),
		UserID:    payment.UserID,
	}
	err = m.CreateBilling(billing)
	if err != nil {
		return nil, err
	}
	return billing, err
}
