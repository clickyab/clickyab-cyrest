package aaa

import (
	"common/assert"
	"common/models/common"
	"fmt"
	"strings"
	"time"
)

// AUTO GENERATED CODE. DO NOT EDIT!

// CreateReservedUser try to save a new ReservedUser in database
func (m *Manager) CreateReservedUser(ru *ReservedUser) error {
	now := time.Now()
	ru.CreatedAt = now
	ru.UpdatedAt = now
	func(in interface{}) {
		if ii, ok := in.(common.Initializer); ok {
			ii.Initialize()
		}
	}(ru)

	return m.GetDbMap().Insert(ru)
}

// UpdateReservedUser try to update ReservedUser in database
func (m *Manager) UpdateReservedUser(ru *ReservedUser) error {
	ru.UpdatedAt = time.Now()
	func(in interface{}) {
		if ii, ok := in.(common.Initializer); ok {
			ii.Initialize()
		}
	}(ru)

	_, err := m.GetDbMap().Update(ru)
	return err
}

// ListReservedUsers try to list all ReservedUsers without pagination
func (m *Manager) ListReservedUsersWithFilter(filter string, params ...interface{}) []ReservedUser {
	filter = strings.Trim(filter, "\n\t ")
	if filter != "" {
		filter = "WHERE " + filter
	}
	var res []ReservedUser
	_, err := m.GetDbMap().Select(
		&res,
		fmt.Sprintf("SELECT * FROM %s %s", ReservedUserTableFull, filter),
		params...,
	)
	assert.Nil(err)

	return res
}

// ListReservedUsers try to list all ReservedUsers without pagination
func (m *Manager) ListReservedUsers() []ReservedUser {
	return m.ListReservedUsersWithFilter("")
}

// CountReservedUsers count entity in ReservedUsers table with valid where filter
func (m *Manager) CountReservedUsersWithFilter(filter string, params ...interface{}) int64 {
	filter = strings.Trim(filter, "\n\t ")
	if filter != "" {
		filter = "WHERE " + filter
	}
	cnt, err := m.GetDbMap().SelectInt(
		fmt.Sprintf("SELECT COUNT(*) FROM %s %s", ReservedUserTableFull, filter),
		params...,
	)
	assert.Nil(err)

	return cnt
}

// CountReservedUsers count entity in ReservedUsers table
func (m *Manager) CountReservedUsers() int64 {
	return m.CountReservedUsersWithFilter("")
}

// ListReservedUsersWithPaginationFilter try to list all ReservedUsers with pagination and filter
func (m *Manager) ListReservedUsersWithPaginationFilter(offset, perPage int, filter string, params ...interface{}) []ReservedUser {
	var res []ReservedUser
	filter = strings.Trim(filter, "\n\t ")
	if filter != "" {
		filter = "WHERE " + filter
	}

	filter += " LIMIT ?, ? "
	params = append(params, offset, perPage)

	// TODO : better pagination without offset and limit
	_, err := m.GetDbMap().Select(
		&res,
		fmt.Sprintf("SELECT * FROM %s %s", ReservedUserTableFull, filter),
		params...,
	)
	assert.Nil(err)

	return res
}

// ListReservedUsersWithPagination try to list all ReservedUsers with pagination
func (m *Manager) ListReservedUsersWithPagination(offset, perPage int) []ReservedUser {
	return m.ListReservedUsersWithPaginationFilter(offset, perPage, "")
}

// FindReservedUserByID return the ReservedUser base on its id
func (m *Manager) FindReservedUserByID(id int64) (*ReservedUser, error) {
	var res ReservedUser
	err := m.GetDbMap().SelectOne(
		&res,
		fmt.Sprintf("SELECT * FROM %s WHERE id=$1", ReservedUserTableFull),
		id,
	)

	if err != nil {
		return nil, err
	}

	return &res, nil
}

// FindReservedUserByContact return the ReservedUser base on its contact
func (m *Manager) FindReservedUserByContact(c string) (*ReservedUser, error) {
	var res ReservedUser
	err := m.GetDbMap().SelectOne(
		&res,
		fmt.Sprintf("SELECT * FROM %s WHERE contact=$1", ReservedUserTableFull),
		c,
	)

	if err != nil {
		return nil, err
	}

	return &res, nil
}
