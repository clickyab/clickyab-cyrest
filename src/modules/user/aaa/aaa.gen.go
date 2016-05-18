package aaa

import (
	"common/models"
	"common/models/common"

	"gopkg.in/gorp.v1"
)

// AUTO GENERATED CODE. DO NOT EDIT!

const (
	// MessageLogSchema is the MessageLog module schema
	MessageLogSchema = "aaa"
	// MessageLogTable is the MessageLog table name
	MessageLogTable = "message_logs"
	// MessageLogTableFull is the MessageLog table name with schema
	MessageLogTableFull = MessageLogSchema + "." + MessageLogTable

	// ReservedUserSchema is the ReservedUser module schema
	ReservedUserSchema = "aaa"
	// ReservedUserTable is the ReservedUser table name
	ReservedUserTable = "reserved_users"
	// ReservedUserTableFull is the ReservedUser table name with schema
	ReservedUserTableFull = ReservedUserSchema + "." + ReservedUserTable

	// RoleSchema is the Role module schema
	RoleSchema = "aaa"
	// RoleTable is the Role table name
	RoleTable = "roles"
	// RoleTableFull is the Role table name with schema
	RoleTableFull = RoleSchema + "." + RoleTable

	// UserRoleSchema is the UserRole module schema
	UserRoleSchema = "aaa"
	// UserRoleTable is the UserRole table name
	UserRoleTable = "user_roles"
	// UserRoleTableFull is the UserRole table name with schema
	UserRoleTableFull = UserRoleSchema + "." + UserRoleTable

	// UserSchema is the User module schema
	UserSchema = "aaa"
	// UserTable is the User table name
	UserTable = "users"
	// UserTableFull is the User table name with schema
	UserTableFull = UserSchema + "." + UserTable
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

	m.AddTableWithNameAndSchema(
		MessageLog{},
		MessageLogSchema,
		MessageLogTable,
	).SetKeys(
		true,
		"ID",
	)

	m.AddTableWithNameAndSchema(
		ReservedUser{},
		ReservedUserSchema,
		ReservedUserTable,
	).SetKeys(
		true,
		"ID",
	)

	m.AddTableWithNameAndSchema(
		Role{},
		RoleSchema,
		RoleTable,
	).SetKeys(
		true,
		"ID",
	)

	m.AddTableWithNameAndSchema(
		UserRole{},
		UserRoleSchema,
		UserRoleTable,
	).SetKeys(
		false,
		"UserID",
		"RoleID",
	)

	m.AddTableWithNameAndSchema(
		User{},
		UserSchema,
		UserTable,
	).SetKeys(
		true,
		"ID",
	)

}
func init() {
	models.Register(NewAaaManager())
}
