package worker

import (
	"common/assert"
	"fmt"
	"modules/telegram/ad/ads"
	"modules/telegram/config"
	"modules/telegram/cyborg/commands"
)

func (mw *MultiWorker) identifyAD(in *commands.IdentifyAD) (bool, error) {
	// first try to resolve the channel
	m := ads.NewAdsManager()
	ad, err := m.FindAdByID(in.AdID)
	assert.Nil(err)
	if !ad.CliMessageID.Valid {
		return false, nil
	}
	_, err = mw.sendMessage(tcfg.Cfg.Telegram.BotID, fmt.Sprintf("/updatead-%d", in.AdID))
	assert.Nil(err)
	_, err = mw.fwdMessage(tcfg.Cfg.Telegram.BotID, ad.CliMessageID.String)
	assert.Nil(err)
	return false, nil
}
