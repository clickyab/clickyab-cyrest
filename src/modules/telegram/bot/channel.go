package bot

import (
	"modules/misc/trans"
	"modules/telegram/ad/ads"
	"modules/telegram/common/tgbot"
	"modules/telegram/teleuser/tlu"
	"time"

	"github.com/Sirupsen/logrus"

	"common/models/common"

	"fmt"

	"gopkg.in/telegram-bot-api.v4"
)

func (bb *bot) addChan(bot *tgbotapi.BotAPI, m *tgbotapi.Message) {
	tu, err := tlu.NewTluManager().FindTeleUserByBotChatID(int64(m.Chat.ID))
	if err != nil {
		send(bot, m.Chat.ID, trans.T("couldn't find your user"))
		return
	}
	send(bot, m.Chat.ID, trans.T("write down your channel tag\nor /cancel to exit process"))
	tgbot.RegisterUserHandler(m.Chat.ID, func(bot *tgbotapi.BotAPI, m *tgbotapi.Message) {
		defer tgbot.UnRegisterUserHandler(m.Chat.ID)

		if m.Text == "/cancel" {
			request := "adding channel process canceled\n" +
				"Enter /getbundle to get a bundle\n" +
				"Enter /addchan to add your new channel\n" +
				"Enter /delchan to delete one of your channels\n" +
				"Enter /report to get your financial report\n"

			keyboard := tgbot.NewKeyboard([]string{"/getbundle", "/addchan", "/delchan", "/report"})
			sendWithKeyboard(bot, keyboard, m.Chat.ID, trans.T(request))
			return
		}

		err := ads.NewAdsManager().CreateChannel(&ads.Channel{
			UserID:        tu.UserID,
			Name:          m.Text,
			Title:         common.NullString{Valid: false},
			AdminStatus:   ads.AdminStatusPending,
			ArchiveStatus: ads.ActiveStatusNo,
			Active:        ads.ActiveStatusYes,
		})
		if err != nil {
			send(bot, m.Chat.ID, trans.T("couldn't add ur channel\ntry /addchan again"))
		}

		request := fmt.Sprintf("ur channel added successfully\n" +
			"Enter /getbundle to get a bundle\n" +
			"Enter /addchan to add your new channel\n" +
			"Enter /delchan to delete one of your channels\n" +
			"Enter /report to get your financial report")
		keyboard := tgbot.NewKeyboard([]string{"/getbundle", "/addchan", "/delchan", "/report"})
		sendWithKeyboard(bot, keyboard, m.Chat.ID, trans.T(request))
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

	var channelsName = make([]string, 0)
	for i := range channels {
		logrus.Warn(channels[i].Name)
		channelsName = append(channelsName, channels[i].Name)
	}

	keyboard := tgbot.NewKeyboard(append(channelsName, "/cancel"))
	sendWithKeyboard(bot, keyboard, m.Chat.ID, trans.T("choose one of these channels to be deleted"))

	tgbot.RegisterUserHandler(m.Chat.ID, func(bot *tgbotapi.BotAPI, m *tgbotapi.Message) {
		defer tgbot.UnRegisterUserHandler(m.Chat.ID)

		if m.Text == "/cancel" {
			request := `deleting channel process canceled\n` +
				`Enter /getbundle to get a bundle\n` +
				`Enter /addchan to add your new channel\n` +
				`Enter /delchan to delete one of your channels\n` +
				`Enter /report to get your financial report\n`

			keyboard := tgbot.NewKeyboard([]string{"/getbundle", "/addchan", "/delchan", "/report"})
			sendWithKeyboard(bot, keyboard, m.Chat.ID, trans.T(request))
			return
		}

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

		request := `Your channel added successfully\n` +
			`Enter /getbundle to get a bundle\n` +
			`Enter /addchan to add your new channel\n` +
			`Enter /delchan to delete one of your channels\n` +
			`Enter /report to get your financial report\n`

		keyboard := tgbot.NewKeyboard([]string{"/getbundle", "/addchan", "/delchan", "/report"})
		sendWithKeyboard(bot, keyboard, m.Chat.ID, trans.T("channel deleted successfully\n"+request))
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
