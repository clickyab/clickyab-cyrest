package aaa

import (
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

// UserProfileCorporation model
// @Model {
//		table = user_profile_personal
//		primary = false, user_id
//		find_by = user_id
// }
type UserProfilePersonal struct {
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
	CreatedAt    time.Time         `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time         `db:"updated_at" json:"updated_at"`
}

// NewUserProfilePersonal is the minimum user profile
func NewUserProfilePersonal(first, last string, gender ProfileGender, cell string) *UserProfilePersonal {
	return &UserProfilePersonal{
		FirstName: first,
		LastName:  last,
		Gender:    gender,
		CellPhone: common.NullString{
			Valid:  len(cell) > 0,
			String: cell,
		},
	}
}
