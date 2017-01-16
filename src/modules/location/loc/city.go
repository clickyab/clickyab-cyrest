package loc

import (
	"common/assert"
	"fmt"
)

// City model
// @Model {
//		table = city
//		primary = true, id
//		find_by = id, name ,province_id
//		list = yes
// }
type City struct {
	ID         int64  `db:"id" json:"id"`
	Name       string `db:"name" json:"name"`
	ProvinceID int64  `db:"province_id" json:"province_id"`
	Latitude   string `db:"latitude" json:"latitude"`
	Longitude  string `db:"longitude" json:"longitude"`
}

// ListCityByProvinceID return the Provinces base on their country_id
func (m *Manager) ListCityByProvinceID(pi int64) []City {
	var res []City
	_, err := m.GetDbMap().Select(
		&res,
		fmt.Sprintf("SELECT * FROM %s WHERE province_id=?", CityTableFull),
		pi,
	)
	assert.Nil(err)

	return res
}
