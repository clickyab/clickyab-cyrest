package bot

import (
	"modules/misc/trans"

	"common/redis"
	"time"

	"common/assert"

	"modules/telegram/ad/ads"
	"strconv"

	"github.com/Sirupsen/logrus"
	"gopkg.in/telegram-bot-api.v4"
)

func (bb *bot) doneORReject(bot *tgbotapi.BotAPI, m *tgbotapi.Message) {
	doneSlice := doneReg.FindStringSubmatch(m.Text)
	if len(doneSlice) != 2 {
		send(bot, m.Chat.ID, trans.T("your command is <b>not valid1</b>"))
		return
	}

	// channelID_bundleID
	value, err := aredis.GetKey(doneSlice[1], true, time.Hour)
	if err != nil || value == "" {
		send(bot, m.Chat.ID, trans.T("u've reached limited time out, try getting new bundle"))
		return
	}

	channelBundleID := channelBundleIDRegex.FindStringSubmatch(value)
	if len(channelBundleID) != 3 {
		logrus.Panic("wrong redis insert in /getBundle command")
	}

	channelID, err := strconv.ParseInt(channelBundleID[1], 10, 0)
	assert.Nil(err)

	/*channel, err := ads.NewAdsManager().FindChannelByID(channelID)
	assert.Nil(err)*/

	bundleID, err := strconv.ParseInt(channelBundleID[2], 10, 0)
	assert.Nil(err)

	bundle, err := ads.NewAdsManager().FindBundlesByID(bundleID)
	assert.Nil(err)

	adID := bundle.TargetAd

	err = ads.NewAdsManager().UpdateBundleChannelAd(&ads.BundleChannelAd{
		ChannelID: channelID,
		BundleID:  bundleID,
		AdID:      adID,
		Active:    ads.ActiveStatusYes,
	})
	assert.Nil(err)
}
