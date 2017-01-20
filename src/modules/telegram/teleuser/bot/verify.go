package bot

import (
	"common/assert"
	"common/initializer"
	"common/models/common"
	"common/redis"
	"fmt"
	"modules/telegram/channel/chn"
	"modules/telegram/common/tgbot"
	"modules/telegram/teleuser/tlu"
	"strconv"
	"strings"
	"time"

	"gopkg.in/telegram-bot-api.v4"
)

const htmlMode = "HTML"

type bot struct {
}

func (bb *bot) verify(bot *tgbotapi.BotAPI, m *tgbotapi.Message) {
	//sample code  /verify-1:12123
	if !strings.Contains(m.Text, "-") && !strings.Contains(m.Text, ":") {
		msg := tgbotapi.NewMessage(m.Chat.ID, "your code is not <b>valid</b>")
		msg.ParseMode = htmlMode
		_, err := bot.Send(msg)
		assert.Nil(err)
		return
	}
	result := strings.Replace(m.Text, "/verify-", "", 1)
	str, err := aredis.GetKey(result, false, time.Hour)
	if str == "" || err != nil {
		msg := tgbotapi.NewMessage(m.Chat.ID, "your code is not <b>valid</b>")
		msg.ParseMode = htmlMode
		_, err := bot.Send(msg)
		assert.Nil(err)
		return
	}
	user := strings.Split(result, ":")
	id, err := strconv.ParseInt(user[0], 0, 10)
	if err == nil {
		msg := tgbotapi.NewMessage(m.Chat.ID, "your account is <b>accepted</b>")
		n := tlu.NewTluManager()
		tl := &tlu.TeleUser{
			UserID:    id,
			BotChatID: m.Chat.ID,
			Username:  common.MakeNullString(m.Chat.UserName),
			Remove:    tlu.RemoveStatusNo,
			Resolve:   tlu.ResolveStatusYes,
		}
		assert.Nil(n.CreateTeleUser(tl))
		msg.ParseMode = htmlMode
		_, err := bot.Send(msg)
		assert.Nil(err)
	}
}

func (bb *bot) wantAD(bot *tgbotapi.BotAPI, m *tgbotapi.Message) {
	//find channels
	chnManger := chn.NewChnManager()
	fmt.Println(m.Chat.ID)
	channels, err := chnManger.FindChannelsByChatID(m.Chat.ID)
	assert.Nil(err)
	if len(channels) == 0 {
		msg := tgbotapi.NewMessage(m.Chat.ID, "no channels for you")
		msg.ParseMode = htmlMode
		_, err := bot.Send(msg)
		assert.Nil(err)
		return
	}
	textMsg := ""
	for i := range channels {
		textMsg += fmt.Sprintf("/ad-%d\n", channels[i].ID)
	}
	msg := tgbotapi.NewMessage(m.Chat.ID, textMsg)
	msg.ParseMode = htmlMode
	_, err = bot.Send(msg)
	assert.Nil(err)
	return

}

func (bb *bot) Initialize() {

	tgbot.RegisterMessageHandler("/verify", bb.verify)
	tgbot.RegisterMessageHandler("/ad", bb.wantAD)
}

func init() {
	initializer.Register(&bot{})
}
