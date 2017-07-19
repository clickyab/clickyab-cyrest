package bot

import (
	"modules/misc/trans"
	"modules/telegram/ad/ads"
	"modules/telegram/common/tgbot"
	"modules/telegram/teleuser/tlu"
	"time"

	"github.com/Sirupsen/logrus"

	"common/assert"

	"common/models/common"

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
		err := ads.NewAdsManager().CreateChannel(&ads.Channel{
			UserID:        tu.ID,
			Name:          m.Text,
			Title:         common.NullString{Valid: false},
			AdminStatus:   ads.AdminStatusPending,
			ArchiveStatus: ads.ActiveStatusNo,
			Active:        ads.ActiveStatusNo,
		})
		assert.Nil(err)

		keyboard := tgbot.NewKeyboard("/get_ad", "/addchan", "/delchan")
		sendWithKeyboard(bot, keyboard, m.Chat.ID, trans.T("use get_ad to get a new ad\nor /delchan to add a new channel"))
	}, func() {
		send(bot, m.Chat.ID, trans.T("times up\nenter /addchan again"))
	}, 20*time.Second)
}

func (bb *bot) delChan(bot *tgbotapi.BotAPI, m *tgbotapi.Message) {
	tu, err := tlu.NewTluManager().FindTeleUserByBotChatID(int64(m.Chat.ID))
	if err != nil {
		send(bot, m.Chat.ID, trans.T("couldn't find your user"))
		return
	}

	channels, err := ads.NewAdsManager().FindActiveChannelsByUserID(tu.UserID)
	if err != nil {
		send(bot, m.Chat.ID, trans.T("couldn't find your channels\ntry again with /addchan"))
		return
	}

	var channelsName []string
	for i := range channels {
		channelsName = append(channelsName, channels[i].Name)
	}
	keyboard := tgbot.NewKeyboard(channelsName...)
	sendWithKeyboard(bot, keyboard, m.Chat.ID, trans.T("choose one of these channels to be deleted"))
	tgbot.RegisterUserHandlerWithExp(m.Chat.ID, func(bot *tgbotapi.BotAPI, m *tgbotapi.Message) {
		defer tgbot.UnRegisterUserHandler(m.Chat.ID)

		if ok := containsString(channelsName, m.Text); !ok {
			send(bot, m.Chat.ID, trans.T("couldn't find this channel\nenter /delchan to try again"))
			return
		}
		_, err = ads.NewAdsManager().DeleteChannelByUserIDChannelName(tu.UserID, m.Text)
		if err != nil {
			logrus.Warn(err.Error())
			send(bot, m.Chat.ID, trans.T("couldn't delete the channel\nenter /delchan to try again"))
			return
		}

		keyboard := tgbot.NewKeyboard("/delchan", "/ad")
		sendWithKeyboard(bot, keyboard, m.Chat.ID, trans.T("channel deleted successfully"))
	}, func() {
		send(bot, m.Chat.ID, trans.T("times up\nenter /delchan to try again"))
	}, 20*time.Second)
}

func containsString(slice []string, target string) bool {
	for i := range slice {
		if slice[i] == target {
			return true
		}
	}
	return false
}
