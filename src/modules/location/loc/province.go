package loc

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
