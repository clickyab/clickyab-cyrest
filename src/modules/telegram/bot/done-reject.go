package bot

import (
	"common/models/common"
	"common/rabbit"
	"modules/telegram/ad/ads"
	"modules/telegram/cyborg/commands"
	"modules/telegram/teleuser/tlu"
	"modules/user/aaa"
	"strconv"
	"time"

	"modules/telegram/config"

	"modules/misc/trans"

	"github.com/Sirupsen/logrus"
	"gopkg.in/telegram-bot-api.v4"
)

func (bb *bot) doneORReject(bot *tgbotapi.BotAPI, m *tgbotapi.Message) {

	b := ads.NewAdsManager()
	doneSlice := doneRejectReg.FindStringSubmatch(m.Text)
	logrus.Warn(doneSlice, len(doneSlice))
	if len(doneSlice) != 3 {
		send(bot, m.Chat.ID, trans.T("your command is <b>not valid1</b>").Translate())
		return
	}
	channelID, err := strconv.ParseInt(doneSlice[2], 10, 0)
	if err != nil {
		send(bot, m.Chat.ID, trans.T("your command is <b>not valid1</b>").Translate())
		return
	}
	tele := tlu.NewTluManager()
	telegramUser, err := tele.FindTeleUserByBotChatID(m.Chat.ID)
	if err != nil {
		send(bot, m.Chat.ID, trans.T("your telegram user is not in our system \n please register!").Translate())
		return
	}
	usr := aaa.NewAaaManager()
	user, err := usr.FindUserByID(telegramUser.UserID)
	if err != nil {
		send(bot, m.Chat.ID, trans.T("your telegram user is not in our system \n please register!").Translate())
		return
	}

	channel, err := b.FindChannelByUserIDChannelID(user.ID, channelID)
	if err != nil {
		send(bot, m.Chat.ID, trans.T("you are not owner this channel").Translate())
		return
	}

	if doneSlice[1] == "done" {
		channelAd, err := b.FindChannelAdActiveByChannelID(channel.ID, ads.ActiveStatusNo)
		if err != nil {
			send(bot, m.Chat.ID, trans.T("your command is <b>not valid</b>").Translate())
			return
		}
		var adS []int64
		for chAd := range channelAd {
			adS = append(adS, channelAd[chAd].AdID)
			channelAd[chAd].Active = ads.ActiveStatusYes
			channelAd[chAd].Start = common.MakeNullTime(time.Now())
			err = b.UpdateChannelAd(&channelAd[chAd])
			if err != nil {
				send(bot, m.Chat.ID, trans.T("your command is <b>not valid</b>").Translate())
				return
			}
		}

		defer func() {
			logrus.Warn("done")
			rabbit.MustPublish(
				commands.DiscoverAd{
					Channel: channelID,
				},
			)
			rabbit.MustPublishAfter(
				commands.ExistChannelAd{
					ChannelID: channel.ID,
					ChatID:    m.Chat.ID,
				},
				tcfg.Cfg.Telegram.TimeReQueUe,
			)
		}()

		send(bot, m.Chat.ID, trans.T("ads active in <b>%s</b> channle", channel.Name).Translate())
		return
	}
	//reject command
	channelAd, err := b.FindChannelAdActiveByChannelID(channel.ID, ads.ActiveStatusYes)
	if err != nil {
		send(bot, m.Chat.ID, trans.T("your command is <b>not valid</b>").Translate())
		return
	}
	for chAd := range channelAd {
		channelAd[chAd].Active = ads.ActiveStatusNo
		channelAd[chAd].End = common.MakeNullTime(time.Now())
		err = b.UpdateChannelAd(&channelAd[chAd])
		if err != nil {
			send(bot, m.Chat.ID, trans.T("your command is <b>not valid</b>").Translate())
			return
		}
	}
	send(bot, m.Chat.ID, trans.T("ads reject in <b>%s</b> channle", channel.Name).Translate())
	return

}
