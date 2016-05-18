package aaa

import (
	"common/assert"
	"common/models/common"
	"fmt"
)

// AUTO GENERATED CODE. DO NOT EDIT!

// CreateUserRole try to save a new UserRole in database
func (m *Manager) CreateUserRole(ur *UserRole) error {

	func(in interface{}) {
		if ii, ok := in.(common.Initializer); ok {
			ii.Initialize()
		}
	}(ur)

	return m.GetDbMap().Insert(ur)
}

// UpdateUserRole try to update UserRole in database
func (m *Manager) UpdateUserRole(ur *UserRole) error {

	func(in interface{}) {
		if ii, ok := in.(common.Initializer); ok {
			ii.Initialize()
		}
	}(ur)

	_, err := m.GetDbMap().Update(ur)
	return err
}

// GetUserRoles return all Roles belong to User (many to many with UserRole)
func (m *Manager) GetUserRoles(u *User) []Role {
	var res []Role
	_, err := m.GetDbMap().Select(
		&res,
		fmt.Sprintf(
			"SELECT r.* FROM %s ur JOIN %s r ON ur.role_id = r.id WHERE ur.user_id=$1",
			UserRoleTableFull,
			RoleTableFull,
		),
		u.ID,
	)

	assert.Nil(err)
	return res
}

// GetRoleUsers return all Users belong to Role (many to many with UserRole)
func (m *Manager) GetRoleUsers(r *Role) []User {
	var res []User
	_, err := m.GetDbMap().Select(
		&res,
		fmt.Sprintf(
			"SELECT u.* FROM %s ur JOIN %s u ON ur.user_id = u.id WHERE ur.role_id=$1",
			UserRoleTableFull,
			UserTableFull,
		),
		r.ID,
	)

	assert.Nil(err)
	return res
}

// CountUserRoles return count Roles belong to User (many to many with UserRole)
func (m *Manager) CountUserRoles(u *User) int64 {
	res, err := m.GetDbMap().SelectInt(
		fmt.Sprintf(
			"SELECT COUNT(*) FROM %s WHERE user_id=$1",
			UserRoleTableFull,
		),
		u.ID,
	)

	assert.Nil(err)
	return res
}

// CountRoleUsers return all Users belong to Role (many to many with UserRole)
func (m *Manager) CountRoleUsers(r *Role) int64 {
	res, err := m.GetDbMap().SelectInt(
		fmt.Sprintf(
			"SELECT COUNT(*) FROM %s WHERE role_id=$1",
			UserRoleTableFull,
		),
		r.ID,
	)

	assert.Nil(err)
	return res
}
