package aaa

import (
	"common/models/common"
	"fmt"
)

// AUTO GENERATED CODE. DO NOT EDIT!

// CreateUserCRM try to save a new UserCRM in database
func (m *Manager) CreateUserCRM(ucrm *UserCRM) error {

	func(in interface{}) {
		if ii, ok := in.(common.Initializer); ok {
			ii.Initialize()
		}
	}(ucrm)

	return m.GetDbMap().Insert(ucrm)
}

// UpdateUserCRM try to update UserCRM in database
func (m *Manager) UpdateUserCRM(ucrm *UserCRM) error {

	func(in interface{}) {
		if ii, ok := in.(common.Initializer); ok {
			ii.Initialize()
		}
	}(ucrm)

	_, err := m.GetDbMap().Update(ucrm)
	return err
}

// FindUserCRMByUserID return the UserCRM base on its user_id
func (m *Manager) FindUserCRMByUserID(ui int64) (*UserCRM, error) {
	var res UserCRM
	err := m.GetDbMap().SelectOne(
		&res,
		fmt.Sprintf("SELECT * FROM %s WHERE user_id=?", UserCRMTableFull),
		ui,
	)

	if err != nil {
		return nil, err
	}

	return &res, nil
}
