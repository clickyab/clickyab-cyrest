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

	//Financial is the key for financial conversation
	Financial = iota
	//Screenshot is the key for screen shot conversation
	Screenshot
)

// ChatHandler ChatHandler
type ChatHandler func(*tgbotapi.Message, map[int]interface{}) (string, bool)

// Chatter consist of request and a func to handle response
type Chatter struct {
	StringRequest   string
	KeyboardRequest *[]string
	Handler         ChatHandler
}

type response struct {
	keyboard *[]string
	isFilled bool
	response string
}

var combination = make(map[int][]*Chatter)

// RegisterConversation adds a chatter to a field
func RegisterConversation(field int, c Chatter) {
	combination[field] = append(combination[field], &c)
}

// StartConversion this will be called for each command
func StartConversion(bot *tgbotapi.BotAPI, chatID int64, mode int) {
	states := make([]response, len(combination[mode]))
	data := map[int]interface{}{}
	qSet := combination[mode]
	if len(qSet) < 1 {
		return
	}

	// asking first Q
	send(bot, chatID, response{response: qSet[0].StringRequest})

	RegisterUserHandler(chatID, func(bot *tgbotapi.BotAPI, m *tgbotapi.Message) {
		for i := range qSet {
			currentState := &states[i]

			// pass if current state is filled
			if currentState.isFilled {
				continue
			}

			// getting response of first Q that haven't completed
			res, valid := qSet[i].Handler(m, data)
			currentState.response = res
			currentState.isFilled = valid

			// add next question to the current mode response
			if i+1 < len(qSet) && currentState.isFilled {
				currentState.response = fmt.Sprintf("%s\n%s", currentState.response, qSet[i+1].StringRequest)
				if qSet[i+1].KeyboardRequest != nil {
					currentState.keyboard = qSet[i+1].KeyboardRequest
				}
			}
			send(bot, chatID, *currentState)
			break
		}

	}, time.Hour)
}

func send(bot *tgbotapi.BotAPI, chatID int64, r response) {
	if r.keyboard == nil {
		sendString(bot, chatID, trans.T(r.response).Translate(trans.PersianLang))
		return
	}
	sendWithKeyboard(bot, NewKeyboard(*r.keyboard...), chatID, trans.T(r.response))
}

func sendString(bot *tgbotapi.BotAPI, chatID int64, message string) {
	translated := trans.T(message).Translate(trans.PersianLang)
	msg := tgbotapi.NewMessage(chatID, translated)
	msg.ParseMode = htmlMode
	_, err := bot.Send(msg)
	assert.Nil(err)
}

func sendWithKeyboard(bot *tgbotapi.BotAPI, keyboard tgbotapi.ReplyKeyboardMarkup, chatID int64, message trans.T9String) {
	msg := tgbotapi.NewMessage(chatID, message.Text)
	msg.ReplyMarkup = keyboard
	bot.Send(msg)
}
