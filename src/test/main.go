package main

import (
	"common/config"
	"common/initializer"
	"common/version"
	"log"
	"modules/bot"

	"gopkg.in/telegram-bot-api.v4"
)

func main() {
	config.Initialize()
	config.InitApplication()

	defer initializer.Initialize().Finalize()

	version.LogVersion().Infof("Application started")

	tBot, err := tgbotapi.NewBotAPI("315976738:AAGC_25yJ4jBN1zHzuR8fGlF_SXXNi6AXjI")
	if err != nil {
		log.Panic(err)
	}

	tBot.Debug = true

	log.Printf("Authorized on account %s", tBot.Self.UserName)
	verify.VerifyBot(tBot)
}
