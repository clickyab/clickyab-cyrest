package aaa

import (
	"common/assert"
	"common/models/common"
	"time"
)

type (
	//ProfileGender is the profile gender
	// @Enum {
	// }
	ProfileGender string
)

const (
	// ProfileGenderMale is the male
	ProfileGenderMale ProfileGender = "male"

	// ProfileGenderFemale is the female
	ProfileGenderFemale ProfileGender = "female"
)

// UserProfile model
// @Model {
//		table = user_profile
//		primary = false, user_id
//		find_by = user_id
// }
type UserProfile struct {
	UserID       int64             `db:"user_id" json:"user_id"`
	FirstName    string            `db:"first_name" json:"first_name"`
	LastName     string            `db:"last_name" json:"last_name"`
	Birthday     common.NullTime   `db:"birthday" json:"birthday"`
	Gender       ProfileGender     `db:"gender" json:"gender"`
	CellPhone    common.NullString `db:"cellphone" json:"cellphone"`
	Phone        common.NullString `db:"phone" json:"phone"`
	Address      common.NullString `db:"address" json:"address"`
	ZipCode      common.NullString `db:"zip_code" json:"zip_code"`
	NationalCode common.NullString `db:"national_code" json:"national_code"`
	CountryID    common.NullInt64  `db:"country_id" json:"country_id"`
	ProvinceID   common.NullInt64  `db:"province_id" json:"province_id"`
	CityID       common.NullInt64  `db:"city_id" json:"city_id"`
	CreatedAt    *time.Time        `db:"created_at" json:"created_at"`
	UpdatedAt    *time.Time        `db:"updated_at" json:"updated_at"`
}

// NewUserProfile is the minimum user profile
func NewUserProfile(first, last string, gender ProfileGender, cell string) *UserProfile {
	return &UserProfile{
		FirstName: first,
		LastName:  last,
		Gender:    gender,
		CellPhone: common.NullString{
			Valid:  len(cell) > 0,
			String: cell,
		},
	}
}

// DeleteProfile delete the user profile
func (m *Manager) DeleteProfile(upp *UserProfile) error {
	_, err := m.GetDbMap().Delete(upp)
	assert.Nil(err)
	return err
}

// RegisterProfile is try to register profile
func (m *Manager) RegisterProfile(userID int64,
	firstName string,
	lastName string,
	birthday time.Time,
	gender ProfileGender,
	cell string,
	phone string,
	address string,
	zipCode string,
	nationalCode string,
	countryID int64,
	provinceID int64,
	cityID int64) (upp *UserProfile, err error) {

	upp = &UserProfile{
		UserID:       userID,
		FirstName:    firstName,
		LastName:     lastName,
		Birthday:     common.NullTime{Valid: !birthday.IsZero(), Time: birthday},
		Gender:       gender,
		CellPhone:    common.NullString{Valid: len(cell) > 0, String: cell},
		Phone:        common.NullString{Valid: len(phone) > 0, String: phone},
		Address:      common.NullString{Valid: len(address) > 0, String: address},
		ZipCode:      common.NullString{Valid: len(zipCode) > 0, String: zipCode},
		NationalCode: common.NullString{Valid: len(nationalCode) > 0, String: nationalCode},
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
			upp = nil
		}
	}()

	//delete user profile
	fupp, err := m.FindUserProfileByUserID(userID)
	if err == nil {
		//delete the user profile row
		err := m.DeleteProfile(fupp)
		if err != nil {
			return nil, err
		}
	}

	//create user profile
	err = m.CreateUserProfile(upp)
	if err != nil {
		upp = nil
		return
	}

	return
}
