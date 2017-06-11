package bot

import (
	"modules/misc/trans"
	"modules/telegram/ad/ads"
	"modules/telegram/common/tgbot"
	"modules/telegram/teleuser/tlu"
	"time"

	"gopkg.in/telegram-bot-api.v4"
)

func (bb *bot) addChan(bot *tgbotapi.BotAPI, m *tgbotapi.Message) {
	tu, err := tlu.NewTluManager().FindTeleUserByBotChatID(int64(m.Chat.ID))
	if err != nil {
		send(bot, m.Chat.ID, trans.T("couldn't find your user"))
		return
	}
	send(bot, m.Chat.ID, trans.T("write down your channel tag"))
	tgbot.RegisterUserHandlerWithExp(m.Chat.ID, func(bot *tgbotapi.BotAPI, m *tgbotapi.Message) {
		defer tgbot.UnRegisterUserHandler(m.Chat.ID)
		ads.NewAdsManager().CreateChannel(&ads.Channel{
			UserID: tu.ID,
			Name:   m.Text,
		})

		keyboard := tgbot.NewKeyboard("/get_ad", "/fff")
		sendWithKeyboard(bot, keyboard, m.Chat.ID, trans.T("use get_ad to get a new ad\nor /fff to add a new channel"))
	}, func() {
		send(bot, m.Chat.ID, trans.T("times up\nenter /fff again"))
	}, 20*time.Second)
}
