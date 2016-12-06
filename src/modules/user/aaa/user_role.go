package aaa

import "time"

// UserRole model
// @Model {
//		table = user_role
//		primary = false, user_id, role_id
// }
type UserRole struct {
	UserID    int64     `db:"user_id" json:"user_id"`
	RoleID    int64     `db:"role_id" json:"role_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
