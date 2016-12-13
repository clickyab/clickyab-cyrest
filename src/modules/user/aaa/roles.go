package aaa

import "time"

// Role model
// @Model {
//		table = roles
//		primary = true, id
//		find_by = id,name
//		list = yes
// }
type Role struct {
	ID        int64     `db:"id" json:"id"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
