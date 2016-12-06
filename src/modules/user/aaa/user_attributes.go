package aaa

// UserAttributes model
// @Model {
//		table = user_attributes
//		primary = true, id
//		find_by = key
// }
type UserAttributes struct {
	ID     int64  `db:"id" json:"id"`
	UserID int64  `db:"user_id" json:"user_id"`
	Key    string `db:"key" json:"key"`
	Value  string `db:"value" json:"value"`
}
