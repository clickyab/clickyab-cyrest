package worker

import (
	"common/assert"
	"common/rabbit"
	"modules/telegram/ad/ads"
	"modules/telegram/cyborg/bot"
	"modules/telegram/cyborg/commands"
	"time"
)

func (mw *MultiWorker) getChanStat(in *commands.GetChanCommand) (bool, error) {
	//find channel
	chnManager := ads.NewAdsManager()
	channel, err := chnManager.FindChannelByID(in.ChannelID)
	assert.Nil(err)
	//check if rhe channel exists in known channel
	knownManger := bot.NewBotManager()
	c, err := knownManger.FindKnownChannelByName(channel.Name)
	if err != nil {
		//known channel not found
		ch, err := mw.discoverChannel(channel.Name)
		if err != nil {
			// Oh crap. can not resolve this :/
			return false, err
		}

		c, err = bot.NewBotManager().CreateChannelByRawData(ch)
		if err != nil {
			return false, err
		}
	}
	var sumView int
	var totalCount int
	h, err := mw.getLastMessages(c.CliTelegramID, in.Count, 0)
	if err != nil {
		return false, err
	}
	for i := range h {
		if h[i].FwdFrom == nil {
			sumView += h[i].Views
			totalCount++
		}

	}
	cd := &ads.ChanDetail{
		Name:       c.Name,
		Title:      c.Title,
		Info:       c.Info,
		UserCount:  c.UserCount,
		TelegramID: c.CliTelegramID,
		AdminCount: c.RawData.AdminsCount,
		PostCount:  totalCount,
		TotalView:  sumView,
		ChannelID:  channel.ID,
	}
	err = ads.NewAdsManager().UpdateOnDuplicateChanDetail(cd)
	assert.Nil(err)
	rabbit.PublishAfter(in, 24*time.Hour)
	//ch, err := mw.discoverChannel(in.Channel)
	return false, nil

}
