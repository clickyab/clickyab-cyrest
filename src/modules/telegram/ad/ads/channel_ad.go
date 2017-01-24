package ads

import (
	"common/models/common"
	"fmt"
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
func (m *Manager) FindChannelIDAdByAdID(c int64, a int64) (*ChannelAd, error) {
	var res ChannelAd
	err := m.GetDbMap().SelectOne(
		&res,
		fmt.Sprintf("SELECT * FROM %s WHERE ad_id=? AND channel_id=?", ChannelAdTableFull),
		a,
		c,
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
