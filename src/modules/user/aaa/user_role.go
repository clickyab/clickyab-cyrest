package aaa

// UserRole model
// @Model {
//		table = user_roles
//		schema = aaa
//		primary = false, user_id, role_id
//		list = false
//		many_to_many = User:user_id, Role:role_id
// }
type UserRole struct {
	UserID int64 `db:"user_id" json:"user_id"`
	RoleID int64 `db:"role_id" json:"role_id"`
}
