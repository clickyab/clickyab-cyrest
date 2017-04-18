package worker

import (
	"common/assert"
	"common/rabbit"
	"modules/misc/trans"
	"modules/telegram/ad/ads"
	bot3 "modules/telegram/ad/worker"
	bot2 "modules/telegram/bot/worker"
	"modules/telegram/cyborg/commands"
	"sort"
)

const (
	promoNonPic selectAdType = "promoNonPic"
	nonPic      selectAdType = "nonPic"
)

type selectAdType string
type selectAd struct {
	ids    []int64
	adType selectAdType
}

func (mw *MultiWorker) selectAd(in *commands.SelectAd) (bool, error) {
	b := ads.NewAdsManager()
	chad, err := b.FindChannelAdActiveByChannelID(in.ChannelID, ads.ActiveStatusYes)
	assert.Nil(err)
	if len(chad) > 0 {
		return false, nil
	}
	chooseAds, err := b.ChooseAd(in.ChannelID)
	assert.Nil(err)
	if len(chooseAds) == 0 {
		rabbit.MustPublish(&bot2.SendWarn{
			AdID:      0,
			ChannelID: in.ChannelID,
			Msg:       trans.T("no ads for you").String(),
			ChatID:    in.ChatID,
		})
		return false, nil
	}
	for k := range chooseAds {
		chooseAds[k].AffectiveView = chooseAds[k].PlanView - chooseAds[k].PossibleView.Int64
	}
	sort.Sort(ads.ByAffectiveView(chooseAds))
	var (
		promoted     int64
		normalPic    int64
		normalNonPic int64
		selectedAd   selectAd
	)
	for i := range chooseAds {

		if promoted == 0 && chooseAds[i].CliMessageID.Valid {
			promoted = chooseAds[i].ID
		}
		if normalPic == 0 && !chooseAds[i].CliMessageID.Valid && chooseAds[i].Src.Valid {
			normalPic = chooseAds[i].ID
		}
		if normalNonPic == 0 && !chooseAds[i].CliMessageID.Valid && !chooseAds[i].Src.Valid {
			normalNonPic = chooseAds[i].ID
		}

		if promoted != 0 && normalNonPic != 0 {
			selectedAd = selectAd{
				ids:    []int64{promoted, normalNonPic},
				adType: promoNonPic,
			}
			break
		}

		if normalPic != 0 {
			selectedAd = selectAd{
				ids:    []int64{normalPic},
				adType: nonPic,
			}
			break
		}
	}
	if len(selectedAd.ids) == 0 {
		rabbit.MustPublish(&bot2.SendWarn{
			AdID:      0,
			ChannelID: in.ChannelID,
			Msg:       "no ads for you",
			ChatID:    in.ChatID,
		})
		return false, nil
	}

	//todo send ad to user
	rabbit.MustPublish(&bot3.AdDelivery{
		AdsID:     selectedAd.ids,
		ChannelID: in.ChannelID,
		ChatID:    in.ChatID,
	})
	return false, nil

}

func (mw *MultiWorker) transaction(m *ads.Manager, chad []ads.ChannelAd, channelAdDetail []*ads.ChannelAdDetail, avg int64) (bool, error) {
	err := m.Begin()
	if err != nil {
		return true, err
	}
	defer func() {
		if err != nil {
			assert.Nil(m.Rollback())
		} else {
			err = m.Commit()
		}

		if err != nil {
			chad = nil
		}
	}()
	err = m.CreateChannelAdDetails(channelAdDetail)
	assert.Nil(err)
	if err != nil {
		return true, err
	}
	err = m.UpdateChannelAds(chad)
	assert.Nil(err)
	if err != nil {
		return true, err
	}
	return false, err
}
