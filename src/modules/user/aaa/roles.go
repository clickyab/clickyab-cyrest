package aaa

import (
	"time"
	"common/assert"
	"common/controllers/base"
)
// Role model
// @Model {
//		table = roles
//		primary = true, id
//		find_by = id,name
//		list = yes
// }
type Role struct {
	ID          int64     `db:"id" json:"id"`
	Name       string    `json:"name" db:"name"`
	Description       string    `json:"description" db:"description"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

// RegisterRole is try to register role
func (m *Manager) RegisterRole(name string,description string,perm map[base.UserScope][]string) (role *Role, err error) {
	role = &Role{
		Name: name,
		Description:description,
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
	err =m.RegisterRolePermission(role.ID,perm)
	if err != nil {
		role = nil
		return
	}

	return
}