package bot

import (
	"common/assert"
	"modules/telegram/common/tgbot"
	"modules/telegram/teleuser/tlu"
	"modules/user/aaa"
	"regexp"

	"common/models/common"

	"gopkg.in/telegram-bot-api.v4"
)

const (
	cardRegex    = `(\d{16})`
	accountRegex = `^IR(\d{24})$`

	cardQ = iota
)

var financialUser *aaa.UserFinancial

func (bb *bot) financial(bot *tgbotapi.BotAPI, m *tgbotapi.Message) {
	user, err := tlu.NewTluManager().FindTeleUserByBotChatID(m.Chat.ID)
	assert.Nil(err)
	financialUser, err = aaa.NewAaaManager().FindUserFinancialByID(user.ID)
	assert.Nil(err)

	tgbot.StartConversion(bot, m.Chat.ID, tgbot.Financial)
}

func init() {

	tgbot.RegisterConversation(tgbot.Financial, tgbot.Chatter{
		StringRequest: "card number added successfully",
		Handler: func(m *tgbotapi.Message, data map[int]interface{}) (string, bool) {
			valid := cardNumberValidator(m.Text)
			if valid {
				data[cardQ] = m.Text
				return "", true
			}
			return "invalid card number, try again", false
		},
	})

	tgbot.RegisterConversation(tgbot.Financial, tgbot.Chatter{
		StringRequest: "Enter your account number please",
		Handler: func(message *tgbotapi.Message, data map[int]interface{}) (string, bool) {
			valid := accountNumberValidator(message.Text)
			if valid {
				cardNumber, ok := data[bundle].(string)
				assert.True(ok, "couldn't cast card number to string")

				financialUser.CardNumber = common.MakeNullString(cardNumber)
				financialUser.AccountNumber = common.MakeNullString(message.Text)

				err := aaa.NewAaaManager().UpdateUserFinancial(financialUser)
				assert.Nil(err)
				return "financial details added successfully", true
			}
			return "invalid account number, try again", false
		},
	})

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
