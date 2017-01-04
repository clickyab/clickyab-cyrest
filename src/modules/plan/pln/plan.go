// Package plan is the models for plan module
package plan

import (
	"common/assert"
	"fmt"
	"modules/misc/base"
	"strings"
	"time"
)

const (
	ActiveStatusYes ActiveStatus = "yes"
	ActiveStatusNo  ActiveStatus = "no"
)

type (
	// ActiveStatus is the plan active
	// @Enum{
	// }
	ActiveStatus string
)

// plan model
// @Model {
//		table = plans
//		primary = true, id
//		find_by = id
//		list = yes
// }
type Plan struct {
	ID          int64        `db:"id" json:"id" sort:"true" title:"ID"`
	Name        string       `json:"name" db:"name" search:"true" title:"Name"`
	Description string       `json:"description" db:"description" search:"true" title:"description"`
	Active      ActiveStatus `json:"active" db:"active" filter:"true" title:"Active"`
	CreatedAt   time.Time    `db:"created_at" json:"created_at" sort:"true" title:"Created at"`
	UpdatedAt   time.Time    `db:"updated_at" json:"updated_at" sort:"true" title:"Updated at"`
}

func (c *Plan) Initialize() {

}

//
//// Create
//func (m *Manager) Create(name, description string, active ActiveStatus) *pLAN {
//
//	pln := &pLAN{
//		Name:        name,
//		Description: description,
//		Active:      active,
//	}
//	assert.Nil(m.CreatePlan(pln))
//	return pln
//}

//planDataTable is the role full data in data table, after join with other field
// @DataTable {
//		url = /list
//		entity = plan
//		view = plan_list:global
//		controller = modules/plan/controllers
//		fill = FillPlanDataTableArray
//		_edit = plan_edit:self
//		_change = plan_manage:global
// }
type PlanDataTable struct {
	Plan
	Email    string `db:"email" json:"email" search:"true" title:"Email"`
	ParentID int64  `db:"-" json:"parent_id" visible:"false"`
	OwnerID  int64  `db:"-" json:"owner_id" visible:"false"`
	Actions  string `db:"-" json:"_actions" visible:"false"`
}

// FillplanDataTableArray is the function to handle
func (m *Manager) FillPlanDataTableArray(u base.PermInterfaceComplete, filters map[string]string, search map[string]string, sort, order string, p, c int) (PlanDataTableArray, int64) {
	var params []interface{}
	var res PlanDataTableArray
	var where []string

	countQuery := fmt.Sprintf("SELECT COUNT(plan.id) FROM %s ", PlanTableFull)
	query := fmt.Sprintf("SELECT plan.*,users.email FROM %s LEFT JOIN %s ON %s.id=%s.user_id", PlanTableFull)
	for field, value := range filters {
		where = append(where, fmt.Sprintf(PlanTableFull+".%s=%s", field, "?"))
		params = append(params, value)
	}

	for column, val := range search {
		where = append(where, fmt.Sprintf("%s LIKE ?", column))
		params = append(params, fmt.Sprintf("%s"+val+"%s", "%", "%"))
	}

	currentUserID := u.GetID()
	highestScope := u.GetCurrentScope()

	if highestScope == base.ScopeSelf {
		where = append(where, fmt.Sprintf("%s.user_id=?", PlanTableFull))
		params = append(params, currentUserID)
	} else if highestScope == base.ScopeParent {
		where = append(where, "users.parent_id=?")
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
