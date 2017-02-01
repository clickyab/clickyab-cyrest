package ads

import (
	"common/assert"
	"common/models/common"
	"fmt"
	"modules/misc/base"
	"modules/user/aaa"
	"strings"
	"time"
)

// ActiveStatus is the channel_ad active
// @Enum {
// }
type ActiveStatus string

// PlanType is the plan type
// @Enum {
// }
type PlanType string

const (
	//ActiveStatusYes is the yes status
	ActiveStatusYes ActiveStatus = "yes"
	// ActiveStatusNo is the no status
	ActiveStatusNo ActiveStatus = "no"
	// PlanTypePromotion is the promotion status
	PlanTypePromotion PlanType = "promotion"
	// PlanTypeIndividual is the individual status
	PlanTypeIndividual PlanType = "individual"
)

// ChannelAd is the list of ad in channel for cyborg
// @Model {
//		table = channel_ad
//		primary = false, channel_id,ad_id
//		find_by = channel_id, ad_id
// }
type ChannelAd struct {
	ChannelID    int64             `db:"channel_id" json:"channel_id"`
	AdID         int64             `db:"ad_id" json:"ad_id"`
	View         int64             `db:"view" json:"view"`
	CliMessageID common.NullString `db:"cli_message_id" json:"cli_message_id"`
	BotChatID    int64             `db:"bot_chat_id" json:"bot_chat_id"`
	BotMessageID int               `db:"bot_message_id" json:"bot_message_id"`
	Active       ActiveStatus      `db:"active" json:"active"`
	Start        common.NullTime   `db:"start" json:"start"`
	End          common.NullTime   `db:"end" json:"end"`
	Warning      int64             `db:"warning" json:"warning"`
	PossibleView int64             `db:"possible_view" json:"possible_view"`
	CreatedAt    time.Time         `db:"created_at" json:"created_at" sort:"true"`
	UpdatedAt    time.Time         `db:"updated_at" json:"updated_at" sort:"true"`
}

//SelectAd choose ad
type SelectAd struct {
	Ad
	View          int64 `db:"view" json:"view"`
	Viewed        int64 `db:"viewed" json:"viewed"`
	PossibleView  int64 `db:"possible_view" json:"possible_view"`
	AffectiveView int64 `json:"affective_view"`
}

//ByAffectiveView sort by AffectiveView
type ByAffectiveView []SelectAd

func (a ByAffectiveView) Len() int {
	return len(a)
}
func (a ByAffectiveView) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a ByAffectiveView) Less(i, j int) bool {
	return a[i].AffectiveView < a[j].AffectiveView
}

// FindChannelIDAdByAdID return the ChannelAd base on its ad_id
func (m *Manager) FindChannelIDAdByAdID(channelID int64, addID int64) (*ChannelAd, error) {
	var res ChannelAd
	err := m.GetDbMap().SelectOne(
		&res,
		fmt.Sprintf("SELECT * FROM %s WHERE channel_id=? AND ad_id=?", ChannelAdTableFull),
		channelID,
		addID,
	)

	if err != nil {
		return nil, err
	}

	return &res, nil
}

// FindChannelAdByAdIDActive return the ChannelAd base on its ad_id
func (m *Manager) FindChannelAdByAdIDActive(a int64) ([]ChannelAd, error) {
	res := []ChannelAd{}
	_, err := m.GetDbMap().Select(
		&res,
		fmt.Sprintf("SELECT * FROM %s WHERE ad_id=? AND active='yes'", ChannelAdTableFull),
		a,
	)

	if err != nil {
		return nil, err
	}

	return res, nil
}

// FindChannelAdActiveByChannelID return the ChannelAd base on its channel_id,active
func (m *Manager) FindChannelAdActiveByChannelID(channelID int64, status ActiveStatus) ([]ChannelAd, error) {
	res := []ChannelAd{}
	_, err := m.GetDbMap().Select(
		&res,
		fmt.Sprintf("SELECT * FROM %s WHERE channel_id=? AND active='?'", ChannelAdTableFull),
		channelID,
		status,
	)

	if err != nil {
		return nil, err
	}

	return res, nil
}

// ChooseAd return the ads
func (m *Manager) ChooseAd(channelID int64) ([]SelectAd, error) {
	res := []SelectAd{}
	_, err := m.GetDbMap().Select(
		&res,
		fmt.Sprintf("SELECT %s.*,sum(%s.possible_view) as possible_view,sum(%s.view) as viewed ,%s.view as view "+
			"FROM %s "+
			"LEFT JOIN %s ON %s.id = %s.plan_id "+
			"INNER JOIN %s on %s.ad_id = %s.id "+
			"WHERE %s.channel_id != ? "+
			"GROUP BY %s.ad_id ",
			AdTableFull,
			ChannelAdTableFull,
			ChannelAdTableFull,
			PlanTableFull,

			AdTableFull,

			PlanTableFull,
			PlanTableFull,
			AdTableFull,

			ChannelAdTableFull,
			ChannelAdTableFull,
			AdTableFull,
			ChannelAdTableFull,
			ChannelAdTableFull,
		),
		channelID,
	)

	if err != nil {
		return nil, err
	}

	return res, nil
}

//ReportAdDataTable is the ad full data in data table, after join with other field
// @DataTable {
//		url = /report
//		entity = adReport
//		view = report_ad:self
//		controller = modules/telegram/ad/adControllers
//		fill = FillAdReportDataTableArray
//		_edit = ad_edit:self
//		_change = ad_manage:global
// }
type ReportAdDataTable struct {
	Ad
	View         int64           `db:"view" json:"view" sort:"true" title:"View"`
	Active       ActiveStatus    `db:"active" json:"active" title:"ActiveStatus" filter:"true"`
	Start        common.NullTime `db:"start" json:"start" sort:"true" title:"Start"`
	End          common.NullTime `db:"end" json:"end" sort:"true" title:"End"`
	Warning      int64           `db:"warning" json:"warning" sort:"true" title:"Warning"`
	PossibleView int64           `db:"possible_view" json:"possible_view" sort:"true" title:"PossibleView"`
	Email        string          `db:"email" json:"email" search:"true" title:"Email"`
	ParentID     int64           `db:"-" json:"parent_id" visible:"false"`
	OwnerID      int64           `db:"-" json:"owner_id" visible:"false"`
	Actions      string          `db:"-" json:"_actions" visible:"false"`
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

	countQuery := fmt.Sprintf("SELECT COUNT(%s.ad_id) FROM %s "+
		"LEFT JOIN %s ON %s.id=%s.ad_id "+
		"LEFT JOIN %s ON %s.id=%s.user_id "+
		"GROUP BY %s.ad_id",
		ChannelAdTableFull,
		ChannelAdTableFull,

		AdTableFull,
		AdTableFull,
		ChannelAdTableFull,

		aaa.UserTableFull,
		aaa.UserTableFull,
		AdTableFull,

		ChannelAdTableFull,
	)
	query := fmt.Sprintf("SELECT %s.*,%s.view,%s.warning,%s.active,%s.start,%s.possible_view,%s.end,%s.email FROM %s "+
		"LEFT JOIN %s ON %s.id=%s.ad_id "+
		"LEFT JOIN %s ON %s.id=%s.user_id ",
		AdTableFull,
		ChannelAdTableFull,
		ChannelAdTableFull,
		ChannelAdTableFull,
		ChannelAdTableFull,
		ChannelAdTableFull,
		ChannelAdTableFull,
		aaa.UserTableFull,
		ChannelAdTableFull,

		AdTableFull,
		AdTableFull,
		ChannelAdTableFull,

		aaa.UserTableFull,
		aaa.UserTableFull,
		AdTableFull,
	)
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
		where = append(where, fmt.Sprintf("%s.user_id=?", ChannelAdTableFull))
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
	query += fmt.Sprintf(" GROUP BY %s.ad_id ", ChannelAdTableFull)
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

//ReportChannelDataTable is the ad full data in data table, after join with other field
// @DataTable {
//		url = /report
//		entity = channelReport
//		view = report_channel:self
//		controller = modules/telegram/ad/chanControllers
//		fill = FillChannelReportDataTableArray
//		_edit = ad_edit:self
//		_change = ad_manage:global
// }
type ReportChannelDataTable struct {
	Channel
	View         int64           `db:"view" json:"view" sort:"true" title:"View"`
	Active       ActiveStatus    `db:"active" json:"active" title:"ActiveStatus" filter:"true"`
	Start        common.NullTime `db:"start" json:"start" sort:"true" title:"Start"`
	End          common.NullTime `db:"end" json:"end" sort:"true" title:"End"`
	Warning      int64           `db:"warning" json:"warning" sort:"true" title:"Warning"`
	PossibleView int64           `db:"possible_view" json:"possible_view" sort:"true" title:"PossibleView"`
	Email        string          `db:"email" json:"email" search:"true" title:"Email"`
	ParentID     int64           `db:"-" json:"parent_id" visible:"false"`
	OwnerID      int64           `db:"-" json:"owner_id" visible:"false"`
	Actions      string          `db:"-" json:"_actions" visible:"false"`
}

// FillChannelReportDataTableArray is the function to handle
func (m *Manager) FillChannelReportDataTableArray(
	u base.PermInterfaceComplete,
	filters map[string]string,
	search map[string]string,
	sort, order string, p, c int) (ReportChannelDataTableArray, int64) {
	var params []interface{}
	var res ReportChannelDataTableArray
	var where []string

	countQuery := fmt.Sprintf("SELECT COUNT(%s.channel_id) FROM %s "+
		"LEFT JOIN %s ON %s.id=%s.channel_id "+
		"LEFT JOIN %s ON %s.id=%s.user_id "+
		"GROUP BY %s.channel_id",
		ChannelAdTableFull,
		ChannelAdTableFull,

		ChannelTableFull,
		ChannelTableFull,
		ChannelAdTableFull,

		aaa.UserTableFull,
		aaa.UserTableFull,
		ChannelTableFull,

		ChannelAdTableFull,
	)
	query := fmt.Sprintf("SELECT %s.*,%s.view,%s.warning,%s.active,%s.start,%s.possible_view,%s.end,%s.email FROM %s "+
		"LEFT JOIN %s ON %s.id=%s.channel_id "+
		"LEFT JOIN %s ON %s.id=%s.user_id ",
		ChannelTableFull,
		ChannelAdTableFull,
		ChannelAdTableFull,
		ChannelAdTableFull,
		ChannelAdTableFull,
		ChannelAdTableFull,
		ChannelAdTableFull,
		aaa.UserTableFull,
		ChannelAdTableFull,

		ChannelTableFull,
		ChannelTableFull,
		ChannelAdTableFull,

		aaa.UserTableFull,
		aaa.UserTableFull,
		ChannelTableFull,
	)
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
		where = append(where, fmt.Sprintf("%s.user_id=?", ChannelAdTableFull))
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
