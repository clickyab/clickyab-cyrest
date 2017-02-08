package ads

import (
	"common/assert"
	"common/models/common"
	"errors"
	"fmt"
	"modules/misc/base"
	"modules/user/aaa"
	"strings"

	"common/utils"
	"database/sql/driver"
	"time"
)

//'pending', 'rejected','accepted','yes','no','yes','no'
const (
	AdAdminStatusPending  AdAdminStatus = "pending"
	AdAdminStatusRejected AdAdminStatus = "rejected"
	AdAdminStatusAccepted AdAdminStatus = "accepted"

	AdArchiveStatusYes AdArchiveStatus = "yes"
	AdArchiveStatusNo  AdArchiveStatus = "no"

	AdPayStatusYes AdPayStatus = "yes"
	AdPayStatusNo  AdPayStatus = "no"

	AdActiveStatusYes AdActiveStatus = "yes"
	AdActiveStatusNo  AdActiveStatus = "no"

	AdTypeIndividual AdType = "individual"
	AdTypePromotion  AdType = "promotion"
)

type (
	// AdAdminStatus is the ad admin status
	// @Enum{
	// }
	AdAdminStatus string

	// AdArchiveStatus is the ad archive status
	// @Enum{
	// }
	AdArchiveStatus string

	// AdPayStatus is the ad pay status
	// @Enum{
	// }
	AdPayStatus string

	// AdActiveStatus is the ad active status
	// @Enum{
	// }
	AdActiveStatus string
	// AdType type ads
	AdType string
)

// Ad model
// @Model {
//		table = ads
//		primary = true, id
//		find_by = id,user_id
//		list = yes
// }
type Ad struct {
	ID              int64             `db:"id" json:"id" sort:"true" title:"ID"`
	UserID          int64             `json:"user_id" db:"user_id" title:"UserID"`
	PlanID          common.NullInt64  `json:"plan_id" db:"plan_id" title:"PlanID"`
	Position        common.NullInt64  `json:"position" db:"position" visible:"false" title:"Position"`
	Name            string            `json:"name" db:"name" title:"Name"`
	Description     common.NullString `json:"description" db:"description" title:"Description"`
	Src             common.NullString `json:"src" db:"src" title:"Src"`
	Mime            common.NullString `json:"mime,ommitempty" db:"-" visible:"false" title:"Mime"`
	CliMessageID    common.NullString `json:"cli_message_id" db:"cli_message_id" visible:"false" title:"CliMessageID"`
	BotChatID       common.NullInt64  `json:"bot_chat_id" db:"bot_chat_id" visible:"false" title:"BotChatID"`
	BotMessageID    common.NullInt64  `json:"bot_message_id" db:"bot_message_id" visible:"false" title:"BotMessageID"`
	View            common.NullInt64  `json:"view" db:"view" visible:"true" title:"View"`
	AdAdminStatus   AdAdminStatus     `json:"admin_status" db:"admin_status" filter:"true" title:"AdminStatus"`
	AdArchiveStatus AdArchiveStatus   `json:"archive_status" db:"archive_status" filter:"true" title:"ArchiveStatus"`
	AdPayStatus     AdPayStatus       `json:"pay_status" db:"pay_status" filter:"true" title:"PayStatus"`
	AdActiveStatus  AdActiveStatus    `json:"active_status" db:"active_status" filter:"true" title:"ActiveStatus"`
	CreatedAt       time.Time         `db:"created_at" json:"created_at" sort:"true" title:"Created at"`
	UpdatedAt       time.Time         `db:"updated_at" json:"updated_at" sort:"true" title:"Updated at"`
}

//AdDataTable is the ad full data in data table, after join with other field
// @DataTable {
//		url = /list
//		entity = ad
//		view = ad_list:self
//		controller = modules/telegram/ad/adControllers
//		fill = FillAdDataTableArray
//		_edit = ad_edit:self
//		_archive_status = change_archive_ad:self
//		_active_status = change_active_ad:global
//		_admin_status = change_admin_ad:parent
// }
type AdDataTable struct {
	Ad
	Email    string `db:"email" json:"email" search:"true" title:"Email"`
	ParentID int64  `db:"-" json:"parent_id" visible:"false"`
	OwnerID  int64  `db:"-" json:"owner_id" visible:"false"`
	Actions  string `db:"-" json:"_actions" visible:"false"`
}

// FillAdDataTableArray is the function to handle
func (m *Manager) FillAdDataTableArray(
	u base.PermInterfaceComplete,
	filters map[string]string,
	search map[string]string,
	sort, order string, p, c int) (AdDataTableArray, int64) {
	var params []interface{}
	var res AdDataTableArray
	var where []string

	countQuery := fmt.Sprintf("SELECT COUNT(ads.id) FROM %s LEFT JOIN %s ON %s.id=%s.user_id", AdTableFull, aaa.UserTableFull, aaa.UserTableFull, AdTableFull)
	query := fmt.Sprintf("SELECT ads.*,users.email FROM %s LEFT JOIN %s ON %s.id=%s.user_id", AdTableFull, aaa.UserTableFull, aaa.UserTableFull, AdTableFull)
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
		where = append(where, fmt.Sprintf("%s.user_id=?", AdTableFull))
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

// LoadNextAd return the next ad in the system
func (m *Manager) LoadNextAd(last int64) (*Ad, error) {
	q := fmt.Sprintf("SELECT * FROM %s WHERE pay_status = ? AND admin_status = ? AND active_status = ? AND id > ? LIMIT 1", AdTableFull)
	res := Ad{}
	err := m.GetDbMap().SelectOne(&res, q, AdPayStatusYes, AdAdminStatusPending, AdActiveStatusYes, last)
	if err != nil && last > 0 {
		return m.LoadNextAd(0)
	}

	return &res, err
}

//ActiveAd selected ad
type ActiveAd struct {
	Ad     Ad
	Viewed int64 `db:"viewed" json:"viewed"`
}

// SelectIndividualActiveAd return the next ad in the system
func (m *Manager) SelectIndividualActiveAd() ([]ActiveAd, error) {
	q := fmt.Sprintf("SELECT %[1]s.*,SUM(%[2]s.view) as viewed FROM %[2]s "+
		" LEFT JOIN %[1]s on %[1]s.id = %[2]s.ad_id "+
		" WHERE %[1]s.cli_message_id = NULL "+
		" AND %[1]s.admin_status = ? "+
		" AND %[1]s.active_status = ? "+
		" AND %[1]s.pay_status = ? "+
		" GROUP BY %[1]s.id",
		AdTableFull,
		ChannelAdTableFull,
	)
	res := []ActiveAd{}
	_, err := m.GetDbMap().Select(&res, q, AdAdminStatusAccepted, AdActiveStatusYes, AdPayStatusYes)

	return res, err
}

// SelectAdsPlan return the next ad in the system
func (m *Manager) SelectAdsPlan() ([]ActiveAd, error) {
	q := fmt.Sprintf("SELECT %[1]s.*,%[2]s.view as viewed FROM %[2]s "+
		" LEFT JOIN %[1]s on %[1]s.id = %[2]s.plan_id "+
		" WHERE %[1]s.admin_status = ? "+
		" AND %[1]s.active_status = ? "+
		" AND %[1]s.pay_status = ? ",
		AdTableFull,
		PlanTableFull,
	)
	res := []ActiveAd{}
	_, err := m.GetDbMap().Select(&res, q, AdAdminStatusAccepted, AdActiveStatusYes, AdPayStatusYes)

	return res, err
}

//UserAdDataTable is the ad full data in data table, after join with other field
// @DataTable {
//		url = /user-ad
//		entity = userAd
//		view = user_ad_list:parent
//		controller = modules/telegram/ad/adControllers
//		fill = FillUserAdDataTableArray
//		_edit = user_ad_edit:parent
//		_change = user_ad_manage:global
// }
type UserAdDataTable struct {
	Ad
	Email    string `db:"email" json:"email" search:"true" title:"Email"`
	ParentID int64  `db:"-" json:"parent_id" visible:"false"`
	OwnerID  int64  `db:"-" json:"owner_id" visible:"false"`
	Actions  string `db:"-" json:"_actions" visible:"false"`
}

// FillUserAdDataTableArray is the function to handle
func (m *Manager) FillUserAdDataTableArray(
	u base.PermInterfaceComplete,
	filters map[string]string,
	search map[string]string,
	sort, order string,
	p, c int) (UserAdDataTableArray, int64) {
	var params []interface{}
	var res UserAdDataTableArray
	var where []string

	currentUserID := u.GetID()
	highestScope := u.GetCurrentScope()

	countQuery := fmt.Sprintf("SELECT COUNT(ads.id) FROM %s LEFT JOIN %s ON %s.id=%s.user_id where (%s.id=%d OR %s.parent_id=%d ) ",
		AdTableFull,
		aaa.UserTableFull,
		aaa.UserTableFull,
		AdTableFull,
		aaa.UserTableFull,
		currentUserID,
		aaa.UserTableFull,
		currentUserID)

	query := fmt.Sprintf("SELECT ads.*,users.email FROM %s LEFT JOIN %s ON %s.id=%s.user_id where (%s.id=%d OR %s.parent_id=%d ) ",
		AdTableFull,
		aaa.UserTableFull,
		aaa.UserTableFull,
		AdTableFull,
		aaa.UserTableFull,
		currentUserID,
		aaa.UserTableFull,
		currentUserID)
	for field, value := range filters {
		where = append(where, fmt.Sprintf("%s.%s=?", AdTableFull, field))
		params = append(params, value)
	}

	for column, val := range search {
		where = append(where, fmt.Sprintf("%s LIKE ?", column))
		params = append(params, "%"+val+"%")
	}

	if highestScope == base.ScopeSelf {
		where = append(where, fmt.Sprintf("%s.user_id=?", AdTableFull))
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

//ReportAdDataTable is the ad full data in data table, after join with other field
// @DataTable {
//		url = /dashboard
//		entity = adReport
//		view = report_ad:self
//		controller = modules/telegram/ad/adControllers
//		fill = FillAdReportDataTableArray
// }
type ReportAdDataTable struct {
	Name     string           `json:"name" db:"name" title:"Name"`
	Type     AdType           `json:"type" db:"type" title:"Type"`
	Start    common.NullTime  `db:"start" json:"start"`
	End      common.NullTime  `db:"end" json:"end" sort:"true" title:"End"`
	PlanView int64            `db:"plan_view" json:"plan_view" sort:"true" title:"PlanView"`
	View     common.NullInt64 `json:"view" db:"view" visible:"true" title:"View"`
	ParentID int64            `db:"-" json:"parent_id" visible:"false"`
	OwnerID  int64            `db:"-" json:"owner_id" visible:"false"`
	Actions  string           `db:"-" json:"_actions" visible:"false"`
}

// IsValid try to validate enum value on ths type
func (e AdType) IsValid() bool {
	return utils.StringInArray(
		string(e),
		string(AdTypeIndividual),
		string(AdTypePromotion),
	)
}

// Scan convert the json array ino string slice
func (e *AdType) Scan(src interface{}) error {
	var b []byte
	switch src.(type) {
	case []byte:
		b = src.([]byte)
	case string:
		b = []byte(src.(string))
	case nil:
		b = make([]byte, 0)
	default:
		return errors.New("unsupported type")
	}
	if string(b) == "0" {
		*e = AdTypeIndividual
		return nil
	}
	if string(b) == "1" {
		*e = AdTypePromotion
		return nil
	}

	return fmt.Errorf("the resualt false %s", string(b))
}

// Value try to get the string slice representation in database
func (e AdType) Value() (driver.Value, error) {
	if !e.IsValid() {
		return nil, errors.New("invaid status")
	}
	return string(e), nil
}

// FillAdReportDataTableArray is the function to handle
func (m *Manager) FillAdReportDataTableArray(
	u base.PermInterfaceComplete,
	filters map[string]string,
	search map[string]string,
	sort, order string, p, c int) (ReportAdDataTableArray, int64) {
	var params []interface{}
	var res ReportAdDataTableArray
	var where []string

	countQuery := fmt.Sprintf("SELECT COUNT(%[2]s.id) FROM %[2]s "+
		"LEFT JOIN %[1]s ON %[2]s.id=%[1]s.ad_id "+
		"LEFT JOIN %[3]s ON %[3]s.id=%[2]s.user_id "+
		"LEFT JOIN %[4]s ON %[4]s.id=%[2]s.plan_id "+
		"GROUP BY %[2]s.id",
		ChannelAdTableFull,
		AdTableFull,
		aaa.UserTableFull,
		PlanTableFull,
	)
	query := fmt.Sprintf("SELECT %[1]s.name as name, %[1]s.cli_message_id IS NULL AS type, %[2]s.start as start, %[2]s.end as end, %[2]s.view as view, %[4]s.view as plan_view FROM %[1]s "+
		"LEFT JOIN %[2]s ON %[1]s.id=%[2]s.ad_id "+
		"LEFT JOIN %[3]s ON %[3]s.id=%[1]s.user_id "+
		"LEFT JOIN %[4]s ON %[4]s.id=%[1]s.plan_id ",
		AdTableFull,
		ChannelAdTableFull,
		aaa.UserTableFull,
		PlanTableFull,
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
		where = append(where, fmt.Sprintf("%s.user_id=?", AdTableFull))
		params = append(params, currentUserID)
	} else if highestScope == base.ScopeParent {
		where = append(where, fmt.Sprintf("%s.parent_id=?", aaa.UserTableFull))
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
	query += fmt.Sprintf(" GROUP BY %s.id ", AdTableFull)
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
