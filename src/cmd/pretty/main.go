package main

import (
	"common/assert"
	"log"

	"github.com/kr/pretty"
	"gopkg.in/telegram-bot-api.v4"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("310314178:AAHVjYxnceAKdo_2cf5XLts0EqJZ7ptz0Wk")
	assert.Nil(err)

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	assert.Nil(err)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, pretty.Sprint(update))
		msg.ReplyToMessageID = update.Message.MessageID

		_, err = bot.Send(msg)
		assert.Nil(err)
	}
}
