package t9n

import (
	"common/models"
	"common/models/common"

	"gopkg.in/gorp.v1"
)

// AUTO GENERATED CODE. DO NOT EDIT!

const (
	// TranslationTableFull is the Translation table name
	TranslationTableFull = "translations"
)

// Manager is the model manager for t9n package
type Manager struct {
	common.Manager
}

// NewT9nManager create and return a manager for this module
func NewT9nManager() *Manager {
	return &Manager{}
}

// NewT9nManagerFromTransaction create and return a manager for this module from a transaction
func NewT9nManagerFromTransaction(tx gorp.SqlExecutor) (*Manager, error) {
	m := &Manager{}
	err := m.Hijack(tx)
	if err != nil {
		return nil, err
	}

	return m, nil
}

// Initialize t9n package
func (m *Manager) Initialize() {

	m.AddTableWithName(
		Translation{},
		TranslationTableFull,
	).SetKeys(
		true,
		"ID",
	)

}
func init() {
	models.Register(NewT9nManager())
}
