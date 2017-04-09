package bot

import (
	"fmt"
	"modules/misc/trans"
	"modules/telegram/ad/ads"
	"modules/telegram/teleuser/tlu"
	"modules/user/aaa"

	"gopkg.in/telegram-bot-api.v4"
)

func (bb *bot) activeAd(bot *tgbotapi.BotAPI, m *tgbotapi.Message) {
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
	//find channelAds
	adManager := ads.NewAdsManager()
	channelAds := adManager.FindActiveChannelAdByUserID(user.ID)
	if len(channelAds) == 0 {
		send(bot, m.Chat.ID, trans.T("no active ad for you"))
	}
	textMsg := trans.T("for complete end show ad & calculate price ad").Translate(trans.PersianLang)
	for i := range channelAds {
		RenderMessage(bot, channelAds[i].ChannelID, &channelAds[i].Ad)
		textMsg += fmt.Sprintf("\n/complete_%d_%d\n", channelAds[i].ChannelID, channelAds[i].ID)
	}
	sendString(bot, m.Chat.ID, textMsg)
	return
}
