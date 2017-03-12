package bot

import (
	"common/assert"
	"common/initializer"
	"modules/telegram/common/tgbot"
	"regexp"

	"modules/misc/trans"

	"gopkg.in/telegram-bot-api.v4"
)

type bot struct {
}

const htmlMode = "HTML"

var doneRejectReg = regexp.MustCompile("/(done|reject)_([0-9]+)")

func sendString(bot *tgbotapi.BotAPI, chatID int64, message string) {
	msg := tgbotapi.NewMessage(chatID, message)
	msg.ParseMode = htmlMode
	_, err := bot.Send(msg)
	assert.Nil(err)
}

func send(bot *tgbotapi.BotAPI, chatID int64, message trans.T9String) {
	sendString(bot, chatID, message.Translate(trans.PersianLang))
}

func (bb *bot) Initialize() {

	tgbot.RegisterMessageHandler("/updatead", bb.updateAD)
	tgbot.RegisterMessageHandler("/ad", bb.wantAD)
	tgbot.RegisterMessageHandler("/confirm", bb.confirm)
	tgbot.RegisterMessageHandler("/done", bb.doneORReject)
	tgbot.RegisterMessageHandler("/reject", bb.doneORReject)
	tgbot.RegisterMessageHandler("/reshot", bb.reshot)
}

func init() {
	initializer.Register(&bot{})
}
