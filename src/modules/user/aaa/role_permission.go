package aaa

import (
	"common/assert"
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
	ID         int64     `db:"id" json:"id"`
	RoleID     int64     `json:"role_id" db:"role_id"`
	Permission string    `json:"permission" db:"permission"`
	Scope      ScopePerm `json:"scope" db:"scope"`
}

// scopePerm is the perm for the role
type (
	// scopePerm is the scope perm status for a single permission
	// @Enum{
	// }
	ScopePerm string
)

const (
	// ScopePermGlobal is the global permission
	ScopePermGlobal ScopePerm = "global"
	// ScopePermParent is the parent permission
	ScopePermParent ScopePerm = "parent"
	// ScopePermOwn is the own permission
	ScopePermOwn ScopePerm = "own"
)

// GetResourceMap return resource map for some roles
func (m *Manager) GetPermissionMap(r ...Role) map[ScopePerm]map[string]bool {
	res := make(map[ScopePerm]map[string]bool)
	res[ScopePermGlobal]=make(map[string]bool)
	res[ScopePermOwn]=make(map[string]bool)
	res[ScopePermParent]=make(map[string]bool)
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
