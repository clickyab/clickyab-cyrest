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

// DailyView struct for dashboard
type DailyView struct {
	View int64 `db:"view" json:"view"`
	End  time.Time
}

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

// FindChannelsByUserID return the Channels owned by a user
func (m *Manager) FindChannelsByUserID(userID int64) ([]Channel, error) {
	channels := []Channel{}
	_, err := m.GetDbMap().Select(
		&channels,
		fmt.Sprintf("SELECT * FROM %s WHERE user_id= ? ", ChannelTableFull),
		userID,
	)
	if err != nil {
		return nil, err
	}
	return channels, nil
}

// GetChanViewByID returns Channels View by day
func (m *Manager) GetChanViewByID(chanID int64) ([]DailyView, error) {
	q := `SELECT sum(view) AS view, end FROM channel_ad
		where channel_id=?
		GROUP BY now() > end`

	var res []DailyView
	_, err := m.GetDbMap().Select(&res, q, chanID)
	if err != nil {
		return nil, err
	}

	return res, nil
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

//ChannelDetailDataTable is the role full data in data table, after join with other field
// @DataTable {
//		url = /detail
//		entity = channel
//		view = channel_list:self
//		controller = modules/telegram/ad/chanControllers
//		fill = FillChannelDetailDataTableArray
//		_edit = edit_channel:self
//		_admin_status = status_channel:parent
//		_archive_status = archive_channel:self
//		_active = active_channel:global
// }
type ChannelDetailDataTable struct {
	ChanID   int64           `json:"name" search:"true" db:"name" visible:"false" title:"Name"`
	View     int64           `db:"view" json:"view"`
	AdName   string          `db:"warning" json:"warning" search:"true"`
	Active   ActiveStatus    `db:"active" json:"active" filter:"true"`
	Start    common.NullTime `db:"start" json:"start"`
	End      common.NullTime `db:"end" json:"end"`
	Warning  int64           `db:"warning" json:"warning"`
	ParentID int64           `db:"-" json:"parent_id" visible:"false"`
	OwnerID  int64           `db:"-" json:"owner_id" visible:"false"`
	Actions  string          `db:"-" json:"_actions" visible:"false"`
}

// FillChannelDetailDataTableArray is the function to handle
func (m *Manager) FillChannelDetailDataTableArray(u base.PermInterfaceComplete,
	filters map[string]string,
	search map[string]string,
	sort,
	order string,
	p, c int) (ChannelDetailDataTableArray, int64) {
	var params []interface{}
	var res ChannelDetailDataTableArray
	var where []string

	countQuery := fmt.Sprintf("SELECT COUNT(channel_ad.id) FROM %[1]s LEFT JOIN %[2]s ON %[1]s.ad_id=%[2]s.id", ChannelAdTableFull, AdTableFull)
	query := fmt.Sprintf("SELECT ads.name, channel_ad.view, active, start, end FROM %[1]s LEFT JOIN %[2]s ON %[1]s.id=%[2]s.ad_id", ChannelAdTableFull, AdTableFull)

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

//ReportActiveAdDataTable is the ad full data in data table, after join with other field
// @DataTable {
//		url = /dashboard/active_ad
//		entity = activeAd
//		view = active_ad:self
//		controller = modules/telegram/ad/chanControllers
//		fill = FillActiveAdReportDataTableArray
// }
type ReportActiveAdDataTable struct {
	Type        AdType          `json:"type" db:"type" title:"Type"`
	AdName      string          `json:"ad_name" db:"ad_name" title:"Ad Name"`
	ChannelName string          `json:"channel_name" db:"channel_name" title:"Channel Name"`
	Active      ActiveStatus    `db:"active" json:"active" title:"Active" filter:"true"`
	Start       common.NullTime `db:"start" json:"start" title:"Start" sort:"true"`
	End         common.NullTime `db:"end" json:"end" title:"End" sort:"true"`
	View        int64           `json:"view" db:"view" visible:"true" title:"View" sort:"true"`
	ParentID    int64           `db:"-" json:"parent_id" visible:"false"`
	OwnerID     int64           `db:"-" json:"owner_id" visible:"false"`
	Actions     string          `db:"-" json:"_actions" visible:"false"`
}

// FillActiveAdReportDataTableArray is the function to handle
func (m *Manager) FillActiveAdReportDataTableArray(
	u base.PermInterfaceComplete,
	filters map[string]string,
	search map[string]string,
	sort, order string, p, c int) (ReportActiveAdDataTableArray, int64) {
	var params []interface{}
	var res ReportActiveAdDataTableArray
	var where []string

	countQuery := fmt.Sprintf("SELECT count(%[1]s.ad_id)  FROM %[1]s "+
		"LEFT JOIN %[2]s ON %[1]s.channel_id = %[2]s.id "+
		"LEFT JOIN  %[3]s ON %[1]s.ad_id = %[3]s.id "+
		"GROUP BY %[1]s.channel_id",
		ChannelAdTableFull,
		ChannelTableFull,
		AdTableFull,
	)
	query := fmt.Sprintf("SELECT %[3]s.cli_message_id is NULL as type,"+
		"%[3]s.name as ad_name, %[2]s.name as channel_name,%[1]s.active as active,"+
		"%[1]s.start as start,%[1]s.end as end,%[1]s.view as view  FROM %[1]s "+
		"LEFT JOIN %[2]s ON %[1]s.channel_id = %[2]s.id "+
		"LEFT JOIN  %[3]s ON %[1]s.ad_id = %[3]s.id ",
		ChannelAdTableFull,
		ChannelTableFull,
		AdTableFull,
	)
	for field, value := range filters {
		where = append(where, fmt.Sprintf("%s.%s=?", AdTableFull, field))
		params = append(params, value)
	}

	for column, val := range search {
		where = append(where, fmt.Sprintf("%s LIKE ?", column))
		params = append(params, "%"+val+"%")
	}

	currentUserID := u.GetID()
	highestScope := u.GetCurrentScope()

	if highestScope == base.ScopeSelf {
		where = append(where, fmt.Sprintf(" %s.user_id=? ", AdTableFull))
		params = append(params, currentUserID)
	} else if highestScope == base.ScopeParent {
		where = append(where, fmt.Sprintf(" %s.parent_id=? ", aaa.UserTableFull))
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
	query += fmt.Sprintf(" GROUP BY %s.channel_id ", ChannelAdTableFull)
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
