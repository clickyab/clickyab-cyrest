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

	"github.com/Sirupsen/logrus"
)

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
		promoted int64
		normal   int64
	)
	for i := range chooseAds {
		if promoted == 0 && chooseAds[i].CliMessageID.Valid {
			promoted = chooseAds[i].ID
		}
		if normal == 0 && !chooseAds[i].CliMessageID.Valid {
			normal = chooseAds[i].ID
		}

		if promoted != 0 && normal != 0 {
			break
		}
	}

	if normal == 0 {
		rabbit.MustPublish(&bot2.SendWarn{
			AdID:      0,
			ChannelID: in.ChannelID,
			Msg:       "no ads for you",
		})
		return false, nil
	}

	adList := []int64{}
	if promoted != 0 {
		adList = append(adList, promoted)
	}
	adList = append(adList, normal)
	logrus.Warn(adList)
	//todo send ad to user
	rabbit.MustPublish(&bot3.AdDelivery{
		AdsID:     adList,
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
