package worker

import (
	"common/assert"
	"common/config"
	"common/models/common"
	"common/rabbit"
	"common/utils"
	"encoding/json"
	"modules/misc/trans"
	"modules/telegram/ad/ads"
	bot2 "modules/telegram/bot/worker"
	"modules/telegram/common/tgo"
	"modules/telegram/config"
	"modules/telegram/cyborg/bot"
	"modules/telegram/cyborg/commands"
	"time"
)

func (mw *MultiWorker) discoverAd(in *commands.DiscoverAd) (bool, error) {
	adsManager := ads.NewAdsManager()
	chads, err := adsManager.FindChannelAdByChannelID(in.Channel)
	assert.Nil(err)

	// first try to resolve the channel
	m := bot.NewBotManager()
	channel, err := adsManager.FindChannelByID(in.Channel)
	assert.Nil(err)
	c, err := m.FindKnownChannelByName(channel.Name)
	if err != nil {
		ch, err := mw.discoverChannel(channel.Name)
		assert.Nil(err)
		c, err = bot.NewBotManager().CreateChannelByRawData(ch)
		assert.Nil(err)
	}
	h, err := mw.getLastMessages(c.CliTelegramID, 10, 0)
	assert.Nil(err)
	// then discover the messages
	found := 0
	var cha []*ads.ChannelAd
bigLoop:
	for i := range chads {
		if chads[i].CliMessageID.Valid {
			found++
			continue
		}
		cha = append(cha, &ads.ChannelAd{
			AdID:         chads[i].AdID,
			ChannelID:    chads[i].ChannelID,
			BotChatID:    chads[i].BotChatID,
			BotMessageID: chads[i].BotMessageID,
			CliMessageID: chads[i].CliMessageID,
			CreatedAt:    chads[i].CreatedAt,
			PossibleView: chads[i].PossibleView,
			View:         chads[i].View,
			Warning:      chads[i].Warning,
			Start:        common.MakeNullTime(time.Now()),
			Active:       ads.ActiveStatusYes,
		})
		var msg *tgo.History
		if chads[i].CliMessageAd.Valid {
			msg = &tgo.History{}
			assert.Nil(json.Unmarshal([]byte(chads[i].PromoteData.String), msg))

		}
		for j := len(h) - 1; j >= 0; j-- {
			// Promotion or individual?
			if h[j].FwdFrom == nil {
				continue
			}
			if msg != nil {
				// promotion
				if comparePromotion(*msg, h[j]) {
					adsManager.SetCLIMessageID(chads[i].ChannelID, chads[i].AdID, h[j].ID)
					found++
					continue bigLoop
				}
			} else {
				if compareIndividual(chads[i], h[j]) {
					adsManager.SetCLIMessageID(chads[i].ChannelID, chads[i].AdID, h[j].ID)
					found++
					continue bigLoop
				}
			}
		}
	}
	// I think we are gonna die, or the user is stupid. replace us, or him/her
	if found != len(chads) {

		data, _ := json.MarshalIndent(struct {
			History []tgo.History
			Chads   []ads.ChannelAdD
		}{
			History: h,
			Chads:   chads,
		}, "\t", "\t")
		if config.Config.Slack.Active {
			go utils.SlackDoMessage(
				"[BUG/USER] can not find messages in channel, nothing special but please check if the user is stupid or we are?",
				":thinking_face:",
				utils.SlackAttachment{Text: string(data), Color: "#AA3939"},
			)
		}
		rabbit.MustPublish(&bot2.SendWarn{
			AdID:      0,
			ChannelID: in.Channel,
			Msg:       trans.T("cant find your ad please make sure the ad is in your channel and press done").String(),
			ChatID:    in.ChatID,
		})
	} else {
		err = adsManager.UpdateActiveEndChannelAds(cha)
		assert.Nil(err)
		rabbit.MustPublishAfter(
			commands.ExistChannelAd{
				ChannelID: channel.ID,
				ChatID:    in.ChatID,
			},
			tcfg.Cfg.Telegram.TimeReQueUe,
		)
		//send ok message
		rabbit.MustPublish(&bot2.SendWarn{
			AdID:      0,
			ChannelID: in.Channel,
			Msg:       trans.T("your add has been successfully activated\nthanks for your cooperation").String(),
			ChatID:    in.ChatID,
		})
	}
	return false, nil

}
