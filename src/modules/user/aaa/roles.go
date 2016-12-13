package aaa

import (
	"common/assert"
	"common/controllers/base"
	"fmt"
	"strings"
	"time"
)

// Role model
// @Model {
//		table = roles
//		primary = true, id
//		find_by = id,name
//		list = yes
// }
type Role struct {
	ID          int64     `db:"id" json:"id" sort:"true"`
	Name        string    `json:"name" db:"name" search:"true"`
	Description string    `db:"description" json:"description"`
	CreatedAt   time.Time `db:"created_at" json:"created_at" sort:"true"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at" sort:"true"`
}

//RoleDataTable is the role full data in data table, after join with other field
// @DataTable {
//		url = /roles
//		entity = role
//		view = role_list:global
//		controller = modules/user/controllers
//		fill = FillRoleDataTableArray
//		_edit = role_edit:global
// }
type RoleDataTable struct {
	Role
	ParentID int64 `db:"-" json:"parent_id" visible:"false"`
	OwnerID  int64 `db:"-" json:"owner_id" visible:"false"`
}

// RegisterRole is try to register role
func (m *Manager) RegisterRole(name string, description string, perm map[base.UserScope][]string) (role *Role, err error) {
	role = &Role{
		Name:        name,
		Description: description,
	}
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

		if err != nil {
			role = nil
		}
	}()
	err = m.CreateRole(role)
	err = m.RegisterRolePermission(role.ID, perm)
	if err != nil {
		role = nil
		return
	}

	return
}

// FillRoleDataTableArray is the function to handle
func (m *Manager) FillRoleDataTableArray(u base.PermInterfaceComplete, filters map[string]string, search map[string]string, sort, order string, p, c int) (RoleDataTableArray, int64) {
	var params []interface{}
	var res RoleDataTableArray
	var where []string

	countQuery := "SELECT COUNT(id) FROM roles"
	query := "SELECT roles.* FROM roles"
	for field, value := range filters {
		where = append(where, fmt.Sprintf("%s=%s", field, "?"))
		params = append(params, value)
	}

	for column, val := range search {
		where = append(where, fmt.Sprintf("%s=%s", column, "?"))
		params = append(params, val)
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
