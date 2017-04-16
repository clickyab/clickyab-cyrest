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
	t, err := mw.returnFwdMessage(tcfg.Cfg.Telegram.BotID, ad.CliMessageID.String)
	assert.Nil(err)
	//update promote data
	b, err := json.Marshal(t)
	assert.Nil(err)
	promote := common.MakeNullString(string(b))
	err = ads.NewAdsManager().UpdateAdPromote(ad.ID, promote)
	assert.Nil(err)
	return false, nil
}
