package bot

import (
	"modules/misc/trans"
	"modules/telegram/ad/ads"
	"modules/telegram/teleuser/tlu"

	"modules/telegram/common/tgbot"

	"common/assert"

	"common/models/common"
	"strings"

	"crypto/sha1"
	"encoding/hex"
	"fmt"

	"common/redis"
	"time"

	"gopkg.in/telegram-bot-api.v4"
)

const (
	bundleData = iota
)

func (bb *bot) sendAd(bot *tgbotapi.BotAPI, m *tgbotapi.Message) {
	mainBot = bot

	teleUser, err := tlu.NewTluManager().FindTeleUserByBotChatID(m.Chat.ID)
	assert.Nil(err)

	userChannels, err := ads.NewAdsManager().FindActiveChannelsByUserID(teleUser.UserID)
	if len(userChannels) < 1 || err != nil {
		send(bot, m.Chat.ID, trans.T("you have no active channel, try another time"))
		return
	}

	for i := range userChannels {
		channelNames[m.Chat.ID] = append(channelNames[m.Chat.ID], userChannels[i].Name)
	}

	tgbot.StartConversion(bot, m.Chat.ID, tgbot.GetBundle)
}

func init() {
	tgbot.RegisterConversation(tgbot.GetBundle, tgbot.Chatter{
		StringRequest: "Enter bundle code you want",
		Handler: func(message *tgbotapi.Message, data map[int]interface{}) (string, bool) {
			bundle, err := ads.NewAdsManager().FindBundlesByCode(message.Text)
			if err != nil {
				return "this bundle code were not found/ntry again", false
			}
			data[bundleData] = bundle
			return "", true
		},
	})

	// TODO not cheching for channel name existence
	tgbot.RegisterConversation(tgbot.GetBundle, tgbot.Chatter{
		StringRequest:   "choose the channel u want to post bundle in",
		KeyboardRequest: channelNames,
		Handler: func(message *tgbotapi.Message, data map[int]interface{}) (string, bool) {
			if message.Text == "" {
				return "wrong input", false
			}
			bc, ok := data[bundleData].(*ads.Bundles)
			assert.True(ok, "couldn't get bundle from data")

			ads := fetchAdsFromArrayString(bc.Ads)
			for i := range ads {
				RenderMessage(mainBot, message.Chat.ID, &ads[i])
			}

			key := fmt.Sprintf("%d,%s", bc.ID, message.Text)
			h := sha1.New()
			_, err := h.Write([]byte(key))
			if err != nil {
				return "couldn't set ur channel, try again", false
			}
			hash := hex.EncodeToString(h.Sum(nil))

			// added to redis
			err = aredis.StoreKey(hash, key, time.Hour)
			assert.Nil(err)

			response := fmt.Sprintf("use /done_%s when you're done", hash)
			return response, true
		},
	})
}

// fetches ads from a comma array string
func fetchAdsFromArrayString(input common.CommaArray) []ads.Ad {
	adIDs := strings.Split(string(input), ",")

	ads, err := ads.NewAdsManager().FindAdsByIDs(adIDs...)
	assert.Nil(err)

	return ads
}
