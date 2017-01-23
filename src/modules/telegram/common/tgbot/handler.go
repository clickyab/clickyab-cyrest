package tgbot

import (
	"common/assert"
	"common/initializer"
	"modules/telegram/config"
	"time"

	"gopkg.in/telegram-bot-api.v4"
)

type tgbotInitializer struct {
}

var handler TelegramBot

func (*tgbotInitializer) Initialize() {
	handler = NewTelegramBot(tcfg.Cfg.Telegram.APIKey)
}

// RegisterMessageHandler try to register a handler in system, the first is the command to match the
// next arg is the handler function
func RegisterMessageHandler(s string, hm HandleMessage) error {
	return handler.RegisterMessageHandler(s, hm)
}

// RegisterUserHandler redirect all user message to a chat
func RegisterUserHandler(i int64, u HandleMessage, t time.Duration) {
	handler.RegisterUserHandler(i, u, t)
}

// UnRegisterUserHandler redirect all user message to a chat
func UnRegisterUserHandler(i int64) {
	handler.UnRegisterUserHandler(i)
}

// Start the handler
func Start() error {
	return handler.Start()
}

// Send a message using the global handler
func Send(c tgbotapi.Chattable) (tgbotapi.Message, error) {
	assert.NotNil(handler, "[BUG] call this to early")
	return handler.Send(c)
}

func init() {
	initializer.Register(&tgbotInitializer{})
}
