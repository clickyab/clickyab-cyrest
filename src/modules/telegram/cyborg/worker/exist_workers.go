package worker

import (
	"common/assert"
	"common/models/common"
	"modules/telegram/ad/ads"
	"modules/telegram/config"
	"modules/telegram/cyborg/bot"
	"time"

	"common/rabbit"
	"fmt"

	bot2 "modules/telegram/bot/worker"

	"github.com/Sirupsen/logrus"
)

func (mw *MultiWorker) existWorker() error {
	mw.lock.Lock()
	defer mw.lock.Unlock()
	channelIDs := []int64{}
	logrus.Debug("Exist Channel AD ...")
	var adsConf = make(map[int64][]channelDetailStat)
	m := ads.NewAdsManager()
	chads, err := m.FindBundleChannelAdActive()
	assert.Nil(err)
	if len(chads) == 0 {
		return nil
	}

	for i := range chads {
		adsConf[chads[i].ChannelID] = append(adsConf[chads[i].ChannelID], channelDetailStat{
			cliChannelAdID: chads[i].CliMessageID,
			frwrd:          chads[i].CliMessageID.Valid,
			adID:           chads[i].AdID,
			channelID:      chads[i].ChannelID,
			bundleID:       chads[i].BundleID,
			botChatID:      chads[i].BotChatID,
			targetView:     chads[i].TargetView,
			code:           chads[i].Code,
		})
		assert.True(chads[i].CliMessageID.Valid, "cli not filled")
		channelIDs = append(channelIDs, chads[i].ChannelID)
	}

	channels, err := m.FindChannelByIDs(channelIDs)
	assert.Nil(err)

	for i := range channels {
		c, err := bot.NewBotManager().FindKnownChannelByName(channels[i].Name)
		if err != nil {
			ch, err := mw.discoverChannel(channels[i].Name)
			assert.Nil(err)
			c, err = bot.NewBotManager().CreateChannelByRawData(ch)
			assert.Nil(err)
		}
		h, err := mw.getLastMessages(c.CliTelegramID, tcfg.Cfg.Telegram.LastPostChannel, 0)
		assert.Nil(err)
		/*channelDetails, err := m.FindChanDetailByChannelID(channel.ID)
		assert.Nil(err)*/

		channelAdStat := mw.existChannelAdFor(h, adsConf[channels[i].ID])
		var ChannelAdDetailArr []*ads.BundleChannelAdDetail
		var ChannelAdArr []ads.BundleChannelAd
		//var reshot bool
		now := time.Now()
		for j := range chads {
			if chads[i].ChannelID != channels[i].ID {
				continue
			}
			defaultPosition := chads[j].Position
			//set warning
			if t, ok := channelAdStat[chads[j].AdID]; !ok || t.pos > defaultPosition {
				ChannelAdDetailArr = append(ChannelAdDetailArr, &ads.BundleChannelAdDetail{
					AdID:      chads[j].AdID,
					ChannelID: chads[j].ChannelID,
					BundleID:  chads[j].BundleID,
					View:      0,
					Position:  common.NullInt64{Valid: t.pos != 0, Int64: t.pos},
					Warning:   1,
					CreatedAt: &now,
				})

				//todo if u wanna tell to publisher send reshot!? for take money??
				if chads[j].Warning >= tcfg.Cfg.Telegram.LimitCountWarning {
					ChannelAdArr = append(ChannelAdArr, ads.BundleChannelAd{
						Warning:   chads[j].Warning + 1,
						View:      chads[j].View,
						AdID:      chads[j].AdID,
						End:       common.MakeNullTime(time.Now()),
						ChannelID: chads[j].ChannelID,
						BundleID:  chads[j].BundleID,
					})

					rabbit.MustPublish(&bot2.SendWarn{
						AdID:      chads[j].AdID,
						ChannelID: chads[j].ChannelID,
						ChatID:    chads[j].BotChatID,
						Msg:       fmt.Sprintf("please remove ad for bundle %s from your channel @%s in true position \n and click on command /screenshot", chads[j].Code, chads[j].ChannelID),
					})
					continue
				}

				ChannelAdArr = append(ChannelAdArr, ads.BundleChannelAd{
					Warning:   chads[j].Warning + 1,
					View:      chads[j].View,
					AdID:      chads[j].AdID,
					ChannelID: chads[j].ChannelID,
					BundleID:  chads[j].BundleID,
				})

				rabbit.MustPublish(&bot2.SendWarn{
					AdID:      chads[j].AdID,
					ChannelID: chads[j].ChannelID,
					ChatID:    chads[j].BotChatID,
					Msg:       fmt.Sprintf("please move ads for bundle %s from your channel @%s to better position", chads[j].Code, chads[j].ChannelID),
				})
				//reshot = true

				continue
			}

			//if channelAdStat[chads[j].AdID].frwrd == true {
			//	currentView = avg
			//	//update ad
			//	assert.Nil(m.UpdateAdView(chads[j].AdID, channelAdStat[chads[j].AdID].view))
			//
			//} else {
			//	currentView = channelAdStat[chads[j].AdID].view
			//}

			currentView := channelAdStat[chads[j].AdID].view
			assert.Nil(m.UpdateAdView(chads[j].AdID, channelAdStat[chads[j].AdID].view))
			ChannelAdDetailArr = append(ChannelAdDetailArr, &ads.BundleChannelAdDetail{
				AdID:      chads[j].AdID,
				ChannelID: chads[j].ChannelID,
				BundleID:  chads[j].BundleID,
				View:      currentView,
				Position:  common.NullInt64{Valid: channelAdStat[chads[j].AdID].pos != 0, Int64: channelAdStat[chads[j].AdID].pos},
				Warning:   channelAdStat[chads[j].AdID].warning,
				CreatedAt: &now,
			})
			if currentView == 0 {
				currentView = chads[j].View
			}

			if currentView > chads[j].TargetView {

				rabbit.MustPublish(&bot2.SendWarn{
					AdID:      chads[j].AdID,
					ChannelID: chads[j].ChannelID,
					ChatID:    chads[j].BotChatID,
					Msg:       fmt.Sprintf("please take screenshot for bundle %s from your channel @%s \n because get limit view for this bundle click on command /screenshot ", chads[j].Code, chads[j].ChannelID),
				})

				ChannelAdArr = append(ChannelAdArr, ads.BundleChannelAd{
					Warning:   chads[j].Warning + channelAdStat[chads[j].AdID].warning,
					AdID:      chads[j].AdID,
					View:      currentView,
					End:       common.MakeNullTime(time.Now()),
					ChannelID: chads[j].ChannelID,
					BundleID:  chads[j].BundleID,
				})

				continue
			}

			ChannelAdArr = append(ChannelAdArr, ads.BundleChannelAd{
				Warning:   chads[j].Warning + channelAdStat[chads[j].AdID].warning,
				AdID:      chads[j].AdID,
				View:      currentView,
				ChannelID: chads[j].ChannelID,
				BundleID:  chads[j].BundleID,
			})
		}
		//if reshot {
		//	//send stop (warn message)
		//	rabbit.MustPublish(&bot2.SendWarn{
		//		AdID:      0,
		//		ChannelID: channels[i].ID,
		//		Msg:       trans.T("please reshot all ads\n%s", fmt.Sprintf("/reshot_%d", channels[i].ID)).String(),
		//		ChatID:    adsConf[channels[i].ID][0].botChatID,
		//	})
		//}

		//transaction
		_, err = mw.transaction(m, ChannelAdArr, ChannelAdDetailArr)
		//if res == true {
		//	continue
		//}
		if err != nil {
			logrus.Debug("have problem in query insert transaction ")
		}

	}
	//check for promotion to be alone or not

	return nil
}
