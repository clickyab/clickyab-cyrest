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
	ID            int64             `db:"id" json:"id,omitempty" sort:"true" title:"ID" perm:"view_channel:global"`
	UserID        int64             `json:"user_id,omitempty" db:"user_id" title:"UserID" perm:"view_channel:global"`
	Name          string            `json:"name" db:"name" search:"true" title:"Name"`
	Title         common.NullString `json:"link" db:"link" search:"true" title:"Title"`
	AdminStatus   AdminStatus       `json:"admin_status" db:"admin_status" filter:"true" title:"AdminStatus"`
	ArchiveStatus ArchiveStatus     `json:"archive_status" db:"archive_status" filter:"true" title:"ArchiveStatus"`
	Active        ActiveStatus      `json:"active" db:"active" filter:"true" title:"Active"`
	CreatedAt     *time.Time        `db:"created_at" json:"created_at,omitempty" sort:"true" title:"Created at" perm:"view_channel:global"`
	UpdatedAt     *time.Time        `db:"updated_at" json:"updated_at,omitempty" sort:"true" title:"Updated at" perm:"view_channel:global"`
}

// ChanStat returns channels and their status by provider
type ChanStat struct {
	Stat  AdminStatus `json:"status" db:"status"`
	Count int         `json:"count" db:"count"`
}

// ChannelCreate a new channel
func (m *Manager) ChannelCreate(link, name string, status AdminStatus, archive ArchiveStatus, active ActiveStatus, userID int64) *Channel {

	ch := &Channel{
		Title:       common.MakeNullString(link),
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
	query := fmt.Sprintf("SELECT %[1]s.* FROM %[1]s"+
		" INNER JOIN %[2]s ON %[2]s.user_id=%[1]s.user_id"+
		" WHERE %[2]s.bot_chat_id=?"+
		" AND %[1]s.name=?"+
		" AND %[2]s.remove=?"+
		" AND %[2]s.resolve=?"+
		" AND %[1]s.admin_status=?"+
		" AND %[1]s.archive_status=?"+
		" AND %[1]s.active=?",
		ChannelTableFull,
		tlu.TeleUserTableFull,
	)
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
	Email    string           `db:"email" json:"email,omitempty" search:"true" title:"Email" perm:"view_channel:global"`
	ParentID common.NullInt64 `db:"parent_id" json:"parent_id" visible:"false"`
	OwnerID  int64            `db:"owner_id" json:"owner_id" visible:"false"`
	Actions  string           `db:"-" json:"_actions" visible:"false"`
}

// FillChannelDataTableArray is the function to handle
func (m *Manager) FillChannelDataTableArray(
	u base.PermInterfaceComplete,
	filters map[string]string,
	search map[string]string,
	contextparams map[string]string,
	sort, order string,
	p, c int) (ChannelDataTableArray, int64) {
	var params []interface{}
	var res ChannelDataTableArray
	var where []string

	countQuery := fmt.Sprintf("SELECT COUNT(%[1]s.id) FROM %[1]s LEFT JOIN %[2]s ON %[2]s.id=%[1]s.user_id ",
		ChannelTableFull,
		aaa.UserTableFull)
	query := fmt.Sprintf("SELECT %[1]s.*,%[2]s.email,%[2]s.id AS owner_id, %[2]s.parent_id as parent_id FROM %[1]s LEFT JOIN %[2]s ON %[2]s.id=%[1]s.user_id ",
		ChannelTableFull,
		aaa.UserTableFull)
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
		where = append(where, "%s.parent_id=?", aaa.UserTableFull)
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
//		url = /detail/:id
//		entity = ChannelDetailDataTable
//		view = channel_list:self
//		controller = modules/telegram/ad/chanControllers
//		fill = FillChannelDetailDataTableArray
//		_edit = edit_channel:self
//		_active = active_channel:global
// }
type ChannelDetailDataTable struct {
	View     int64            `db:"view" json:"view"`
	AdName   string           `db:"name" json:"adname" search:"true"`
	Active   ActiveStatus     `db:"active" json:"active" filter:"true"`
	Start    common.NullTime  `db:"start" json:"start"`
	End      common.NullTime  `db:"end" json:"end"`
	Warning  int64            `db:"warning" json:"warning"`
	ParentID common.NullInt64 `db:"-" json:"parent_id" visible:"false"`
	OwnerID  int64            `db:"-" json:"owner_id" visible:"false"`
	Actions  string           `db:"-" json:"_actions" visible:"false"`
}

// FillChannelDetailDataTableArray is the function to handle
func (m *Manager) FillChannelDetailDataTableArray(u base.PermInterfaceComplete,
	filters map[string]string,
	search map[string]string,
	contextparams map[string]string,
	sort,
	order string,
	p, c int) (ChannelDetailDataTableArray, int64) {
	var params []interface{}
	var res ChannelDetailDataTableArray
	var where []string
	id := contextparams["id"]

	countQuery := fmt.Sprintf("SELECT COUNT(channel_id) FROM %[1]s "+
		"LEFT JOIN %[2]s ON %[1]s.ad_id=%[2]s.id "+
		"LEFT JOIN %[3]s ON %[2]s.user_id = %[3]s.id "+
		"WHERE channel_id=?",
		ChannelAdTableFull,
		AdTableFull,
		aaa.UserTableFull)

	query := fmt.Sprintf("SELECT ads.name, channel_ad.view, active, start, end FROM %[1]s "+
		"LEFT JOIN %[2]s ON %[2]s.id=%[1]s.ad_id "+
		"LEFT JOIN %[3]s ON %[2]s.user_id = %[3]s.id "+
		"WHERE channel_id=?",
		ChannelAdTableFull,
		AdTableFull,
		aaa.UserTableFull)
	params = append(params, id)

	for field, value := range filters {
		where = append(where, fmt.Sprintf(ChannelTableFull+".%s=?", field))
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
		where = append(where, fmt.Sprintf("%[1]s.parent_id=? OR %[1]s.user_id=?", aaa.UserTableFull))
		params = append(params, currentUserID, currentUserID)
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

//ChannelSpecificDataTable is the role full data in data table, after join with other field
// @DataTable {
//		url = /specific/:id
//		entity = ChannelSpecificDataTable
//		view = channel_list:self
//		controller = modules/telegram/ad/chanControllers
//		fill = FillChannelSpecificDataTableArray
//		_edit = edit_channel:self
// }
type ChannelSpecificDataTable struct {
	AdName   string           `db:"name" json:"adname" search:"true"`
	Start    common.NullTime  `db:"start" json:"start"`
	End      common.NullTime  `db:"end" json:"end"`
	Price    int64            `db:"-" json:"price"`
	ParentID common.NullInt64 `db:"-" json:"parent_id" visible:"false"`
	OwnerID  int64            `db:"-" json:"owner_id" visible:"false"`
	Actions  string           `db:"-" json:"_actions" visible:"false"`
}

// FillChannelSpecificDataTableArray is the function to handle
func (m *Manager) FillChannelSpecificDataTableArray(u base.PermInterfaceComplete,
	filters map[string]string,
	search map[string]string,
	contextparams map[string]string,
	sort,
	order string,
	p, c int) (ChannelSpecificDataTableArray, int64) {
	var params []interface{}
	var res ChannelSpecificDataTableArray
	var where []string
	id := contextparams["id"]

	countQuery := fmt.Sprintf("SELECT COUNT(channel_id) FROM %[1]s "+
		"LEFT JOIN %[2]s ON %[1]s.ad_id=%[2]s.id "+
		"LEFT JOIN %[3]s ON %[2]s.user_id = %[3]s.id "+
		"WHERE channel_id=?",
		ChannelAdTableFull,
		AdTableFull,
		aaa.UserTableFull)

	query := fmt.Sprintf("SELECT ads.name, start, end FROM %[1]s "+
		"LEFT JOIN %[2]s ON %[2]s.id=%[1]s.ad_id "+
		"LEFT JOIN %[3]s ON %[2]s.user_id = %[3]s.id "+
		"WHERE channel_id=?",
		ChannelAdTableFull,
		AdTableFull,
		aaa.UserTableFull)
	params = append(params, id)

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
		where = append(where, fmt.Sprintf("%[1]s.parent_id=? OR %[1]s.user_id=?", aaa.UserTableFull))
		params = append(params, currentUserID, currentUserID)
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
	// TODO price algorithm
	return res, count
}

// EditChannel function for channel editing
func (m *Manager) EditChannel(link, name string, status AdminStatus, activeStatus ActiveStatus, userID int64, createdAt time.Time, id int64) *Channel {

	ch := &Channel{
		ID:          id,
		UserID:      userID,
		Title:       common.NullString{Valid: link != "", String: link},
		AdminStatus: status,
		Active:      activeStatus,
		Name:        name,
		CreatedAt:   &createdAt,
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
		Title:       common.NullString{Valid: link != "", String: link},
		AdminStatus: status,
		CreatedAt:   &createdAt,
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
	Type        AdType           `json:"type" db:"type" title:"Type"`
	AdName      string           `json:"ad_name" db:"ad_name" title:"Ad Name"`
	ChannelName string           `json:"channel_name" db:"channel_name" title:"Channel Name"`
	Active      ActiveStatus     `db:"active" json:"active" title:"Active" filter:"true"`
	Start       common.NullTime  `db:"start" json:"start" title:"Start" sort:"true"`
	End         common.NullTime  `db:"end" json:"end" title:"End" sort:"true"`
	View        int64            `json:"view" db:"view" visible:"true" title:"View" sort:"true"`
	ParentID    common.NullInt64 `db:"parent_id" json:"parent_id" visible:"false"`
	OwnerID     int64            `db:"owner_id" json:"owner_id" visible:"false"`
	Actions     string           `db:"-" json:"_actions" visible:"false"`
}

// FillActiveAdReportDataTableArray is the function to handle
func (m *Manager) FillActiveAdReportDataTableArray(
	u base.PermInterfaceComplete,
	filters map[string]string,
	search map[string]string,
	contextparams map[string]string,
	sort, order string, p, c int) (ReportActiveAdDataTableArray, int64) {
	var params []interface{}
	var res ReportActiveAdDataTableArray
	var where []string

	countQuery := fmt.Sprintf("SELECT count(%[1]s.ad_id)  FROM %[1]s "+
		"LEFT JOIN %[2]s ON %[1]s.channel_id = %[2]s.id "+
		"LEFT JOIN  %[3]s ON %[1]s.ad_id = %[3]s.id "+
		"GROUP BY %[1]s.channel_id"+
		"LEFT JOIN  %[4]s ON %[4]s.id = %[2]s.user_id ",
		ChannelAdTableFull,
		ChannelTableFull,
		AdTableFull,
		aaa.UserTableFull,
	)
	query := fmt.Sprintf("SELECT %[3]s.cli_message_id is NULL as type,"+
		"%[3]s.name as ad_name, %[2]s.name as channel_name,%[1]s.active as active,"+
		"%[1]s.start as start,%[1]s.end as end,%[1]s.view as view"+
		", %[4]s.id AS owner_id,  %[4]s.parent_id as parent_id "+
		"FROM %[1]s "+
		"LEFT JOIN %[2]s ON %[1]s.channel_id = %[2]s.id "+
		"LEFT JOIN  %[3]s ON %[1]s.ad_id = %[3]s.id "+
		"LEFT JOIN  %[4]s ON %[4]s.id = %[2]s.user_id ",
		ChannelAdTableFull,
		ChannelTableFull,
		AdTableFull,
		aaa.UserTableFull,
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

// GetChanStat returns channels and admin statuss'
func (m *Manager) GetChanStat(userID int64, scope base.UserScope) (result []ChanStat) {
	var params []interface{}
	where := ""
	if scope == base.ScopeSelf {
		where = "WHERE user_id=? "
		params = append(params, userID)
	} else if scope == base.ScopeParent {
		where = "WHERE user_id=? OR parent_id=? "
		params = append(params, userID, userID)
	}

	q := fmt.Sprintf("SELECT COUNT(%[1]s.id) AS count, admin_status AS status from %[1]s "+
		"LEFT JOIN %[2]s ON %[1]s.user_id = %[2]s.id "+
		"%[3]s"+
		" GROUP BY admin_status",
		ChannelTableFull,
		aaa.UserTableFull,
		where)

	_, err := m.GetDbMap().Select(
		&result,
		q,
		params...,
	)
	assert.Nil(err)

	return
}

// CountActiveChannel dashboard query count active and wait channel
func (m *Manager) CountActiveChannel(userID int64, scope base.UserScope) (int64, int64) {
	var where string
	var params1 []interface{}
	var params2 []interface{}
	switch scope {
	case base.ScopeGlobal:
		where = ""
		params1 = []interface{}{ActiveStatusYes}
		params2 = []interface{}{ActiveStatusNo}
	case base.ScopeParent:
		where = fmt.Sprintf(" AND ( %[1]s.id = ? OR %[1]s.parent_id = ? )", aaa.UserTableFull)
		params1 = []interface{}{ActiveStatusYes, userID, userID}
		params2 = []interface{}{ActiveStatusNo, userID, userID}
	case base.ScopeSelf:
		where = fmt.Sprintf(" AND %[1]s.id = ?", aaa.UserTableFull)
		params1 = []interface{}{ActiveStatusYes, userID}
		params2 = []interface{}{ActiveStatusNo, userID}
	}
	q1 := fmt.Sprintf(`SELECT COUNT(%[1]s.channel_id) as count_active
				   FROM %[1]s
				   INNER JOIN %[2]s ON %[1]s.channel_id = %[2]s.id
				   INNER JOIN %[3]s ON %[2]s.user_id = %[3]s.id
				   WHERE %[1]s.active = ?
				   AND %[1]s.start IS NOT NULL
				   AND %[1]s.end IS NULL
				   %[4]s`,
		ChannelAdTableFull,
		ChannelTableFull,
		aaa.UserTableFull,
		where,
	)
	active, err := m.GetDbMap().SelectInt(q1, params1...)
	assert.Nil(err)

	q2 := fmt.Sprintf(`SELECT COUNT(%[1]s.channel_id) as count_wait
				   FROM %[1]s
				   INNER JOIN %[2]s ON %[1]s.channel_id = %[2]s.id
				   INNER JOIN %[3]s ON %[2]s.user_id = %[3]s.id
				   WHERE %[1]s.active = ?
				   AND %[1]s.start IS NULL
				   AND %[1]s.end IS NULL
				   %[4]s`,
		ChannelAdTableFull,
		ChannelTableFull,
		aaa.UserTableFull,
		where,
	)

	wait, err := m.GetDbMap().SelectInt(q2, params2...)
	assert.Nil(err)
	return active, wait
}
