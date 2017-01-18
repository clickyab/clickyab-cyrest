package worker

import (
	"common/assert"
	"common/config"
	"common/rabbit"
	"common/redis"
	"common/tgo"
	"encoding/json"
	"errors"
	"fmt"
	"modules/ad/ads"
	"modules/channel/chn"
	"modules/cyborg/bot"
	"modules/cyborg/commands"
	"net"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
)

// MultiWorker is a worker for all commands but share a single tcli client
type MultiWorker struct {
	client tgo.TelegramCli
	lock   *sync.Mutex
}

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
	return mw.client.History(telegramID, count, offset)
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
	h, err := mw.getLastMessages(c.TelegramID, in.Count, 0)
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
	ad, err := m.FindAdByID(in.AddID)
	assert.Nil(err)
	if !ad.CliMessageID.Valid {
		return false, nil
	}
	_, err = mw.sendMessage(config.Config.Telegram.BotID, fmt.Sprintf("/updatead-%d", in.AddID))
	assert.Nil(err)
	_, err = mw.fwdMessage(config.Config.Telegram.BotID, ad.CliMessageID.String)
	assert.Nil(err)
	return false, nil
}

func (mw *MultiWorker) getChanStat(in *commands.GetChanCommand) (bool, error) {
	//find channel
	chnManager := chn.NewChnManager()
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
	h, err := mw.getLastMessages(c.TelegramID, in.Count, 0)
	if err != nil {
		return false, err
	}
	for i := range h {
		if h[i].FwdFrom == nil {
			sumView += h[i].Views
			totalCount++
		}

	}
	cd := &bot.ChanDetail{
		Name:       c.Name,
		Title:      c.Title,
		Info:       c.Info,
		UserCount:  c.UserCount,
		TelegramID: c.TelegramID,
		AdminCount: c.RawData.AdminsCount,
		Num:        totalCount,
		TotalView:  sumView,
		ChannelID:  channel.ID,
	}
	err = knownManger.CreateChanDetail(cd)
	assert.Nil(err)
	rabbit.PublishAfter(in, 24*time.Hour)
	//ch, err := mw.discoverChannel(in.Channel)
	return false, nil

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
	go rabbit.RunWorker(&commands.IdentifyAD{}, res.identifyAD, 1)
	go rabbit.RunWorker(&commands.GetChanCommand{}, res.getChanStat, 1)
	return res, nil
}
