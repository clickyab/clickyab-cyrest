package aaa

import (
	"common/models/common"
	"time"
)

// AUTO GENERATED CODE. DO NOT EDIT!

// CreateUserRole try to save a new UserRole in database
func (m *Manager) CreateUserRole(ur *UserRole) error {
	now := time.Now()
	ur.CreatedAt = now

	func(in interface{}) {
		if ii, ok := in.(common.Initializer); ok {
			ii.Initialize()
		}
	}(ur)

	return m.GetDbMap().Insert(ur)
}

// UpdateUserRole try to update UserRole in database
func (m *Manager) UpdateUserRole(ur *UserRole) error {

	func(in interface{}) {
		if ii, ok := in.(common.Initializer); ok {
			ii.Initialize()
		}
	}(ur)

	_, err := m.GetDbMap().Update(ur)
	return err
}
