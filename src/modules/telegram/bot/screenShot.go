package bot

import (
	"modules/telegram/common/tgbot"

	"common/assert"
	"modules/telegram/ad/ads"
	"modules/telegram/teleuser/tlu"

	"common/models/common"

	"modules/file/fila"

	"github.com/Sirupsen/logrus"
	"gopkg.in/telegram-bot-api.v4"
)

const (
	bundle = iota
	channelName
)

var (
	user    *tlu.TeleUser
	mainBot *tgbotapi.BotAPI
)

func (bb *bot) uploadSS(bot *tgbotapi.BotAPI, m *tgbotapi.Message) {
	mainBot = bot

	var err error
	user, err = tlu.NewTluManager().FindTeleUserByBotChatID(m.Chat.ID)
	assert.Nil(err)

	var channels []ads.Channel
	channels, err = ads.NewAdsManager().FindActiveChannelsByUserID(user.UserID)
	assert.Nil(err)

	for i := range channels {
		channelNames[m.Chat.ID] = append(channelNames[m.Chat.ID], channels[i].Name)
	}

	tgbot.StartConversion(bot, m.Chat.ID, tgbot.Screenshot)
}

func init() {
	// asking bundle code
	tgbot.RegisterConversation(tgbot.Screenshot, tgbot.Chatter{
		StringRequest: "Enter your bundle code",
		Handler: func(m *tgbotapi.Message, data map[int]interface{}) (string, bool) {
			if m.Text == "" {
				return "invalid input, try again", false
			}
			fetchedBundle, err := ads.NewAdsManager().FindBundlesByCode(m.Text)
			if err != nil {
				println(err.Error())
				return "no bundle code with this data was found, try again", false
			}
			data[bundle] = fetchedBundle
			return "bundle code added successfully", true
		},
	})

	// asking for channel name
	tgbot.RegisterConversation(tgbot.Screenshot, tgbot.Chatter{
		StringRequest:   "choose one of these as your channel",
		KeyboardRequest: channelNames,
		Handler: func(m *tgbotapi.Message, data map[int]interface{}) (string, bool) {
			data[channelName] = m.Text
			return "channel name added successfully", true
		},
	})

	// asking for screen shot
	tgbot.RegisterConversation(tgbot.Screenshot, tgbot.Chatter{
		StringRequest: "upload your screen shot",
		Handler: func(m *tgbotapi.Message, data map[int]interface{}) (string, bool) {
			channelName, ok := data[channelName].(string)
			assert.True(ok, "couldn't parse channel name to string")

			Fbundle, ok := data[bundle].(*ads.Bundles)
			assert.True(ok, "couldn't cast bundle from data")

			channel, err := ads.NewAdsManager().FindChannelByUserIDChannelName(user.UserID, channelName)
			if err != nil {
				return "could upload ur ss, try again", false
			}

			//TODO take care of slice
			fp, err := mainBot.GetFileDirectURL((*m.Photo)[3].FileID)
			logrus.Warn(fp)
			assert.Nil(err)
			url, err := fila.UploadFromURL(fp, channel.UserID)
			if err != nil {
				return "could upload ur ss, try again", false
			}

			bundleChannelAd := ads.NewAdsManager().FindBundleChannelAd(channel.ID, Fbundle.ID, Fbundle.TargetAd)
			bundleChannelAd.Shot = common.NullString{Valid: true, String: url}
			err = ads.NewAdsManager().UpdateBundleChannelAd(bundleChannelAd)
			if err != nil {
				return "could upload ur ss, try again", false
			}

			return "photo uploaded successfully", true
		},
	})
}
