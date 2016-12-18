package loc

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
