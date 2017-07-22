package worker

import (
	"common/assert"
	"common/config"
	"common/mail"
	"common/models/common"
	"common/rabbit"
	"modules/billing/bil"
	"modules/misc/trans"
	"modules/telegram/ad/ads"
	bot2 "modules/telegram/bot/worker"
	"modules/telegram/config"
	"modules/user/aaa"
	"time"
)

// CronReview cron review for finished ads
func (mw *MultiWorker) cronReview() error {
	m := ads.NewAdsManager()
	m.UpdateIndividualViewCount()

	finishedAds := m.FinishedActiveAds()
	// TODO : transaction
	for key := range finishedAds {
		ad := finishedAds[key].Ad
		ad.AdActiveStatus = ads.ActiveStatusNo
		assert.Nil(m.UpdateAd(&ad))
		ca := m.FinishedActiveChannels(finishedAds[key].ID, tcfg.Cfg.Telegram.LimitCountWarning)
		b := bil.NewBilManager()
		err := b.ChannelAdBilling(ca, finishedAds[key])
		assert.Nil(err)

		channelAd, err := m.FindChannelAdActiveByAdID(finishedAds[key].ID, ads.ActiveStatusYes)
		assert.Nil(err)
		for c := range channelAd {
			//todo send message to admin channel
			channelAds := m.FindActiveChannelAdByChannelID(channelAd[c].ChannelID)
			for i := range channelAds {
				channelAds[i].Active = ads.ActiveStatusNo
				channelAds[i].End = common.MakeNullTime(time.Now())
				err = m.UpdateChannelAd(&channelAds[i])
				assert.Nil(err)
				str := trans.T("you have %d view in channel \n please remove all ads from your channel", channelAds[i].View)
				rabbit.MustPublish(&bot2.SendWarn{
					AdID:      0,
					ChannelID: channelAds[i].ChannelID,
					Msg:       str.String(),
					ChatID:    channelAds[i].BotChatID,
				})
			}
		}
		//send mail
		currentUser, err := aaa.NewAaaManager().FindUserByID(finishedAds[key].UserID)
		assert.Nil(err)
		go func() {
			err = mail.SendByTemplateName(trans.T("your ad is end").Translate("fa_IR"), "end-ad", struct {
				Name      string
				Ad        string
				StartDate time.Time
				EndDate   time.Time
				EndTime   int
			}{
				Name:      currentUser.Email,
				Ad:        finishedAds[key].Name,
				StartDate: *finishedAds[key].CreatedAt,
				EndDate:   *finishedAds[key].UpdatedAt,
				EndTime:   finishedAds[key].UpdatedAt.Hour(),
			}, config.Config.Mail.From, currentUser.Email)
			assert.Nil(err)
		}()
	}

	for _, chad := range m.GetWarningLimited(tcfg.Cfg.Telegram.LimitCountWarning) {
		chad.Active = ads.ActiveStatusNo
		chad.End = common.MakeNullTime(time.Now())
		assert.Nil(m.UpdateChannelAd(&chad))
		str := trans.T("Sorry, but since there is a lot of error, please remove the ads from your channel")
		rabbit.MustPublish(&bot2.SendWarn{
			AdID:      0,
			ChannelID: chad.ChannelID,
			Msg:       str.String(),
			ChatID:    chad.BotChatID,
		})
	}

	return nil
}
