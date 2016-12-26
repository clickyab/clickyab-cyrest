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

const (
	AdStatusPending  AdStatus = "pending"
	AdStatusRejected AdStatus = "rejected"
	AdStatusAccepted AdStatus = "accepted"
	AdStatusArchive  AdStatus = "archive"

	AdTypeImg      AdType = "img"
	AdTypeDocument AdType = "document"
	AdTypeVideo    AdType = "video"
)

type (
	// AdStatus is the ad status
	// @Enum{
	// }
	AdStatus string

	// AdType is the ad status
	// @Enum{
	// }
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
	ID          int64             `db:"id" json:"id" sort:"true" title:"ID"`
	UserID      int64             `json:"user_id" db:"user_id" title:"UserID"`
	Name        string            `json:"name" db:"name" search:"true" title:"Name"`
	Media       common.NullString `json:"media" db:"media" title:"Media"`
	Description common.NullString `json:"description" db:"description" title:"Description"`
	Link        common.NullString `json:"link" db:"link" search:"true" title:"Link"`
	Type        AdType            `json:"type" db:"type" filter:"true" title:"Type"`
	Status      AdStatus          `json:"status" db:"status" filter:"true" title:"Status"`
	CreatedAt   time.Time         `db:"created_at" json:"created_at" sort:"true" title:"Created at"`
	UpdatedAt   time.Time         `db:"updated_at" json:"updated_at" sort:"true" title:"Updated at"`
}

//AdDataTable is the ad full data in data table, after join with other field
// @DataTable {
//		url = /list
//		entity = ad
//		view = ad_list:self
//		controller = modules/ad/controllers
//		fill = FillAdDataTableArray
//		_edit = ad_edit:self
//		_change = ad_manage:global
// }
type AdDataTable struct {
	Ad
	Email    string `db:"email" json:"email" search:"true" title:"Email"`
	ParentID int64  `db:"-" json:"parent_id" visible:"false"`
	OwnerID  int64  `db:"-" json:"owner_id" visible:"false"`
	Actions  string `db:"-" json:"_actions" visible:"false"`
}

// FillAdDataTableArray is the function to handle
func (m *Manager) FillAdDataTableArray(u base.PermInterfaceComplete, filters map[string]string, search map[string]string, sort, order string, p, c int) (AdDataTableArray, int64) {
	var params []interface{}
	var res AdDataTableArray
	var where []string

	countQuery := fmt.Sprintf("SELECT COUNT(ads.id) FROM %s LEFT JOIN %s ON %s.id=%s.user_id", AdTableFull, aaa.UserTableFull, aaa.UserTableFull, AdTableFull)
	query := fmt.Sprintf("SELECT ads.*,users.email FROM %s LEFT JOIN %s ON %s.id=%s.user_id", AdTableFull, aaa.UserTableFull, aaa.UserTableFull, AdTableFull)
	for field, value := range filters {
		where = append(where, fmt.Sprintf(AdTableFull+".%s=%s", field, "?"))
		params = append(params, value)
	}

	for column, val := range search {
		where = append(where, fmt.Sprintf("%s LIKE ?", column))
		params = append(params, fmt.Sprintf("%s"+val+"%s", "%", "%"))
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

// EditAd function for ad editing
func (m *Manager) EditAd(usr *aaa.User, ad *Ad) *Ad {

	ads := &Ad{
		ID:          ad.ID,
		UserID:      ad.UserID,
		Name:        ad.Name,
		Media:       ad.Media,
		Description: ad.Description,
		Link:        ad.Link,
		Type:        ad.Type,
		Status:      ad.Status,
		CreatedAt:   ad.CreatedAt,
		UpdatedAt:   ad.UpdatedAt,
	}

	assert.Nil(m.UpdateAd(ads))
	return ads
}
