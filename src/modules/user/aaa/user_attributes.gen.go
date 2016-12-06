package aaa

import (
	"common/models/common"
	"fmt"
)

// AUTO GENERATED CODE. DO NOT EDIT!

// CreateUserAttributes try to save a new UserAttributes in database
func (m *Manager) CreateUserAttributes(ua *UserAttributes) error {

	func(in interface{}) {
		if ii, ok := in.(common.Initializer); ok {
			ii.Initialize()
		}
	}(ua)

	return m.GetDbMap().Insert(ua)
}

// UpdateUserAttributes try to update UserAttributes in database
func (m *Manager) UpdateUserAttributes(ua *UserAttributes) error {

	func(in interface{}) {
		if ii, ok := in.(common.Initializer); ok {
			ii.Initialize()
		}
	}(ua)

	_, err := m.GetDbMap().Update(ua)
	return err
}

// FindUserAttributesByKey return the UserAttributes base on its key
func (m *Manager) FindUserAttributesByKey(k string) (*UserAttributes, error) {
	var res UserAttributes
	err := m.GetDbMap().SelectOne(
		&res,
		fmt.Sprintf("SELECT * FROM %s WHERE key=?", UserAttributesTableFull),
		k,
	)

	if err != nil {
		return nil, err
	}

	return &res, nil
}
