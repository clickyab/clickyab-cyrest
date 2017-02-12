package bot

import (
	"common/models/common"
	"common/rabbit"
	"fmt"
	"modules/telegram/ad/ads"
	"modules/telegram/cyborg/commands"
	"modules/telegram/teleuser/tlu"
	"modules/user/aaa"
	"strconv"
	"time"

	"modules/telegram/config"

	"github.com/Sirupsen/logrus"
	"gopkg.in/telegram-bot-api.v4"
)

func (bb *bot) doneORReject(bot *tgbotapi.BotAPI, m *tgbotapi.Message) {

	b := ads.NewAdsManager()
	doneSlice := doneRejectReg.FindStringSubmatch(m.Text)
	logrus.Warn(doneSlice, len(doneSlice))
	if len(doneSlice) != 3 {
		send(bot, m.Chat.ID, "your command is <b>not valid1</b>")
		return
	}
	channelID, err := strconv.ParseInt(doneSlice[2], 10, 0)
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

	channel, err := b.FindChannelByUserIDChannelID(user.ID, channelID)
	if err != nil {
		send(bot, m.Chat.ID, "you are not owner this channel")
		return
	}

	if doneSlice[1] == "done" {
		channelAd, err := b.FindChannelAdActiveByChannelID(channel.ID, ads.ActiveStatusNo)
		if err != nil {
			send(bot, m.Chat.ID, "your command is <b>not valid</b>")
			return
		}
		var adS []int64
		for chAd := range channelAd {
			adS = append(adS, channelAd[chAd].AdID)
			channelAd[chAd].Active = ads.ActiveStatusYes
			channelAd[chAd].Start = common.MakeNullTime(time.Now())
			err = b.UpdateChannelAd(&channelAd[chAd])
			if err != nil {
				send(bot, m.Chat.ID, "your command is <b>not valid</b>")
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

		send(bot, m.Chat.ID, fmt.Sprintf("ads active in <b>%s</b> channle", channel.Name))
		return
	}
	//reject command
	channelAd, err := b.FindChannelAdActiveByChannelID(channel.ID, ads.ActiveStatusYes)
	if err != nil {
		send(bot, m.Chat.ID, "your command is <b>not valid</b>")
		return
	}
	for chAd := range channelAd {
		channelAd[chAd].Active = ads.ActiveStatusNo
		channelAd[chAd].End = common.MakeNullTime(time.Now())
		err = b.UpdateChannelAd(&channelAd[chAd])
		if err != nil {
			send(bot, m.Chat.ID, "your command is <b>not valid</b>")
			return
		}
	}
	send(bot, m.Chat.ID, fmt.Sprintf("ads reject in <b>%s</b> channle", channel.Name))
	return

}
