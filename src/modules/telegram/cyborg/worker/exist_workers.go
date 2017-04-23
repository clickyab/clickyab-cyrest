package worker

import (
	"common/assert"
	"common/models/common"
	"common/rabbit"
	"fmt"
	"modules/misc/trans"
	"modules/telegram/ad/ads"
	bot2 "modules/telegram/bot/worker"
	"modules/telegram/config"
	"modules/telegram/cyborg/bot"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
)

var existLock sync.Mutex

func (mw *MultiWorker) existWorker() error {
	existLock.Lock()
	defer existLock.Unlock()
	channelIDs := []int64{}
	logrus.Debug("Exist Channel AD ...")
	var adsConf = make(map[int64][]channelDetailStat)
	m := ads.NewAdsManager()
	chads, err := m.FindChannelAdActive()
	assert.Nil(err)
	if len(chads) == 0 {
		return nil
	}

	for i := range chads {
		adsConf[chads[i].ChannelID] = append(adsConf[chads[i].ChannelID], channelDetailStat{
			cliChannelAdID: chads[i].CliMessageID,
			frwrd:          chads[i].CliMessageAd.Valid,
			adID:           chads[i].AdID,
			channelID:      chads[i].ChannelID,
			botChatID:      chads[i].BotChatID,
		})
		assert.True(chads[i].CliMessageID.Valid, "cli not filled")
		channelIDs = append(channelIDs, chads[i].ChannelID)
	}

	channels, err := m.FindChannelByIDs(channelIDs)
	assert.Nil(err)

	for i := range channels {
		var promotionCount int
		var individualCount int
		for adCon := range adsConf[channels[i].ID] {
			if adsConf[channels[i].ID][adCon].frwrd {
				promotionCount++
			}
			individualCount++
		}

		if individualCount == 0 {

			for adCon := range adsConf[channels[i].ID] {
				//send stop (warn message)
				rabbit.MustPublish(&bot2.SendWarn{
					AdID:      adsConf[channels[i].ID][adCon].adID,
					ChannelID: adsConf[channels[i].ID][adCon].channelID,
					Msg:       "please remove the following ad",
				})

			}
			continue
		}

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

		channelAdStat, avg := mw.existChannelAdFor(channels[i].ID, adsConf[channels[i].ID][0].botChatID, h, adsConf[channels[i].ID])
		var ChannelAdDetailArr []*ads.ChannelAdDetail
		var ChannelAdArr []ads.ChannelAd
		var reshot bool
		now := time.Now()
		for j := range chads {
			defaultPosition := chads[j].PlanPosition
			if t, ok := channelAdStat[chads[j].AdID]; !ok || t.pos > defaultPosition {
				ChannelAdDetailArr = append(ChannelAdDetailArr, &ads.ChannelAdDetail{
					AdID:      chads[j].AdID,
					ChannelID: chads[j].ChannelID,
					View:      0,
					Position:  common.NullInt64{Valid: t.pos != 0, Int64: t.pos},
					Warning:   1,
					CreatedAt: &now,
				})
				ChannelAdArr = append(ChannelAdArr, ads.ChannelAd{

					Warning:   chads[j].Warning + 1,
					View:      chads[j].View,
					AdID:      chads[j].AdID,
					ChannelID: chads[j].ChannelID,
				})
				reshot = true

				continue
			}
			var currentView int64
			if channelAdStat[chads[j].AdID].frwrd == true {
				currentView = avg
				//update ad
				assert.Nil(m.UpdateAdView(chads[j].AdID, channelAdStat[chads[j].AdID].view))

			} else {
				currentView = channelAdStat[chads[j].AdID].view
			}
			ChannelAdDetailArr = append(ChannelAdDetailArr, &ads.ChannelAdDetail{
				AdID:      chads[j].AdID,
				ChannelID: chads[j].ChannelID,
				View:      currentView,
				Position:  common.NullInt64{Valid: channelAdStat[chads[j].AdID].pos != 0, Int64: channelAdStat[chads[j].AdID].pos},
				Warning:   channelAdStat[chads[j].AdID].warning,
				CreatedAt: &now,
			})
			if currentView == 0 {
				ChannelAdArr = append(ChannelAdArr, ads.ChannelAd{

					Warning:   chads[j].Warning + channelAdStat[chads[j].AdID].warning,
					AdID:      chads[j].AdID,
					View:      chads[j].View,
					ChannelID: chads[j].ChannelID,
				})
			} else {
				ChannelAdArr = append(ChannelAdArr, ads.ChannelAd{

					Warning:   chads[j].Warning + channelAdStat[chads[j].AdID].warning,
					View:      currentView,
					AdID:      chads[j].AdID,
					ChannelID: chads[j].ChannelID,
				})
			}

		}
		if reshot {
			//send stop (warn message)
			rabbit.MustPublish(&bot2.SendWarn{
				AdID:      0,
				ChannelID: channels[i].ID,
				Msg:       trans.T("please reshot all ads\n%s", fmt.Sprintf("/reshot_%d", channels[i].ID)).String(),
				ChatID:    adsConf[channels[i].ID][0].botChatID,
			})
		}

		//transaction
		res, _ := mw.transaction(m, ChannelAdArr, ChannelAdDetailArr, avg)
		if res == true {
			continue
		}

	}
	//check for promotion to be alone or not

	return nil
}
