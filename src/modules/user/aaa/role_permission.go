package aaa

import (
	"common/assert"
	"errors"
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
	ID         int64           `db:"id" json:"id"`
	RoleID     int64           `json:"role_id" db:"role_id"`
	Permission base.Permission `json:"permission" db:"permission"`
	Scope      base.UserScope  `json:"scope" db:"scope"`
	CreatedAt  *time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt  *time.Time      `json:"updated_at" db:"updated_at"`
}

// GetPermissionMap return resource map for some roles
func (m *Manager) GetPermissionMap(r ...Role) map[base.UserScope]map[base.Permission]bool {
	res := make(map[base.UserScope]map[base.Permission]bool)
	res[base.ScopeGlobal] = make(map[base.Permission]bool)
	res[base.ScopeSelf] = make(map[base.Permission]bool)
	res[base.ScopeParent] = make(map[base.Permission]bool)
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
		res[rr[i].Scope][base.Permission(rr[i].Permission)] = true
	}
	return res
}

// RegisterRolePermission register role with permission assign
func (m *Manager) RegisterRolePermission(roleID int64, perm map[base.UserScope][]base.Permission) error {
	now := time.Now()
	var rolePermission []interface{}
	for scope, val := range perm {
		for i := range val {
			ok := base.PermissionCheckRegistered(val[i])
			if !ok {
				return errors.New("perm not exists")
			}
			role := &RolePermission{
				Permission: val[i],
				RoleID:     roleID,
				Scope:      scope,
				CreatedAt:  &now,
				UpdatedAt:  &now,
			}
			rolePermission = append(rolePermission, role)
		}

	}
	return m.GetDbMap().Insert(rolePermission...)
}

// DeleteRolePermissionByRoleID delete role permission by role id
func (m *Manager) DeleteRolePermissionByRoleID(roleID int64) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE role_id=?", RolePermissionTableFull)
	_, err := m.GetDbMap().Exec(
		query,
		roleID,
	)
	return err
}
