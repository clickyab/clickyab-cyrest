package bil

import (
	"common/assert"
	"common/models/common"
	"encoding/json"
	"fmt"
	"modules/misc/base"
	"modules/telegram/ad/ads"
	"modules/user/aaa"
	"strings"
	"time"
)

const (
	//BilTypeWithdrawal const withdrawal type withdrawal
	BilTypeWithdrawal BillingType = "withdrawal"
	// BilTypeBilling const billing type billing
	BilTypeBilling BillingType = "billing"
	// BilTypeIncome const income type billing
	BilTypeIncome BillingType = "income"
	// BilTypeCampaign const campaign type billing
	BilTypeCampaign BillingType = "campaign"

	//BilStatusAccepted const status accepted
	BilStatusAccepted BillingStatus = "accepted"
	//BilStatusRejected const status rejected
	BilStatusRejected BillingStatus = "rejected"
	//BilStatusPending const billing status pending
	BilStatusPending BillingStatus = "pending"

	//BilDepositYes deposit yes
	BilDepositYes BillingDeposit = "yes"
	//BilDepositNo deposit No
	BilDepositNo BillingDeposit = "no"
)

type (
	//BillingType type billing
	//@Enum{
	//}
	BillingType string

	//BillingStatus status billing
	//@Enum{
	//}
	BillingStatus string

	//BillingDeposit deposit billing
	//@Enum{
	//}
	BillingDeposit string
)

// Billing model
// @Model {
//		table = billings
//		primary = true, id
//		find_by = id,payment_id,user_id
//		list = yes
// }
type Billing struct {
	ID        int64             `db:"id" json:"id" sort:"true" title:"ID" visible:"false"`
	UserID    int64             `json:"user_id" db:"user_id" title:"UserID" visible:"false"`
	ChannelID common.NullInt64  `db:"channel_id" json:"channel_id"  title:"ChannelID"`
	AdID      common.NullInt64  `db:"ad_id" json:"ad_id"   title:"AdID"`
	PaymentID common.NullInt64  `json:"payment_id" db:"payment_id" title:"PaymentID"`
	Amount    int64             `json:"amount" db:"amount" sort:"true" title:"Amount"`
	Reason    common.NullString `json:"reason" db:"reason" title:"Reason"`
	Type      BillingType       `json:"type" db:"type" title:"Type" filter:"true"`
	Status    BillingStatus     `json:"status" db:"status" title:"Status" filter:"true"`
	Deposit   BillingDeposit    `json:"deposit" db:"deposit" title:"Deposit" filter:"true"`
	CreatedAt *time.Time        `json:"created_at" db:"created_at" sort:"true" title:"Created at"`
	UpdatedAt *time.Time        `json:"updated_at" db:"updated_at" sort:"true" title:"Updated at" visible:"false"`
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
func (m *Manager) RegisterBilling(currentUser *aaa.User, authority string, refID int64, price int64, statusCode int64, adID int64) (*Billing, error) {
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
		AdID:      common.MakeNullInt64(adID),
		Type:      BilTypeBilling,
		Deposit:   BilDepositNo,
		Status:    BilStatusPending,
	}
	err = m.CreateBilling(billing)
	if err != nil {
		return nil, err
	}
	b, err := json.Marshal(billing)
	if err != nil {
		return nil, err
	}
	bilDetail := &BillingDetail{
		BillingID: billing.ID,
		UserID:    currentUser.ID,
		Reason:    common.MakeNullString(string(b)),
	}
	err = m.CreateBillingDetail(bilDetail)
	if err != nil {
		return nil, err
	}
	return billing, err
}

// ChannelAdBilling insert income & campaign billing
func (m *Manager) ChannelAdBilling(channelAds []ads.FinishedActiveChannels, ad ads.FinishedActiveAds) error {
	err := m.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			assert.Nil(m.Rollback())
		} else {
			err = m.Commit()
		}

	}()
	//insert campaign into billing
	billing := Billing{
		Status:  BilStatusPending,
		Deposit: BilDepositNo,
		Type:    BilTypeCampaign,
		AdID:    common.MakeNullInt64(ad.ID),
		Amount:  ad.Price * -1,
		UserID:  ad.UserID,
	}
	err = m.CreateBilling(&billing)
	if err != nil {
		return err
	}

	jsonBilling, err := json.Marshal(billing)
	if err != nil {
		return err
	}
	bilDetail := BillingDetail{
		UserID:    ad.UserID,
		BillingID: billing.ID,
		Reason:    common.MakeNullString(string(jsonBilling)),
	}
	err = m.CreateBillingDetail(&bilDetail)

	//insert income billing
	for k := range channelAds {
		//calculate amount
		amount := channelAds[k].View * ad.Share
		fAmount := float64(amount) * 0.1
		amount = int64(fAmount)

		billing.UserID = channelAds[k].UserID
		billing.ChannelID = common.MakeNullInt64(channelAds[k].ChannelID)
		billing.AdID = common.MakeNullInt64(ad.ID)
		billing.Amount = amount
		billing.Type = BilTypeIncome
		billing.Status = BilStatusPending
		billing.Deposit = BilDepositNo

		err = m.CreateBilling(&billing)
		if err != nil {
			return err
		}

		jsonBilling, err = json.Marshal(billing)
		if err != nil {
			return err
		}
		bilDetail = BillingDetail{
			UserID:    channelAds[k].UserID,
			BillingID: billing.ID,
			Reason:    common.MakeNullString(string(jsonBilling)),
		}
		err = m.CreateBillingDetail(&bilDetail)
		if err != nil {
			return err
		}

	}

	return err
}

// ChannelBilling insert income & channel billing
func (m *Manager) ChannelBilling(adActive *ads.ChannelAdActive) error {
	err := m.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			assert.Nil(m.Rollback())
		} else {
			err = m.Commit()
		}

	}()
	//insert campaign into billing
	billing := Billing{}
	bilDetail := BillingDetail{}
	//calculate amount
	amount := adActive.View * adActive.Share
	fAmount := float64(amount) * 0.1
	amount = int64(fAmount)

	billing.UserID = adActive.UserID
	billing.ChannelID = common.MakeNullInt64(adActive.ChannelID)
	billing.AdID = common.MakeNullInt64(adActive.AdID)
	billing.Amount = amount
	billing.Type = BilTypeIncome
	billing.Status = BilStatusPending
	billing.Deposit = BilDepositNo

	err = m.CreateBilling(&billing)
	if err != nil {
		return err
	}

	jsonBilling, err := json.Marshal(billing)
	if err != nil {
		return err
	}
	bilDetail = BillingDetail{
		UserID:    adActive.UserID,
		BillingID: billing.ID,
		Reason:    common.MakeNullString(string(jsonBilling)),
	}
	err = m.CreateBillingDetail(&bilDetail)
	return err
}

//BillingDataTable is the ad full data in data table, after join with other field
// @DataTable {
//		url = /list
//		entity = billingList
//		view = billing_list:self
//		controller = modules/billing/controllers
//		fill = FillBillingDataTableArray
//		_status = change_status_billing:global
//		_deposit = change_deposit_billing:global
// }
type BillingDataTable struct {
	Billing
	Email    string           `db:"email" json:"email" search:"true" title:"Email"`
	ParentID common.NullInt64 `db:"parent_id" json:"parent_id" visible:"false"`
	OwnerID  int64            `db:"owner_id" json:"owner_id" visible:"false"`
	Actions  string           `db:"-" json:"_actions" visible:"false"`
}

// FillBillingDataTableArray is the function to handle
func (m *Manager) FillBillingDataTableArray(
	u base.PermInterfaceComplete,
	filters map[string]string,
	search map[string]string,
	contextparams map[string]string,
	sort, order string, p, c int) (BillingDataTableArray, int64) {
	var params []interface{}
	var res BillingDataTableArray
	var where []string

	countQuery := fmt.Sprintf("SELECT COUNT(%[1]s.id) FROM %[1]s "+
		"LEFT JOIN %[2]s ON %[2]s.id=%[1]s.user_id ",
		BillingTableFull,
		aaa.UserTableFull,
	)
	query := fmt.Sprintf("SELECT %[1]s.*,%[2]s.email,%[2]s.id AS owner_id, %[2]s.parent_id as parent_id FROM %[1]s "+
		"LEFT JOIN %[2]s ON %[2]s.id=%[1]s.user_id ",
		BillingTableFull,
		aaa.UserTableFull,
	)
	for field, value := range filters {
		where = append(where, fmt.Sprintf("%s.%s=?", BillingTableFull, field))
		params = append(params, value)
	}

	for column, val := range search {
		where = append(where, fmt.Sprintf("%s LIKE ?", column))
		params = append(params, "%"+val+"%")
	}

	currentUserID := u.GetID()
	highestScope := u.GetCurrentScope()

	if highestScope == base.ScopeSelf {
		where = append(where, fmt.Sprintf("%s.user_id=?", BillingTableFull))
		params = append(params, currentUserID)
	} else if highestScope == base.ScopeParent {
		where = append(where, fmt.Sprintf("%s.parent_id=?", aaa.UserTableFull))
		params = append(params, currentUserID)
	}

	//check for perm
	if len(where) > 0 {
		query += " WHERE "
		countQuery += " WHERE "
	}
	query += strings.Join(where, " AND ")
	countQuery += strings.Join(where, " AND ")
	limit := c
	offset := (p - 1) * c
	if sort != "" {
		query += fmt.Sprintf(" ORDER BY %s %s ", sort, order)
	}
	query += fmt.Sprintf(" LIMIT %d OFFSET %d ", limit, offset)
	count, err := m.GetDbMap().SelectInt(countQuery, params...)
	assert.Nil(err)

	_, err = m.GetDbMap().Select(&res, query, params...)
	assert.Nil(err)
	return res, count
}

// SumBil SumBil
type SumBil struct {
	Balance common.ZeroNullInt64 `json:"balance" db:"balance"`
}

// SumBilling sum billing
func (m *Manager) SumBilling(userID int64) SumBil {
	q := fmt.Sprintf("SELECT SUM(amount) AS balance FROM %s WHERE user_id = ?", BillingTableFull)
	var res SumBil
	_, err := m.GetDbMap().Select(&res, q, userID)
	assert.Nil(err)
	return res
}
