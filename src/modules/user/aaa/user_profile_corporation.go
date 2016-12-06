package aaa

import (
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
	RegisterCode string            `db:"register_code" json:"register_code"`
	Phone        string            `db:"phone" json:"phone"`
	Address      common.NullString `db:"address" json:"address"`
	CountryID    int64             `db:"country_id" json:"country_id"`
	ProvinceID   int64             `db:"province_id" json:"province_id"`
	CityID       int64             `db:"city_id" json:"city_id"`
	CreatedAt    time.Time         `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time         `db:"updated_at" json:"updated_at"`
}
