package main

import (
	"common/assert"
	"common/config"
	"common/initializer"
	"common/redis"
	"common/tgbot"
	"modules/teleuser/tlu"
	"strings"
	"time"

	"strconv"

	"gopkg.in/telegram-bot-api.v4"
)

func main() {
	config.Initialize()
	config.InitApplication()

	defer initializer.Initialize().Finalize()

	//b := tgbot.NewTelegramBot("273335144:AAEv4uPeo68X7Scc3MLKxwMO1YI3JFkWiJM")
	b := tgbot.NewTelegramBot("232630313:AAHRVcaQxFvs3u2-VGAAlsD3Xe1TIUr5rhk")

	b.RegisterMessageHandler("/verify", func(bot *tgbotapi.BotAPI, m *tgbotapi.Message) {
		//sample code  /verify-1:12123

		if !strings.Contains(m.Text, "-") && !strings.Contains(m.Text, ":") {
			msg := tgbotapi.NewMessage(m.Chat.ID, "your code is not <b>valid</b>")
			msg.ParseMode = "HTML"
			_, err := bot.Send(msg)
			assert.Nil(err)
			return
		}
		result := strings.Replace(m.Text, "/verify-", "", 1)
		str, err := aredis.GetKey(result, false, time.Hour)
		if str == "" || err != nil {
			msg := tgbotapi.NewMessage(m.Chat.ID, "your code is not <b>valid</b>")
			msg.ParseMode = "HTML"
			_, err := bot.Send(msg)
			assert.Nil(err)
			return
		}
		user := strings.Split(result, ":")
		id, err := strconv.ParseInt(user[0], 0, 10)
		if err == nil {
			msg := tgbotapi.NewMessage(m.Chat.ID, "your account is <b>accepted</b>")
			n := tlu.NewTluManager()
			tl := &tlu.Teleuser{
				UserID:     id,
				TelegramID: m.Chat.ID,
				Username:   m.Chat.UserName,
				Remove:     tlu.RemoveStatusNo,
				Resolve:    tlu.ResolveStatusYes,
			}
			assert.Nil(n.CreateTeleuser(tl))
			msg.ParseMode = "HTML"
			_, err := bot.Send(msg)
			assert.Nil(err)
		}
	})
	assert.Nil(b.Start())
}
