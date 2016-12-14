// Package cat is the models for category module
package cat

import (
	"time"

	"github.com/Sirupsen/logrus"
)

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
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

func (c *Category) Initialize() {
	if !IsValidScope(c.Scope) {
		logrus.Panic("[BUG] you try to use a scope that is not valid in this app")
	}
}
