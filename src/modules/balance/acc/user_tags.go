package acc

// UnitTags is my try to achieve tags in units
// @Model {
//		table = user_tags
//		schema = acc
//		primary = true, id
//		belong_to = aaa.User:user_id
// }
type UserTags struct {
	ID     int64  `db:"id" json:"id"`
	UserID int64  `db:"user_id" json:"user_id"`
	Tag    string `db:"tag" json:"tag"`
	Count  int    `db:"count" json:"count"`
}
