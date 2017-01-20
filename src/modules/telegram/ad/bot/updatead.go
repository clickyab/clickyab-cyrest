package bot

import (
	"common/assert"
	"common/initializer"
	"common/models/common"
	"modules/telegram/ad/ads"
	"modules/telegram/common/tgbot"
	"strconv"
	"strings"
	"time"

	"gopkg.in/telegram-bot-api.v4"
)

type bot struct {
}

func (bb *bot) updateAD(bot *tgbotapi.BotAPI, m *tgbotapi.Message) {
	result := strings.Replace(m.Text, "/updatead-", "", 1)
	id, err := strconv.ParseInt(result, 0, 10)
	if err == nil {
		tgbot.RegisterUserHandler(m.Chat.ID, func(bot *tgbotapi.BotAPI, m *tgbotapi.Message) {
			defer tgbot.UnRegisterUserHandler(m.Chat.ID)
			botChatID := strconv.FormatInt(m.Chat.ID, 10)
			botMsgID := strconv.Itoa(m.MessageID)
			n := ads.NewAdsManager()
			currentAd, err := n.FindAdByID(id)
			assert.Nil(err)
			currentAd.BotChatID = common.MakeNullString(botChatID)
			currentAd.BotMessageID = common.MakeNullString(botMsgID)
			assert.Nil(n.UpdateAd(currentAd))

		}, time.Minute)
	}

}

func (bb *bot) Initialize() {

	tgbot.RegisterMessageHandler("/updatead", bb.updateAD)
	//assert.Nil(b.Start())
}

func init() {
	initializer.Register(&bot{})
}
