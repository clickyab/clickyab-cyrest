package worker

import (
	"common/assert"
	"common/rabbit"
	"common/redis"
	"encoding/json"
	"errors"
	"fmt"
	"modules/telegram/ad/ads"
	"modules/telegram/common/tgo"
	"modules/telegram/config"
	"modules/telegram/cyborg/bot"
	"modules/telegram/cyborg/commands"
	"net"
	"sync"
	"time"

	"common/models/common"

	"sort"

	"regexp"
	"strconv"

	"github.com/Sirupsen/logrus"
)

// MultiWorker is a worker for all commands but share a single tcli client
type MultiWorker struct {
	client tgo.TelegramCli
	lock   *sync.Mutex
}

//ChnAdPattern is a pattern for message
var ChnAdPattern = regexp.MustCompile(`^([0-9]+)/([0-9]+)$`)

// Ping command verify if the client is alive
func (mw *MultiWorker) Ping() error {
	mw.lock.Lock()
	defer mw.lock.Unlock()

	u, err := mw.client.GetSelf()
	if err != nil {
		return err
	}

	logrus.Debugf("%+v", *u)
	return nil
}

func (mw *MultiWorker) discoverChannel(c string) (*tgo.ChannelInfo, error) {
	mw.lock.Lock()
	defer mw.lock.Unlock()

	ch, err := mw.client.ResolveUsername(c)
	if err != nil {
		return nil, err
	}
	if ch.PeerType != "channel" {
		return nil, errors.New("invalid channel type")
	}
	return mw.client.ChannelInfo(ch.ID)
}

func (mw *MultiWorker) getLastMessages(telegramID string, count int, offset int) ([]tgo.History, error) {
	mw.lock.Lock()
	defer mw.lock.Unlock()

	if count > 99 {
		count = 99
	}

	if count < 1 {
		count = 20
	}
	logrus.Warn(count, offset)
	return mw.client.MessageHistory(telegramID, count, offset)
}
func (mw *MultiWorker) fwdMessage(telegramID string, messageID string) (*tgo.SuccessResp, error) {
	mw.lock.Lock()
	defer mw.lock.Unlock()
	return mw.client.FwdMsg(telegramID, messageID)
}
func (mw *MultiWorker) sendMessage(telegramID string, messageID string) (*tgo.SuccessResp, error) {
	mw.lock.Lock()
	defer mw.lock.Unlock()
	return mw.client.Msg(telegramID, messageID)
}

func (mw *MultiWorker) getLast(in *commands.GetLastCommand) (bool, error) {
	// first try to resolve the channel
	m := bot.NewBotManager()
	c, err := m.FindKnownChannelByName(in.Channel)
	if err != nil {
		ch, err := mw.discoverChannel(in.Channel)
		if err != nil {
			// Oh crap. can not resolve this :/
			assert.Nil(aredis.StoreHashKey(in.HashKey, "STATUS", "failed", time.Hour))
			return false, err
		}

		c, err = bot.NewBotManager().CreateChannelByRawData(ch)
		if err != nil {
			assert.Nil(aredis.StoreHashKey(in.HashKey, "STATUS", "failed", time.Hour))
			return false, err
		}
	}
	h, err := mw.getLastMessages(c.CliTelegramID, in.Count, 0)
	if err != nil {
		assert.Nil(aredis.StoreHashKey(in.HashKey, "STATUS", "failed", time.Hour))
		return false, err
	}
	b, err := json.Marshal(h)
	assert.Nil(err)
	assert.Nil(aredis.StoreHashKey(in.HashKey, "DATA", string(b), time.Hour))
	assert.Nil(aredis.StoreHashKey(in.HashKey, "STATUS", "done", time.Hour))
	return false, nil
}

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

func (mw *MultiWorker) selectAd(in *commands.SelectAd) (bool, error) {
	b := ads.NewAdsManager()
	chad, err := b.FindChannelAdByAdIDActive(in.ChannelID)
	assert.Nil(err)
	if len(chad) > 0 {
		return false, nil
	}
	chooseAds, err := b.ChooseAd(in.ChannelID)
	logrus.Info("chooseAds", chooseAds)
	assert.Nil(err)
	if len(chooseAds) == 0 {
		return false, nil
	}
	for k := range chooseAds {
		chooseAds[k].AffectiveView = chooseAds[k].View - chooseAds[k].PossibleView
	}
	sort.Sort(ads.ByAffectiveView(chooseAds))
	//todo send ad to user
	return false, nil

}

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
		Num:        totalCount,
		TotalView:  sumView,
		ChannelID:  channel.ID,
	}
	err = ads.NewAdsManager().CreateChanDetail(cd)
	assert.Nil(err)
	rabbit.PublishAfter(in, 24*time.Hour)
	//ch, err := mw.discoverChannel(in.Channel)
	return false, nil

}
func (mw *MultiWorker) transaction(m *ads.Manager, chad *ads.ChannelAd, channelAdDetail *ads.ChannelAdDetail) (bool, error) {
	err := m.Begin()
	if err != nil {
		return true, err
	}
	defer func() {
		if err != nil {
			assert.Nil(m.Rollback())
		} else {
			err = m.Commit()
		}

		if err != nil {
			chad = nil
		}
	}()
	err = m.UpdateChannelAd(chad)
	if err != nil {
		return true, err
	}
	err = m.CreateChannelAdDetail(channelAdDetail)
	if err != nil {
		return true, err
	}
	return false, err
}

func (mw *MultiWorker) existChannelAdFor(h []tgo.History, cliMessageID common.NullString, len int) (int64, int64, int64) {
	var pos int64
	var view int64
	var warning int64
	warning = 1
	for k := range h {
		if h[k].Event == "message" && h[k].Service == false {
			if h[k].FwdFrom != nil {
				if h[k].ID == cliMessageID.String {
					warning = 0
					view = int64(h[k].Views)
					pos = int64(len - k)
				}
			}
		}
	}
	return pos, view, warning
}

func (mw *MultiWorker) existChannelAd(in *commands.ExistChannelAd) (bool, error) {
	m := ads.NewAdsManager()
	chad, err := m.FindChannelIDAdByAdID(in.ChannelID, in.AdID)
	assert.Nil(err)
	if chad.Active == "no" || !chad.Active.IsValid() {
		return false, nil
	}
	if chad.CliMessageID.Valid {
		rabbit.PublishAfter(&commands.ExistChannelAd{
			ChannelID: in.ChannelID,
			AdID:      in.AdID,
		}, tcfg.Cfg.Telegram.TimeReQueUe)
	}
	if !chad.End.Valid {
		return false, nil
	}
	channel, err := m.FindChannelByID(in.ChannelID)
	assert.Nil(err)
	c, err := bot.NewBotManager().FindKnownChannelByName(channel.Name)
	assert.Nil(err)
	a := ads.NewAdsManager()
	ad, err := a.FindAdByID(chad.AdID)
	assert.Nil(err)
	if err != nil {
		ch, err := mw.discoverChannel(channel.Name)
		assert.Nil(err)
		c, err = bot.NewBotManager().CreateChannelByRawData(ch)
		assert.Nil(err)
	}
	h, err := mw.getLastMessages(c.CliTelegramID, tcfg.Cfg.Telegram.LastPostChannel, 0)
	assert.Nil(err)
	chde, err := m.FindChanDetailByChannelID(chad.ChannelID)
	assert.Nil(err)
	//var len int64
	len := len(h)
	depos := tcfg.Cfg.Telegram.PositionAdDefault
	if ad.Position.Valid {
		depos = ad.Position.Int64
	}
	/*var pos int64
	var view int64
	var warning int64
	warning = 1*/
	var possible int64
	pos, view, warning := mw.existChannelAdFor(h, chad.CliMessageID, len)

	possible = int64(chde.TotalView / chde.Num)
	channelAdDetail := &ads.ChannelAdDetail{}
	chad.View = view
	chad.PossibleView = possible
	if pos < depos {
		warning = 1
	}
	chad.Warning = chad.Warning + warning
	if chad.Warning >= tcfg.Cfg.Telegram.LimitCountWarning {
		//todo send message to user
		chad.End = common.NullTime{Valid: true, Time: time.Now()}
	}

	channelAdDetail.AdID = chad.AdID
	channelAdDetail.AdID = chad.ChannelID
	channelAdDetail.Warning = warning
	channelAdDetail.Position = common.NullInt64{Valid: pos != 0, Int64: pos}
	channelAdDetail.View = view
	//transaction
	res, err := mw.transaction(m, chad, channelAdDetail)
	if res == true {
		return true, err
	}

	rabbit.PublishAfter(&commands.ExistChannelAd{
		ChannelID: in.ChannelID,
		AdID:      in.AdID,
	}, tcfg.Cfg.Telegram.TimeReQueUe)

	return false, nil
}

//updateMessage get channel id and read each post on it then if not save on db,
//save it
func (mw *MultiWorker) updateMessage(in *commands.UpdateMessage) (bool, error) {
	defer rabbit.MustPublishAfter(in, 2*time.Minute)
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
		codes := ChnAdPattern.FindStringSubmatch(h.Text)
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

// NewMultiWorker create a multi worker that listen on all commands
func NewMultiWorker(ip net.IP, port int) (*MultiWorker, error) {
	t, err := tgo.NewTelegramCli(ip, port)
	if err != nil {
		return nil, err
	}
	res := &MultiWorker{
		client: t,
		lock:   &sync.Mutex{},
	}
	if err := res.Ping(); err != nil {
		return nil, err
	}
	go rabbit.RunWorker(&commands.GetLastCommand{}, res.getLast, 1)
	go rabbit.RunWorker(&commands.GetChanCommand{}, res.getChanStat, 1)
	go rabbit.RunWorker(&commands.IdentifyAD{}, res.identifyAD, 1)
	go rabbit.RunWorker(&commands.ExistChannelAd{}, res.existChannelAd, 1)
	go rabbit.RunWorker(&commands.SelectAd{}, res.selectAd, 1)
	go rabbit.RunWorker(&commands.UpdateMessage{}, res.updateMessage, 1)
	return res, nil
}
