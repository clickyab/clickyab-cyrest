package aaa

import (
	"common/assert"
	"common/controllers/base"
	"fmt"
	"strings"
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
