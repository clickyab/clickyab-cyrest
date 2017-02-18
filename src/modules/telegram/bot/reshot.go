package bot

import (
	"common/rabbit"
	"modules/telegram/ad/ads"
	"modules/telegram/cyborg/commands"
	"modules/telegram/teleuser/tlu"
	"modules/user/aaa"
	"strconv"
	"strings"

	"gopkg.in/telegram-bot-api.v4"
)

func (bb *bot) reshot(bot *tgbotapi.BotAPI, m *tgbotapi.Message) {
	parts := strings.Split(m.Text, "_")
	if len(parts) != 2 {
		send(bot, m.Chat.ID, "your command is <b>not valid1</b>")
		return
	}

	channelID, err := strconv.ParseInt(parts[1], 10, 0)
	if err != nil {
		send(bot, m.Chat.ID, "your command is <b>not valid2</b>")
		return
	}

	tele := tlu.NewTluManager()
	telegramUser, err := tele.FindTeleUserByBotChatID(m.Chat.ID)
	if err != nil {
		send(bot, m.Chat.ID, "your telegram user is not in our system \n please register!")
		return
	}
	usr := aaa.NewAaaManager()
	user, err := usr.FindUserByID(telegramUser.UserID)
	if err != nil {
		send(bot, m.Chat.ID, "your telegram user is not in our system \n please register!")
		return
	}
	b := ads.NewAdsManager()
	channel, err := b.FindChannelByUserIDChannelID(user.ID, channelID)
	if err != nil {
		send(bot, m.Chat.ID, "you are not owner this channel")
		return
	}
	channelAd, err := b.FindChannelAdActiveByChannelID(channel.ID, ads.ActiveStatusYes)
	if err != nil || len(channelAd) == 0 {
		send(bot, m.Chat.ID, "your command is <b>not valid</b>")
		return
	}

	rabbit.MustPublish(
		commands.DiscoverAd{
			Channel: channelID,
			ChatID:  m.Chat.ID,
			Reshot:  true,
		},
	)
}
