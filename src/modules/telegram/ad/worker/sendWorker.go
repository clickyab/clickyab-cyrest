package worker

import (
	"common/assert"
	"modules/telegram/ad/ads"
	"modules/telegram/bot"
	"modules/telegram/common/tgbot"

	"fmt"
	bot2 "modules/telegram/bot/worker"
	"modules/telegram/config"

	"common/rabbit"

	"gopkg.in/telegram-bot-api.v4"
)

//AdDelivery is type of ad delivery
type AdDelivery struct {
	ChatID    int64
	AdsID     []int64
	ChannelID int64
}

// GetTopic return this message topic
func (AdDelivery) GetTopic() string {
	return "cy.rubik.AdDelivery"
}

// GetQueue is the request queue
func (AdDelivery) GetQueue() string {
	return "cy_rubik_AdDelivery"
}

//AdDeliveryAction is a function that send ad and channel data and metadata
func AdDeliveryAction(in *AdDelivery) (bool, error) {

	for adID := range in.AdsID {
		adManager := ads.NewAdsManager()
		ad, err := adManager.FindAdByID(in.AdsID[adID])

		if err != nil {
			continue
		}

		rabbit.MustPublish(bot2.SendWarn{
			ChannelID: in.ChannelID,
			AdID:      0,
			ChatID:    in.ChatID,
			Msg:       "please forward the following ad to your channel",
		})
		res := bot.RenderMessage(tgbot.GetBot(), in.ChatID, ad)
		msgx := fmt.Sprintf("please forward these as to your channel and press done otherwise press reject\n/done_%[1]d\n/reject_%[1]d\n", in.ChannelID)
		userMsg := tgbotapi.NewMessage(in.ChatID, msgx)
		userMsg.ParseMode = "HTML"
		_, err = tgbot.Send(userMsg)
		if err != nil {
			continue
		}
		if !ad.CliMessageID.Valid {
			fwd := tgbotapi.NewForward(tcfg.Cfg.Telegram.ChannelID, res.Chat.ID, res.MessageID)
			_, err := tgbot.Send(fwd)
			if err != nil {
				continue
			}

			msgTxt := fmt.Sprintf("%d/%d", ad.ID, in.ChannelID)
			msg := tgbotapi.NewMessage(tcfg.Cfg.Telegram.ChannelID, msgTxt)
			_, err = tgbot.Send(msg)
			if err != nil {
				continue
			}
		}
		var cha ads.ChannelAd
		cha.ChannelID = in.ChannelID
		cha.CliMessageID = ad.CliMessageID
		cha.AdID = ad.ID
		cha.BotChatID = res.Chat.ID
		cha.BotMessageID = res.MessageID
		cha.Active = ads.ActiveStatusNo
		cha.View = 0

		err = ads.NewAdsManager().CreateChannelAd(&cha)
		assert.Nil(err)

	}
	return false, nil
}
