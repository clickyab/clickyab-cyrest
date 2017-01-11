package tgbot

import "common/initializer"
import "common/config"

type tgbotInitializer struct {
}

var handler TelegramBot

func (*tgbotInitializer) Initialzie() {

}

// RegisterMessageHandler try to register a handler in system, the first is the command to match the
// next arg is the handler function
func RegisterMessageHandler(s string, hm HandleMessage) error {
	return handler.RegisterMessageHandler(s, hm)
}

// Start the handler
func Start() error {
	return handler.Start()
}

func init() {
	handler = NewTelegramBot(config.Config.Telegram.APIKey)
	initializer.Register(&tgbotInitializer{})
}
