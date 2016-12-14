// Package cat is the models for category module
package cat

import "time"

// Category model
// @Model {
//		table = categories
//		primary = true, id
//		find_by = id, title
//		list = yes
// }
type Category struct {
	ID        int64     `db:"id" json:"id"`
	Scope     string    `db:"scope" json:"scope"`
	Title     string    `db:"title" json:"title"`
	CreatedAt time.Time `db:"created_at" json:"created_at" sort:"true"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at" sort:"true"`
}
