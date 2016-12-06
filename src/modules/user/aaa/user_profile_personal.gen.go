package aaa

import (
	"common/models/common"
	"fmt"
	"time"
)

// AUTO GENERATED CODE. DO NOT EDIT!

// CreateUserProfileCorporation try to save a new UserProfileCorporation in database
func (m *Manager) CreateUserProfileCorporation(upc *UserProfileCorporation) error {
	now := time.Now()
	upc.CreatedAt = now
	upc.UpdatedAt = now
	func(in interface{}) {
		if ii, ok := in.(common.Initializer); ok {
			ii.Initialize()
		}
	}(upc)

	return m.GetDbMap().Insert(upc)
}

// UpdateUserProfileCorporation try to update UserProfileCorporation in database
func (m *Manager) UpdateUserProfileCorporation(upc *UserProfileCorporation) error {
	upc.UpdatedAt = time.Now()
	func(in interface{}) {
		if ii, ok := in.(common.Initializer); ok {
			ii.Initialize()
		}
	}(upc)

	_, err := m.GetDbMap().Update(upc)
	return err
}

// FindUserProfileCorporationByUserID return the UserProfileCorporation base on its user_id
func (m *Manager) FindUserProfileCorporationByUserID(ui int64) (*UserProfileCorporation, error) {
	var res UserProfileCorporation
	err := m.GetDbMap().SelectOne(
		&res,
		fmt.Sprintf("SELECT * FROM %s WHERE user_id=?", UserProfileCorporationTableFull),
		ui,
	)

	if err != nil {
		return nil, err
	}

	return &res, nil
}
