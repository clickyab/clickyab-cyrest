package bot

import (
	"common/assert"
	"common/initializer"
	"common/models/common"
	"modules/telegram/ad/ads"
	"modules/telegram/common/tgbot"
	"strconv"
	"strings"
	"time"

	"fmt"
	"modules/telegram/channel/chn"

	"gopkg.in/telegram-bot-api.v4"
)

type bot struct {
}

const htmlMode = "HTML"

func (bb *bot) updateAD(bot *tgbotapi.BotAPI, m *tgbotapi.Message) {
	result := strings.Replace(m.Text, "/updatead-", "", 1)
	id, err := strconv.ParseInt(result, 0, 10)
	if err == nil {
		tgbot.RegisterUserHandler(m.Chat.ID, func(bot *tgbotapi.BotAPI, m *tgbotapi.Message) {
			defer tgbot.UnRegisterUserHandler(m.Chat.ID)
			botChatID := strconv.FormatInt(m.Chat.ID, 10)
			botMsgID := strconv.Itoa(m.MessageID)
			n := ads.NewAdsManager()
			currentAd, err := n.FindAdByID(id)
			assert.Nil(err)
			currentAd.BotChatID = common.MakeNullString(botChatID)
			currentAd.BotMessageID = common.MakeNullString(botMsgID)
			assert.Nil(n.UpdateAd(currentAd))

		}, time.Minute)
	}

}

func (bb *bot) wantAD(bot *tgbotapi.BotAPI, m *tgbotapi.Message) {
	//find channels
	chnManger := chn.NewChnManager()
	if strings.Contains(m.Text, "_") {
		res := strings.Split(m.Text, "_")
		if len(res) != 2 {
			return
		}
		//find channel by chat ID and channel_name
		channel, err := chnManger.FindChannelsByChatIDName(m.Chat.ID, res[1])
		if err != nil {
			msg := tgbotapi.NewMessage(m.Chat.ID, "channel not found for you")
			msg.ParseMode = htmlMode
			_, err := bot.Send(msg)
			assert.Nil(err)
			return
		}
		//everything ok publish a job TODO:
		fmt.Println(channel.Name)
		return
	}
	channels, err := chnManger.FindChannelsByChatID(m.Chat.ID)
	assert.Nil(err)
	if len(channels) == 0 {
		msg := tgbotapi.NewMessage(m.Chat.ID, "no channels for you")
		msg.ParseMode = htmlMode
		_, err := bot.Send(msg)
		assert.Nil(err)
		return
	}
	textMsg := "please choose one of the below channels\n"
	for i := range channels {
		textMsg += fmt.Sprintf("/ad_%s\n", channels[i].Name)
	}
	msg := tgbotapi.NewMessage(m.Chat.ID, textMsg)
	msg.ParseMode = htmlMode
	_, err = bot.Send(msg)
	assert.Nil(err)
	return

}

func (bb *bot) Initialize() {

	tgbot.RegisterMessageHandler("/updatead", bb.updateAD)
	tgbot.RegisterMessageHandler("/ad", bb.wantAD)
	//assert.Nil(b.Start())
}

func init() {
	initializer.Register(&bot{})
}
