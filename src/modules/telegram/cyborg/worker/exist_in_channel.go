package worker

import (
	bot2 "modules/telegram/bot/worker"
	"modules/telegram/common/tgo"

	"common/rabbit"

	"fmt"

	"github.com/Sirupsen/logrus"
)

func (mw *MultiWorker) existChannelAdFor(h []tgo.History, adConfs []channelDetailStat) map[int64]channelViewStat {
	var finalResult = make(map[int64]channelViewStat)
	var sumIndividualView int64
	var countIndividual int64
	var found int
	historyLen := len(h)
bigloop:
	for k := historyLen - 1; k >= 0; k-- {
		if h[k].Event == "message" && !h[k].Service {
			if h[k].FwdFrom != nil {
				for i := range adConfs {
					if h[k].ID == adConfs[i].cliChannelAdID.String {
						finalResult[adConfs[i].adID] = channelViewStat{
							view:    int64(h[k].Views),
							warning: 0,
							pos:     int64(historyLen - k),
							frwrd:   adConfs[i].frwrd,
							adID:    adConfs[i].adID,
						}
						found++
						if !adConfs[i].frwrd { //the ad is  not forward type
							sumIndividualView += int64(h[k].Views)
							countIndividual++
						}
					}
					if found == len(adConfs) {
						break bigloop
					}
				}
			}
		}
	}

	for i := range adConfs {
		if _, ok := finalResult[adConfs[i].adID]; !ok {
			logrus.Infof("%+v", finalResult[adConfs[i].adID])
			finalResult[adConfs[i].adID] = channelViewStat{
				view:    0,
				warning: 1,
				frwrd:   adConfs[i].frwrd,
				adID:    adConfs[i].adID,
				pos:     0,
			}
			//send stop (warn message)
			rabbit.MustPublish(&bot2.SendWarn{
				AdID:      adConfs[i].adID,
				ChannelID: adConfs[i].channelID,
				ChatID:    adConfs[i].botChatID,
				Msg:       fmt.Sprintf("please remove the following ad for bundle %s from your channel @%s", adConfs[i].code, adConfs[i].channelID),
			})
		}
	}

	logrus.Warnf("%+v", finalResult, sumIndividualView, countIndividual)
	return finalResult
}
