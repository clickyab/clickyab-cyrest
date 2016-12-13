package aaa

import (
	"common/models"
	"common/models/common"

	"gopkg.in/gorp.v1"
)

// AUTO GENERATED CODE. DO NOT EDIT!

const (
	// RolePermissionTableFull is the RolePermission table name
	RolePermissionTableFull = "role_permission"

	// UserAttributesTableFull is the UserAttributes table name
	UserAttributesTableFull = "user_attributes"

	// UserProfileCorporationTableFull is the UserProfileCorporation table name
	UserProfileCorporationTableFull = "user_profile_corporation"

	// UserProfilePersonalTableFull is the UserProfilePersonal table name
	UserProfilePersonalTableFull = "user_profile_personal"

	// DomainTableFull is the Domain table name
	DomainTableFull = "domains"

	// RoleTableFull is the Role table name
	RoleTableFull = "roles"

	// UserFinancialTableFull is the UserFinancial table name
	UserFinancialTableFull = "user_financial"

	// UserRoleTableFull is the UserRole table name
	UserRoleTableFull = "user_role"

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
		RolePermission{},
		RolePermissionTableFull,
	).SetKeys(
		true,
		"ID",
	)

	m.AddTableWithName(
		UserAttributes{},
		UserAttributesTableFull,
	).SetKeys(
		true,
		"ID",
	)

	m.AddTableWithName(
		UserProfileCorporation{},
		UserProfileCorporationTableFull,
	).SetKeys(
		false,
		"UserID",
	)

	m.AddTableWithName(
		UserProfilePersonal{},
		UserProfilePersonalTableFull,
	).SetKeys(
		false,
		"UserID",
	)

	m.AddTableWithName(
		Domain{},
		DomainTableFull,
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
		UserFinancial{},
		UserFinancialTableFull,
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
