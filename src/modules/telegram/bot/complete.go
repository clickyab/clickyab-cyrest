package bot

import (
	"common/models/common"
	"modules/misc/trans"
	"modules/telegram/ad/ads"
	"modules/telegram/teleuser/tlu"
	"modules/user/aaa"
	"time"

	"strconv"

	"modules/billing/bil"

	"github.com/Sirupsen/logrus"
	"gopkg.in/telegram-bot-api.v4"
)

func (bb *bot) complete(botAPI *tgbotapi.BotAPI, m *tgbotapi.Message) {
	regular := completeAd.FindStringSubmatch(m.Text)
	logrus.Warn(regular, len(regular))
	if len(regular) != 3 {
		send(botAPI, m.Chat.ID, trans.T("your command is <botAPI>not valid1</botAPI>"))
		return
	}
	channelID, err := strconv.ParseInt(regular[1], 10, 0)
	if err != nil {
		send(botAPI, m.Chat.ID, trans.T("your command is <botAPI>not valid2</botAPI>"))
		return
	}
	adID, err := strconv.ParseInt(regular[2], 10, 0)
	if err != nil {
		send(botAPI, m.Chat.ID, trans.T("your command is <botAPI>not valid2</botAPI>"))
		return
	}
	tele := tlu.NewTluManager()
	telegramUser, err := tele.FindTeleUserByBotChatID(m.Chat.ID)
	if err != nil {
		send(botAPI, m.Chat.ID, trans.T("your telegram user is not in our system \n please register!"))
		return
	}
	usr := aaa.NewAaaManager()
	user, err := usr.FindUserByID(telegramUser.UserID)
	if err != nil {
		send(botAPI, m.Chat.ID, trans.T("your telegram user is not in our system \n please register!"))
		return
	}
	//find channelAdActive
	adManager := ads.NewAdsManager()
	channelAdActive, err := adManager.FindChannelIDAdByAdIDByActive(channelID, adID, user.ID)
	if err != nil {
		send(botAPI, m.Chat.ID, trans.T("no active ad for you"))
	}

	ch := ads.ChannelAd{}
	ch.Active = ads.ActiveStatusNo
	ch.End = common.MakeNullTime(time.Now())
	err = adManager.UpdateChannelAd(&ch)

	if err != nil {
		send(botAPI, m.Chat.ID, trans.T("have problem exist in action please try later"))
	}

	bi := bil.NewBilManager()
	err = bi.ChannelBilling(channelAdActive)
	if err != nil {
		send(botAPI, m.Chat.ID, trans.T("have problem exist in action please try later"))
	}

	textMsg := trans.T("you have %d view in channel \n please remove all ads from your channel", channelAdActive.View).Translate(trans.PersianLang)
	sendString(botAPI, m.Chat.ID, textMsg)
	return
}
