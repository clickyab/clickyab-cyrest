// Package cat is the models for category module
package cat

import (
	"common/assert"
	"fmt"
	"modules/misc/base"
	"strings"
	"time"

	"common/models/common"

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
	ID          int64      `db:"id" json:"id" sort:"true" title:"ID"`
	Title       string     `db:"title" json:"title" search:"true" title:"Title"`
	Description string     `db:"description" json:"description" title:"Description"`
	CreatedAt   *time.Time `db:"created_at" json:"created_at" sort:"true" title:"Created at"`
	UpdatedAt   *time.Time `db:"updated_at" json:"updated_at" sort:"true" title:"Updated at"`
}

// CategoryAd model
// @Model {
//		table = ad_category
// }
type CategoryAd struct {
	AdID  int64 `db:"ad_id"`
	CatID int64 `db:"category_id"`
}

//type (
//	CategoryDataTableArray []CategoryDataTable
//
//)

//CategoryDataTable is the role full data in data table, after join with other field
// @DataTable {
//		url = /list
//		entity = category
//		view = category_list:global
//		controller = modules/category/controllers
//		fill = FillCategoryDataTableArray
//		_edit = category_edit:global
//		_change = category_manage:global
// }
type CategoryDataTable struct {
	Category
	ParentID common.NullInt64 `db:"-" json:"parent_id" visible:"false"`
	OwnerID  int64            `db:"-" json:"owner_id" visible:"false"`
	Actions  string           `db:"-" json:"_actions" visible:"false"`
}

// Create is for create category
func (m *Manager) Create(title string, description string) *Category {
	c := &Category{
		Title:       title,
		Description: description,
	}
	assert.Nil(m.CreateCategory(c))
	return c
}

//FetchCategory get all categories query
func (m *Manager) FetchCategory() []Category {
	var res []Category
	query := "SELECT * FROM categories"
	_, err := m.GetDbMap().Select(&res, query)
	assert.Nil(err)
	return res
}

// FillCategoryDataTableArray is the function to handle
func (m *Manager) FillCategoryDataTableArray(
	u base.PermInterfaceComplete,
	filters map[string]string,
	search map[string]string,
	contextparams map[string]string,
	sort, order string,
	p, c int) (CategoryDataTableArray, int64) {
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
		where = append(where, fmt.Sprintf("%s LIKE ?", column))
		params = append(params, fmt.Sprintf("%s"+val+"%s", "%", "%"))
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

// DeleteAdCatByAdID delete pivot by adID
func (m *Manager) DeleteAdCatByAdID(id int64) error {
	q := fmt.Sprintf("DELETE FROM %s WHERE ad_id=?", CategoryAdTableFull)
	_, err := m.GetDbMap().Exec(q, id)
	return err
}

// AddCat add category
func (m *Manager) AddCat(adID int64, catIDs []int64) error {
	var err error
	for i := range catIDs {
		catAd := &CategoryAd{
			AdID:  adID,
			CatID: catIDs[i],
		}
		err = m.GetDbMap().Insert(catAd)
		logrus.Warn(err)
		if err != nil {
			return err
		}
	}
	return err
}

// AssignCat assign category to ad (transaction)
func (m *Manager) AssignCat(id int64, catIDs []int64) error {
	err := m.Begin()
	assert.Nil(err)
	defer func() {
		if err != nil {
			err = m.Rollback()
			assert.Nil(err)
		} else {
			err = m.Commit()
			assert.Nil(err)
		}
	}()
	err = m.DeleteAdCatByAdID(id)
	if err != nil {
		return err
	}
	err = m.AddCat(id, catIDs)
	return err
}
