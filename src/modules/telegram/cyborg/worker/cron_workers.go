package worker

import (
	"common/assert"
	"common/models/common"
	"fmt"
	"modules/telegram/ad/ads"
	bot2 "modules/telegram/bot/worker"
	"time"

	"github.com/Sirupsen/logrus"
)

// CronReview cron review for finished ads
func (mw *MultiWorker) cronReview() error {
	logrus.Debug("SSS")
	m := ads.NewAdsManager()
	activeIndividualAds, err := m.SelectIndividualActiveAd()
	assert.Nil(err)
	if len(activeIndividualAds) <= 0 {
		return nil
	}

	for k := range activeIndividualAds {
		activeIndividualAds[k].Ad.View.Int64 = activeIndividualAds[k].Viewed
		activeIndividualAds[k].Ad.View.Valid = true
		err = m.UpdateAd(&activeIndividualAds[k].Ad)
		assert.Nil(err)
	}
	allActiveAds, err := m.SelectAdsPlan()
	assert.Nil(err)
	for key := range allActiveAds {
		if allActiveAds[key].Viewed < allActiveAds[key].Ad.View.Int64 {
			channelAd, err := m.FindChannelAdActiveByAdID(allActiveAds[key].Ad.ID, ads.ActiveStatusYes)
			assert.Nil(err)
			for c := range channelAd {
				//todo send message to admin channel
				channelAd[c].ChannelAd.Active = ads.ActiveStatusNo
				channelAd[c].ChannelAd.End = common.MakeNullTime(time.Now())
				assert.Nil(m.UpdateChannelAd(&channelAd[c].ChannelAd))
				str := fmt.Sprintf("you have have <%d> in channel \n please remove ad from <%s> channel", channelAd[c].ChannelAd.View, channelAd[c].Channel.Name)
				bot2.SendWarnAction(&bot2.SendWarn{
					AdID:      channelAd[c].ChannelAd.AdID,
					ChannelID: channelAd[c].ChannelAd.ChannelID,
					Msg:       str,
					ChatID:    channelAd[c].ChannelAd.BotChatID,
				})

			}
		}
	}

	return nil

}
