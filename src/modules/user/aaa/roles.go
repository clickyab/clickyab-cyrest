package aaa

import (
	"common/assert"
	"common/config"
	"common/models/common"
	"common/utils"
	"fmt"
	"modules/misc/base"
	"modules/misc/trans"
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
	ID          int64             `db:"id" json:"id" sort:"true"`
	Name        string            `json:"name" db:"name" search:"true"`
	Description common.NullString `db:"description" json:"description"`
	CreatedAt   *time.Time        `db:"created_at" json:"created_at" sort:"true"`
	UpdatedAt   *time.Time        `db:"updated_at" json:"updated_at" sort:"true"`
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
	ParentID common.NullInt64 `db:"" json:"parent_id" visible:"false"`
	OwnerID  int64            `db:"-" json:"owner_id" visible:"false"`
	Actions  string           `db:"-" json:"_actions" visible:"false"`
}

// RegisterRole is try to register role
func (m *Manager) RegisterRole(name string, description string, perm map[base.UserScope][]base.Permission) (role *Role, err error) {
	role = &Role{
		Name:        name,
		Description: common.NullString{String: description, Valid: true},
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
	if err != nil {
		role = nil
		return
	}
	err = m.RegisterRolePermission(role.ID, perm)
	if err != nil {
		role = nil
		return
	}

	return
}

// FillRoleDataTableArray is the function to handle
func (m *Manager) FillRoleDataTableArray(
	u base.PermInterfaceComplete,
	filters map[string]string,
	search map[string]string,
	contextparams map[string]string,
	sort, order string, p, c int) (RoleDataTableArray, int64) {
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
		where = append(where, fmt.Sprintf("%s LIKE ?", column))
		params = append(params, fmt.Sprintf("%s"+val+"%s", "%", "%"))
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

// CountRoleUserByID count the role user by id
func (m *Manager) CountRoleUserByID(roleID int64) (int64, error) {
	query := fmt.Sprintf("SELECT COUNT(role_id) FROM %s WHERE role_id=?", UserRoleTableFull)
	return m.GetDbMap().SelectInt(
		query,
		roleID,
	)
}

// DeleteRoleByID delete a role by id
func (m *Manager) DeleteRoleByID(roleID int64) (*Role, error) {
	role, err := m.FindRoleByID(roleID)
	if err != nil {
		return nil, trans.E("no role found")
	}
	_, err = m.GetDbMap().Delete(role)
	return role, err
}

// DeleteRole in transaction try delete role
func (m *Manager) DeleteRole(ID int64) (r *Role, err error) {
	if ID == config.Config.Role.Default {
		return r, trans.E("you cannot remove Default role")
	}
	r, err = m.FindRoleByID(ID)
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
			r = nil
		}

	}()

	if err != nil {
		r = nil
		return
	}

	//delete role_permission
	err = m.DeleteRolePermissionByRoleID(ID)

	//delete role
	_, err = m.DeleteRoleByID(ID)
	return
}

// UpdateRoleWithPerm try to save a new Role in database
func (m *Manager) UpdateRoleWithPerm(ID int64, name, description string, perm map[base.UserScope][]base.Permission) (r *Role, err error) {
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
			r = nil
		}
	}()

	if err != nil {
		r = nil
		return
	}

	r, err = m.FindRoleByID(ID)
	if err != nil {
		return
	}
	r.Name = name
	r.Description = common.MakeNullString(description)

	err = m.UpdateRole(r)
	if err != nil {
		err = trans.E("cant update role")
	}

	//delete role_permission
	err = m.DeleteRolePermissionByRoleID(ID)
	if err != nil {
		err = trans.E("cant delete role permission")
	}

	err = m.RegisterRolePermission(ID, perm)
	if err != nil {
		err = trans.E("cant register role permission")
	}
	return

}

// GetAllRole return all registered roles in system
func (m *Manager) GetAllRole() ([]Role, error) {
	var res []Role
	query := fmt.Sprintf("SELECT * FROM %s", RoleTableFull)
	_, err := m.GetDbMap().Select(
		&res,
		query,
	)

	return res, err

}

// FindRoleByIDs return the Role base on its id
func (m *Manager) FindRoleByIDs(IDs []int64) ([]Role, error) {
	var res []Role
	var t []interface{}
	for i := range IDs {
		t = append(t, IDs[i])
	}
	_, err := m.GetDbMap().Select(
		&res,
		fmt.Sprintf("SELECT * FROM %s WHERE id IN(%s)", RoleTableFull, utils.DBImplode(IDs)),
		t...,
	)

	if err != nil {
		return nil, err
	}

	return res, nil
}
