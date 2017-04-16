package worker

import (
	"common/assert"
	"common/models/common"
	"encoding/json"
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
	res, err := mw.getLastMessages(tcfg.Cfg.Telegram.BotID, 1, 0)
	assert.Nil(err)
	//update promote data
	b, err := json.Marshal(res[0])
	assert.Nil(err)
	ad.PromoteData = common.MakeNullString(string(b))
	assert.Nil(err)
	err = ads.NewAdsManager().UpdateAd(ad)
	assert.Nil(err)
	return false, nil
}
