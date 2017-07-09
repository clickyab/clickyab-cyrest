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

const htmlMode string = "HTML"

var (
	doneReg    = regexp.MustCompile("/(done|reject)_([0-9]+)")
	rejectReg  = regexp.MustCompile("/(done|reject)_([0-9]+)_([0-9]+)")
	completeAd = regexp.MustCompile("/complete_([0-9]+)_([0-9]+)")
	//send bundle to channel first bundle id & second channel id
	sendAd = regexp.MustCompile("/send_([0-9]+)_([0-9]+)")
)

func sendString(bot *tgbotapi.BotAPI, chatID int64, message string) {
	msg := tgbotapi.NewMessage(chatID, message)
	msg.ParseMode = htmlMode
	_, err := bot.Send(msg)
	assert.Nil(err)
}

func send(bot *tgbotapi.BotAPI, chatID int64, message trans.T9String) {
	sendString(bot, chatID, message.Translate(trans.PersianLang))
}

func sendWithKeyboard(bot *tgbotapi.BotAPI, keyboard tgbotapi.ReplyKeyboardMarkup, chatID int64, message trans.T9String) {
	msg := tgbotapi.NewMessage(chatID, message.Text)
	msg.ReplyMarkup = keyboard
	bot.Send(msg)
}

func (bb *bot) Initialize() {
	tgbot.RegisterMessageHandler("/updatead", bb.updateAD)
	tgbot.RegisterMessageHandler("/ad", bb.wantAD)
	tgbot.RegisterMessageHandler("/confirm", bb.confirm)
	tgbot.RegisterMessageHandler("/done", bb.doneORReject)
	tgbot.RegisterMessageHandler("/reject", bb.doneORReject)
	tgbot.RegisterMessageHandler("/reshot", bb.reshot)
	tgbot.RegisterMessageHandler("/activead", bb.activeAd)
	tgbot.RegisterMessageHandler("/complete", bb.complete)

	tgbot.RegisterMessageHandler("/name", bb.getName)
	tgbot.RegisterMessageHandler("/channel", bb.getChannel)
	tgbot.RegisterMessageHandler("/secret", CheckUserExisted(bb.test))

	tgbot.RegisterMessageHandler("/addchan", bb.addChan)
	tgbot.RegisterMessageHandler("/delchan", bb.delChan)
	tgbot.RegisterMessageHandler("/addCard", CheckUserExisted(bb.financial))

	// lint hack
	if false {
		bb.getCard("")
		bb.getAccount("")
	}
	tgbot.RegisterMessageHandler("/send", bb.sendAd)
}

func init() {
	initializer.Register(&bot{})
}
