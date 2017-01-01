// Package cat is the models for category module
package cat

import (
	"time"

	"common/assert"

	"fmt"
	"modules/misc/base"
	"strings"

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

//type (
//	CategoryDataTableArray []CategoryDataTable
//
//)

//CategoryDataTable is the role full data in data table, after join with other field
// @DataTable {
//		url = /list
//		entity = categories
//		view = category_list:global
//		controller = modules/category/controllers
//		fill = FillCategoryDataTableArray
//		_edit = category_edit:global
// }
type CategoryDataTable struct {
	Category
	ParentID int64 `db:"-" json:"parent_id" visible:"false"`
	OwnerID  int64 `db:"-" json:"owner_id" visible:"false"`
	Actions	string `db:"-" json:"_actions" visible:"false"`
}

func (c *Category) Initialize() {
	if !IsValidScope(c.Scope) {
		logrus.Panic("[BUG] you try to use a scope that is not valid in this app")
	}
}

// Create is for create category
func (m *Manager) Create(title string, description string, scope string) *Category {
	c := &Category{
		Title:       title,
		Description: description,
		Scope:       scope,
	}
	assert.Nil(m.CreateCategory(c))
	return c
}

// FillRoleDataTableArray is the function to handle
func (m *Manager) FillCategoryDataTableArray(u base.PermInterfaceComplete, filters map[string]string, search map[string]string, sort, order string, p, c int) (CategoryDataTableArray, int64) {
	var params []interface{}
	var res CategoryDataTableArray
	var where []string

	countQuery := "SELECT COUNT(id) FROM categories"
	query := "SELECT * FROM categories"
	for field, value := range filters {
		where = append(where, fmt.Sprintf("%s=%s", field, "?"))
		params = append(params, value)
	}

	for column, val := range search {
		where = append(where, fmt.Sprintf("%s=%s", column, "?"))
		params = append(params, val)
	}

	//check for perm
	if len(where) > 0 {
		query += " WHERE "
		countQuery += " WHERE "
	}
	query += strings.Join(where, " AND ")
	countQuery += strings.Join(where, " AND ")
	limit := c
	offset := (p - 1) * c
	if sort != "" {
		query += fmt.Sprintf(" ORDER BY %s %s ", sort, order)
	}
	query += fmt.Sprintf(" LIMIT %d OFFSET %d ", limit, offset)
	count, err := m.GetDbMap().SelectInt(countQuery, params...)
	assert.Nil(err)

	_, err = m.GetDbMap().Select(&res, query, params...)
	assert.Nil(err)
	return res, count
}
