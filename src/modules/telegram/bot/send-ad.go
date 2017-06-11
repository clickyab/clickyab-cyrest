package bot

//
//import (
//	"modules/misc/trans"
//	"modules/telegram/ad/ads"
//	"modules/telegram/teleuser/tlu"
//	"modules/user/aaa"
//
//	"strconv"
//
//	"github.com/Sirupsen/logrus"
//	"gopkg.in/telegram-bot-api.v4"
//)
//
//func (bb *bot) sendAd(bot *tgbotapi.BotAPI, m *tgbotapi.Message) {
//	regular := sendAd.FindStringSubmatch(m.Text)
//	logrus.Warn(regular, len(regular))
//	if len(regular) != 3 {
//		send(bot, m.Chat.ID, trans.T("your command is <botAPI>not valid1</botAPI>"))
//		return
//	}
//	bundleID, err := strconv.ParseInt(regular[1], 10, 0)
//	if err != nil {
//		send(bot, m.Chat.ID, trans.T("your command is <botAPI>not valid2</botAPI>"))
//		return
//	}
//	channelID, err := strconv.ParseInt(regular[2], 10, 0)
//	if err != nil {
//		send(bot, m.Chat.ID, trans.T("your command is <botAPI>not valid2</botAPI>"))
//		return
//	}
//	tele := tlu.NewTluManager()
//	telegramUser, err := tele.FindTeleUserByBotChatID(m.Chat.ID)
//	if err != nil {
//		send(bot, m.Chat.ID, trans.T("your telegram user is not in our system \n please register!"))
//		return
//	}
//	usr := aaa.NewAaaManager()
//	user, err := usr.FindUserByID(telegramUser.UserID)
//	if err != nil {
//		send(bot, m.Chat.ID, trans.T("your telegram user is not in our system \n please register!"))
//		return
//	}
//	//find bundleChannel
//	adManager := ads.NewAdsManager()
//	_, err = adManager.FindChannelByUserIDChannelID(telegramUser.UserID, channelID)
//	if err != nil {
//		send(bot, m.Chat.ID, trans.T("this channel not for you"))
//	}
//	bundleChannel := adManager.FindActiveChannelBundleByUserID(user.ID, bundleID, channelID)
//	if len(bundleChannel) > 0 {
//		send(bot, m.Chat.ID, trans.T("you have active bundle ad for your channel"))
//	}
//	send(bot, m.Chat.ID, trans.T("please send below ads to your channel"))
//	for i := range bundleChannel {
//		RenderMessage(bot, bundleChannel[i].ChannelID, &bundleChannel[i].Ad)
//		//textMsg += fmt.Sprintf("\n/complete_%d_%d\n", bundleChannel[i].ChannelID, bundleChannel[i].ID)
//	}
//	send(bot, m.Chat.ID, trans.T("after send ads enter /done_%d_%d", bundleID, channelID))
//	return
//}
