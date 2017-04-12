package bot

import (
	"common/rabbit"
	"modules/misc/trans"
	"modules/telegram/ad/ads"
	"modules/telegram/cyborg/commands"
	"modules/telegram/teleuser/tlu"
	"modules/user/aaa"
	"strconv"

	"common/config"
	"common/mail"

	"github.com/Sirupsen/logrus"
	"gopkg.in/telegram-bot-api.v4"
)

func (bb *bot) doneORReject(bot *tgbotapi.BotAPI, m *tgbotapi.Message) {

	b := ads.NewAdsManager()
	doneSlice := doneRejectReg.FindStringSubmatch(m.Text)
	logrus.Warn(doneSlice, len(doneSlice))
	if len(doneSlice) != 3 {
		send(bot, m.Chat.ID, trans.T("your command is <b>not valid1</b>"))
		return
	}
	channelID, err := strconv.ParseInt(doneSlice[2], 10, 0)
	if err != nil {
		send(bot, m.Chat.ID, trans.T("your command is <b>not valid2</b>"))
		return
	}
	tele := tlu.NewTluManager()
	telegramUser, err := tele.FindTeleUserByBotChatID(m.Chat.ID)
	if err != nil {
		send(bot, m.Chat.ID, trans.T("your telegram user is not in our system \n please register!"))
		return
	}
	usr := aaa.NewAaaManager()
	user, err := usr.FindUserByID(telegramUser.UserID)
	if err != nil {
		send(bot, m.Chat.ID, trans.T("your telegram user is not in our system \n please register!"))
		return
	}

	channel, err := b.FindChannelByUserIDChannelID(user.ID, channelID)
	if err != nil {
		send(bot, m.Chat.ID, trans.T("you are not owner this channel"))
		return
	}

	if doneSlice[1] == "done" {
		channelAd, err := b.FindChannelAdActiveByChannelID(channel.ID, ads.ActiveStatusNo)
		if err != nil || len(channelAd) == 0 {
			send(bot, m.Chat.ID, trans.T("your command is <b>not valid</b>"))
			return
		}

		rabbit.MustPublish(
			commands.DiscoverAd{
				Channel: channelID,
				ChatID:  m.Chat.ID,
			},
		)

		return
	}
	//reject command
	err = b.DeleteChannelAdByChannelID(channel.ID)
	if err != nil {
		send(bot, m.Chat.ID, trans.T("your command is <b>not valid</b>"))
		return
	}
	send(bot, m.Chat.ID, trans.T("ads reject in <b>%s</b> channel", channel.Name))
	//send mail
	owner, err := aaa.NewAaaManager().FindUserByID(channel.UserID)
	if err != nil {
		return
	}
	go func() {
		mail.SendByTemplateName(trans.T("channel rejected").Translate("fa_IR"), "reject-channel", struct {
			Name    string
			Channel string
		}{
			Name:    owner.Email,
			Channel: channel.Name,
		}, config.Config.Mail.From, owner.Email)
	}()
	return

}
