package aaa

import (
	"common/assert"
	"common/models/common"
	"time"
)

// UserProfileCorporation model
// @Model {
//		table = user_profile_corporation
//		primary = false, user_id
//		find_by = user_id
// }
type UserProfileCorporation struct {
	UserID       int64             `db:"user_id" json:"user_id"`
	Title        string            `db:"title" json:"title"`
	EconomicCode common.NullString `db:"economic_code" json:"economic_code"`
	RegisterCode common.NullString `db:"register_code" json:"register_code"`
	Phone        common.NullString `db:"phone" json:"phone"`
	Address      common.NullString `db:"address" json:"address"`
	CountryID    common.NullInt64  `db:"country_id" json:"country_id"`
	ProvinceID   common.NullInt64  `db:"province_id" json:"province_id"`
	CityID       common.NullInt64  `db:"city_id" json:"city_id"`
	CreatedAt    time.Time         `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time         `db:"updated_at" json:"updated_at"`
}

// NewUserProfileCorporation is the minimum user profile
func NewUserProfileCorporation(title string, phone string) *UserProfileCorporation {
	return &UserProfileCorporation{
		Title: title,
		Phone: common.NullString{
			Valid:  len(phone) > 0,
			String: phone,
		},
	}
}

// DeleteCorporation delete the selected user profile corporation
func (m *Manager) DeleteCorporation(upc *UserProfileCorporation) error {
	_, err := m.GetDbMap().Delete(upc)
	assert.Nil(err)
	return err
}

// RegisterCorporation is try to register personal
func (m *Manager) RegisterCorporation(userID int64,
	title string,
	economicCode,
	registerCode string,
	phone string,
	address string,
	countryID int64,
	provinceID int64,
	cityID int64) (cpp *UserProfileCorporation, err error) {

	cpp = &UserProfileCorporation{
		UserID:       userID,
		Title:        title,
		EconomicCode: common.NullString{Valid: len(economicCode) > 0, String: economicCode},
		RegisterCode: common.NullString{Valid: len(registerCode) > 0, String: registerCode},
		Phone:        common.NullString{Valid: len(phone) > 0, String: phone},
		Address:      common.NullString{Valid: len(address) > 0, String: address},
		CountryID:    common.NullInt64{Valid: countryID > 0, Int64: countryID},
		ProvinceID:   common.NullInt64{Valid: provinceID > 0, Int64: provinceID},
		CityID:       common.NullInt64{Valid: cityID > 0, Int64: cityID},
	}

	err = m.Begin()
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			assert.Nil(m.Rollback())
		} else {
			err = m.Commit()
		}

		if err != nil {
			cpp = nil
		}
	}()

	//delete user profile
	fupp, err := m.FindUserProfilePersonalByUserID(userID)
	if err == nil {
		//delete the user profile row
		m.DeletePersonal(fupp)
	}

	//delete corporation profile row
	fucp, err := m.FindUserProfileCorporationByUserID(userID)
	if err == nil {
		//delete the user corporation profile row
		m.DeleteCorporation(fucp)
	}

	//create user profile personal
	err = m.CreateUserProfileCorporation(cpp)
	if err != nil {
		cpp = nil
		return
	}

	//update user type
	err = m.ChangeUserType(userID, UserTypeCorporation)
	if err != nil {
		cpp = nil
		return
	}

	return
}
