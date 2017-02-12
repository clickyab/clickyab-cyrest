package worker

import (
	"modules/telegram/ad/ads"
	"modules/telegram/common/tgo"
	"modules/telegram/config"
)

func compareIndividual(src ads.ChannelAdD, dst tgo.History) bool {
	// the chat is forwarded from our bot or not?
	if dst.FwdFrom.Username != tcfg.Cfg.Telegram.BotName {
		return false
	}

	// TODO : compare text? if the text is comparable, and there is no link. the problem is that the telegram change
	// the links

	return true
}

func comparePromotion(src tgo.History, dst tgo.History) bool {
	if src.From.Username != dst.FwdFrom.Username {
		if src.FwdFrom.Username != dst.FwdFrom.Username {
			return false
		}
	}

	return src.Text == dst.Text
}
