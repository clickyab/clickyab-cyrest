package bil

import (
	"common/assert"
	"common/models/common"
	"fmt"
	"modules/misc/base"
	"modules/user/aaa"
	"strings"
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
	ID        int64             `db:"id" json:"id" sort:"true" title:"ID" title:"ID"`
	UserID    int64             `json:"user_id" db:"user_id" title:"UserID"  visible:"false"`
	Amount    int64             `json:"amount" db:"amount" title:"Amount" title:"Amount" sort:"true"`
	Reason    common.NullString `json:"reason" db:"reason" title:"Reason" title:"Reason"`
	CreatedAt time.Time         `db:"created_at" json:"created_at" sort:"true" title:"Created at" sort:"true"`
	UpdatedAt time.Time         `db:"updated_at" json:"updated_at" sort:"true" title:"Updated at" sort:"true"`
}

//BillingDataTable is the ad full data in data table, after join with other field
// @DataTable {
//		url = /
//		entity = billingList
//		view = billing_list:self
//		controller = modules/billing/controllers
//		fill = FillBillingDataTableArray
//		_change = billing_manage:global
// }
type BillingDataTable struct {
	Billing
	Email    string `db:"email" json:"email" search:"true" title:"Email"`
	ParentID int64  `db:"-" json:"parent_id" visible:"false"`
	OwnerID  int64  `db:"-" json:"owner_id" visible:"false"`
	Actions  string `db:"-" json:"_actions" visible:"false"`
}

// FillBillingDataTableArray is the function to handle
func (m *Manager) FillBillingDataTableArray(
	u base.PermInterfaceComplete,
	filters map[string]string,
	search map[string]string,
	sort, order string, p, c int) (BillingDataTableArray, int64) {
	var params []interface{}
	var res BillingDataTableArray
	var where []string

	countQuery := fmt.Sprintf("SELECT COUNT(%[1]s.id) FROM %[1]s "+
		"LEFT JOIN %[2]s ON %[2]s.id=%[1]s.user_id ",
		BillingTableFull,
		aaa.UserTableFull,
	)
	query := fmt.Sprintf("SELECT %[1]s.*,%[2]s.email FROM %[1]s "+
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
