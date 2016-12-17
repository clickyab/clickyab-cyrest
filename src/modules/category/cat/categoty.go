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
	ID          int64     `db:"id" json:"id"`
	Scope       string    `db:"scope" json:"scope"`
	Title       string    `db:"title" json:"title"`
	Description string    `db:"Description" json:"Description"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

func (c *Category) Initialize() {
	if !IsValidScope(c.Scope) {
		logrus.Panic("[BUG] you try to use a scope that is not valid in this app")
	}
}

func (m *Manager) Create(title string, description string, scope string) (Category, error) {
	c := &Category{Title: title, Description: description, Scope: scope}
	err := m.CreateCategory(c)
	return c, err
}
func (m *Manager) Update(title string, description string, scope string, id string) (Category, error) {
	c := &Category{Title: title, Description: description, Scope: scope, ID: id}
	err := m.UpdateCategory(c)
	return c, err
}
