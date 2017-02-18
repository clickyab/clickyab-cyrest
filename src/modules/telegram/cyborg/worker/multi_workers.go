package worker

import (
	"common/assert"
	"common/models/common"
	"common/rabbit"
	"common/utils"
	"errors"
	"modules/telegram/common/tgo"
	"modules/telegram/cyborg/commands"
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

var (
	once = &sync.Once{}
)

type channelDetailStat struct {
	frwrd bool
	cliID common.NullString
	adID  int64
}

type channelViewStat struct {
	adID    int64
	pos     int64
	view    int64
	warning int64
	frwrd   bool
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
	go rabbit.RunWorker(&commands.DiscoverAd{}, res.discoverAd, 1)

	once.Do(func() {

		go utils.SafeGO(func() {
			for {
				assert.Nil(res.cronReview())
				<-time.After(1 * time.Minute)
			}
		}, true)
	})

	return res, nil
}
