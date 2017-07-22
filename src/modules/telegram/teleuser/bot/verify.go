package bot

import (
	"common/assert"
	"common/initializer"
	"common/models/common"
	"common/redis"
	"modules/telegram/common/tgbot"
	"modules/telegram/teleuser/tlu"
	"strconv"
	"strings"
	"time"

	"modules/misc/trans"

	"gopkg.in/telegram-bot-api.v4"
)

const htmlMode = "HTML"

type bot struct {
}

func (bb *bot) verify(bot *tgbotapi.BotAPI, m *tgbotapi.Message) {
	//sample code  /verify-1:12123
	if !strings.Contains(m.Text, "-") && !strings.Contains(m.Text, ":") {
		msg := tgbotapi.NewMessage(m.Chat.ID, trans.T("your command is <b>not valid</b>").String())
		msg.ParseMode = htmlMode
		_, err := bot.Send(msg)
		assert.Nil(err)
		return
	}
	result := strings.Replace(m.Text, "/verify-", "", 1)
	str, err := aredis.GetKey(result, false, time.Hour)
	if str == "" || err != nil {
		msg := tgbotapi.NewMessage(m.Chat.ID, trans.T("your command is <b>not valid</b>").String())
		msg.ParseMode = htmlMode
		_, err := bot.Send(msg)
		assert.Nil(err)
		return
	}
	user := strings.Split(result, ":")
	id, err := strconv.ParseInt(user[0], 0, 10)
	if err == nil {
		msg := tgbotapi.NewMessage(m.Chat.ID, trans.T("your account is <b>accepted</b>").String())
		n := tlu.NewTluManager()
		tl := &tlu.TeleUser{
			UserID:    id,
			BotChatID: m.Chat.ID,
			Username:  common.MakeNullString(m.Chat.UserName),
			Remove:    tlu.RemoveStatusNo,
			Resolve:   tlu.ResolveStatusYes,
		}
		err := n.CreateTeleUser(tl)
		if err != nil {
			msg1 := tgbotapi.NewMessage(m.Chat.ID, trans.T("your account <b>cant</b> be accepted").String())
			msg1.ParseMode = htmlMode
			_, err = bot.Send(msg1)
			assert.Nil(err)
		} else {
			msg.ParseMode = htmlMode
			_, err = bot.Send(msg)
			assert.Nil(err)
		}
	}
}

func (bb *bot) Initialize() {

	err := tgbot.RegisterMessageHandler("/verify", bb.verify)
	assert.Nil(err)
}

func init() {
	initializer.Register(&bot{})
}
