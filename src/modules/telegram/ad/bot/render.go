package bot

import (
	"common/assert"
	"modules/telegram/ad/ads"

	"common/config"
	"path/filepath"
	"strings"

	"gopkg.in/telegram-bot-api.v4"
)

func forwardCli(bot *tgbotapi.BotAPI, chatID int64, ad *ads.Ad) tgbotapi.Message {
	assert.True(ad.BotChatID.Valid, "[BUG] not yet checked by the bot")
	assert.True(ad.BotMessageID.Valid, "[BUG] not yet checked by the bot")
	msg := tgbotapi.NewForward(chatID, ad.BotChatID.Int64, int(ad.BotMessageID.Int64))
	x, err := bot.Send(msg)
	assert.Nil(err)
	return x
}

func createMessage(bot *tgbotapi.BotAPI, chatID int64, ad *ads.Ad) tgbotapi.Message {
	f := filepath.Join(config.Config.StaticRoot, ad.Src.String)
	ext := strings.ToLower(filepath.Ext(f))
	var chat tgbotapi.Chattable
	if ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif" {
		ph := tgbotapi.NewPhotoUpload(chatID, f)
		ph.Caption = ad.Description.String
		chat = ph
	} else if ext == ".mov" || ext == ".mp4" {
		vd := tgbotapi.NewVideoUpload(chatID, f)
		vd.Caption = ad.Description.String
		chat = vd
	}

	assert.NotNil(chat, "[BUG] Unhandled ext ")

	x, err := bot.Send(chat)
	assert.Nil(err)
	return x

}

//RenderMessage  is a sender that send message depend type of message forward or create
func RenderMessage(bot *tgbotapi.BotAPI, chatID int64, ad *ads.Ad) tgbotapi.Message {
	if ad.CliMessageID.Valid {
		return forwardCli(bot, chatID, ad)
	}
	return createMessage(bot, chatID, ad)
}
