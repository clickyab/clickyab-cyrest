package ads

import (
	"common/assert"
	"common/models/common"
	"fmt"
	"modules/misc/base"
	"strings"
	"time"
)

// Plan model
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
	Price       int64        `db:"price" json:"price" title:"Price"`
	View        int64        `db:"view" json:"view" title:"View"`
	Position    int64        `db:"position" json:"position" title:"Position"`
	Share       int64        `db:"share" json:"share" title:"Share" perm:"plan_list:global"`
	Type        PlanType     `json:"type" db:"type" filter:"true" title:"Type"`
	Active      ActiveStatus `json:"active" db:"active" filter:"true" title:"Active"`
	CreatedAt   *time.Time   `db:"created_at" json:"created_at" sort:"true" title:"Created at"`
	UpdatedAt   *time.Time   `db:"updated_at" json:"updated_at" sort:"true" title:"Updated at"`
}

// GetAllActivePlans return all the active plans
func (m *Manager) GetAllActivePlans() ([]Plan, error) {
	var res []Plan
	query := fmt.Sprintf("SELECT * FROM %s WHERE active=?", PlanTableFull)
	_, err := m.GetDbMap().Select(
		&res,
		query,
		ActiveStatusYes,
	)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetAllIndividualActivePlans return all the individual active plans
func (m *Manager) GetAllIndividualActivePlans() ([]Plan, error) {
	var res []Plan
	query := fmt.Sprintf("SELECT * FROM %s WHERE active=? AND type=?", PlanTableFull)
	_, err := m.GetDbMap().Select(
		&res,
		query,
		ActiveStatusYes,
		PlanTypeIndividual,
	)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetAllPromotionActivePlans return all the promotion active plans
func (m *Manager) GetAllPromotionActivePlans() ([]Plan, error) {
	var res []Plan
	query := fmt.Sprintf("SELECT * FROM %s WHERE active=? AND type=?", PlanTableFull)
	_, err := m.GetDbMap().Select(
		&res,
		query,
		ActiveStatusYes,
		PlanTypePromotion,
	)
	if err != nil {
		return nil, err
	}

	return res, nil
}

//PlanDataTable is the role full data in data table, after join with other field
// @DataTable {
//		url = /list
//		entity = plan
//		view = plan_list:global
//		controller = modules/telegram/ad/planControllers
//		fill = FillPlanDataTableArray
//		_edit = plan_edit:self
//		_change = plan_manage:global
// }
type PlanDataTable struct {
	Plan
	ParentID common.NullInt64 `db:"-" json:"parent_id" visible:"false"`
	OwnerID  int64            `db:"-" json:"owner_id" visible:"false"`
	Actions  string           `db:"-" json:"_actions" visible:"false"`
}

// FillPlanDataTableArray is the function to handle
func (m *Manager) FillPlanDataTableArray(
	u base.PermInterfaceComplete,
	filters map[string]string,
	search map[string]string,
	contextparams map[string]string,
	sort, order string, p, c int) (PlanDataTableArray, int64) {
	var params []interface{}
	var res PlanDataTableArray
	var where []string

	countQuery := fmt.Sprintf("SELECT COUNT(plans.id) FROM %s ", PlanTableFull)
	query := fmt.Sprintf("SELECT plans.* FROM %s", PlanTableFull)
	for field, value := range filters {
		where = append(where, fmt.Sprintf(PlanTableFull+".%s=?", field))
		params = append(params, value)
	}

	for column, val := range search {
		where = append(where, fmt.Sprintf("%s LIKE ?", column))
		params = append(params, "%"+val+"%")
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
