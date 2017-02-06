package ads

import (
	"common/assert"
	"common/models/common"
	"fmt"
	"modules/misc/base"
	"modules/telegram/teleuser/tlu"
	"modules/user/aaa"
	"strings"
	"time"
)

//'pending', 'rejected','accepted','archive','yes','no'
const (
	// AdminStatusPending is the pending status
	AdminStatusPending AdminStatus = "pending"
	// AdminStatusRejected is the rejected status
	AdminStatusRejected AdminStatus = "rejected"
	// AdminStatusAccepted is the accepted status
	AdminStatusAccepted AdminStatus = "accepted"

	ArchiveStatusYes ArchiveStatus = "yes"
	ArchiveStatusNo  ArchiveStatus = "no"
)

type (
	// AdminStatus is the channel status
	// @Enum{
	// }
	AdminStatus string

	// ArchiveStatus is the channel active
	// @Enum{
	// }
	ArchiveStatus string
)

// Channel model
// @Model {
//		table = channels
//		primary = true, id
//		find_by = id
//		list = yes
// }
type Channel struct {
	ID            int64             `db:"id" json:"id" sort:"true" title:"ID"`
	UserID        int64             `json:"user_id" db:"user_id" title:"UserID"`
	Name          string            `json:"name" db:"name" search:"true" title:"Name"`
	Link          common.NullString `json:"link" db:"link" search:"true" title:"Link"`
	AdminStatus   AdminStatus       `json:"admin_status" db:"admin_status" filter:"true" title:"AdminStatus"`
	ArchiveStatus ArchiveStatus     `json:"archive_status" db:"archive_status" filter:"true" title:"ArchiveStatus"`
	Active        ActiveStatus      `json:"active" db:"active" filter:"true" title:"Active"`
	CreatedAt     time.Time         `db:"created_at" json:"created_at" sort:"true" title:"Created at"`
	UpdatedAt     time.Time         `db:"updated_at" json:"updated_at" sort:"true" title:"Updated at"`
}

// ChannelCreate a new channel
func (m *Manager) ChannelCreate(link, name string, status AdminStatus, archive ArchiveStatus, active ActiveStatus, userID int64) *Channel {

	ch := &Channel{
		Link:        common.MakeNullString(link),
		Name:        name,
		AdminStatus: status,
		Active:      active,
		UserID:      userID,
	}
	assert.Nil(m.CreateChannel(ch))
	return ch
}

// FindChannelsByChatID return the Channel base on its user_id
func (m *Manager) FindChannelsByChatID(chatID int64) ([]Channel, error) {
	var res []Channel
	query := "SELECT channels.* FROM channels INNER JOIN telegram_users ON telegram_users.user_id=channels.user_id WHERE telegram_users.bot_chat_id=? AND channels.admin_status=?"
	_, err := m.GetDbMap().Select(
		&res,
		query,
		chatID,
		AdminStatusAccepted,
	)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// FindChannelByUserIDChannelID return the Channel base on its user_id
func (m *Manager) FindChannelByUserIDChannelID(userID int64, channelID int64) (*Channel, error) {
	var res Channel
	err := m.GetDbMap().SelectOne(
		&res,
		fmt.Sprintf("SELECT * FROM %s WHERE user_id = ? AND id = ?", ChannelTableFull),
		userID,
		channelID,
	)

	if err != nil {
		return nil, err
	}

	return &res, nil
}

// FindChannelsByChatIDName return the Channel base on its chatId and name
func (m *Manager) FindChannelsByChatIDName(chatID int64, name string) (*Channel, error) {
	var res Channel
	query := "SELECT channels.* FROM channels" +
		" INNER JOIN telegram_users ON telegram_users.user_id=channels.user_id" +
		" WHERE telegram_users.bot_chat_id=?" +
		" AND channels.name=?" +
		" AND telegram_users.remove=?" +
		" AND telegram_users.resolve=?" +
		" AND channels.admin_status=?" +
		" AND channels.archive_status=?" +
		" AND channels.active=?"
	err := m.GetDbMap().SelectOne(
		&res,
		query,
		chatID,
		name,
		tlu.RemoveStatusNo,
		tlu.ResolveStatusYes,
		AdminStatusAccepted,
		ArchiveStatusNo,
		ActiveStatusYes,
	)

	if err != nil {
		return nil, err
	}

	return &res, nil
}

//ChannelDataTable is the role full data in data table, after join with other field
// @DataTable {
//		url = /list
//		entity = channel
//		view = channel_list:self
//		controller = modules/telegram/ad/chanControllers
//		fill = FillChannelDataTableArray
//		_edit = edit_channel:self
//		_admin_status = status_channel:parent
//		_archive_status = archive_channel:self
//		_active = active_channel:global
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
	fmt.Println(countQuery)
	count, err := m.GetDbMap().SelectInt(countQuery, params...)
	assert.Nil(err)

	_, err = m.GetDbMap().Select(&res, query, params...)
	assert.Nil(err)
	return res, count
}

// EditChannel function for channel editing
func (m *Manager) EditChannel(link, name string, status AdminStatus, activeStatus ActiveStatus, userID int64, createdAt time.Time, id int64) *Channel {

	ch := &Channel{
		ID:          id,
		UserID:      userID,
		Link:        common.NullString{Valid: link != "", String: link},
		AdminStatus: status,
		Active:      activeStatus,
		Name:        name,
		CreatedAt:   createdAt,
	}
	assert.Nil(m.UpdateChannel(ch))
	return ch
}

// ChangeActive toggle between active status
func (m *Manager) ChangeActive(ID int64, userID int64, name string, link string, status AdminStatus, currentStat ActiveStatus, createdAt time.Time) *Channel {
	ch := &Channel{
		ID:          ID,
		UserID:      userID,
		Name:        name,
		Link:        common.NullString{Valid: link != "", String: link},
		AdminStatus: status,
		CreatedAt:   createdAt,
	}
	if currentStat == ActiveStatusYes {
		ch.Active = ActiveStatusNo
	} else {
		ch.Active = ActiveStatusYes
	}
	assert.Nil(m.UpdateChannel(ch))
	return ch
}
