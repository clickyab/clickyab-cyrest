package acc

import (
	"common/assert"
	"common/models/common"
	"fmt"
	"strings"
	"time"

	"gopkg.in/gorp.v1"
)

// AUTO GENERATED CODE. DO NOT EDIT!

// CreateTransactionTags try to save a new TransactionTags in database
func (m *Manager) CreateTransactionTags(tt *TransactionTags) error {

	func(in interface{}) {
		if ii, ok := in.(common.Initializer); ok {
			ii.Initialize()
		}
	}(tt)

	return m.GetDbMap().Insert(tt)
}

// UpdateTransactionTags try to update TransactionTags in database
func (m *Manager) UpdateTransactionTags(tt *TransactionTags) error {

	func(in interface{}) {
		if ii, ok := in.(common.Initializer); ok {
			ii.Initialize()
		}
	}(tt)

	_, err := m.GetDbMap().Update(tt)
	return err
}

// FindTransactionTagsByTransactionID return the TransactionTags base on its transaction_id
func (m *Manager) FindTransactionTagsByTransactionID(ti int64) (*TransactionTags, error) {
	var res TransactionTags
	err := m.GetDbMap().SelectOne(
		&res,
		fmt.Sprintf("SELECT * FROM %s WHERE transaction_id=$1", TransactionTagsTableFull),
		ti,
	)

	if err != nil {
		return nil, err
	}

	return &res, nil
}

// PreInsert is gorp hook to prevent Insert without transaction
func (tt *TransactionTags) PreInsert(q gorp.SqlExecutor) error {
	if _, ok := q.(*gorp.Transaction); !ok {
		// This is not good. gorp is not in transaction
		return fmt.Errorf("Insert TransactionTags must be in transaction")
	}
	return nil
}

// PreUpdate is gorp hook to prevent Update without transaction
func (tt *TransactionTags) PreUpdate(q gorp.SqlExecutor) error {
	if _, ok := q.(*gorp.Transaction); !ok {
		// This is not good. gorp is not in transaction
		return fmt.Errorf("Update TransactionTags must be in transaction")
	}
	return nil
}

// CreateTransaction try to save a new Transaction in database
func (m *Manager) CreateTransaction(t *Transaction) error {
	now := time.Now()
	t.CreatedAt = now
	t.UpdatedAt = now
	func(in interface{}) {
		if ii, ok := in.(common.Initializer); ok {
			ii.Initialize()
		}
	}(t)

	return m.GetDbMap().Insert(t)
}

// UpdateTransaction try to update Transaction in database
func (m *Manager) UpdateTransaction(t *Transaction) error {
	t.UpdatedAt = time.Now()
	func(in interface{}) {
		if ii, ok := in.(common.Initializer); ok {
			ii.Initialize()
		}
	}(t)

	_, err := m.GetDbMap().Update(t)
	return err
}

// ListTransactions try to list all Transactions without pagination
func (m *Manager) ListTransactionsWithFilter(filter string, params ...interface{}) []Transaction {
	filter = strings.Trim(filter, "\n\t ")
	if filter != "" {
		filter = "WHERE " + filter
	}
	var res []Transaction
	_, err := m.GetDbMap().Select(
		&res,
		fmt.Sprintf("SELECT * FROM %s %s", TransactionTableFull, filter),
		params...,
	)
	assert.Nil(err)

	return res
}

// ListTransactions try to list all Transactions without pagination
func (m *Manager) ListTransactions() []Transaction {
	return m.ListTransactionsWithFilter("")
}

// CountTransactions count entity in Transactions table with valid where filter
func (m *Manager) CountTransactionsWithFilter(filter string, params ...interface{}) int64 {
	filter = strings.Trim(filter, "\n\t ")
	if filter != "" {
		filter = "WHERE " + filter
	}
	cnt, err := m.GetDbMap().SelectInt(
		fmt.Sprintf("SELECT COUNT(*) FROM %s %s", TransactionTableFull, filter),
		params...,
	)
	assert.Nil(err)

	return cnt
}

// CountTransactions count entity in Transactions table
func (m *Manager) CountTransactions() int64 {
	return m.CountTransactionsWithFilter("")
}

// ListTransactionsWithPaginationFilter try to list all Transactions with pagination and filter
func (m *Manager) ListTransactionsWithPaginationFilter(offset, perPage int, filter string, params ...interface{}) []Transaction {
	var res []Transaction
	filter = strings.Trim(filter, "\n\t ")
	if filter != "" {
		filter = "WHERE " + filter
	}

	filter += fmt.Sprintf(" OFFSET $%d LIMIT $%d", len(params)+1, len(params)+2)
	params = append(params, offset, perPage)

	// TODO : better pagination without offset and limit
	_, err := m.GetDbMap().Select(
		&res,
		fmt.Sprintf("SELECT * FROM %s %s", TransactionTableFull, filter),
		params...,
	)
	assert.Nil(err)

	return res
}

// ListTransactionsWithPagination try to list all Transactions with pagination
func (m *Manager) ListTransactionsWithPagination(offset, perPage int) []Transaction {
	return m.ListTransactionsWithPaginationFilter(offset, perPage, "")
}

// FindTransactionByID return the Transaction base on its id
func (m *Manager) FindTransactionByID(id int64) (*Transaction, error) {
	var res Transaction
	err := m.GetDbMap().SelectOne(
		&res,
		fmt.Sprintf("SELECT * FROM %s WHERE id=$1", TransactionTableFull),
		id,
	)

	if err != nil {
		return nil, err
	}

	return &res, nil
}

// FilterTransactionsByAccountID return all Transactions base on its account_id, panic on query error
func (m *Manager) FilterTransactionsByAccountID(ai int64) []Transaction {
	var res []Transaction
	_, err := m.GetDbMap().Select(
		&res,
		fmt.Sprintf("SELECT * FROM %s WHERE account_id=$1", TransactionTableFull),
		ai,
	)
	assert.Nil(err)

	return res
}

// GetAccountTransactions return all Transactions belong to Account
func (m *Manager) GetAccountTransactions(a *Account) []Transaction {
	var res []Transaction
	_, err := m.GetDbMap().Select(
		&res,
		fmt.Sprintf(
			"SELECT * FROM %s WHERE account_id=$1",
			TransactionTableFull,
		),
		a.ID,
	)

	assert.Nil(err)
	return res
}

// CountAccountTransactions return count Transactions belong to Account
func (m *Manager) CountAccountTransactions(a *Account) int64 {
	res, err := m.GetDbMap().SelectInt(
		fmt.Sprintf(
			"SELECT COUNT(*) FROM %s WHERE account_id=$1",
			TransactionTableFull,
		),
		a.ID,
	)

	assert.Nil(err)
	return res
}

// PreInsert is gorp hook to prevent Insert without transaction
func (t *Transaction) PreInsert(q gorp.SqlExecutor) error {
	if _, ok := q.(*gorp.Transaction); !ok {
		// This is not good. gorp is not in transaction
		return fmt.Errorf("Insert Transaction must be in transaction")
	}
	return nil
}

// PreUpdate is gorp hook to prevent Update without transaction
func (t *Transaction) PreUpdate(q gorp.SqlExecutor) error {
	if _, ok := q.(*gorp.Transaction); !ok {
		// This is not good. gorp is not in transaction
		return fmt.Errorf("Update Transaction must be in transaction")
	}
	return nil
}
