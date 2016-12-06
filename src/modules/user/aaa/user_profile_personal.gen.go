package aaa

import (
	"common/models/common"
	"fmt"
	"time"
)

// AUTO GENERATED CODE. DO NOT EDIT!

// CreateUserProfilePersonal try to save a new UserProfilePersonal in database
func (m *Manager) CreateUserProfilePersonal(upp *UserProfilePersonal) error {
	now := time.Now()
	upp.CreatedAt = now
	upp.UpdatedAt = now
	func(in interface{}) {
		if ii, ok := in.(common.Initializer); ok {
			ii.Initialize()
		}
	}(upp)

	return m.GetDbMap().Insert(upp)
}

// UpdateUserProfilePersonal try to update UserProfilePersonal in database
func (m *Manager) UpdateUserProfilePersonal(upp *UserProfilePersonal) error {
	upp.UpdatedAt = time.Now()
	func(in interface{}) {
		if ii, ok := in.(common.Initializer); ok {
			ii.Initialize()
		}
	}(upp)

	_, err := m.GetDbMap().Update(upp)
	return err
}

// FindUserProfilePersonalByUserID return the UserProfilePersonal base on its user_id
func (m *Manager) FindUserProfilePersonalByUserID(ui int64) (*UserProfilePersonal, error) {
	var res UserProfilePersonal
	err := m.GetDbMap().SelectOne(
		&res,
		fmt.Sprintf("SELECT * FROM %s WHERE user_id=?", UserProfilePersonalTableFull),
		ui,
	)

	if err != nil {
		return nil, err
	}

	return &res, nil
}
