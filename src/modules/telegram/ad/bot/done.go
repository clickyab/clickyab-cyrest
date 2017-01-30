package bot

import (
	"common/assert"

	"regexp"

	"strconv"

	"modules/telegram/ad/ads"

	"time"

	"fmt"

	"common/rabbit"
	"modules/telegram/cyborg/commands"

	"common/models/common"

	"github.com/Sirupsen/logrus"
	"gopkg.in/telegram-bot-api.v4"
)

func (bb *bot) done(bot *tgbotapi.BotAPI, m *tgbotapi.Message) {
	var doneReg = regexp.MustCompile("/done_([0-9]+)_([0-9]+)")
	doneSlice := doneReg.FindStringSubmatch(m.Text)
	logrus.Infof("%+v", doneSlice)
	if len(doneSlice) != 3 {
		msg := tgbotapi.NewMessage(m.Chat.ID, "your command is <b>not valid</b>")
		msg.ParseMode = htmlMode
		_, err := bot.Send(msg)
		assert.Nil(err)
		return
	}
	channelID, err := strconv.ParseInt(doneSlice[1], 10, 0)
	if err != nil {
		msg := tgbotapi.NewMessage(m.Chat.ID, "your command is <b>not valid</b>")
		msg.ParseMode = htmlMode
		_, err := bot.Send(msg)
		assert.Nil(err)
		return
	}
	b := ads.NewAdsManager()
	channel, err := b.FindChannelByID(channelID)
	if err != nil {
		msg := tgbotapi.NewMessage(m.Chat.ID, "your command is <b>not valid</b>")
		msg.ParseMode = htmlMode
		_, err := bot.Send(msg)
		assert.Nil(err)
		return
	}
	adID, err := strconv.ParseInt(doneSlice[2], 10, 0)
	if err != nil {
		msg := tgbotapi.NewMessage(m.Chat.ID, "your command is <b>not valid</b>")
		msg.ParseMode = htmlMode
		_, err := bot.Send(msg)
		assert.Nil(err)
		return
	}
	ad, err := b.FindAdByID(adID)
	if err != nil {
		msg := tgbotapi.NewMessage(m.Chat.ID, "your command is <b>not valid</b>")
		msg.ParseMode = htmlMode
		_, err := bot.Send(msg)
		assert.Nil(err)
		return
	}
	channelAd, err := b.FindChannelIDAdByAdID(channel.ID, ad.ID)
	if err != nil {
		msg := tgbotapi.NewMessage(m.Chat.ID, "your command is <b>not valid</b>")
		msg.ParseMode = htmlMode
		_, err := bot.Send(msg)
		assert.Nil(err)
		return
	}
	channelAd.Active = ads.ActiveStatusYes
	channelAd.Start = common.MakeNullTime(time.Now())
	err = b.UpdateChannelAd(channelAd)
	if err != nil {
		msg := tgbotapi.NewMessage(m.Chat.ID, "your command is <b>not valid</b>")
		msg.ParseMode = htmlMode
		_, err := bot.Send(msg)
		assert.Nil(err)
		return
	}
	defer func() {
		rabbit.MustPublish(
			commands.ExistChannelAd{
				AdID:      ad.ID,
				ChannelID: channel.ID,
			},
		)
	}()

	msg := tgbotapi.NewMessage(m.Chat.ID, fmt.Sprintf("ad active in <b>%s</b> channle", channel.Name))
	msg.ParseMode = htmlMode
	_, err = bot.Send(msg)
	assert.Nil(err)
	return

}
