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

// PlanType is the plan type
// @Enum {
// }
type PlanType string

const (
	// PlanTypePromotion is the promotion status
	PlanTypePromotion PlanType = "promotion"
	// PlanTypeIndividual is the individual status
	PlanTypeIndividual PlanType = "individual"
)

// ChannelAd is the list of ad in channel for cyborg
// @Model {
//		table = channel_ad
//		primary = false, channel_id,ad_id
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
	CreatedAt    *time.Time        `db:"created_at" json:"created_at" sort:"true"`
	UpdatedAt    *time.Time        `db:"updated_at" json:"updated_at" sort:"true"`
}

//SelectAd choose ad
type SelectAd struct {
	Ad
	Type          AdType           `json:"type" db:"type" title:"Type"`
	PlanView      int64            `db:"plan_view" json:"plan_view"`
	Viewed        common.NullInt64 `db:"viewed" json:"viewed"`
	PossibleView  common.NullInt64 `db:"possible_view" json:"possible_view"`
	AffectiveView int64            `json:"affective_view"`
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

// ChannelAdActive type active ad
type ChannelAdActive struct {
	ChannelAd
	UserID int64 `db:"user_id"`
	Share  int64 `db:"share"`
}

// FindChannelIDAdByAdIDByActive return the ChannelAd base on its ad_id
func (m *Manager) FindChannelIDAdByAdIDByActive(channelID int64, addID int64, userID int64) (*ChannelAdActive, error) {
	var res ChannelAdActive
	err := m.GetDbMap().SelectOne(
		&res,
		fmt.Sprintf("SELECT ca.*,u.id AS user_id, p.share AS share"+
			"FROM %s AS ca "+
			"INNER JOIN %s AS c "+
			"ON ca.channel_id = c.id "+
			"INNER JOIN %s AS u "+
			"ON u.id = c.user_id"+
			"INNER JOIN %s as a "+
			"ON ca.ad_id = a.id  "+
			"INNER JOIN %s AS p "+
			"ON p.id = a.plan_id"+
			"WHERE c.user_id = u.id "+
			"AND ca.active = ? "+
			"AND u.id = ? "+
			"AND ca.channel_id = ? "+
			"AND ca.ad_id = ?", ChannelAdTableFull, ChannelTableFull, aaa.UserTableFull, AdTableFull, PlanTableFull),
		ActiveStatusYes,
		userID,
		channelID,
		addID,
	)

	if err != nil {
		return nil, err
	}

	return &res, nil
}

// FindChannelAdActiveByChannelID return the ChannelAd base on its channel_id,active
func (m *Manager) FindChannelAdActiveByChannelID(channelID int64, status ActiveStatus) ([]ChannelAd, error) {
	res := []ChannelAd{}
	_, err := m.GetDbMap().Select(
		&res,
		fmt.Sprintf("SELECT * FROM %s WHERE channel_id=? AND active=? AND end IS NULL", ChannelAdTableFull),
		channelID,
		status,
	)

	if err != nil {
		return nil, err
	}

	return res, nil
}

// DeleteChannelAdByChannelID delete the ChannelAd
func (m *Manager) DeleteChannelAdByChannelID(channelID int64) error {
	q := fmt.Sprintf("DELETE FROM %s WHERE channel_id=? AND start IS NULL AND end IS NULL", ChannelAdTableFull)
	_, err := m.GetDbMap().Exec(
		q,
		channelID,
	)

	return err
}

// ChannelAdD ChannelAdD
type ChannelAdD struct {
	ChannelID    int64             `db:"channel_id" json:"channel_id"`
	AdID         int64             `db:"ad_id" json:"ad_id"`
	View         int64             `db:"view" json:"view"`
	Src          common.NullString `json:"src" db:"src"`
	CliMessageID common.NullString `db:"cli_message_id" json:"cli_message_id"`
	CliMessageAd common.NullString `db:"cli_message_ad" json:"cli_message_ad"`
	PromoteData  common.NullString `db:"promote_data" json:"promote_data"`
	PlanView     int64             `db:"plan_view" json:"plan_view"`
	PlanPosition int64             `json:"plan_position" db:"plan_position"`
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

// FindChannelAdByChannelIDActive return the ChannelAd base on its channel_id
func (m *Manager) FindChannelAdByChannelIDActive(a int64) ([]ChannelAdD, error) {
	res := []ChannelAdD{}
	_, err := m.GetDbMap().Select(
		&res,
		fmt.Sprintf("SELECT %[1]s.*,%[2]s.cli_message_id AS cli_message_ad,%[2]s.promote_data,%[2]s.src, "+
			"%[3]s.view AS plan_view,%[3]s.position AS plan_position"+
			" FROM %[1]s INNER JOIN %[2]s ON %[2]s.id=%[1]s.ad_id "+
			" INNER JOIN %[3]s ON %[3]s.id=%[2]s.plan_id "+
			" WHERE channel_id=? "+
			" AND %[1]s.active='yes' "+
			" AND %[2]s.active_status='yes' "+
			" AND end IS NULL",
			ChannelAdTableFull,
			AdTableFull,
			PlanTableFull,
		),
		a,
	)

	if err != nil {
		return nil, err
	}

	return res, nil
}

// FindChannelAdByChannelID return the ChannelAd base on its channel_id
func (m *Manager) FindChannelAdByChannelID(a int64) ([]ChannelAdD, error) {
	res := []ChannelAdD{}
	_, err := m.GetDbMap().Select(
		&res,
		fmt.Sprintf("SELECT %[1]s.*,%[2]s.cli_message_id AS cli_message_ad,%[2]s.promote_data,%[2]s.src, "+
			"%[3]s.position AS plan_position,"+
			"%[3]s.view AS plan_view"+
			" FROM %[1]s INNER JOIN %[2]s ON %[2]s.id=%[1]s.ad_id "+
			" INNER JOIN %[3]s ON %[3]s.id=%[2]s.plan_id "+
			" WHERE channel_id=? "+
			" AND %[1]s.active='no' "+
			" AND start IS NULL"+
			" AND end IS NULL",
			ChannelAdTableFull,
			AdTableFull,
			PlanTableFull,
		),
		a,
	)

	if err != nil {
		return nil, err
	}

	return res, nil
}

// FindReshotChannelAdByChannelID return the ChannelAd base on its channel_id
func (m *Manager) FindReshotChannelAdByChannelID(a int64) ([]ChannelAdD, error) {
	res := []ChannelAdD{}
	_, err := m.GetDbMap().Select(
		&res,
		fmt.Sprintf("SELECT %[1]s.*,%[2]s.cli_message_id AS cli_message_ad,%[2]s.promote_data,%[2]s.src, "+
			"%[3]s.position AS plan_position,"+
			"%[3]s.view AS plan_view"+
			" FROM %[1]s INNER JOIN %[2]s ON %[2]s.id=%[1]s.ad_id "+
			" INNER JOIN %[3]s ON %[3]s.id=%[2]s.plan_id "+
			" WHERE channel_id=? "+
			" AND %[1]s.active='yes' "+
			" AND start IS NOT NULL"+
			" AND end IS NULL",
			ChannelAdTableFull,
			AdTableFull,
			PlanTableFull,
		),
		a,
	)

	if err != nil {
		return nil, err
	}

	return res, nil
}

// ChannelAdChan struct channel ad
type ChannelAdChan struct {
	ChannelAd ChannelAd
	Channel   Channel
}

// FindChannelAdActiveByAdID return the adID base on its ad_id,ActiveStatus
func (m *Manager) FindChannelAdActiveByAdID(adID int64, status ActiveStatus) ([]ChannelAd, error) {
	res := []ChannelAd{}
	_, err := m.GetDbMap().Select(
		&res,
		fmt.Sprintf("SELECT * FROM %s WHERE ad_id=? AND active=?", ChannelAdTableFull),
		adID,
		status,
	)

	if err != nil {
		return nil, err
	}

	return res, nil
}

// FindChannelAdActive return the adID base on its ad_id,ActiveStatus
func (m *Manager) FindChannelAdActive() ([]ChannelAdD, error) {
	res := []ChannelAdD{}
	_, err := m.GetDbMap().Select(
		&res,
		fmt.Sprintf("SELECT %[1]s.*,%[2]s.cli_message_id AS cli_message_ad,%[2]s.promote_data,%[2]s.src, "+
			"%[3]s.view AS plan_view,%[3]s.position AS plan_position"+
			" FROM %[1]s INNER JOIN %[2]s ON %[2]s.id=%[1]s.ad_id "+
			" INNER JOIN %[3]s ON %[3]s.id=%[2]s.plan_id "+
			" AND %[1]s.active='yes' "+
			" AND %[2]s.active_status='yes' "+
			" AND end IS NULL",
			ChannelAdTableFull,
			AdTableFull,
			PlanTableFull,
		),
	)

	if err != nil {
		return nil, err
	}

	return res, nil
}

// UpdateChannelAds update channel ads
func (m *Manager) UpdateChannelAds(ca []ChannelAd) error {
	var q = fmt.Sprintf("UPDATE %s SET  updated_at=?,warning=? , view=? WHERE channel_id=? AND ad_id=?", ChannelAdTableFull)
	for i := range ca {
		_, err := m.GetDbMap().Exec(
			q,
			time.Now(),
			ca[i].Warning,
			ca[i].View,
			ca[i].ChannelID,
			ca[i].AdID,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

// UpdateActiveEndChannelAds update channel ads
func (m *Manager) UpdateActiveEndChannelAds(ca []*ChannelAd) error {
	for i := range ca {
		err := m.UpdateChannelAd(ca[i])
		if err != nil {
			return err
		}
	}
	return nil
}

// ChooseAd return the ads
func (m *Manager) ChooseAd(channelID int64) ([]SelectAd, error) {
	res := []SelectAd{}
	_, err := m.GetDbMap().Select(
		&res,
		fmt.Sprintf("SELECT %[1]s.*,sum(%[2]s.possible_view) as possible_view,"+
			"%[3]s.view as plan_view, %[1]s.cli_message_id IS NULL as type ,sum(%[2]s.view) as viewed "+
			"FROM %[1]s "+
			"LEFT JOIN %[3]s ON %[3]s.id = %[1]s.plan_id "+
			"LEFT JOIN %[2]s on %[2]s.ad_id = %[1]s.id "+
			"WHERE  ( %[2]s.channel_id != ? OR %[2]s.channel_id IS NULL ) "+
			" AND %[1]s.plan_id IS NOT NULL "+
			" AND %[1]s.admin_status = ? "+
			" AND %[1]s.active_status = ? "+
			" AND %[1]s.pay_status = ? "+
			"GROUP BY %[1]s.id ",
			AdTableFull,
			ChannelAdTableFull,
			PlanTableFull,
		),
		channelID,
		AdminStatusAccepted,
		ActiveStatusYes,
		ActiveStatusYes,
	)

	if err != nil {
		return nil, err
	}

	return res, nil
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
	View         int64            `db:"view" json:"view" sort:"true" title:"View"`
	Active       ActiveStatus     `db:"active" json:"active" title:"ActiveStatus" filter:"true"`
	Start        common.NullTime  `db:"start" json:"start" sort:"true" title:"Start"`
	End          common.NullTime  `db:"end" json:"end" sort:"true" title:"End"`
	Warning      int64            `db:"warning" json:"warning" sort:"true" title:"Warning"`
	PossibleView int64            `db:"possible_view" json:"possible_view" sort:"true" title:"PossibleView"`
	Email        string           `db:"email" json:"email" search:"true" title:"Email"`
	ParentID     common.NullInt64 `db:"parent_id" json:"parent_id" visible:"false"`
	OwnerID      int64            `db:"owner_id" json:"owner_id" visible:"false"`
	Actions      string           `db:"-" json:"_actions" visible:"false"`
}

// FillChannelReportDataTableArray is the function to handle
func (m *Manager) FillChannelReportDataTableArray(
	u base.PermInterfaceComplete,
	filters map[string]string,
	search map[string]string,
	contextparams map[string]string,
	sort, order string, p, c int) (ReportChannelDataTableArray, int64) {
	var params []interface{}
	var res ReportChannelDataTableArray
	var where []string

	countQuery := fmt.Sprintf("SELECT COUNT(%[1]s.channel_id) FROM %[1]s "+
		"INNER JOIN %[2]s ON %[2]s.id=%[1]s.channel_id "+
		"INNER JOIN %[3]s ON %[3]s.id=%[2]s.user_id "+
		"GROUP BY %[1]s.channel_id",
		ChannelAdTableFull,
		ChannelTableFull,
		aaa.UserTableFull,
	)
	query := fmt.Sprintf("SELECT %[1]s.*,%[2]s.view,%[2]s.warning,%[2]s.active,%[2]s.start,%[2]s.possible_view,%[2]s.end,%[3]s.email,%[3]s.id AS owner_id, %[3]s.parent_id as parent_id FROM %[2]s "+
		"INNER JOIN %[1]s ON %[1]s.id=%[2]s.channel_id "+
		"INNER JOIN %[3]s ON %[3]s.id=%[1]s.user_id ",
		ChannelTableFull,
		ChannelAdTableFull,
		aaa.UserTableFull,
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

// SetCLIMessageID try to set CLI message id on channel_ad table
func (m *Manager) SetCLIMessageID(channelID, adID int64, cliMessageID string) error {
	q := fmt.Sprintf("UPDATE %s SET cli_message_id = ? WHERE channel_id=? AND ad_id = ?", ChannelAdTableFull)
	_, err := m.GetDbMap().Exec(q, cliMessageID, channelID, adID)
	return err
}

// UpdateOnDuplicateChanDetail try to save a new ChanDetail or update in database
func (m *Manager) UpdateOnDuplicateChanDetail(cd *ChanDetail) error {
	now := time.Now()
	cd.CreatedAt = &now
	cd.UpdatedAt = &now
	_, err := m.GetDbMap().Exec(fmt.Sprintf("INSERT INTO %s "+
		"(id,name, channel_id, title, info, cli_telegram_id, user_count, admin_count, post_count, total_view,created_at,updated_at)"+
		" VALUES (NULL,?,?,?,?,?,?,?,?,?,?,?)"+
		" ON DUPLICATE KEY UPDATE "+
		"title=VALUES(title)"+
		",info=VALUES(info)"+
		",user_count=VALUES(user_count)"+
		",admin_count=VALUES(admin_count)"+
		",post_count=VALUES(post_count)"+
		",total_view=VALUES(total_view)"+
		",updated_at=VALUES(updated_at)", ChanDetailTableFull),
		cd.Name,
		cd.ChannelID,
		cd.Title,
		cd.Info,
		cd.TelegramID,
		cd.UserCount,
		cd.AdminCount,
		cd.PostCount,
		cd.TotalView,
		cd.CreatedAt,
		cd.UpdatedAt,
	)
	return err
}

// FindActiveChannelAdByChannelID try to find all channel ad for specific channel
func (m *Manager) FindActiveChannelAdByChannelID(channelID int64) []ChannelAd {
	var res []ChannelAd
	q := fmt.Sprintf("SELECT * FROM %s WHERE channel_id=? AND active = ?", ChannelAdTableFull)
	_, err := m.GetDbMap().Select(
		&res,
		q,
		channelID,
		ActiveStatusYes,
	)
	assert.Nil(err)
	return res
}

// ActiveAdUser struct active ad user
type ActiveAdUser struct {
	Ad
	ChannelID int64 `db:"channel_id"`
}

// FindActiveChannelAdByUserID try to find by user id
func (m *Manager) FindActiveChannelAdByUserID(userID int64) []ActiveAdUser {
	var res []ActiveAdUser
	q := fmt.Sprintf(
		"SELECT a.*,c.id AS channel_id "+
			"FROM %s AS ca "+
			"INNER JOIN %s AS c "+
			"ON ca.channel_id = c.id "+
			"INNER JOIN %s AS u "+
			"ON u.id = c.user_id"+
			"INNER JOIN %s as a "+
			"ON ca.ad_id = a.id "+
			"WHERE c.user_id = u.id "+
			"AND ca.active = ? "+
			"AND u.id = ?", ChannelAdTableFull, ChannelTableFull, aaa.UserTableFull, AdTableFull)
	_, err := m.GetDbMap().Select(
		&res,
		q,
		ActiveStatusYes,
		userID,
	)
	assert.Nil(err)
	return res
}
