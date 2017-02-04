package bot

import (
	"strconv"

	"modules/telegram/ad/ads"

	"time"

	"fmt"

	"common/rabbit"
	"modules/telegram/cyborg/commands"

	"common/models/common"

	"modules/telegram/teleuser/tlu"

	"modules/user/aaa"

	"gopkg.in/telegram-bot-api.v4"
)

func (bb *bot) doneORReject(bot *tgbotapi.BotAPI, m *tgbotapi.Message) {

	b := ads.NewAdsManager()
	doneSlice := doneRejectReg.FindStringSubmatch(m.Text)
	if len(doneSlice) != 2 {
		send(bot, m.Chat.ID, "your command is <b>not valid</b>")
		return
	}
	channelID, err := strconv.ParseInt(doneSlice[1], 10, 0)
	if err != nil {
		send(bot, m.Chat.ID, "your command is <b>not valid</b>")
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

	if doneSlice[0] == "done" {
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
			rabbit.MustPublish(
				commands.ExistChannelAd{
					ChannelID: channel.ID,
					ChatID:    m.Chat.ID,
				},
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
