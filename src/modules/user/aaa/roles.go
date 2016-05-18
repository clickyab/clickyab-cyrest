package aaa

import (
	"common/models/common"
	"time"
)

// Role model
// @Model {
//		table = roles
//		schema = aaa
//		primary = true, id
//		find_by = id,name
//		list = yes
// }
type Role struct {
	ID          int64              `db:"id" json:"id"`
	Name        string             `db:"name" json:"name"`
	Description string             `db:"description" json:"description"`
	Resources   common.StringSlice `db:"resources" json:"resources"`
	CreatedAt   time.Time          `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `db:"updated_at" json:"updated_at"`
}
