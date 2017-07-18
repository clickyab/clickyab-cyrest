package bot

import (
	"modules/telegram/common/tgbot"
	"modules/telegram/teleuser/tlu"
	"modules/user/aaa"

	"regexp"

	"common/assert"

	"gopkg.in/telegram-bot-api.v4"
)

const (
	financial = `financial`
)

var (
	cardRegex    = `(\d{16})`
	accountRegex = `^IR(\d{24})$`
)

func (bb *bot) financial(bot *tgbotapi.BotAPI, m *tgbotapi.Message) {
	user, err := tlu.NewTluManager().FindTeleUserByBotChatID(m.Chat.ID)
	assert.Nil(err)
	_, err = aaa.NewAaaManager().FindUserFinancialByID(user.ID)
	if err != nil {
		tgbot.RegisterConversation(financial, tgbot.Chatter{
			Request: "Enter your card number please",
			Handler: bb.getCard,
		})

		tgbot.RegisterConversation(financial, tgbot.Chatter{
			Request: "Enter your account number please",
			Handler: bb.getAccount,
		})

		tgbot.NewConversation(financial, bot, m.Chat.ID)
	} /*else {
		keyboard := tgbot.NewKeyboard("/agree", "/change")
		sendWithKeyboard(bot, keyboard, m.Chat.ID, trans.T("you already have an account, agree or change it"))
		var cardNumber string
		getCardNumber(bot, m, user, &cardNumber)
		getAccountNumber(bot, m, user, &cardNumber)
	}*/
}

func (bb *bot) getCard(response string) (string, bool) {
	valid := cardNumberValidator(response)
	if valid {
		return "card number added successfully", true
	}
	return "invalid card number, try again", false

}

func (bb *bot) getAccount(response string) (string, bool) {
	valid := accountNumberValidator(response)
	if valid {
		return "Enter your account number please", true
	}
	return "invalid account number, try again", false
}

// TODO take care
func cardNumberValidator(number string) bool {
	valid, err := regexp.MatchString(cardRegex, number)
	return valid && err == nil
}

func accountNumberValidator(number string) bool {
	valid, err := regexp.MatchString(accountRegex, number)
	return valid && err == nil
}
