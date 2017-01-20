package chn

import (
	"common/assert"
	"common/models/common"
	"fmt"
	"modules/misc/base"
	"modules/user/aaa"
	"strings"
	"time"
)

//'pending', 'rejected','accepted','archive','yes','no'
const (
	// ChannelStatusPending is the pending status
	ChannelStatusPending ChannelStatus = "pending"
	// ChannelStatusRejected is the rejected status
	ChannelStatusRejected ChannelStatus = "rejected"
	// ChannelStatusAccepted is the accepted status
	ChannelStatusAccepted ChannelStatus = "accepted"
	// ChannelStatusArchive is the archive status
	ChannelStatusArchive ChannelStatus = "archive"

	//ActiveStatusYes is the yes status
	ActiveStatusYes ActiveStatus = "yes"
	// ActiveStatusNo is the no status
	ActiveStatusNo ActiveStatus = "no"
)

type (
	// ChannelStatus is the channel status
	// @Enum{
	// }
	ChannelStatus string

	// ActiveStatus is the channel active
	// @Enum{
	// }
	ActiveStatus string
)

// Channel model
// @Model {
//		table = channels
//		primary = true, id
//		find_by = id,user_id
//		list = yes
// }
type Channel struct {
	ID        int64             `db:"id" json:"id" sort:"true" title:"ID"`
	UserID    int64             `json:"user_id" db:"user_id" title:"UserID"`
	Name      string            `json:"name" db:"name" search:"true" title:"Name"`
	Link      common.NullString `json:"link" db:"link" search:"true" title:"Link"`
	Admin     common.NullString `json:"admin" db:"admin" search:"true" title:"Admin"`
	Code      string            `json:"admin" db:"code" title:"code"`
	Status    ChannelStatus     `json:"status" db:"status" filter:"true" title:"Status"`
	Active    ActiveStatus      `json:"active" db:"active" filter:"true" title:"Active"`
	CreatedAt time.Time         `db:"created_at" json:"created_at" sort:"true" title:"Created at"`
	UpdatedAt time.Time         `db:"updated_at" json:"updated_at" sort:"true" title:"Updated at"`
}

// ChannelCreate a new channel
func (m *Manager) ChannelCreate(admin, link, name string, status ChannelStatus, active ActiveStatus, userID int64) *Channel {

	ch := &Channel{
		Admin:  common.MakeNullString(admin),
		Link:   common.MakeNullString(link),
		Name:   name,
		Status: status,
		Active: active,
		UserID: userID,
	}
	assert.Nil(m.CreateChannel(ch))
	return ch
}

//ChannelDataTable is the role full data in data table, after join with other field
// @DataTable {
//		url = /list
//		entity = channel
//		view = channel_list:self
//		controller = modules/telegram/channel/controllers
//		fill = FillChannelDataTableArray
//		_edit = channel_edit:self
//		_change = channel_manage:global
// }
type ChannelDataTable struct {
	Channel
	Email    string `db:"email" json:"email" search:"true" title:"Email"`
	ParentID int64  `db:"-" json:"parent_id" visible:"false"`
	OwnerID  int64  `db:"-" json:"owner_id" visible:"false"`
	Actions  string `db:"-" json:"_actions" visible:"false"`
}

// FillChannelDataTableArray is the function to handle
func (m *Manager) FillChannelDataTableArray(u base.PermInterfaceComplete, filters map[string]string, search map[string]string, sort, order string, p, c int) (ChannelDataTableArray, int64) {
	var params []interface{}
	var res ChannelDataTableArray
	var where []string

	countQuery := fmt.Sprintf("SELECT COUNT(channels.id) FROM %s LEFT JOIN %s ON %s.id=%s.user_id", ChannelTableFull, aaa.UserTableFull, aaa.UserTableFull, ChannelTableFull)
	query := fmt.Sprintf("SELECT channels.*,users.email FROM %s LEFT JOIN %s ON %s.id=%s.user_id", ChannelTableFull, aaa.UserTableFull, aaa.UserTableFull, ChannelTableFull)
	for field, value := range filters {
		where = append(where, fmt.Sprintf(ChannelTableFull+".%s=%s", field, "?"))
		params = append(params, value)
	}

	for column, val := range search {
		where = append(where, fmt.Sprintf("%s LIKE ?", column))
		params = append(params, fmt.Sprintf("%s"+val+"%s", "%", "%"))
	}

	currentUserID := u.GetID()
	highestScope := u.GetCurrentScope()

	if highestScope == base.ScopeSelf {
		where = append(where, fmt.Sprintf("%s.user_id=?", ChannelTableFull))
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

// EditChannel function for channel editing
func (m *Manager) EditChannel(admin, link, name string, status ChannelStatus, activeStatus ActiveStatus, userID int64, createdAt time.Time, id int64) *Channel {

	ch := &Channel{
		ID:        id,
		UserID:    userID,
		Admin:     common.NullString{Valid: admin != "", String: admin},
		Link:      common.NullString{Valid: link != "", String: link},
		Status:    status,
		Active:    activeStatus,
		Name:      name,
		CreatedAt: createdAt,
	}
	assert.Nil(m.UpdateChannel(ch))
	return ch
}

// ChangeActive toggle between active status
func (m *Manager) ChangeActive(ID int64, userID int64, name string, link string, admin string, status ChannelStatus, currentStat ActiveStatus, createdAt time.Time) *Channel {
	ch := &Channel{
		ID:        ID,
		UserID:    userID,
		Name:      name,
		Link:      common.NullString{Valid: link != "", String: link},
		Admin:     common.NullString{Valid: admin != "", String: admin},
		Status:    status,
		CreatedAt: createdAt,
	}
	if currentStat == ActiveStatusYes {
		ch.Active = ActiveStatusNo
	} else {
		ch.Active = ActiveStatusYes
	}
	assert.Nil(m.UpdateChannel(ch))
	return ch
}
