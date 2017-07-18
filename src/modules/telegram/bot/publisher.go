package bot

import (
	"common/assert"
	"common/redis"
	"fmt"
	"modules/misc/trans"
	"modules/telegram/ad/ads"
	"modules/telegram/common/tgbot"
	"modules/telegram/config"
	"modules/telegram/teleuser/tlu"
	"modules/user/aaa"
	"time"

	"gopkg.in/telegram-bot-api.v4"
)

func (bb *bot) getName(bot *tgbotapi.BotAPI, m *tgbotapi.Message) {
	name, _ := aredis.GetHashKey(fmt.Sprintf("p_%d", m.Chat.ID), "name", false, 2*time.Minute)
	if name != "" {
		//redirect to get channel command

		bb.getChannel(bot, m)

		return
	}
	send(bot, m.Chat.ID, trans.T("Please Enter Your Name "))
	tgbot.RegisterUserHandler(m.Chat.ID, func(bot *tgbotapi.BotAPI, m *tgbotapi.Message) {
		if len(m.Text) < 3 {
			tgbot.UnRegisterUserHandler(m.Chat.ID)
			return
		}
		err := aredis.StoreHashKey(fmt.Sprintf("p_%d", m.Chat.ID), "name", m.Text, 10*time.Minute)
		if err != nil {
			tgbot.UnRegisterUserHandler(m.Chat.ID)
			return
		}
		tgbot.UnRegisterUserHandler(m.Chat.ID)
		bb.getChannel(bot, m)
	}, time.Minute)
}

func (bb *bot) getChannel(bot *tgbotapi.BotAPI, m *tgbotapi.Message) {
	name, _ := aredis.GetHashKey(fmt.Sprintf("p_%d", m.Chat.ID), "name", false, 2*time.Minute)
	if name == "" {
		bb.getName(bot, m)
		return
	}
	send(bot, m.Chat.ID, trans.T("Please Enter Your Channel Username "))
	tgbot.RegisterUserHandler(m.Chat.ID, func(bot *tgbotapi.BotAPI, m *tgbotapi.Message) {
		if len(m.Text) < 3 {
			tgbot.UnRegisterUserHandler(m.Chat.ID)
			return
		}
		var err error
		userManager := aaa.NewAaaManager()
		err = userManager.Begin()
		if err != nil {
			return
		}

		defer func() {
			if err != nil {
				tgbot.UnRegisterUserHandler(m.Chat.ID)
				assert.Nil(userManager.Rollback())
			} else {
				val := fmt.Sprintf("%d", m.Chat.ID)
				err = aredis.StoreKey(val, val, tcfg.Cfg.Telegram.PublisherAuthRedis)
				assert.Nil(err)
				err = userManager.Commit()
				assert.Nil(err)
			}
		}()
		user, err := userManager.RegisterPublisherUser(name, m.Chat.ID)
		if err != nil {
			return
		}
		tluManager, err := tlu.NewTluManagerFromTransaction(userManager.GetDbMap())

		if err != nil {
			return
		}
		err = tluManager.RegisterPublisherTeleUser(m.Chat.ID, user.ID, name)
		if err != nil {
			return
		}

		channelManager, err := ads.NewAdsManagerFromTransaction(userManager.GetDbMap())
		if err != nil {
			return
		}
		err = channelManager.RegisterPublisherChannel(name, user.ID)
		if err != nil {
			return
		}
		aredis.Client.Del(fmt.Sprintf("p_%d", m.Chat.ID))
		tgbot.UnRegisterUserHandler(m.Chat.ID)
	}, time.Minute)
}
