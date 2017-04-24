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

	"common/assert"
	"common/redis"
	"fmt"
	"time"

	"gopkg.in/telegram-bot-api.v4"
)

func (bb *bot) doneORReject(bot *tgbotapi.BotAPI, m *tgbotapi.Message) {

	b := ads.NewAdsManager()
	doneSlice := doneReg.FindStringSubmatch(m.Text)
	rejectSlice := rejectReg.FindStringSubmatch(m.Text)
	if len(doneSlice) != 3 && len(rejectSlice) != 4 {
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
		assert.Nil(err)
		if len(channelAd) == 0 {
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
	channelOwner, err := aaa.NewAaaManager().FindUserByID(channel.UserID)
	if err != nil || channelOwner.ID != user.ID {
		return
	}
	err = b.DeleteChannelAdByChannelID(channel.ID)
	assert.Nil(err)
	adID, err := strconv.ParseInt(rejectSlice[3], 10, 0)
	if err != nil {
		return
	}
	currentAd, err := ads.NewAdsManager().FindAdByID(adID)
	assert.Nil(err)
	err = aredis.StoreHashKey(fmt.Sprintf("REJECT_%d", channel.ID), fmt.Sprintf("AD_%d", currentAd.ID), fmt.Sprintf("%d", currentAd.ID), 2*time.Hour)
	assert.Nil(err)
	if err != nil {
		send(bot, m.Chat.ID, trans.T("your command is <b>not valid</b>"))
		return
	}
	if err != nil {
		send(bot, m.Chat.ID, trans.T("your command is <b>not valid2</b>"))
		return
	}
	//dont select the same rejected ad for him
	send(bot, m.Chat.ID, trans.T("ads reject in <b>%s</b> channel", channel.Name))
	//send mail

	go func() {
		mail.SendByTemplateName(trans.T("channel rejected").Translate("fa_IR"), "reject-channel", struct {
			Name    string
			Channel string
		}{
			Name:    channelOwner.Email,
			Channel: channel.Name,
		}, config.Config.Mail.From, channelOwner.Email)
	}()
	return

}
