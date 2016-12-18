package loc

import (
	"fmt"
	"common/assert"
)

// Province model
// @Model {
//		table = province
//		primary = true, id
//		find_by = id, name ,country_id
//		list = yes
// }
type Province struct {
	ID        int64  `db:"id" json:"id"`
	Name      string `db:"name" json:"name"`
	CountryID int64  `db:"country_id" json:"country_id"`
	Latitude  string `db:"latitude" json:"latitude"`
	Longitude string `db:"longitude" json:"longitude"`
}

// ListProvinceByCountryID return the Provinces base on their country_id
func (m *Manager) ListProvinceByCountryID(ci int64) ([]Province) {
	var res []Province
	_, err := m.GetDbMap().Select(
		&res,
		fmt.Sprintf("SELECT * FROM %s WHERE country_id=?", ProvinceTableFull),
		ci,
	)
	assert.Nil(err)

	return res
}
