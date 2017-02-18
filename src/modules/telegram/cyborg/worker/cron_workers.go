package worker

import (
	"common/assert"
	"common/models/common"
	"common/rabbit"
	"modules/misc/trans"
	"modules/telegram/ad/ads"
	bot2 "modules/telegram/bot/worker"
	"modules/telegram/config"
	"time"
)

// CronReview cron review for finished ads
func (mw *MultiWorker) cronReview() error {
	m := ads.NewAdsManager()
	m.UpdateIndividualViewCount()

	finishedAds := m.FinishedActiveAds()
	// TODO : transaction
	for key := range finishedAds {
		finishedAds[key].AdActiveStatus = ads.AdActiveStatusNo
		assert.Nil(m.UpdateAd(&finishedAds[key]))
		channelAd, err := m.FindChannelAdActiveByAdID(finishedAds[key].ID, ads.ActiveStatusYes)
		assert.Nil(err)
		for c := range channelAd {
			//todo send message to admin channel
			channelAd[c].ChannelAd.Active = ads.ActiveStatusNo
			channelAd[c].ChannelAd.End = common.MakeNullTime(time.Now())
			assert.Nil(m.UpdateChannelAd(&channelAd[c].ChannelAd))
			str := trans.T("you have have <%d> in channel \n please remove ad from <%s> channel", channelAd[c].ChannelAd.View, channelAd[c].Channel.Name)
			rabbit.MustPublish(&bot2.SendWarn{
				AdID:      channelAd[c].ChannelAd.AdID,
				ChannelID: channelAd[c].ChannelAd.ChannelID,
				Msg:       str.String(),
				ChatID:    channelAd[c].ChannelAd.BotChatID,
			})
		}
	}

	for _, chad := range m.GetWarningLimited(tcfg.Cfg.Telegram.LimitCountWarning) {
		chad.Active = ads.ActiveStatusNo
		chad.End = common.MakeNullTime(time.Now())
		assert.Nil(m.UpdateChannelAd(&chad))
		str := trans.T("Sorry, but since there is a lot of error, please remove the ads from your channel")
		rabbit.MustPublish(&bot2.SendWarn{
			AdID:      chad.AdID,
			ChannelID: chad.ChannelID,
			Msg:       str.String(),
			ChatID:    chad.BotChatID,
		})
	}

	return nil
}
