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
	// BTypeBanner type banner
	BTypeBanner BType = "banner"
	// BTypeBannerRep type banner rep
	BTypeBannerRep BType = "banner+rep"
	// BTypeRepBanner type rep banner
	BTypeRepBanner BType = "rep+banner"
	// BTypeRepBannerRep type rep banner rep
	BTypeRepBannerRep BType = "rep+banner+rep"
)

// BType bundle type
//	@Enum{
//	}
type BType string

// Bundles model
// @Model {
//		table = bundles
//		primary = true, id
//		find_by = id,user_id
//		list = yes
//	}
type Bundles struct {
	ID            int64             `db:"id" json:"id" sort:"true" title:"ID"`
	UserID        int64             `db:"user_id" json:"user_id" sort:"true" title:"UserID"`
	Position      int64             `db:"position" json:"position"  title:"Position"`
	View          int64             `db:"view" json:"view" sort:"true"  title:"View"`
	Price         int64             `db:"price" json:"price" sort:"true"  title:"Price"`
	PercentFinish int64             `db:"percent_finish" json:"percent_finish" sort:"true"  title:"PercentFinish"`
	BundleType    BType             `db:"bundle_type" json:"bundle_type"  title:"BundleType" filter:"true"`
	Rules         common.NullString `json:"rules" db:"rules"  title:"Rules" `
	AdminStatus   ActiveStatus      `json:"admin_status" db:"admin_status"  title:"AdminStatus" filter:"true"`
	ActiveStatus  ActiveStatus      `json:"active_status" db:"active_status"  title:"ActiveStatus" filter:"true"`
	Ads           common.CommaArray `json:"ads" db:"ads"  title:"Ads"`
	TargetAd      int64             `json:"target_ad" db:"target_ad"  title:"TargetAd"`
	Start         common.NullTime   `json:"start" db:"start" sort:"true" title:"Start"`
	End           common.NullTime   `json:"end" db:"end" sort:"true" title:"End"`
	CreatedAt     *time.Time        `json:"created_at" db:"created_at" sort:"true" title:"CreatedAt"`
	UpdatedAt     *time.Time        `json:"updated_at" db:"updated_at" sort:"true" title:"UpdatedAt"`
}

//BundleDataTable is the ad full data in data table, after join with other field
// @DataTable {
//		url = /bundle_list
//		entity = bundle_ad
//		view = bundle_list:self
//		controller = modules/telegram/ad/adControllers
//		fill = FillBundleDataTableArray
//		_edit = bundle_edit:self
//		_active_status = change_active_bundle:self
//		_admin_status = change_admin_bundle:self
// }
type BundleDataTable struct {
	Bundles
	Email    string           `db:"email" json:"email" search:"true" title:"Email" visible:"false"`
	ParentID common.NullInt64 `db:"parent_id" json:"parent_id" visible:"false"`
	OwnerID  int64            `db:"owner_id" json:"owner_id" visible:"false"`
	Actions  string           `db:"-" json:"_actions" visible:"false"`
}

// FillBundleDataTableArray is the function to handle
func (m *Manager) FillBundleDataTableArray(
	u base.PermInterfaceComplete,
	filters map[string]string,
	search map[string]string,
	contextparams map[string]string,
	sort, order string, p, c int) (BundleDataTableArray, int64) {
	var params []interface{}
	var res BundleDataTableArray
	var where []string

	countQuery := fmt.Sprintf("SELECT COUNT(%[1]s.id) FROM %[1]s",
		BundlesTableFull)
	query := fmt.Sprintf("SELECT %[1]s.*,%[2]s.email,%[2]s.id AS owner_id FROM %[1]s LEFT JOIN %[2]s ON %[2]s.id=%[1]s.user_id",
		BundlesTableFull,
		aaa.UserTableFull)
	for field, value := range filters {
		where = append(where, fmt.Sprintf("%s.%s=?", BundlesTableFull, field))
		params = append(params, value)
	}

	for column, val := range search {
		where = append(where, fmt.Sprintf("%s LIKE ?", column))
		params = append(params, "%"+val+"%")
	}

	currentUserID := u.GetID()
	highestScope := u.GetCurrentScope()

	if highestScope == base.ScopeSelf {
		where = append(where, fmt.Sprintf("%s.user_id=?", BundlesTableFull))
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

// CheckAdActiveAdBundle return all finished ads
func (m *Manager) CheckAdActiveAdBundle(ads common.CommaArray) ([]BundleChannelAd, error) {
	q := `SELECT * FROM %s AS b
		WHERE b.active = ?
		AND b.start NOT NULL
		AND b.ad_id IN (?)`
	q = fmt.Sprintf(q, BundleChannelAdTableFull)
	var res []BundleChannelAd
	_, err := m.GetDbMap().Select(&res, q, ActiveStatusYes, ads)
	return res, err
}
