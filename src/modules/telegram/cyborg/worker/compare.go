package worker

import (
	"common/utils"
	"modules/telegram/ad/ads"
	"modules/telegram/common/tgo"
	"modules/telegram/config"
)

func compareIndividual(src ads.ChannelAdD, dst tgo.History) bool {
	// the chat is forwarded from our bot or not?
	return dst.FwdFrom.Username == tcfg.Cfg.Telegram.BotName

	// TODO : compare text? if the text is comparable, and there is no link. the problem is that the telegram change
	// the links
}

func comparePromotion(src tgo.History, dst tgo.History) bool {
	if src.From.Username != dst.FwdFrom.Username {
		if src.FwdFrom != nil {
			if src.FwdFrom.Username != dst.FwdFrom.Username {
				return false
			}
		}
	}
	if src.Media != nil && dst.Media != nil {
		return src.Media.Caption == utils.RemoveEmojis(dst.Media.Caption)
	}
	return src.Text == utils.RemoveEmojis(dst.Text)
}
