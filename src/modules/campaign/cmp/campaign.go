// Package cmp is the models for campaign module
package cmp

import (
	"common/assert"
	"common/models/common"
	"fmt"
	"modules/misc/base"
	"modules/user/aaa"
	"strings"
	"time"
)

// Campaign model
// @Model {
//		table = campaigns
//		primary = true, id
//		find_by = id,user_id
//		list = yes
// }
type Campaign struct {
	ID        int64           `db:"id" json:"id" sort:"true" title:"ID" map:"campaigns.id"`
	UserID    int64           `json:"user_id" db:"user_id" title:"UserID"`
	Name      string          `json:"name" db:"name" search:"true" title:"Name"`
	Active    CampaignActive  `json:"active" db:"active" filter:"true" title:"Status"`
	Start     common.NullTime `json:"start" db:"start" title:"Start"`
	End       common.NullTime `json:"end" db:"end" title:"End"`
	CreatedAt time.Time       `db:"created_at" json:"created_at" sort:"true" title:"Created at"`
	UpdatedAt time.Time       `db:"updated_at" json:"updated_at" sort:"true" title:"Updated at"`
}

const (
	/*	CampaignStatusPending  CampaignStatus = "pending"
		CampaignStatusRejected CampaignStatus = "rejected"
		CampaignStatusAccepted CampaignStatus = "accepted"
		CampaignStatusArchive  CampaignStatus = "archive"*/

	CampaignActiveStart CampaignActive = "yes"
	CampaignActiveStop  CampaignActive = "no"
)

type (

	// CampaignActive is the campaign active
	// @Enum{
	// }
	CampaignActive string
)

//CampaignDataTable is the role full data in data table, after join with other field
// @DataTable {
//		url = /list
//		entity = campaign
//		view = campaign_list:self
//		controller = modules/campaign/controllers
//		fill = FillCampaignDataTableArray
//		_edit = campaign_edit:self
//		_change = campaign_manage:global
// }
type CampaignDataTable struct {
	Campaign
	Email    string `db:"email" json:"email" search:"true" title:"Email"`
	ParentID int64  `db:"parent_id_dt" json:"parent_id" visible:"false"`
	OwnerID  int64  `db:"owner_id_dt" json:"owner_id" visible:"false"`
}

// FillCampaignDataTableArray is the function to handle
func (m *Manager) FillCampaignDataTableArray(u base.PermInterfaceComplete, filters map[string]string, search map[string]string, sort, order string, p, c int) (CampaignDataTableArray, int64) {
	var params []interface{}
	var res CampaignDataTableArray
	var where []string

	countQuery := fmt.Sprintf("SELECT COUNT(campaigns.id) FROM %s LEFT JOIN %s ON %s.id=%s.user_id", CampaignTableFull, aaa.UserTableFull, aaa.UserTableFull, CampaignTableFull)
	query := fmt.Sprintf("SELECT campaigns.*,campaigns.user_id AS owner_id_dt,CASE WHEN users.parent_id IS NOT NULL THEN users.parent_id ELSE 0 END AS parent_id_dt,users.email FROM %s LEFT JOIN %s ON %s.id=%s.user_id", CampaignTableFull, aaa.UserTableFull, aaa.UserTableFull, CampaignTableFull)
	for field, value := range filters {
		where = append(where, fmt.Sprintf(CampaignTableFull+".%s=%s", field, "?"))
		params = append(params, value)
	}

	for column, val := range search {
		where = append(where, fmt.Sprintf("%s LIKE ?", column))
		params = append(params, fmt.Sprintf("%s"+val+"%s", "%", "%"))
	}

	currentUserID := u.GetID()
	highestScope := u.GetCurrentScope()

	if highestScope == base.ScopeSelf {
		where = append(where, fmt.Sprintf("%s.user_id=?", CampaignTableFull))
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

func (c *Campaign) Initialize() {

}

// Create campaign
func (m *Manager) Create(user *aaa.User, name string, start, end time.Time) *Campaign {

	c := &Campaign{
		UserID: user.ID,
		Name:   name,
		Active: CampaignActiveStop,
		Start:  common.NullTime{Valid: !start.IsZero(), Time: start},
		End:    common.NullTime{Valid: !end.IsZero(), Time: end},
	}
	assert.Nil(m.CreateCampaign(c))
	return c
}

// ChangeActive toggle between active status
func (m *Manager) ChangeActive(ID int64, userID int64, currentStat CampaignActive, createdAt time.Time) *Campaign {
	ch := &Campaign{
		ID:        ID,
		UserID:    userID,
		CreatedAt: createdAt,
	}
	if currentStat == CampaignActiveStart {
		ch.Active = CampaignActiveStop
	} else {
		ch.Active = CampaignActiveStart
	}
	assert.Nil(m.UpdateCampaign(ch))
	return ch
}
