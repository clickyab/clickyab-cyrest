package aaa

import (
	"common/assert"
	"fmt"
	"modules/misc/base"
	"strings"
	"time"
)

// RolePermission model
// @Model {
//		table = role_permission
//		primary = true, id
//		find_by = id
//		belong_to = Role:role_id
//		list = true
// }
type RolePermission struct {
	ID         int64          `db:"id" json:"id"`
	RoleID     int64          `json:"role_id" db:"role_id"`
	Permission string         `json:"permission" db:"permission"`
	Scope      base.UserScope `json:"scope" db:"scope"`
	CreatedAt  time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at" db:"updated_at"`
}

// GetResourceMap return resource map for some roles
func (m *Manager) GetPermissionMap(r ...Role) map[base.UserScope]map[string]bool {
	res := make(map[base.UserScope]map[string]bool)
	res[base.ScopeGlobal] = make(map[string]bool)
	res[base.ScopeSelf] = make(map[string]bool)
	res[base.ScopeParent] = make(map[string]bool)
	if len(r) == 0 {
		return res
	}

	var roleIDs []string
	for i := range r {
		roleIDs = append(roleIDs, fmt.Sprintf("%d", r[i].ID))
	}

	var rr []RolePermission
	query := fmt.Sprintf("SELECT * FROM role_permission WHERE role_id IN (%s)", strings.Join(roleIDs, ","))
	_, err := m.GetDbMap().Select(
		&rr,
		query,
	)
	assert.Nil(err)
	for i := range rr {
		res[rr[i].Scope][rr[i].Permission] = true
	}
	return res
}

// RegisterRolePermission register role with permission assign
func (m *Manager) RegisterRolePermission(roleID int64, perm map[base.UserScope][]string) error {
	var rolePermission []interface{}
	for scope, val := range perm {
		for i := range val {
			role := &RolePermission{
				Permission: val[i],
				RoleID:     roleID,
				Scope:      scope,
			}
			rolePermission = append(rolePermission, role)
		}

	}
	return m.GetDbMap().Insert(rolePermission...)
}

func (m *Manager) DeleteRolePermissionByRoleID(roleID int64) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE role_id=?", RolePermissionTableFull)
	_, err := m.GetDbMap().Exec(
		query,
		roleID,
	)
	return err
}
