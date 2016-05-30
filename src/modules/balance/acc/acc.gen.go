package acc

import (
	"common/models"
	"common/models/common"

	"gopkg.in/gorp.v1"
)

// AUTO GENERATED CODE. DO NOT EDIT!

const (
	// AccountSchema is the Account module schema
	AccountSchema = "acc"
	// AccountTable is the Account table name
	AccountTable = "accounts"
	// AccountTableFull is the Account table name with schema
	AccountTableFull = AccountSchema + "." + AccountTable

	// TransactionTagsSchema is the TransactionTags module schema
	TransactionTagsSchema = "acc"
	// TransactionTagsTable is the TransactionTags table name
	TransactionTagsTable = "transaction_tags"
	// TransactionTagsTableFull is the TransactionTags table name with schema
	TransactionTagsTableFull = TransactionTagsSchema + "." + TransactionTagsTable

	// TransactionSchema is the Transaction module schema
	TransactionSchema = "acc"
	// TransactionTable is the Transaction table name
	TransactionTable = "transactions"
	// TransactionTableFull is the Transaction table name with schema
	TransactionTableFull = TransactionSchema + "." + TransactionTable

	// UserTagsSchema is the UserTags module schema
	UserTagsSchema = "acc"
	// UserTagsTable is the UserTags table name
	UserTagsTable = "user_tags"
	// UserTagsTableFull is the UserTags table name with schema
	UserTagsTableFull = UserTagsSchema + "." + UserTagsTable
)

// Manager is the model manager for acc package
type Manager struct {
	common.Manager
}

// NewAccManager create and return a manager for this module
func NewAccManager() *Manager {
	return &Manager{}
}

// NewAccManagerFromTransaction create and return a manager for this module from a transaction
func NewAccManagerFromTransaction(tx gorp.SqlExecutor) (*Manager, error) {
	m := &Manager{}
	err := m.Hijack(tx)
	if err != nil {
		return nil, err
	}

	return m, nil
}

// Initialize acc package
func (m *Manager) Initialize() {

	m.AddTableWithNameAndSchema(
		Account{},
		AccountSchema,
		AccountTable,
	).SetKeys(
		true,
		"ID",
	)

	m.AddTableWithNameAndSchema(
		TransactionTags{},
		TransactionTagsSchema,
		TransactionTagsTable,
	).SetKeys(
		false,
		"TransactionID",
	)

	m.AddTableWithNameAndSchema(
		Transaction{},
		TransactionSchema,
		TransactionTable,
	).SetKeys(
		true,
		"ID",
	)

	m.AddTableWithNameAndSchema(
		UserTags{},
		UserTagsSchema,
		UserTagsTable,
	).SetKeys(
		true,
		"ID",
	)

}
func init() {
	models.Register(NewAccManager())
}
