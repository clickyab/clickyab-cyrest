package bot

import (
	"common/assert"
	"common/initializer"
	"common/models/common"
	"common/redis"
	"common/tgbot"
	"modules/teleuser/tlu"
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
			UserID:     id,
			TelegramID: m.Chat.ID,
			Username:   common.MakeNullString(m.Chat.UserName),
			Remove:     tlu.RemoveStatusNo,
			Resolve:    tlu.ResolveStatusYes,
		}
		assert.Nil(n.CreateTeleuser(tl))
		msg.ParseMode = htmlMode
		_, err := bot.Send(msg)
		assert.Nil(err)
	}
}

func (bb *bot) Initialize() {

	tgbot.RegisterMessageHandler("/verify", bb.verify)
	//assert.Nil(b.Start())
}

func init() {
	initializer.Register(&bot{})
}
