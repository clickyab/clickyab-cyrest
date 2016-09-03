package aaa

import (
	"common/models"
	"common/models/common"

	"gopkg.in/gorp.v1"
)

// AUTO GENERATED CODE. DO NOT EDIT!

const (
	// MessageLogTableFull is the MessageLog table name
	MessageLogTableFull = "message_logs"

	// ReservedUserTableFull is the ReservedUser table name
	ReservedUserTableFull = "reserved_users"

	// RoleTableFull is the Role table name
	RoleTableFull = "roles"

	// UserRoleTableFull is the UserRole table name
	UserRoleTableFull = "user_roles"

	// UserTableFull is the User table name
	UserTableFull = "users"
)

// Manager is the model manager for aaa package
type Manager struct {
	common.Manager
}

// NewAaaManager create and return a manager for this module
func NewAaaManager() *Manager {
	return &Manager{}
}

// NewAaaManagerFromTransaction create and return a manager for this module from a transaction
func NewAaaManagerFromTransaction(tx gorp.SqlExecutor) (*Manager, error) {
	m := &Manager{}
	err := m.Hijack(tx)
	if err != nil {
		return nil, err
	}

	return m, nil
}

// Initialize aaa package
func (m *Manager) Initialize() {

	m.AddTableWithName(
		MessageLog{},
		MessageLogTableFull,
	).SetKeys(
		true,
		"ID",
	)

	m.AddTableWithName(
		ReservedUser{},
		ReservedUserTableFull,
	).SetKeys(
		true,
		"ID",
	)

	m.AddTableWithName(
		Role{},
		RoleTableFull,
	).SetKeys(
		true,
		"ID",
	)

	m.AddTableWithName(
		UserRole{},
		UserRoleTableFull,
	).SetKeys(
		false,
		"UserID",
		"RoleID",
	)

	m.AddTableWithName(
		User{},
		UserTableFull,
	).SetKeys(
		true,
		"ID",
	)

}
func init() {
	models.Register(NewAaaManager())
}
