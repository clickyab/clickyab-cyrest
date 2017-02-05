package bot

import (
	"modules/telegram/ad/ads"

	"modules/telegram/common/tgbot"

	"fmt"

	"gopkg.in/telegram-bot-api.v4"
)

// SendWarn is the command for
type SendWarn struct {
	// The msg
	Msg string
	// The channel ID
	ChannelID int64
	// AdID
	AdID int64
	// ChatID
	ChatID int64
}

// GetTopic return this message topic
func (SendWarn) GetTopic() string {
	return "cy.rubik.sendWarn"
}

// GetQueue is the request queue
func (SendWarn) GetQueue() string {
	return "cy_rubik_sendWarn"
}

// SendWarnAction worker
func SendWarnAction(in *SendWarn) (bool, error) {
	m := ads.NewAdsManager()
	//find channel
	channel, err := m.FindChannelByID(in.ChannelID)
	if err != nil {
		return false, err
	}
	baseMSg := fmt.Sprintf("Dear Admin of the <b>%s</b> channel:\n", channel.Name)
	x := tgbotapi.NewMessage(in.ChatID, fmt.Sprintf("%s%s", baseMSg, in.Msg))
	tgbot.Send(x)
	if in.AdID != 0 { //forward the ad
		channelAd, err := m.FindChannelIDAdByAdID(in.ChannelID, in.AdID)
		if err != nil {
			return false, err
		}
		msg := tgbotapi.NewForward(in.ChatID, channelAd.BotChatID, int(channelAd.BotMessageID))
		tgbot.Send(msg)
	}
	return false, nil
}