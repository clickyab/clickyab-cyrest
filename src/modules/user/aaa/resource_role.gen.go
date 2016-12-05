package aaa

import (
	"common/assert"
	"common/models/common"
	"fmt"
	"strings"
)

// AUTO GENERATED CODE. DO NOT EDIT!

// CreateRolePermission try to save a new RolePermission in database
func (m *Manager) CreateRolePermission(rp *RolePermission) error {

	func(in interface{}) {
		if ii, ok := in.(common.Initializer); ok {
			ii.Initialize()
		}
	}(rp)

	return m.GetDbMap().Insert(rp)
}

// UpdateRolePermission try to update RolePermission in database
func (m *Manager) UpdateRolePermission(rp *RolePermission) error {

	func(in interface{}) {
		if ii, ok := in.(common.Initializer); ok {
			ii.Initialize()
		}
	}(rp)

	_, err := m.GetDbMap().Update(rp)
	return err
}

// ListRolePermissionsWithFilter try to list all RolePermissions without pagination
func (m *Manager) ListRolePermissionsWithFilter(filter string, params ...interface{}) []RolePermission {
	filter = strings.Trim(filter, "\n\t ")
	if filter != "" {
		filter = "WHERE " + filter
	}
	var res []RolePermission
	_, err := m.GetDbMap().Select(
		&res,
		fmt.Sprintf("SELECT * FROM %s %s", RolePermissionTableFull, filter),
		params...,
	)
	assert.Nil(err)

	return res
}

// ListRolePermissions try to list all RolePermissions without pagination
func (m *Manager) ListRolePermissions() []RolePermission {
	return m.ListRolePermissionsWithFilter("")
}

// CountRolePermissionsWithFilter count entity in RolePermissions table with valid where filter
func (m *Manager) CountRolePermissionsWithFilter(filter string, params ...interface{}) int64 {
	filter = strings.Trim(filter, "\n\t ")
	if filter != "" {
		filter = "WHERE " + filter
	}
	cnt, err := m.GetDbMap().SelectInt(
		fmt.Sprintf("SELECT COUNT(*) FROM %s %s", RolePermissionTableFull, filter),
		params...,
	)
	assert.Nil(err)

	return cnt
}

// CountRolePermissions count entity in RolePermissions table
func (m *Manager) CountRolePermissions() int64 {
	return m.CountRolePermissionsWithFilter("")
}

// ListRolePermissionsWithPaginationFilter try to list all RolePermissions with pagination and filter
func (m *Manager) ListRolePermissionsWithPaginationFilter(
	offset, perPage int, filter string, params ...interface{}) []RolePermission {
	var res []RolePermission
	filter = strings.Trim(filter, "\n\t ")
	if filter != "" {
		filter = "WHERE " + filter
	}

	filter += " LIMIT ?, ? "
	params = append(params, offset, perPage)

	// TODO : better pagination without offset and limit
	_, err := m.GetDbMap().Select(
		&res,
		fmt.Sprintf("SELECT * FROM %s %s", RolePermissionTableFull, filter),
		params...,
	)
	assert.Nil(err)

	return res
}

// ListRolePermissionsWithPagination try to list all RolePermissions with pagination
func (m *Manager) ListRolePermissionsWithPagination(offset, perPage int) []RolePermission {
	return m.ListRolePermissionsWithPaginationFilter(offset, perPage, "")
}

// FindRolePermissionByID return the RolePermission base on its id
func (m *Manager) FindRolePermissionByID(id int64) (*RolePermission, error) {
	var res RolePermission
	err := m.GetDbMap().SelectOne(
		&res,
		fmt.Sprintf("SELECT * FROM %s WHERE id=?", RolePermissionTableFull),
		id,
	)

	if err != nil {
		return nil, err
	}

	return &res, nil
}

// GetRoleRolePermissions return all RolePermissions belong to Role
func (m *Manager) GetRoleRolePermissions(r *Role) []RolePermission {
	var res []RolePermission
	_, err := m.GetDbMap().Select(
		&res,
		fmt.Sprintf(
			"SELECT * FROM %s WHERE role_id=?",
			RolePermissionTableFull,
		),
		r.ID,
	)

	assert.Nil(err)
	return res
}

// CountRoleRolePermissions return count RolePermissions belong to Role
func (m *Manager) CountRoleRolePermissions(r *Role) int64 {
	res, err := m.GetDbMap().SelectInt(
		fmt.Sprintf(
			"SELECT COUNT(*) FROM %s WHERE role_id=?",
			RolePermissionTableFull,
		),
		r.ID,
	)

	assert.Nil(err)
	return res
}
