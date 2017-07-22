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
	doneReg              = regexp.MustCompile(`/done_(\w+)`)
	channelBundleIDRegex = regexp.MustCompile(`(\d+)_(\d+)`)
	completeAd           = regexp.MustCompile("/complete_([0-9]+)_([0-9]+)")
	//send bundle to channel first bundle id & second channel id
	//sendAd = regexp.MustCompile("/send_([0-9]+)_([0-9]+)")

	channelNames map[int64][]string
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
	_, err := bot.Send(msg)
	assert.Nil(err)
}

func (bb *bot) Initialize() {
	err := tgbot.RegisterMessageHandler("/start", bb.start)
	assert.Nil(err)
	err = tgbot.RegisterMessageHandler("/updatead", bb.updateAD)
	assert.Nil(err)
	err = tgbot.RegisterMessageHandler("/ad", bb.wantAD)
	assert.Nil(err)
	err = tgbot.RegisterMessageHandler("/confirm", bb.confirm)
	assert.Nil(err)
	err = tgbot.RegisterMessageHandler("/done", CheckUserExisted(bb.doneORReject))
	assert.Nil(err)
	err = tgbot.RegisterMessageHandler("/reshot", bb.reshot)
	assert.Nil(err)
	err = tgbot.RegisterMessageHandler("/activead", bb.activeAd)
	assert.Nil(err)
	err = tgbot.RegisterMessageHandler("/complete", bb.complete)
	assert.Nil(err)
	err = tgbot.RegisterMessageHandler("/name", bb.getName)
	assert.Nil(err)
	err = tgbot.RegisterMessageHandler("/channel", bb.getChannel)
	assert.Nil(err)
	err = tgbot.RegisterMessageHandler("/secret", CheckUserExisted(bb.test))
	assert.Nil(err)
	err = tgbot.RegisterMessageHandler("/addchan", bb.addChan)
	assert.Nil(err)
	err = tgbot.RegisterMessageHandler("/delchan", bb.delChan)
	assert.Nil(err)
	err = tgbot.RegisterMessageHandler("/addCard", CheckUserExisted(bb.financial))
	assert.Nil(err)
	err = tgbot.RegisterMessageHandler("/upscreen", CheckUserExisted(bb.uploadSS))
	assert.Nil(err)
	err = tgbot.RegisterMessageHandler("/getbundle", CheckUserExisted(bb.getBundle))
	assert.Nil(err)
}

func init() {
	channelNames = map[int64][]string{}
	initializer.Register(&bot{})
}
