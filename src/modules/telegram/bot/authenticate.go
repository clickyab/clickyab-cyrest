package bot

import (
	"modules/telegram/common/tgbot"
	"modules/telegram/teleuser/tlu"

	"common/assert"
	"common/redis"
	"fmt"
	"modules/telegram/config"

	"gopkg.in/telegram-bot-api.v4"
)

// CheckUserExisted try to check if user already registered or not
func CheckUserExisted(a tgbot.HandleMessage) tgbot.HandleMessage {
	return func(bot1 *tgbotapi.BotAPI, m *tgbotapi.Message) {
		val := fmt.Sprintf("%d", m.Chat.ID)
		data, err := aredis.GetKey(val, true, tcfg.Cfg.Telegram.PublisherAuthRedis)
		if err == nil && data != "" {
			if data == val {
				a(bot1, m)
				return
			}
		}
		_, err = tlu.NewTluManager().FindTeleUserByBotChatID(m.Chat.ID)
		if err != nil {
			var g *bot
			g.getName(bot1, m)
			return
		}
		err = aredis.StoreKey(val, val, tcfg.Cfg.Telegram.PublisherAuthRedis)
		assert.Nil(err)
		a(bot1, m)
	}
}
