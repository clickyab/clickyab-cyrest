package bot

import (
	"common/assert"
	"modules/misc/trans"
	"modules/telegram/teleuser/tlu"

	"modules/telegram/common/tgbot"

	"fmt"

	"gopkg.in/telegram-bot-api.v4"
)

func (bb *bot) start(bot *tgbotapi.BotAPI, m *tgbotapi.Message) {
	user, err := tlu.NewTluManager().FindTeleUserByBotChatID(m.Chat.ID)
	assert.Nil(err)

	k := tgbot.NewKeyboard(`/getBundle`, `/addchan`, `/delchan`, `/report`)
	request := fmt.Sprintf("hello %s, how are u today bitch ?\n"+
		"Enter /sendAd to get a bundle\n"+
		"Enter /addchan to add your new channel\n"+
		"Enter /delchan to delete one of your channels\n"+
		"Enter /report to get your financial report", user.Username.String)
	sendWithKeyboard(bot, k, m.Chat.ID, trans.T(request))

}
