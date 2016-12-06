package aaa

import (
	"common/assert"
	"common/models/common"
	"fmt"
	"strings"
	"time"
)

// AUTO GENERATED CODE. DO NOT EDIT!

// CreateUserFinancial try to save a new UserFinancial in database
func (m *Manager) CreateUserFinancial(uf *UserFinancial) error {
	now := time.Now()
	uf.CreatedAt = now
	uf.UpdatedAt = now
	func(in interface{}) {
		if ii, ok := in.(common.Initializer); ok {
			ii.Initialize()
		}
	}(uf)

	return m.GetDbMap().Insert(uf)
}

// UpdateUserFinancial try to update UserFinancial in database
func (m *Manager) UpdateUserFinancial(uf *UserFinancial) error {
	uf.UpdatedAt = time.Now()
	func(in interface{}) {
		if ii, ok := in.(common.Initializer); ok {
			ii.Initialize()
		}
	}(uf)

	_, err := m.GetDbMap().Update(uf)
	return err
}

// ListUserFinancialsWithFilter try to list all UserFinancials without pagination
func (m *Manager) ListUserFinancialsWithFilter(filter string, params ...interface{}) []UserFinancial {
	filter = strings.Trim(filter, "\n\t ")
	if filter != "" {
		filter = "WHERE " + filter
	}
	var res []UserFinancial
	_, err := m.GetDbMap().Select(
		&res,
		fmt.Sprintf("SELECT * FROM %s %s", UserFinancialTableFull, filter),
		params...,
	)
	assert.Nil(err)

	return res
}

// ListUserFinancials try to list all UserFinancials without pagination
func (m *Manager) ListUserFinancials() []UserFinancial {
	return m.ListUserFinancialsWithFilter("")
}

// CountUserFinancialsWithFilter count entity in UserFinancials table with valid where filter
func (m *Manager) CountUserFinancialsWithFilter(filter string, params ...interface{}) int64 {
	filter = strings.Trim(filter, "\n\t ")
	if filter != "" {
		filter = "WHERE " + filter
	}
	cnt, err := m.GetDbMap().SelectInt(
		fmt.Sprintf("SELECT COUNT(*) FROM %s %s", UserFinancialTableFull, filter),
		params...,
	)
	assert.Nil(err)

	return cnt
}

// CountUserFinancials count entity in UserFinancials table
func (m *Manager) CountUserFinancials() int64 {
	return m.CountUserFinancialsWithFilter("")
}

// ListUserFinancialsWithPaginationFilter try to list all UserFinancials with pagination and filter
func (m *Manager) ListUserFinancialsWithPaginationFilter(
	offset, perPage int, filter string, params ...interface{}) []UserFinancial {
	var res []UserFinancial
	filter = strings.Trim(filter, "\n\t ")
	if filter != "" {
		filter = "WHERE " + filter
	}

	filter += " LIMIT ?, ? "
	params = append(params, offset, perPage)

	// TODO : better pagination without offset and limit
	_, err := m.GetDbMap().Select(
		&res,
		fmt.Sprintf("SELECT * FROM %s %s", UserFinancialTableFull, filter),
		params...,
	)
	assert.Nil(err)

	return res
}

// ListUserFinancialsWithPagination try to list all UserFinancials with pagination
func (m *Manager) ListUserFinancialsWithPagination(offset, perPage int) []UserFinancial {
	return m.ListUserFinancialsWithPaginationFilter(offset, perPage, "")
}

// FindUserFinancialByID return the UserFinancial base on its id
func (m *Manager) FindUserFinancialByID(id int64) (*UserFinancial, error) {
	var res UserFinancial
	err := m.GetDbMap().SelectOne(
		&res,
		fmt.Sprintf("SELECT * FROM %s WHERE id=?", UserFinancialTableFull),
		id,
	)

	if err != nil {
		return nil, err
	}

	return &res, nil
}

// FindUserFinancialByUserID return the UserFinancial base on its user_id
func (m *Manager) FindUserFinancialByUserID(ui int64) (*UserFinancial, error) {
	var res UserFinancial
	err := m.GetDbMap().SelectOne(
		&res,
		fmt.Sprintf("SELECT * FROM %s WHERE user_id=?", UserFinancialTableFull),
		ui,
	)

	if err != nil {
		return nil, err
	}

	return &res, nil
}
