package bot

import (
	"modules/misc/trans"

	"common/redis"
	"time"

	"common/assert"

	"modules/telegram/ad/ads"
	"strconv"

	"strings"

	"modules/telegram/teleuser/tlu"

	"github.com/Sirupsen/logrus"
	"gopkg.in/telegram-bot-api.v4"
)

func (bb *bot) doneORReject(bot *tgbotapi.BotAPI, m *tgbotapi.Message) {
	teleUser, err := tlu.NewTluManager().FindTeleUserByBotChatID(m.Chat.ID)
	assert.Nil(err)

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

	// bundleID,ChannelName
	channelBundleID := strings.Split(value, ",")
	if len(channelBundleID) != 2 {
		logrus.Panic("wrong redis insert in /getBundle command")
	}

	channel, err := ads.NewAdsManager().FindChannelByUserIDChannelName(teleUser.UserID, channelBundleID[1])
	assert.Nil(err)

	bundleID, err := strconv.ParseInt(channelBundleID[0], 10, 0)
	assert.Nil(err)

	bundle, err := ads.NewAdsManager().FindBundlesByID(bundleID)
	assert.Nil(err)

	adID := bundle.TargetAd

	err = ads.NewAdsManager().UpdateBundleChannelAd(&ads.BundleChannelAd{
		ChannelID: channel.ID,
		BundleID:  bundleID,
		AdID:      adID,
		Active:    ads.ActiveStatusNo,
	})
	assert.Nil(err)

	send(bot, m.Chat.ID, trans.T("successfully done"))
}
