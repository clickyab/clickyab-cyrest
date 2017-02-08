package worker

import (
	"common/assert"
	"common/models/common"
	"fmt"
	"modules/telegram/ad/ads"
	bot2 "modules/telegram/ad/bot/worker"
	"modules/telegram/cyborg/bot"
	"modules/telegram/cyborg/commands"
	"strconv"
	"time"
)

//UpdateMessage get channel id and read each post on it then if not save on db,
//save it
func (mw *MultiWorker) UpdateMessage(in *commands.UpdateMessage) (bool, error) {
	knownManger := bot.NewBotManager()
	c, err := knownManger.FindKnownChannelByName(in.CLiChannelName)
	if err != nil {
		//known channel not found
		ch, err := mw.discoverChannel(in.CLiChannelName)

		if err != nil {
			// Oh crap. can not resolve this :/
			return false, err
		}
		c, err = bot.NewBotManager().CreateChannelByRawData(ch)
		if err != nil {
			return false, err
		}
	}
	caManager := ads.NewAdsManager()

	history, err := mw.getLastMessages(c.CliTelegramID, in.Count, in.Offset)
	assert.Nil(err)

	if len(history) == 0 {
		return true, nil
	}
	for i, h := range history {
		codes := chnAdPattern.FindStringSubmatch(h.Text)
		if len(codes) == 0 {
			continue
		}
		adID, err := strconv.ParseInt(codes[1], 10, 0)
		if err != nil {
			//logrus.Warn(err)
			continue
		}
		channelID, err := strconv.ParseInt(codes[2], 10, 0)
		if err != nil {
			//logrus.Warn(err)
			continue
		}

		chn, err := caManager.FindChannelIDAdByAdID(adID, channelID)
		if err != nil {
			//logrus.Warn(err)
			continue
		}
		if chn.CliMessageID.Valid && chn.CliMessageID.String == h.ID {
			break

		}
		chn.CliMessageID = common.MakeNullString(history[i-1].ID)

		assert.Nil(caManager.UpdateChannelAd(chn))

	}
	return false, err
}

// CronReview cron review for finished ads
func (mw *MultiWorker) CronReview(in *commands.UpdateMessage) (bool, error) {
	m := ads.NewAdsManager()
	activeIndividualAds, err := m.SelectIndividualActiveAd()
	assert.Nil(err)
	if len(activeIndividualAds) <= 0 {
		return false, nil
	}

	for k := range activeIndividualAds {
		activeIndividualAds[k].Ad.View.Int64 = activeIndividualAds[k].Viewed
		activeIndividualAds[k].Ad.View.Valid = true
		err = m.UpdateAd(&activeIndividualAds[k].Ad)
		assert.Nil(err)
	}
	allActiveAds, err := m.SelectAdsPlan()
	assert.Nil(err)
	for key := range allActiveAds {
		if allActiveAds[key].Viewed < allActiveAds[key].Ad.View.Int64 {
			channelAd, err := m.FindChannelAdActiveByAdID(allActiveAds[key].Ad.ID, ads.ActiveStatusYes)
			assert.Nil(err)
			for c := range channelAd {
				//todo send message to admin channel
				channelAd[c].ChannelAd.Active = ads.ActiveStatusNo
				channelAd[c].ChannelAd.End = common.MakeNullTime(time.Now())
				assert.Nil(m.UpdateChannelAd(&channelAd[c].ChannelAd))
				str := fmt.Sprintf("you have have <%d> in channel \n please remove ad from <%s> channel", channelAd[c].ChannelAd.View, channelAd[c].Channel.Name)
				bot2.SendWarnAction(&bot2.SendWarn{
					AdID:      channelAd[c].ChannelAd.AdID,
					ChannelID: channelAd[c].ChannelAd.ChannelID,
					Msg:       str,
					ChatID:    channelAd[c].ChannelAd.BotChatID,
				})

			}
		}
	}

	return false, nil

}