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

	"github.com/Sirupsen/logrus"
)

func (mw *MultiWorker) existChannelAdFor(h []tgo.History, adConfs []channelDetailStat) (map[int64]channelViewStat, int64) {
	var finalResult = make(map[int64]channelViewStat)
	var sumNotpromotionView int64
	var countNotPromotion int64
	historyLen := len(h)
	for k := range h {
		if h[k].Event == "message" && h[k].Service == false {
			if h[k].FwdFrom != nil {
				for i := range adConfs {
					if adConfs[i].frwrd && h[k].ID == adConfs[i].cliID.String { //the ad is forward type
						finalResult[adConfs[i].adID] = channelViewStat{
							view:    int64(h[k].Views),
							warning: 0,
							pos:     int64(historyLen - k),
							frwrd:   true,
							adID:    adConfs[i].adID,
						}
					} else if !adConfs[i].frwrd && h[k].ID == adConfs[i].cliID.String { //the ad is  not forward type
						finalResult[adConfs[i].adID] = channelViewStat{
							view:    int64(h[k].Views),
							warning: 0,
							pos:     int64(historyLen - k),
							frwrd:   false,
							adID:    adConfs[i].adID,
						}
						sumNotpromotionView += int64(h[k].Views)
						countNotPromotion++
					} else { //don't find ad in the history
						finalResult[adConfs[i].adID] = channelViewStat{
							view:    0,
							warning: 1,
							frwrd:   adConfs[i].frwrd,
							adID:    adConfs[i].adID,
							pos:     0,
						}
					}
				}
			}
		}
	}
	if countNotPromotion == 0 {
		return finalResult, 0
	}
	return finalResult, (sumNotpromotionView) / (countNotPromotion)
}

func (mw *MultiWorker) existChannelAd(in *commands.ExistChannelAd) (bool, error) {
	logrus.Warn("existChannelAdFor")
	var adsConf []channelDetailStat
	m := ads.NewAdsManager()

	chads, err := m.FindChannelAdByChannelIDActive(in.ChannelID)
	assert.Nil(err)
	for i := range chads {
		adsConf = append(adsConf, channelDetailStat{
			cliID: chads[i].CliMessageID,
			frwrd: chads[i].CliMessageAd.Valid,
			adID:  chads[i].AdID,
		})
		assert.True(chads[i].CliMessageID.Valid, "cli not filled")
		if !chads[i].CliMessageID.Valid {
			rabbit.PublishAfter(&commands.ExistChannelAd{
				ChannelID: in.ChannelID,
				ChatID:    in.ChatID,
			}, tcfg.Cfg.Telegram.TimeReQueUe)
			return false, nil
		}
	}

	//check for promotion to be alone or not
	var promotionCount int
	var notPromotionCount int
	for adConf := range adsConf {
		if adsConf[adConf].frwrd {
			promotionCount++
		}
		notPromotionCount++
	}

	if notPromotionCount == 0 {

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
	for j := range chads {
		var currentView int64
		depos := tcfg.Cfg.Telegram.PositionAdDefault
		if chads[j].AdPosition.Valid {
			depos = chads[j].AdPosition.Int64
		}
		if channelAdStat[chads[j].AdID].pos < depos {
			channelAdStat[chads[j].AdID] = channelViewStat{
				warning: 1,
				adID:    chads[j].AdID,
				frwrd:   channelAdStat[chads[j].AdID].frwrd,
				pos:     channelAdStat[chads[j].AdID].pos,
				view:    channelAdStat[chads[j].AdID].view,
			}
		}
		if channelAdStat[chads[j].AdID].frwrd == true {
			currentView = avg
		} else {
			currentView = channelAdStat[chads[j].AdID].view
		}
		ChannelAdDetailArr = append(ChannelAdDetailArr, &ads.ChannelAdDetail{
			AdID:      chads[j].AdID,
			ChannelID: chads[j].ChannelID,
			View:      currentView,
			Position:  common.NullInt64{Valid: channelAdStat[chads[j].AdID].pos != 0, Int64: channelAdStat[chads[j].AdID].pos},
			Warning:   channelAdStat[chads[j].AdID].warning,
		})
	}

	var ChannelAdArr []ads.ChannelAd

	for chad := range chads {
		var currentView int64
		if channelAdStat[chads[chad].AdID].frwrd == true {
			currentView = avg
		} else {
			currentView = channelAdStat[chads[chad].AdID].view
		}
		ChannelAdArr = append(ChannelAdArr, ads.ChannelAd{

			Warning:   chads[chad].Warning + channelAdStat[chads[chad].AdID].warning,
			View:      currentView,
			AdID:      chads[chad].AdID,
			ChannelID: chads[chad].ChannelID,
		})
		if chads[chad].Warning >= tcfg.Cfg.Telegram.LimitCountWarning {
			//send stop (warn message)
			rabbit.MustPublish(&bot2.SendWarn{
				AdID:      chads[chad].AdID,
				ChannelID: in.ChannelID,
				Msg:       "please reshot the following ad",
			})
			ChannelAdArr = append(ChannelAdArr, ads.ChannelAd{
				End: common.NullTime{Valid: true, Time: time.Now()},
			})
		}
	}

	//transaction
	res, err := mw.transaction(m, ChannelAdArr, ChannelAdDetailArr, avg)
	if res == true {
		return true, err
	}

	return false, nil
}
