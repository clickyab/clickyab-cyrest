package worker

import (
	"common/assert"
	"common/models/common"
	"common/rabbit"
	"encoding/json"
	"errors"
	"fmt"
	"modules/telegram/ad/ads"
	bot3 "modules/telegram/ad/worker"
	bot2 "modules/telegram/bot/worker"
	"modules/telegram/common/tgo"
	"modules/telegram/cyborg/bot"
	"net"
	"sort"
	"sync"
	"time"

	"modules/telegram/cyborg/commands"

	"common/redis"

	"modules/telegram/config"

	"common/config"
	"common/utils"
	"modules/misc/trans"

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

func (mw *MultiWorker) discoverMessage(msgID string) (*tgo.History, error) {
	mw.lock.Lock()
	defer mw.lock.Unlock()

	return mw.client.GetMessage(msgID)
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
	h, err := mw.getLastMessages(c.CliTelegramID, 99, 0)
	if len(h) > in.Count {
		h = h[:in.Count]
	}
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
	chad, err := b.FindChannelAdActiveByChannelID(in.ChannelID, ads.ActiveStatusYes)
	assert.Nil(err)
	if len(chad) > 0 {
		return false, nil
	}
	chooseAds, err := b.ChooseAd(in.ChannelID)
	assert.Nil(err)
	if len(chooseAds) == 0 {
		rabbit.MustPublish(&bot2.SendWarn{
			AdID:      0,
			ChannelID: in.ChannelID,
			Msg:       trans.T("no ads for you"),
			ChatID:    in.ChatID,
		})
		return false, nil
	}
	for k := range chooseAds {
		chooseAds[k].AffectiveView = chooseAds[k].PlanView - chooseAds[k].PossibleView.Int64
	}
	sort.Sort(ads.ByAffectiveView(chooseAds))
	logrus.Infof("%+v", chooseAds)
	var (
		promoted int64
		normal   int64
	)

	for i := range chooseAds {
		if promoted == 0 && chooseAds[i].CliMessageID.Valid {
			promoted = chooseAds[i].ID
		}
		if normal == 0 && !chooseAds[i].CliMessageID.Valid {
			normal = chooseAds[i].ID
		}

		if promoted != 0 && normal != 0 {
			break
		}
	}

	if normal == 0 {
		rabbit.MustPublish(&bot2.SendWarn{
			AdID:      0,
			ChannelID: in.ChannelID,
			Msg:       trans.T("you have an active ad"),
		})
		return false, nil
	}

	adList := []int64{}
	if promoted != 0 {
		adList = append(adList, promoted)
	}
	adList = append(adList, normal)
	logrus.Info(adList)
	//todo send ad to user
	rabbit.MustPublish(&bot3.AdDelivery{
		AdsID:     adList,
		ChannelID: in.ChannelID,
		ChatID:    in.ChatID,
	})
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
func (mw *MultiWorker) transaction(m *ads.Manager, chad []ads.ChannelAd, channelAdDetail []*ads.ChannelAdDetail, avg int64) (bool, error) {
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
	err = m.CreateChannelAdDetails(channelAdDetail)
	assert.Nil(err)
	if err != nil {
		return true, err
	}
	err = m.UpdateChannelAds(chad)
	assert.Nil(err)
	//update ads table field view for promotion ad
	for ch := range chad {
		if chad[ch].CliMessageID.Valid {
			//get ad
			ad, err := m.FindAdByID(chad[ch].AdID)
			if err != nil {
				return true, err
			}
			//ad.View = common.NullInt64{Valid: avg != 0, Int64: avg}
			ad.View = common.NullInt64{Valid: avg != 0, Int64: avg}
			assert.Nil(m.UpdateAd(ad))
		}
	}

	if err != nil {
		return true, err
	}
	return false, err
}

func (mw *MultiWorker) existChannelAdFor(h []tgo.History, adConfs []channelDetailStat) (map[int64]channelViewStat, int64) {
	var finalResult = make(map[int64]channelViewStat)
	var sumNotpromotionView int64
	var countNotPromotion int64
	historyLen := len(h)
	for k := range h {
		if h[k].Event == "message" && h[k].Service == false {
			if h[k].FwdFrom != nil {
				for i := range adConfs {
					if adConfs[i].frwrd && h[k].ID == adConfs[i].cliID.String { //the ad is forward type
						finalResult[adConfs[i].adID] = channelViewStat{
							view:    int64(h[k].Views),
							warning: 0,
							pos:     int64(historyLen - k),
							frwrd:   true,
							adID:    adConfs[i].adID,
						}
					} else if !adConfs[i].frwrd && h[k].ID == adConfs[i].cliID.String { //the ad is  not forward type
						finalResult[adConfs[i].adID] = channelViewStat{
							view:    int64(h[k].Views),
							warning: 0,
							pos:     int64(historyLen - k),
							frwrd:   false,
							adID:    adConfs[i].adID,
						}
						sumNotpromotionView += int64(h[k].Views)
						countNotPromotion++
					} else { //don't find ad in the history
						finalResult[adConfs[i].adID] = channelViewStat{
							view:    0,
							warning: 1,
							frwrd:   adConfs[i].frwrd,
							adID:    adConfs[i].adID,
							pos:     0,
						}
					}
				}
			}
		}
	}
	if countNotPromotion == 0 {
		return finalResult, 0
	}
	return finalResult, (sumNotpromotionView) / (countNotPromotion)
}

func (mw *MultiWorker) existChannelAd(in *commands.ExistChannelAd) (bool, error) {
	logrus.Warn("existChannelAdFor")
	var adsConf []channelDetailStat
	m := ads.NewAdsManager()

	chads, err := m.FindChannelAdByChannelIDActive(in.ChannelID)
	assert.Nil(err)
	for i := range chads {
		adsConf = append(adsConf, channelDetailStat{
			cliID: chads[i].CliMessageID,
			frwrd: chads[i].CliMessageAd.Valid,
			adID:  chads[i].AdID,
		})
		assert.True(chads[i].CliMessageID.Valid, "cli not filled")
		if !chads[i].CliMessageID.Valid {
			rabbit.PublishAfter(&commands.ExistChannelAd{
				ChannelID: in.ChannelID,
				ChatID:    in.ChatID,
			}, tcfg.Cfg.Telegram.TimeReQueUe)
			return false, nil
		}
	}

	//check for promotion to be alone or not
	var promotionCount int
	var notPromotionCount int
	for adConf := range adsConf {
		if adsConf[adConf].frwrd {
			promotionCount++
		}
		notPromotionCount++
	}

	if notPromotionCount == 0 {

		for adConf := range adsConf {
			//send stop (warn message)
			rabbit.MustPublish(&bot2.SendWarn{
				AdID:      adsConf[adConf].adID,
				ChannelID: in.ChannelID,
				Msg:       trans.T("please remove the following ad"),
			})

		}
		return false, nil
	}

	defer rabbit.PublishAfter(&commands.ExistChannelAd{
		ChannelID: in.ChannelID,
		ChatID:    in.ChatID,
	}, tcfg.Cfg.Telegram.TimeReQueUe)

	channel, err := m.FindChannelByID(in.ChannelID)
	assert.Nil(err)
	c, err := bot.NewBotManager().FindKnownChannelByName(channel.Name)
	if err != nil {
		ch, err := mw.discoverChannel(channel.Name)
		assert.Nil(err)
		c, err = bot.NewBotManager().CreateChannelByRawData(ch)
		assert.Nil(err)
	}
	h, err := mw.getLastMessages(c.CliTelegramID, tcfg.Cfg.Telegram.LastPostChannel, 0)
	assert.Nil(err)
	/*channelDetails, err := m.FindChanDetailByChannelID(channel.ID)
	assert.Nil(err)*/
	channelAdStat, avg := mw.existChannelAdFor(h, adsConf)

	var ChannelAdDetailArr []*ads.ChannelAdDetail
	for j := range chads {
		var currentView int64
		depos := tcfg.Cfg.Telegram.PositionAdDefault
		if chads[j].AdPosition.Valid {
			depos = chads[j].AdPosition.Int64
		}
		if channelAdStat[chads[j].AdID].pos < depos {
			channelAdStat[chads[j].AdID] = channelViewStat{
				warning: 1,
				adID:    chads[j].AdID,
				frwrd:   channelAdStat[chads[j].AdID].frwrd,
				pos:     channelAdStat[chads[j].AdID].pos,
				view:    channelAdStat[chads[j].AdID].view,
			}
		}
		if channelAdStat[chads[j].AdID].frwrd == true {
			currentView = avg
		} else {
			currentView = channelAdStat[chads[j].AdID].view
		}
		ChannelAdDetailArr = append(ChannelAdDetailArr, &ads.ChannelAdDetail{
			AdID:      chads[j].AdID,
			ChannelID: chads[j].ChannelID,
			View:      currentView,
			Position:  common.NullInt64{Valid: channelAdStat[chads[j].AdID].pos != 0, Int64: channelAdStat[chads[j].AdID].pos},
			Warning:   channelAdStat[chads[j].AdID].warning,
		})
	}

	var ChannelAdArr []ads.ChannelAd

	for chad := range chads {
		var currentView int64
		if channelAdStat[chads[chad].AdID].frwrd == true {
			currentView = avg
		} else {
			currentView = channelAdStat[chads[chad].AdID].view
		}
		ChannelAdArr = append(ChannelAdArr, ads.ChannelAd{

			Warning:   chads[chad].Warning + channelAdStat[chads[chad].AdID].warning,
			View:      currentView,
			AdID:      chads[chad].AdID,
			ChannelID: chads[chad].ChannelID,
		})
		if chads[chad].Warning >= tcfg.Cfg.Telegram.LimitCountWarning {
			//send stop (warn message)
			rabbit.MustPublish(&bot2.SendWarn{
				AdID:      chads[chad].AdID,
				ChannelID: in.ChannelID,
				Msg:       trans.T("please reshot the following ad"),
			})
			ChannelAdArr = append(ChannelAdArr, ads.ChannelAd{
				End: common.NullTime{Valid: true, Time: time.Now()},
			})
		}
	}

	//transaction
	res, err := mw.transaction(m, ChannelAdArr, ChannelAdDetailArr, avg)
	if res == true {
		return true, err
	}

	return false, nil
}

func (mw *MultiWorker) discoverAd(in *commands.DiscoverAd) (bool, error) {
	adsManager := ads.NewAdsManager()
	chads, err := adsManager.FindChannelAdByChannelIDActive(in.Channel)
	assert.Nil(err)

	// first try to resolve the channel
	m := bot.NewBotManager()
	channel, err := adsManager.FindChannelByID(in.Channel)
	assert.Nil(err)
	c, err := m.FindKnownChannelByName(channel.Name)
	if err != nil {
		ch, err := mw.discoverChannel(channel.Name)
		assert.Nil(err)
		c, err = bot.NewBotManager().CreateChannelByRawData(ch)
		assert.Nil(err)
	}
	h, err := mw.getLastMessages(c.CliTelegramID, 10, 0)
	assert.Nil(err)
	// then discover the messages
	found := 0
bigLoop:
	for i := range chads {
		if chads[i].CliMessageID.Valid {
			found++
			continue
		}
		var msg *tgo.History
		if chads[i].CliMessageAd.Valid {
			msg, err = mw.discoverMessage(chads[i].CliMessageAd.String)
			assert.Nil(err)
		}
		for j := len(h) - 1; j >= 0; j-- {
			// Promotion or individual?
			if h[j].FwdFrom == nil {
				continue
			}
			if msg != nil {
				// promotion
				if comparePromotion(*msg, h[j]) {
					adsManager.SetCLIMessageID(chads[i].ChannelID, chads[i].AdID, h[j].ID)
					found++
					continue bigLoop
				}
			} else {
				if compareIndividual(chads[i], h[j]) {
					adsManager.SetCLIMessageID(chads[i].ChannelID, chads[i].AdID, h[j].ID)
					found++
					continue bigLoop
				}
			}
		}
	}

	// I think we are gonna die, or the user is stupid. replace us, or him/her
	if found != len(chads) {

		data, _ := json.MarshalIndent(struct {
			History []tgo.History
			Chads   []ads.ChannelAdD
		}{
			History: h,
			Chads:   chads,
		}, "\t", "\t")
		logrus.Warn(data)
		if config.Config.Slack.Active {
			go utils.SlackDoMessage(
				"[BUG/USER] can not find messages in channel, nothing special but please check if the user is stupid or we are?",
				":thinking_face:",
				utils.SlackAttachment{Text: string(data), Color: "#AA3939"},
			)
		}
	} else {
		rabbit.MustPublishAfter(
			commands.ExistChannelAd{
				ChannelID: channel.ID,
				ChatID:    in.ChatID,
			},
			tcfg.Cfg.Telegram.TimeReQueUe,
		)
		//send ok message
		rabbit.MustPublish(&bot2.SendWarn{
			AdID:      0,
			ChannelID: in.Channel,
			Msg:       trans.T("your add has been successfully activated\nthanks for your cooperation"),
			ChatID:    in.ChatID,
		})
	}
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
