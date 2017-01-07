package tlu

import (
	"common/assert"
	"fmt"
	"modules/misc/base"
	"modules/user/aaa"
	"strings"

	"time"
)

//'yes', 'no','yes','no'
const (
	ResolveStatusYes ResolveStatus = "yes"
	ResolveStatusNo  ResolveStatus = "no"

	RemoveStatusYes RemoveStatus = "yes"
	RemoveStatusNo  RemoveStatus = "no"
)

type (
	// ResolveStatus is the telegram user
	// @Enum{
	// }
	ResolveStatus string

	// RemoveStatus is the telegram user
	// @Enum{
	// }
	RemoveStatus string
)

// TeleUser model
// @Model {
//		table = telegram_users
//		primary = true, id
//		find_by = id,user_id
//		list = yes
// }
type Teleuser struct {
	ID         int64         `db:"id" json:"id" sort:"true" title:"ID"`
	UserID     int64         `json:"user_id" db:"user_id" title:"UserID"`
	TelegramID string        `json:"telegram_id" db:"telegram_id" title:"TelegramID"`
	Username   string        `json:"username" db:"username" title:"UserName"`
	Resolve    ResolveStatus `json:"resolve" db:"resolve" title:"Resolve"`
	Remove     RemoveStatus  `json:"remove" db:"remove" title:"Remove"`
	CreatedAt  time.Time     `db:"created_at" json:"created_at" sort:"true" title:"Created at"`
	UpdatedAt  time.Time     `db:"updated_at" json:"updated_at" sort:"true" title:"Updated at"`
}

//TeleuserDataTable is the teleuser full data in data table, after join with other field
// @DataTable {
//		url = /list
//		entity = teleuser
//		view = teleuser_list:self
//		controller = modules/teleuser/controllers
//		fill = FillTeleuserDataTableArray
//		_edit = teleuser_edit:self
// }
type TeleuserDataTable struct {
	Teleuser
	Email    string `db:"email" json:"email" search:"true" title:"Email"`
	ParentID int64  `db:"-" json:"parent_id" visible:"false"`
	OwnerID  int64  `db:"-" json:"owner_id" visible:"false"`
	Actions  string `db:"-" json:"_actions" visible:"false"`
}

// FillTeleUserDataTableArray is the function to handle
func (m *Manager) FillTeleuserDataTableArray(u base.PermInterfaceComplete, filters map[string]string, search map[string]string, sort, order string, p, c int) (TeleuserDataTableArray, int64) {
	var params []interface{}
	var res TeleuserDataTableArray
	var where []string

	countQuery := fmt.Sprintf("SELECT COUNT(telegram_users.id) FROM %s LEFT JOIN %s ON %s.id=%s.user_id", TeleuserTableFull, aaa.UserTableFull, aaa.UserTableFull, TeleuserTableFull)
	query := fmt.Sprintf("SELECT telegram_users.*,users.email FROM %s LEFT JOIN %s ON %s.id=%s.user_id", TeleuserTableFull, aaa.UserTableFull, aaa.UserTableFull, TeleuserTableFull)
	for field, value := range filters {
		where = append(where, fmt.Sprintf(TeleuserTableFull+".%s=%s", field, "?"))
		params = append(params, value)
	}

	for column, val := range search {
		where = append(where, fmt.Sprintf("%s LIKE ?", column))
		params = append(params, fmt.Sprintf("%s"+val+"%s", "%", "%"))
	}

	currentUserID := u.GetID()
	highestScope := u.GetCurrentScope()

	if highestScope == base.ScopeSelf {
		where = append(where, fmt.Sprintf("%s.user_id=?", TeleuserTableFull))
		params = append(params, currentUserID)
	} else if highestScope == base.ScopeParent {
		where = append(where, "users.parent_id=?")
		params = append(params, currentUserID)
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

func (c *Teleuser) Initialize() {

}
