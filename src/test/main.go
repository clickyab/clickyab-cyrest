package main

import (
	"common/assert"
	"common/config"
	"common/initializer"

	"common/tgbot"

	"gopkg.in/telegram-bot-api.v4"
)

func main() {
	config.Initialize()
	config.InitApplication()

	defer initializer.Initialize().Finalize()

	b := tgbot.NewTelegramBot("273335144:AAEv4uPeo68X7Scc3MLKxwMO1YI3JFkWiJM")

	b.RegisterMessageHandler("/test", func(bot *tgbotapi.BotAPI, m *tgbotapi.Message) {
		msg := tgbotapi.NewMessage(m.Chat.ID, "<b>Hi</b> Dude! use /test_for_this")
		msg.ParseMode = "HTML"
		_, err := bot.Send(msg)
		assert.Nil(err)
	})
	assert.Nil(b.Start())
}
