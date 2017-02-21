package worker

import (
	"common/assert"
	"common/models/common"
	"common/rabbit"
	"modules/telegram/ad/ads"
	bot2 "modules/telegram/bot/worker"
	"modules/telegram/common/tgo"
	"modules/telegram/config"
	"modules/telegram/cyborg/bot"
	"modules/telegram/cyborg/commands"
	"time"

	"fmt"
	"modules/misc/trans"

	"github.com/Sirupsen/logrus"
)

func (mw *MultiWorker) existChannelAdFor(h []tgo.History, adConfs []channelDetailStat) (map[int64]channelViewStat, int64) {
	var finalResult = make(map[int64]channelViewStat)
	var sumIndividualView int64
	var countIndividual int64
	var found int
	historyLen := len(h)
bigloop:
	for k := historyLen - 1; k >= 0; k-- {
		if h[k].Event == "message" && h[k].Service == false {
			if h[k].FwdFrom != nil {
				for i := range adConfs {
					if h[k].ID == adConfs[i].cliChannelAdID.String {
						finalResult[adConfs[i].adID] = channelViewStat{
							view:    int64(h[k].Views),
							warning: 0,
							pos:     int64(historyLen - k),
							frwrd:   adConfs[i].frwrd,
							adID:    adConfs[i].adID,
						}
						found++
						if !adConfs[i].frwrd { //the ad is  not forward type
							sumIndividualView += int64(h[k].Views)
							countIndividual++
						}
					}
					if found == len(adConfs) {
						break bigloop
					}
				}
			}
		}
	}

	for i := range adConfs {
		if _, ok := finalResult[adConfs[i].adID]; !ok {
			logrus.Infof("%+v", finalResult[adConfs[i].adID])
			finalResult[adConfs[i].adID] = channelViewStat{
				view:    0,
				warning: 1,
				frwrd:   adConfs[i].frwrd,
				adID:    adConfs[i].adID,
				pos:     0,
			}
		}
	}

	if countIndividual == 0 {
		return finalResult, 0
	}
	logrus.Warnf("%+v", finalResult, sumIndividualView, countIndividual)
	return finalResult, (sumIndividualView) / (countIndividual)
}

func (mw *MultiWorker) existChannelAd(in *commands.ExistChannelAd) (bool, error) {
	var adsConf []channelDetailStat
	m := ads.NewAdsManager()

	chads, err := m.FindChannelAdByChannelIDActive(in.ChannelID)
	if len(chads) == 0 {
		return false, nil
	}
	assert.Nil(err)
	for i := range chads {
		adsConf = append(adsConf, channelDetailStat{
			cliChannelAdID: chads[i].CliMessageID,
			frwrd:          chads[i].CliMessageAd.Valid,
			adID:           chads[i].AdID,
		})
		assert.True(chads[i].CliMessageID.Valid, "cli not filled")
	}

	//check for promotion to be alone or not
	var promotionCount int
	var individualCount int
	for adConf := range adsConf {
		if adsConf[adConf].frwrd {
			promotionCount++
		}
		individualCount++
	}

	if individualCount == 0 {

		for adConf := range adsConf {
			//send stop (warn message)
			rabbit.MustPublish(&bot2.SendWarn{
				AdID:      adsConf[adConf].adID,
				ChannelID: in.ChannelID,
				Msg:       "please remove the following ad",
			})

		}
		return false, nil
	}

	defer rabbit.PublishAfter(&commands.ExistChannelAd{
		ChannelID: in.ChannelID,
		ChatID:    in.ChatID,
	}, tcfg.Cfg.Telegram.TimeReQueUe)

	channel, err := m.FindChannelByID(in.ChannelID)
	assert.Nil(err)
	c, err := bot.NewBotManager().FindKnownChannelByName(channel.Name)
	if err != nil {
		ch, err := mw.discoverChannel(channel.Name)
		assert.Nil(err)
		c, err = bot.NewBotManager().CreateChannelByRawData(ch)
		assert.Nil(err)
	}
	h, err := mw.getLastMessages(c.CliTelegramID, tcfg.Cfg.Telegram.LastPostChannel, 0)
	assert.Nil(err)
	/*channelDetails, err := m.FindChanDetailByChannelID(channel.ID)
	assert.Nil(err)*/

	channelAdStat, avg := mw.existChannelAdFor(h, adsConf)
	var ChannelAdDetailArr []*ads.ChannelAdDetail
	var ChannelAdArr []ads.ChannelAd
	var reshot bool
	for j := range chads {
		defaultPosition := chads[j].PlanPosition
		if t, ok := channelAdStat[chads[j].AdID]; !ok || t.pos > defaultPosition {
			ChannelAdDetailArr = append(ChannelAdDetailArr, &ads.ChannelAdDetail{
				AdID:      chads[j].AdID,
				ChannelID: chads[j].ChannelID,
				View:      0,
				Position:  common.NullInt64{Valid: t.pos != 0, Int64: t.pos},
				Warning:   1,
				CreatedAt: time.Now(),
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
			CreatedAt: time.Now(),
		})
		ChannelAdArr = append(ChannelAdArr, ads.ChannelAd{

			Warning:   chads[j].Warning + channelAdStat[chads[j].AdID].warning,
			View:      currentView,
			AdID:      chads[j].AdID,
			ChannelID: chads[j].ChannelID,
		})

	}
	if reshot {
		//send stop (warn message)
		rabbit.MustPublish(&bot2.SendWarn{
			AdID:      0,
			ChannelID: in.ChannelID,
			Msg:       trans.T("please reshot all ads\n%s", fmt.Sprintf("/reshot_%d", in.ChannelID)).String(),
			ChatID:    in.ChatID,
		})
	}

	//transaction
	res, err := mw.transaction(m, ChannelAdArr, ChannelAdDetailArr, avg)
	if res == true {
		return true, err
	}

	return false, nil
}
