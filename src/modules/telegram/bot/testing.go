package bot

import (
	"modules/misc/trans"

	"gopkg.in/telegram-bot-api.v4"
)

func (bb *bot) test(bot *tgbotapi.BotAPI, m *tgbotapi.Message) {
	send(bot, m.Chat.ID, trans.T("you registered"))
}
