package acc

import (
	"common/assert"
	"common/models/common"
	"fmt"
	"modules/user/aaa"
	"strings"
	"time"

	"gopkg.in/gorp.v1"
)

// AUTO GENERATED CODE. DO NOT EDIT!

// CreateAccount try to save a new Account in database
func (m *Manager) CreateAccount(a *Account) error {
	now := time.Now()
	a.CreatedAt = now
	a.UpdatedAt = now
	func(in interface{}) {
		if ii, ok := in.(common.Initializer); ok {
			ii.Initialize()
		}
	}(a)

	return m.GetDbMap().Insert(a)
}

// UpdateAccount try to update Account in database
func (m *Manager) UpdateAccount(a *Account) error {
	a.UpdatedAt = time.Now()
	func(in interface{}) {
		if ii, ok := in.(common.Initializer); ok {
			ii.Initialize()
		}
	}(a)

	_, err := m.GetDbMap().Update(a)
	return err
}

// ListAccounts try to list all Accounts without pagination
func (m *Manager) ListAccountsWithFilter(filter string, params ...interface{}) []Account {
	filter = strings.Trim(filter, "\n\t ")
	if filter != "" {
		filter = "WHERE " + filter
	}
	var res []Account
	_, err := m.GetDbMap().Select(
		&res,
		fmt.Sprintf("SELECT * FROM %s %s", AccountTableFull, filter),
		params...,
	)
	assert.Nil(err)

	return res
}

// ListAccounts try to list all Accounts without pagination
func (m *Manager) ListAccounts() []Account {
	return m.ListAccountsWithFilter("")
}

// CountAccounts count entity in Accounts table with valid where filter
func (m *Manager) CountAccountsWithFilter(filter string, params ...interface{}) int64 {
	filter = strings.Trim(filter, "\n\t ")
	if filter != "" {
		filter = "WHERE " + filter
	}
	cnt, err := m.GetDbMap().SelectInt(
		fmt.Sprintf("SELECT COUNT(*) FROM %s %s", AccountTableFull, filter),
		params...,
	)
	assert.Nil(err)

	return cnt
}

// CountAccounts count entity in Accounts table
func (m *Manager) CountAccounts() int64 {
	return m.CountAccountsWithFilter("")
}

// ListAccountsWithPaginationFilter try to list all Accounts with pagination and filter
func (m *Manager) ListAccountsWithPaginationFilter(offset, perPage int, filter string, params ...interface{}) []Account {
	var res []Account
	filter = strings.Trim(filter, "\n\t ")
	if filter != "" {
		filter = "WHERE " + filter
	}

	filter += fmt.Sprintf(" OFFSET $%d LIMIT $%d", len(params)+1, len(params)+2)
	params = append(params, offset, perPage)

	// TODO : better pagination without offset and limit
	_, err := m.GetDbMap().Select(
		&res,
		fmt.Sprintf("SELECT * FROM %s %s", AccountTableFull, filter),
		params...,
	)
	assert.Nil(err)

	return res
}

// ListAccountsWithPagination try to list all Accounts with pagination
func (m *Manager) ListAccountsWithPagination(offset, perPage int) []Account {
	return m.ListAccountsWithPaginationFilter(offset, perPage, "")
}

// FindAccountByID return the Account base on its id
func (m *Manager) FindAccountByID(id int64) (*Account, error) {
	var res Account
	err := m.GetDbMap().SelectOne(
		&res,
		fmt.Sprintf("SELECT * FROM %s WHERE id=$1", AccountTableFull),
		id,
	)

	if err != nil {
		return nil, err
	}

	return &res, nil
}

// GetUserAccounts return all Accounts belong to User
func (m *Manager) GetUserAccounts(au *aaa.User) []Account {
	var res []Account
	_, err := m.GetDbMap().Select(
		&res,
		fmt.Sprintf(
			"SELECT * FROM %s WHERE owner_id=$1",
			AccountTableFull,
		),
		au.ID,
	)

	assert.Nil(err)
	return res
}

// CountUserAccounts return count Accounts belong to User
func (m *Manager) CountUserAccounts(au *aaa.User) int64 {
	res, err := m.GetDbMap().SelectInt(
		fmt.Sprintf(
			"SELECT COUNT(*) FROM %s WHERE owner_id=$1",
			AccountTableFull,
		),
		au.ID,
	)

	assert.Nil(err)
	return res
}

// PreInsert is gorp hook to prevent Insert without transaction
func (a *Account) PreInsert(q gorp.SqlExecutor) error {
	if _, ok := q.(*gorp.Transaction); !ok {
		// This is not good. gorp is not in transaction
		return fmt.Errorf("Insert Account must be in transaction")
	}
	return nil
}

// PreUpdate is gorp hook to prevent Update without transaction
func (a *Account) PreUpdate(q gorp.SqlExecutor) error {
	if _, ok := q.(*gorp.Transaction); !ok {
		// This is not good. gorp is not in transaction
		return fmt.Errorf("Update Account must be in transaction")
	}
	return nil
}

// PreDelete is gorp hook to prevent Delete without transaction
func (a *Account) PreDelete(q gorp.SqlExecutor) error {
	if _, ok := q.(*gorp.Transaction); !ok {
		// This is not good. gorp is not in transaction
		return fmt.Errorf("Delete Account must be in transaction")
	}
	return nil
}
