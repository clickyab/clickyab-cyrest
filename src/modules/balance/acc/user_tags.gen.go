package acc

import (
	"common/assert"
	"common/models/common"
	"fmt"
	"modules/user/aaa"
)

// AUTO GENERATED CODE. DO NOT EDIT!

// CreateUserTags try to save a new UserTags in database
func (m *Manager) CreateUserTags(ut *UserTags) error {

	func(in interface{}) {
		if ii, ok := in.(common.Initializer); ok {
			ii.Initialize()
		}
	}(ut)

	return m.GetDbMap().Insert(ut)
}

// UpdateUserTags try to update UserTags in database
func (m *Manager) UpdateUserTags(ut *UserTags) error {

	func(in interface{}) {
		if ii, ok := in.(common.Initializer); ok {
			ii.Initialize()
		}
	}(ut)

	_, err := m.GetDbMap().Update(ut)
	return err
}

// GetUserUserTags return all UserTags belong to User
func (m *Manager) GetUserUserTags(au *aaa.User) []UserTags {
	var res []UserTags
	_, err := m.GetDbMap().Select(
		&res,
		fmt.Sprintf(
			"SELECT * FROM %s WHERE user_id=$1",
			UserTagsTableFull,
		),
		au.ID,
	)

	assert.Nil(err)
	return res
}

// CountUserUserTags return count UserTags belong to User
func (m *Manager) CountUserUserTags(au *aaa.User) int64 {
	res, err := m.GetDbMap().SelectInt(
		fmt.Sprintf(
			"SELECT COUNT(*) FROM %s WHERE user_id=$1",
			UserTagsTableFull,
		),
		au.ID,
	)

	assert.Nil(err)
	return res
}
