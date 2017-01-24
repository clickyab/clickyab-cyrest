package worker

import (
	"common/assert"
	"modules/telegram/ad/ads"
	"modules/telegram/ad/bot"
	"modules/telegram/common/tgbot"

	"fmt"

	"modules/telegram/config"

	"gopkg.in/telegram-bot-api.v4"
)

//AdDelivery is type of ad delivery
type AdDelivery struct {
	chatID    int64
	adsID     []int64
	channelID int64
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

	for adID := range in.adsID {
		adManager := ads.NewAdsManager()
		ad, err := adManager.FindAdByID(in.adsID[adID])

		if err != nil {
			continue
		}
		res := bot.RenderMessage(tgbot.GetBot(), in.chatID, ad)
		if !ad.CliMessageID.Valid {
			fwd := tgbotapi.NewForward(tcfg.Cfg.Telegram.ChannelID, res.Chat.ID, res.MessageID)
			_, err := tgbot.Send(fwd)
			if err != nil {
				continue
			}

			msgTxt := fmt.Sprintf("%d/%d", ad.ID, in.channelID)
			msg := tgbotapi.NewMessage(tcfg.Cfg.Telegram.ChannelID, msgTxt)
			_, err = tgbot.Send(msg)
			if err != nil {
				continue
			}
		}
		var cha ads.ChannelAd
		cha.ChannelID = in.channelID
		cha.CliMessageID = ad.CliMessageID
		cha.AdID = ad.ID
		cha.BotChatID = res.Chat.ID
		cha.BotMessageID = res.MessageID

		err = ads.NewAdsManager().CreateChannelAd(&cha)
		assert.Nil(err)

	}
	return false, nil
}
