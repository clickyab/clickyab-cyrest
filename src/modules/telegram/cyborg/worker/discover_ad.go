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

	"github.com/Sirupsen/logrus"
)

func (mw *MultiWorker) discoverAd(in *commands.DiscoverAd) (bool, error) {
	logrus.Warn("discover ad 1")
	adsManager := ads.NewAdsManager()
	var chads []ads.ChannelAdD
	var err error
	if in.Reshot {
		chads, err = adsManager.FindReshotChannelAdByChannelID(in.Channel)
		assert.Nil(err)
	} else {
		chads, err = adsManager.FindChannelAdByChannelID(in.Channel)
		assert.Nil(err)
	}

	// first try to resolve the channel
	m := bot.NewBotManager()
	channel, err := adsManager.FindChannelByID(in.Channel)
	assert.Nil(err)
	c, err := m.FindKnownChannelByName(channel.Name)
	if err != nil {
		ch, err := mw.discoverChannel(channel.Name)
		if err != nil {
			rabbit.MustPublish(&bot2.SendWarn{
				AdID:      0,
				ChannelID: channel.ID,
				Msg:       trans.T("cant find your channel").String(),
				ChatID:    in.ChatID,
			})
			return false, nil
		}
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
		logrus.Warn("discover ad")
		if chads[i].CliMessageID.Valid && !in.Reshot {
			found++
			continue
		}
		logrus.Warn("discover ad")
		cha = append(cha, &ads.ChannelAd{
			AdID:         chads[i].AdID,
			ChannelID:    chads[i].ChannelID,
			BotChatID:    chads[i].BotChatID,
			BotMessageID: chads[i].BotMessageID,
			CreatedAt:    &chads[i].CreatedAt,
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
					err = adsManager.SetCLIMessageID(chads[i].ChannelID, chads[i].AdID, h[j].ID)
					assert.Nil(err)
					cha[i].CliMessageID = common.MakeNullString(h[j].ID)
					found++
					continue bigLoop
				}
			} else {
				if compareIndividual(chads[i], h[j]) {
					err = adsManager.SetCLIMessageID(chads[i].ChannelID, chads[i].AdID, h[j].ID)
					assert.Nil(err)
					cha[i].CliMessageID = common.MakeNullString(h[j].ID)
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
		if in.Reshot {
			rabbit.MustPublish(&bot2.SendWarn{
				AdID:      0,
				ChannelID: in.Channel,
				Msg:       trans.T("your reshot process begin\nthanks for your cooperation").String(),
				ChatID:    in.ChatID,
			})
		} else {
			rabbit.MustPublish(&bot2.SendWarn{
				AdID:      0,
				ChannelID: in.Channel,
				Msg:       trans.T("your add has been successfully activated\nthanks for your cooperation").String(),
				ChatID:    in.ChatID,
			})
		}

	}
	return false, nil

}
