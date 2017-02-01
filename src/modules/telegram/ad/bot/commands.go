package bot

import (
	"common/initializer"
	"modules/telegram/common/tgbot"
)

type bot struct {
}

const htmlMode = "HTML"

func (bb *bot) Initialize() {

	tgbot.RegisterMessageHandler("/updatead", bb.updateAD)
	tgbot.RegisterMessageHandler("/ad", bb.wantAD)
	tgbot.RegisterMessageHandler("/confirm", bb.confirm)
	tgbot.RegisterMessageHandler("/done", bb.done)
}

func init() {
	initializer.Register(&bot{})
}
