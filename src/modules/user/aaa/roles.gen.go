package aaa

import (
	"common/assert"
	"common/models/common"
	"fmt"
	"strings"
	"time"
)

// AUTO GENERATED CODE. DO NOT EDIT!

// CreateRole try to save a new Role in database
func (m *Manager) CreateRole(r *Role) error {
	now := time.Now()
	r.CreatedAt = now
	r.UpdatedAt = now
	func(in interface{}) {
		if ii, ok := in.(common.Initializer); ok {
			ii.Initialize()
		}
	}(r)

	return m.GetDbMap().Insert(r)
}

// UpdateRole try to update Role in database
func (m *Manager) UpdateRole(r *Role) error {
	r.UpdatedAt = time.Now()
	func(in interface{}) {
		if ii, ok := in.(common.Initializer); ok {
			ii.Initialize()
		}
	}(r)

	_, err := m.GetDbMap().Update(r)
	return err
}

// ListRoles try to list all Roles without pagination
func (m *Manager) ListRolesWithFilter(filter string, params ...interface{}) []Role {
	filter = strings.Trim(filter, "\n\t ")
	if filter != "" {
		filter = "WHERE " + filter
	}
	var res []Role
	_, err := m.GetDbMap().Select(
		&res,
		fmt.Sprintf("SELECT * FROM %s %s", RoleTableFull, filter),
		params...,
	)
	assert.Nil(err)

	return res
}

// ListRoles try to list all Roles without pagination
func (m *Manager) ListRoles() []Role {
	return m.ListRolesWithFilter("")
}

// CountRoles count entity in Roles table with valid where filter
func (m *Manager) CountRolesWithFilter(filter string, params ...interface{}) int64 {
	filter = strings.Trim(filter, "\n\t ")
	if filter != "" {
		filter = "WHERE " + filter
	}
	cnt, err := m.GetDbMap().SelectInt(
		fmt.Sprintf("SELECT COUNT(*) FROM %s %s", RoleTableFull, filter),
		params...,
	)
	assert.Nil(err)

	return cnt
}

// CountRoles count entity in Roles table
func (m *Manager) CountRoles() int64 {
	return m.CountRolesWithFilter("")
}

// ListRolesWithPaginationFilter try to list all Roles with pagination and filter
func (m *Manager) ListRolesWithPaginationFilter(offset, perPage int, filter string, params ...interface{}) []Role {
	var res []Role
	filter = strings.Trim(filter, "\n\t ")
	if filter != "" {
		filter = "WHERE " + filter
	}

	filter += fmt.Sprintf(" OFFSET $%d LIMIT $%d", len(params)+1, len(params)+2)
	params = append(params, offset, perPage)

	// TODO : better pagination without offset and limit
	_, err := m.GetDbMap().Select(
		&res,
		fmt.Sprintf("SELECT * FROM %s %s", RoleTableFull, filter),
		params...,
	)
	assert.Nil(err)

	return res
}

// ListRolesWithPagination try to list all Roles with pagination
func (m *Manager) ListRolesWithPagination(offset, perPage int) []Role {
	return m.ListRolesWithPaginationFilter(offset, perPage, "")
}

// FindRoleByID return the Role base on its id
func (m *Manager) FindRoleByID(id int64) (*Role, error) {
	var res Role
	err := m.GetDbMap().SelectOne(
		&res,
		fmt.Sprintf("SELECT * FROM %s WHERE id=$1", RoleTableFull),
		id,
	)

	if err != nil {
		return nil, err
	}

	return &res, nil
}

// FindRoleByName return the Role base on its name
func (m *Manager) FindRoleByName(n string) (*Role, error) {
	var res Role
	err := m.GetDbMap().SelectOne(
		&res,
		fmt.Sprintf("SELECT * FROM %s WHERE name=$1", RoleTableFull),
		n,
	)

	if err != nil {
		return nil, err
	}

	return &res, nil
}
