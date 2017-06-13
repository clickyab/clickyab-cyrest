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
	"database/sql"
	"database/sql/driver"
	"strconv"
	"time"
)

//'pending', 'rejected','accepted','yes','no','yes','no'
const (
	AdminStatusPending  AdminStatus = "pending"
	AdminStatusRejected AdminStatus = "rejected"
	AdminStatusAccepted AdminStatus = "accepted"

	ActiveStatusYes ActiveStatus = "yes"
	ActiveStatusNo  ActiveStatus = "no"

	AdTypeIndividual AdType = "individual"
	AdTypePromotion  AdType = "promotion"
)

type (
	// AdminStatus is the ad admin status
	// @Enum{
	// }
	AdminStatus string

	// ActiveStatus is the ad active status
	// @Enum{
	// }
	ActiveStatus string
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
	UserID          int64             `json:"user_id" db:"user_id" title:"UserID" visible:"false"`
	PlanID          common.NullInt64  `json:"plan_id" db:"plan_id" title:"PlanID" visible:"false"`
	Position        common.NullInt64  `json:"position" db:"position" visible:"false" title:"Position"`
	Name            string            `json:"name" db:"name" title:"Name" search:"true"`
	Description     common.MB4String  `json:"description" db:"description" title:"Description"`
	Src             common.NullString `json:"src" db:"src" title:"Src"`
	PromoSrc        common.NullString `json:"promo_src" db:"promo_src" title:"PromoSrc"`
	Extension       common.NullString `json:"extension,ommitempty" db:"-" visible:"false" title:"Extension"`
	CliMessageID    common.NullString `json:"cli_message_id" db:"cli_message_id" visible:"false" title:"CliMessageID"`
	PromoteData     common.NullString `json:"promote_data" db:"promote_data" visible:"false" title:"PromoteData"`
	BotChatID       common.NullInt64  `json:"bot_chat_id" db:"bot_chat_id" visible:"false" title:"BotChatID"`
	BotMessageID    common.NullInt64  `json:"bot_message_id" db:"bot_message_id" visible:"false" title:"BotMessageID"`
	View            common.NullInt64  `json:"view" db:"view" visible:"true" title:"View"`
	AdAdminStatus   AdminStatus       `json:"admin_status" db:"admin_status" filter:"true" title:"AdminStatus"`
	AdArchiveStatus ActiveStatus      `json:"archive_status" db:"archive_status" filter:"true" title:"AdminStatus"`
	AdPayStatus     ActiveStatus      `json:"pay_status" db:"pay_status" filter:"true" title:"PayStatus"`
	AdActiveStatus  ActiveStatus      `json:"active_status" db:"active_status" filter:"true" title:"ActiveStatus"`
	CreatedAt       *time.Time        `db:"created_at" json:"created_at" sort:"true" title:"Created at"`
	UpdatedAt       *time.Time        `db:"updated_at" json:"updated_at" sort:"true" title:"Updated at" visible:"false"`
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
//		_admin_status = change_admin_ad:parent
// }
type AdDataTable struct {
	Ad
	Type     AdType           `json:"type" db:"type" title:"Type"`
	Email    string           `db:"email" json:"email" search:"true" title:"Email" visible:"false"`
	ParentID common.NullInt64 `db:"parent_id" json:"parent_id" visible:"false"`
	OwnerID  int64            `db:"owner_id" json:"owner_id" visible:"false"`
	Actions  string           `db:"-" json:"_actions" visible:"false"`
}

// FillAdDataTableArray is the function to handle
func (m *Manager) FillAdDataTableArray(
	u base.PermInterfaceComplete,
	filters map[string]string,
	search map[string]string,
	contextparams map[string]string,
	sort, order string, p, c int) (AdDataTableArray, int64) {
	var params []interface{}
	var res AdDataTableArray
	var where []string

	countQuery := fmt.Sprintf("SELECT COUNT(%[1]s.id) FROM %[1]s LEFT JOIN %[2]s ON %[2]s.id=%[1]s.user_id",
		AdTableFull,
		aaa.UserTableFull)
	query := fmt.Sprintf("SELECT %[1]s.*,%[2]s.email,%[2]s.id AS owner_id, %[2]s.parent_id as parent_id,%[1]s.cli_message_id IS NULL AS type FROM %[1]s LEFT JOIN %[2]s ON %[2]s.id=%[1]s.user_id",
		AdTableFull,
		aaa.UserTableFull)
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
	err := m.GetDbMap().SelectOne(&res, q, ActiveStatusYes, AdminStatusPending, ActiveStatusYes, last)
	if err != nil && last > 0 {
		return m.LoadNextAd(0)
	}

	return &res, err
}

//ActiveAd selected ad
type ActiveAd struct {
	Ad
	Viewed int64 `db:"viewed" json:"viewed"`
}

// SelectIndividualActiveAd return the next ad in the system
func (m *Manager) SelectIndividualActiveAd() ([]ActiveAd, error) {
	q := fmt.Sprintf("SELECT %[1]s.*,SUM(%[2]s.view) as viewed FROM %[2]s "+
		" LEFT JOIN %[1]s on %[1]s.id = %[2]s.ad_id "+
		" WHERE %[1]s.cli_message_id IS NULL "+
		" AND %[1]s.admin_status = ? "+
		" AND %[1]s.active_status = ? "+
		" AND %[1]s.pay_status = ? "+
		" GROUP BY %[1]s.id",
		AdTableFull,
		ChannelAdTableFull,
	)
	res := []ActiveAd{}
	_, err := m.GetDbMap().Select(&res, q, AdminStatusAccepted, ActiveStatusYes, ActiveStatusYes)

	return res, err
}

// UpdateIndividualViewCount try to fill the view of individuals base on its sub view
func (m *Manager) UpdateIndividualViewCount() {
	q := `UPDATE %s
		LEFT JOIN (SELECT ad_id, SUM(view) AS dv FROM %s GROUP BY ad_id ) AS chad ON chad.ad_id = ads.id
		SET ads.view = chad.dv
		WHERE ads.cli_message_id IS NULL AND ads.active_status = ? AND ads.admin_status = ? AND ads.pay_status = ? `
	q = fmt.Sprintf(q, AdTableFull, ChannelAdTableFull)
	_, err := m.GetDbMap().Exec(q, ActiveStatusYes, AdminStatusAccepted, ActiveStatusYes)
	assert.Nil(err)
}

// FinishedActiveAds struct finished ad
type FinishedActiveAds struct {
	Ad
	Price int64 `db:"price" json:"price"`
	Share int64 `db:"share" json:"share"`
}

// FinishedActiveAds return all finished ads
func (m *Manager) FinishedActiveAds() []FinishedActiveAds {
	q := `SELECT a.*,p.price,p.share FROM %s AS a LEFT JOIN %s AS p ON a.plan_id = p.id
		WHERE p.view < a.view AND a.admin_status = ?
		AND a.active_status = ? AND a.pay_status = ?`
	q = fmt.Sprintf(q, AdTableFull, PlanTableFull)
	var res []FinishedActiveAds
	_, err := m.GetDbMap().Select(&res, q, AdminStatusAccepted, ActiveStatusYes, ActiveStatusYes)
	assert.Nil(err)
	return res
}

//FinishedActiveChannels struct finished ad
type FinishedActiveChannels struct {
	ChannelAd
	UserID int64 `db:"user_id" json:"user_id"`
}

// FinishedActiveChannels return all finished ChannelAd
func (m *Manager) FinishedActiveChannels(adID int64, warning int64) []FinishedActiveChannels {
	q := `SELECT ca.*,c.user_id FROM %s AS ca
		LEFT JOIN %s AS c ON ca.channel_id = c.id
		WHERE ca.ad_id = ?
		AND ca.warning < ?
		AND ca.view > 0`
	q = fmt.Sprintf(q, ChannelAdTableFull, ChannelTableFull)
	var res []FinishedActiveChannels
	_, err := m.GetDbMap().Select(&res, q, adID, warning)
	assert.Nil(err)
	return res
}

// GetWarningLimited return all passed warnings
func (m *Manager) GetWarningLimited(warning int64) []ChannelAd {
	q := fmt.Sprintf("SELECT * FROM %s WHERE warning > ? AND active=? AND end IS NULL", ChannelAdTableFull)
	var ch []ChannelAd
	_, err := m.GetDbMap().Select(&ch, q, warning, ActiveStatusYes)
	assert.Nil(err)

	return ch
}

// SelectAdsPlan return the next ad in the system
func (m *Manager) SelectAdsPlan() ([]ActiveAd, error) {
	q := fmt.Sprintf("SELECT %[1]s.*,%[2]s.view as viewed FROM %[2]s "+
		" LEFT JOIN %[1]s on %[1]s.plan_id = %[2]s.id "+
		" WHERE %[1]s.admin_status = ? "+
		" AND %[1]s.active_status = ? "+
		" AND %[1]s.pay_status = ? ",
		AdTableFull,
		PlanTableFull,
	)
	res := []ActiveAd{}
	_, err := m.GetDbMap().Select(&res, q, AdminStatusAccepted, ActiveStatusYes, ActiveStatusYes)

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
	Email    string           `db:"email" json:"email" search:"true" title:"Email"`
	ParentID common.NullInt64 `db:"parent_id" json:"parent_id" visible:"false"`
	OwnerID  int64            `db:"owner_id" json:"owner_id" visible:"false"`
	Actions  string           `db:"-" json:"_actions" visible:"false"`
}

// FillUserAdDataTableArray is the function to handle
func (m *Manager) FillUserAdDataTableArray(
	u base.PermInterfaceComplete,
	filters map[string]string,
	search map[string]string,
	contextparams map[string]string,
	sort, order string,
	p, c int) (UserAdDataTableArray, int64) {
	var params []interface{}
	var res UserAdDataTableArray
	var where []string

	currentUserID := u.GetID()
	highestScope := u.GetCurrentScope()

	countQuery := fmt.Sprintf("SELECT COUNT(%[1]s.id) FROM %[1]s LEFT JOIN %[2]s ON %[2]s.id=%[1]s.user_id ",
		AdTableFull,
		aaa.UserTableFull)

	query := fmt.Sprintf("SELECT %[1]s.*,%[2]s.email,%[2]s.id AS owner_id, %[2]s.parent_id as parent_id FROM %[1]s LEFT JOIN %[2]s ON %[2]s.id=%[1]s.user_id ",
		AdTableFull,
		aaa.UserTableFull)
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
	ParentID common.NullInt64 `db:"parent_id" json:"parent_id" visible:"false"`
	OwnerID  int64            `db:"owner_id" json:"owner_id" visible:"false"`
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
	case int64:
		b = []byte(fmt.Sprint(src.(int64)))
	default:
		return fmt.Errorf("unsupported type %T", src)
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
	contextparams map[string]string,
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
	query := fmt.Sprintf("SELECT %[1]s.name as name, %[1]s.cli_message_id IS NULL AS type,"+
		" %[2]s.start as start, %[2]s.end as end, %[2]s.view as view,"+
		" %[4]s.view as plan_view ,%[3]s.id AS owner_id, %[3]s.parent_id as parent_id FROM %[1]s "+
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

// AdReport shows the report for ad
type AdReport struct {
	Name  string          `json:"name"`
	View  int64           `json:"view"`
	End   common.NullTime `json:"end"`
	Price int             `json:"price"`
}

// PieChart struct type
type PieChart struct {
	Name string           `json:"name" db:"name" title:"Name"`
	View common.NullInt64 `json:"view" db:"view" visible:"true" title:"View"`
}

// PieChartAdvertiser return the ads
func (m *Manager) PieChartAdvertiser(userID int64) ([]PieChart, error) {

	res := []PieChart{}
	_, err := m.GetDbMap().Select(
		&res,
		fmt.Sprintf("SELECT %[1]s.name,%[1]s.view "+
			" FROM %[1]s "+
			" LEFT JOIN %[2]s ON %[1]s.user_id = %[2]s.id "+
			" WHERE ( %[2]s.id = ? OR %[2]s.parent_id = ? )",
			AdTableFull,
			aaa.UserTableFull,
		),
		userID,
		userID,
	)

	if err != nil {
		return nil, err
	}

	return res, nil
}

// AdDashboard AdDashboard
type AdDashboard struct {
	AdName    string               `json:"ad_name" db:"name"`
	Viewed    common.ZeroNullInt64 `json:"viewed" db:"viewed"`
	Remaining common.ZeroNullInt64 `json:"remaining" db:"remaining"`
}

// PieChartAd return the ads
func (m *Manager) PieChartAd(userID int64, scope base.UserScope) []AdDashboard {
	var where string
	var params []interface{}
	switch scope {
	case base.ScopeGlobal:
		where = ""
		params = []interface{}{ActiveStatusYes}
	case base.ScopeParent:
		where = " AND ( u.id = ? OR u.parent_id = ? )"
		params = []interface{}{ActiveStatusYes, userID, userID}
	case base.ScopeSelf:
		where = " AND u.id = ?"
		params = []interface{}{ActiveStatusYes, userID}
	}
	res := []AdDashboard{}
	_, err := m.GetDbMap().Select(
		&res,
		fmt.Sprintf("SELECT a.name,a.view AS viewed,(p.view-a.view) AS remaining FROM %s AS a "+
			"INNER JOIN %s AS u ON u.id=a.user_id "+
			"INNER JOIN %s AS p ON p.id=a.plan_id "+
			"WHERE a.pay_status=? %s", AdTableFull, aaa.UserTableFull, PlanTableFull, where),
		params...,
	)
	assert.Nil(err)
	return res
}

// AdDashboardPerChannel AdDashboardPerChannel
type AdDashboardPerChannel struct {
	ChannelName string  `json:"channel_name" db:"channel_name"`
	Spent       float64 `json:"spent" db:"spent"`
}

// PieChartAdPerChannel return the ads
func (m *Manager) PieChartAdPerChannel(userID int64, scope base.UserScope) []AdDashboardPerChannel {
	var where string
	var params []interface{}
	switch scope {
	case base.ScopeGlobal:
		where = ""
		params = []interface{}{ActiveStatusYes}
	case base.ScopeParent:
		where = " AND ( u.id = ? OR u.parent_id = ? )"
		params = []interface{}{ActiveStatusYes, userID, userID}
	case base.ScopeSelf:
		where = " AND u.id = ?"
		params = []interface{}{ActiveStatusYes, userID}
	}

	res := []AdDashboardPerChannel{}
	_, err := m.GetDbMap().Select(
		&res,
		fmt.Sprintf("SELECT SUM(cha.view)*(p.price/p.view) AS spent,"+
			"ch.name AS channel_name "+
			"FROM %s AS cha "+
			"INNER JOIN %s AS a ON a.id=cha.ad_id "+
			"INNER JOIN %s AS u ON u.id=a.user_id "+
			"INNER JOIN %s AS p ON p.id=a.plan_id "+
			"INNER JOIN %s AS ch ON ch.id=cha.channel_id "+
			"WHERE a.pay_status=? %s "+
			"GROUP BY cha.channel_id", ChannelAdTableFull, AdTableFull, aaa.UserTableFull, PlanTableFull, ChannelTableFull, where),
		params...,
	)
	assert.Nil(err)
	return res
}

// UpdateAdView update ad view
func (m *Manager) UpdateAdView(ID, view int64) error {
	q := fmt.Sprintf("UPDATE %s SET view = ? WHERE id = ?", AdTableFull)
	_, err := m.GetDbMap().Exec(
		q,
		view,
		ID,
	)
	return err
}

// PubDashboardTotalView AdDashboard
type PubDashboardTotalView struct {
	ChannelName string `json:"channel_name" db:"name"`
	Viewed      int64  `json:"viewed" db:"viewed"`
}

// PubDashboardTotalView return the pub ad
func (m *Manager) PubDashboardTotalView(userID int64, scope base.UserScope) []PubDashboardTotalView {
	var where string
	var params []interface{}
	switch scope {
	case base.ScopeGlobal:
		where = ""
		params = []interface{}{}
	case base.ScopeParent:
		where = " AND ( u.id = ? OR u.parent_id = ? )"
		params = []interface{}{userID, userID}
	case base.ScopeSelf:
		where = " AND u.id = ?"
		params = []interface{}{userID}
	}
	res := []PubDashboardTotalView{}
	_, err := m.GetDbMap().Select(
		&res,
		fmt.Sprintf("SELECT ch.name AS name,SUM(cha.view) AS viewed FROM %s AS cha "+
			"INNER JOIN %s AS ch ON ch.id=cha.channel_id "+
			"INNER JOIN %s AS u ON u.id=ch.user_id "+
			"WHERE cha.start IS NOT NULL %s GROUP BY cha.channel_id",
			ChannelAdTableFull,
			ChannelTableFull,
			aaa.UserTableFull,
			where,
		),
		params...,
	)
	assert.Nil(err)
	return res
}

//AdDetailDataTable is the ad full data in data table, after join with other field
// @DataTable {
//		url = /detailtable/:id
//		entity = AdDetailDataTable
//		view = ad_detail:self
//		controller = modules/telegram/ad/adControllers
//		fill = FillAdDetailDataTableArray
// }
type AdDetailDataTable struct {
	Name        common.NullString `json:"name" db:"name" search:"true"`
	Active      ActiveStatus      `json:"active" db:"active" search:"true"`
	Src         common.NullString `db:"src" json:"src"`
	Description common.NullString `db:"description" json:"description"`
	Paid        sql.NullBool      `db:"paid" json:"paid" search:"true"`
	PlanType    PlanType          `db:"type" json:"plantype" search:"true"`
	Start       common.NullTime   `db:"start" json:"start"`
	PlanView    common.NullInt64  `db:"planview" json:"plan_view" sort:"true" title:"PlanView"`
	View        common.NullInt64  `json:"view" db:"view" visible:"true" title:"View"`
	ParentID    common.NullInt64  `db:"-" json:"parent_id" visible:"false"`
	OwnerID     common.NullInt64  `db:"-" json:"owner_id" visible:"false"`
	Actions     string            `db:"-" json:"_actions" visible:"false"`
}

// FillAdDetailDataTableArray is the function to handle
func (m *Manager) FillAdDetailDataTableArray(u base.PermInterfaceComplete,
	filters map[string]string,
	search map[string]string,
	contextparams map[string]string,
	sort, order string, p, c int) (AdDetailDataTableArray, int64) {
	var params []interface{}
	var res AdDetailDataTableArray
	var where []string

	id, err := strconv.ParseInt(contextparams["id"], 0, 64)
	assert.Nil(err)

	countQuery := fmt.Sprintf("SELECT COUNT(%[2]s.name) from %[1]s "+
		"LEFT JOIN %[2]s ON %[1]s.ad_id = %[2]s.id "+
		"LEFT JOIN %[4]s ON %[2]s.plan_id = %[4]s.id "+
		"LEFT JOIN %[3]s ON %[2]s.user_id = %[3]s.id "+
		"WHERE %[2]s.id=? ",
		ChannelAdTableFull,
		AdTableFull,
		aaa.UserTableFull,
		PlanTableFull,
	)
	query := fmt.Sprintf("SELECT %[2]s.name, %[1]s.active, src, %[2]s.description, pay_status AS paid,"+
		"type, start, %[4]s.view as planview, %[1]s.view from %[1]s "+
		"LEFT JOIN %[2]s ON %[1]s.ad_id = %[2]s.id "+
		"LEFT JOIN %[4]s ON %[2]s.plan_id = %[4]s.id "+
		"LEFT JOIN %[3]s ON %[2]s.user_id = %[3]s.id "+
		"WHERE %[2]s.id=? ",
		ChannelAdTableFull,
		AdTableFull,
		aaa.UserTableFull,
		PlanTableFull,
	)
	params = append(params, id)

	for field, value := range filters {
		where = append(where, fmt.Sprintf("%s.%s=?", ChannelAdTableFull, field))
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

//SpecificAdDataTable is the ad full data in data table, after join with other field
// @DataTable {
//		url = /detail/:id
//		entity = SpecificAdDataTable
//		view = get_ad:self
//		controller = modules/telegram/ad/adControllers
//		fill = FillSpecificAdDataTableArray
// }
type SpecificAdDataTable struct {
	Name     string           `json:"name" db:"name" search:"true"`
	View     common.NullInt64 `json:"view" db:"view" title:"View"`
	End      common.NullTime  `db:"end" json:"end"`
	Price    int              `json:"price" db:"-"`
	ParentID common.NullInt64 `db:"parent_id" json:"parent_id" visible:"false"`
	OwnerID  int64            `db:"owner_id" json:"owner_id" visible:"false"`
	Actions  string           `db:"-" json:"_actions" visible:"false"`
}

// FillSpecificAdDataTableArray is the function to handle
func (m *Manager) FillSpecificAdDataTableArray(u base.PermInterfaceComplete,
	filters map[string]string,
	search map[string]string,
	contextparams map[string]string,
	sort, order string, p, c int) (SpecificAdDataTableArray, int64) {
	var params []interface{}
	var res SpecificAdDataTableArray
	var where []string

	id, err := strconv.ParseInt(contextparams["id"], 0, 64)
	assert.Nil(err)

	countQuery := fmt.Sprintf("SELECT count(channel_id) from %[1]s "+
		"LEFT JOIN %[2]s ON %[1]s.ad_id = %[2]s.id",
		ChannelAdTableFull,
		AdTableFull)

	query := "SELECT %[3]s.name, %[1]s.view, end from %[2]s " +
		"LEFT JOIN %[1]s on %[2]s.ad_id = %[1]s.id " +
		"LEFT JOIN %[3]s ON %[2]s.channel_id = %[3]s.id "
	query = fmt.Sprintf(query, AdTableFull, ChannelAdTableFull, ChannelTableFull)

	for field, value := range filters {
		where = append(where, fmt.Sprintf("%s.%s=?", ChannelAdTableFull, field))
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
		where = append(where, fmt.Sprintf("%[1]s.parent_id=? OR %[1]s.user_id=?", aaa.UserTableFull))
		params = append(params, currentUserID, currentUserID)
	}

	where = append(where, "ad_id=?", "NOW() > end")

	//check for perm
	if len(where) > 0 {
		query += " WHERE "
		countQuery += " WHERE "
		params = append(params, id)
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

	//todo
	// price algorithm

	return res, count
}

// UpdateAdPromote update ad view
func (m *Manager) UpdateAdPromote(ID int64, data common.NullString) error {
	q := fmt.Sprintf("UPDATE %s SET promote_data = ? WHERE id = ?", AdTableFull)
	_, err := m.GetDbMap().Exec(
		q,
		data.String,
		ID,
	)
	return err
}
