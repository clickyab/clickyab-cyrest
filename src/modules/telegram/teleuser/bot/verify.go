package bot

import (
	"common/assert"
	"common/initializer"
	"common/models/common"
	"common/redis"
	"common/tgbot"
	"modules/telegram/teleuser/tlu"
	"strconv"
	"strings"
	"time"

	"modules/telegram/ad/ads"

	"gopkg.in/telegram-bot-api.v4"
)

const htmlMode = "HTML"

type bot struct {
}

func (bb *bot) verify(bot *tgbotapi.BotAPI, m *tgbotapi.Message) {
	//sample code  /verify-1:12123

	if !strings.Contains(m.Text, "-") && !strings.Contains(m.Text, ":") {
		msg := tgbotapi.NewMessage(m.Chat.ID, "your code is not <b>valid</b>")
		msg.ParseMode = htmlMode
		_, err := bot.Send(msg)
		assert.Nil(err)
		return
	}
	result := strings.Replace(m.Text, "/verify-", "", 1)
	str, err := aredis.GetKey(result, false, time.Hour)
	if str == "" || err != nil {
		msg := tgbotapi.NewMessage(m.Chat.ID, "your code is not <b>valid</b>")
		msg.ParseMode = htmlMode
		_, err := bot.Send(msg)
		assert.Nil(err)
		return
	}
	user := strings.Split(result, ":")
	id, err := strconv.ParseInt(user[0], 0, 10)
	if err == nil {
		msg := tgbotapi.NewMessage(m.Chat.ID, "your account is <b>accepted</b>")
		n := tlu.NewTluManager()
		tl := &tlu.TeleUser{
			UserID:     id,
			TelegramID: m.Chat.ID,
			Username:   common.MakeNullString(m.Chat.UserName),
			Remove:     tlu.RemoveStatusNo,
			Resolve:    tlu.ResolveStatusYes,
		}
		assert.Nil(n.CreateTeleUser(tl))
		msg.ParseMode = htmlMode
		_, err := bot.Send(msg)
		assert.Nil(err)
	}
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

	tgbot.RegisterMessageHandler("/verify", bb.verify)
	tgbot.RegisterMessageHandler("/updatead", bb.updateAD)
	//assert.Nil(b.Start())
}

func init() {
	initializer.Register(&bot{})
}
