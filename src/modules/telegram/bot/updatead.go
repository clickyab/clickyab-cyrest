package bot

import (
	"common/assert"
	"common/models/common"
	"fmt"
	"modules/telegram/ad/ads"
	"modules/telegram/common/tgbot"
	"strconv"
	"strings"
	"time"

	"common/rabbit"
	"modules/telegram/cyborg/commands"

	"gopkg.in/telegram-bot-api.v4"
	"modules/misc/trans"
)

func (bb *bot) updateAD(bot *tgbotapi.BotAPI, m *tgbotapi.Message) {
	result := strings.Replace(m.Text, "/updatead-", "", 1)
	id, err := strconv.ParseInt(result, 0, 10)
	if err == nil {
		tgbot.RegisterUserHandler(m.Chat.ID, func(bot *tgbotapi.BotAPI, m *tgbotapi.Message) {
			defer tgbot.UnRegisterUserHandler(m.Chat.ID)

			n := ads.NewAdsManager()
			currentAd, err := n.FindAdByID(id)
			assert.Nil(err)
			currentAd.BotChatID = common.NullInt64{Valid: true, Int64: m.Chat.ID}
			currentAd.BotMessageID = common.NullInt64{Valid: true, Int64: int64(m.MessageID)}
			assert.Nil(n.UpdateAd(currentAd))

		}, time.Minute)
	}

}

func (bb *bot) wantAD(bot *tgbotapi.BotAPI, m *tgbotapi.Message) {
	//find channels
	chnManger := ads.NewAdsManager()
	if strings.Contains(m.Text, "_") {
		res := strings.Split(m.Text, "_")
		if len(res) != 2 {
			return
		}
		//find channel by chat ID and channel_name
		channel, err := chnManger.FindChannelsByChatIDName(m.Chat.ID, res[1])
		if err != nil {
			msg := tgbotapi.NewMessage(m.Chat.ID, trans.T("channel not found for you").Translate())
			msg.ParseMode = htmlMode
			_, err := bot.Send(msg)
			assert.Nil(err)
			return
		}
		//everything ok publish a job TODO:
		rabbit.MustPublish(&commands.SelectAd{
			ChannelID: channel.ID,
			ChatID:    m.Chat.ID,
		})
		fmt.Println(channel.Name)
		return
	}
	channels, err := chnManger.FindChannelsByChatID(m.Chat.ID)
	assert.Nil(err)
	if len(channels) == 0 {
		msg := tgbotapi.NewMessage(m.Chat.ID, trans.T("no channels for you").Translate())
		msg.ParseMode = htmlMode
		_, err := bot.Send(msg)
		assert.Nil(err)
		return
	}
	textMsg := trans.T("please choose one of the below channels\n").Translate()
	for i := range channels {
		textMsg += fmt.Sprintf("/ad_%s\n", channels[i].Name)
	}
	msg := tgbotapi.NewMessage(m.Chat.ID, textMsg)
	msg.ParseMode = htmlMode
	_, err = bot.Send(msg)
	assert.Nil(err)
	return

}
