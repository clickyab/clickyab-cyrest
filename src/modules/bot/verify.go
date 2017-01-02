package verify

import (
	"modules/channel/chn"
	"strconv"
	"strings"

	"regexp"

	"gopkg.in/telegram-bot-api.v4"
)

var Verify = regexp.MustCompile(`\/verify (\d)+-(\d)+`)

func VerifyBot(bot *tgbotapi.BotAPI) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	message := "<b>your channel has been successfully verified</b>"

	updates, _ := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() && Verify.MatchString(update.Message.Text) {

			code := update.Message.CommandArguments()
			params := strings.Split(code, "-")
			if len(params) != 2 {
				continue

			}
			channel_id, err := strconv.ParseInt(params[0], 10, 0)
			if err != nil {
				continue
			}

			//get the channel from database
			chManager := chn.NewChnManager()
			channel, err := chManager.FindChannelByID(channel_id)
			if err != nil {
				continue
			}
			//check if the channel belongs to user and the code correct
			if channel.Admin.String == update.Message.From.UserName && channel.Code == code {
				//check if the channel status is pending
				if channel.Status == chn.ChannelStatusPending {
					channel.Status = chn.ChannelStatusAccepted
					err = chManager.UpdateChannel(channel)
					if err != nil {
						continue
					}
				} else {
					continue
				}
			} else {
				continue
			}
			//get
			/*username:=update.Message.From.UserName*/

			//check if the channel username exists in the channel table

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, message)
			msg.ParseMode = "HTML"
			bot.Send(msg)
		} else {
			continue
		}

	}
}
