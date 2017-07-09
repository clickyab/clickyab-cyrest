package tgbot

import (
	"time"

	"common/assert"
	"modules/misc/trans"

	"fmt"

	"gopkg.in/telegram-bot-api.v4"
)

const (
	htmlMode = "HTML"
)

// Chatter consist of request and a func to handle response
type Chatter struct {
	Request  string
	IsFilled bool
	Handler  ChatHandler
}

// ChatHandler ChatHandler
type ChatHandler func(string) (string, bool)

var (
	combination = make(map[string][]Chatter)
	clientState = make(map[string]int)
	counter     int
)

// RegisterConversation adds a chatter to a field
func RegisterConversation(field string, c Chatter) {
	combination[field] = append(combination[field], c)
	combination[field][counter] = c
	counter++
}

// NewConversation makes a new conversation for each chatID
func NewConversation(field string, bot *tgbotapi.BotAPI, chatID int64) {
	if len(combination[field]) == 0 {
		return
	}
	send(bot, chatID, trans.T(combination[field][0].Request))
	RegisterUserHandler(chatID, func(bot *tgbotapi.BotAPI, m *tgbotapi.Message) {
		for i := range combination[field] {
			if !combination[field][i].IsFilled {
				res, valid := combination[field][i].Handler(m.Text)
				if !valid {
					send(bot, chatID, trans.T(res))
					return
				}
				a := &combination[field][clientState[string(chatID)]]
				a.IsFilled = true
				clientState[string(chatID)]++
				if clientState[string(chatID)] < len(combination[field]) {
					res = fmt.Sprintf("%s\n%s", res, combination[field][clientState[string(chatID)]].Request)
				}
				send(bot, chatID, trans.T(res))
				clientState[string(chatID)] = 0
				return

			}
		}
		UnRegisterUserHandler(chatID)
	}, time.Minute)
	//TODO ^_^
}

func sendString(bot *tgbotapi.BotAPI, chatID int64, message string) {
	msg := tgbotapi.NewMessage(chatID, message)
	msg.ParseMode = htmlMode
	_, err := bot.Send(msg)
	assert.Nil(err)
}

func send(bot *tgbotapi.BotAPI, chatID int64, message trans.T9String) {
	sendString(bot, chatID, message.Translate(trans.PersianLang))
}
