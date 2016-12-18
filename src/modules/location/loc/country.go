package loc

// Country model
// @Model {
//		table = country
//		primary = true, id
//		find_by = id, name
//		list = yes
// }
type Country struct {
	ID        int64  `db:"id" json:"id"`
	Name      string `db:"name" json:"name"`
	Latitude  string `db:"latitude" json:"latitude"`
	Longitude string `db:"longitude" json:"longitude"`
}
