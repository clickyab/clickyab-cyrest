package bot

import (
	"common/assert"

	"modules/telegram/ad/ads"

	"regexp"

	"strconv"

	"fmt"

	"modules/misc/base"
	"modules/telegram/teleuser/tlu"

	"gopkg.in/telegram-bot-api.v4"
)

var (
	conMsg = regexp.MustCompile("/confirm_([a-z]+)_([0-9]+)")
)

func doMessage(bot *tgbotapi.BotAPI, chatID int64, m string) {
	msg := tgbotapi.NewMessage(chatID, m)
	msg.ParseMode = "HTML"
	_, _ = bot.Send(msg)
	return
}

func (bb *bot) confirm(bot *tgbotapi.BotAPI, m *tgbotapi.Message) {
	u, err := tlu.NewTluManager().GetUser(m.Chat.ID)
	if err != nil {
		// No action
		doMessage(bot, m.Chat.ID, "<b>Not authorized</b> #1001")
		return
	}
	if _, ok := u.HasPerm(base.ScopeGlobal, "confirm_ad"); !ok {
		// No action
		doMessage(bot, m.Chat.ID, "<b>Not authorized</b> #1002")
		return
	}
	var (
		param   int64
		command string
		resp    string
	)
	cns := conMsg.FindStringSubmatch(m.Text)
	if len(cns) == 3 {
		command = cns[1]
		param, _ = strconv.ParseInt(cns[2], 10, 0)
	}
	mm := ads.NewAdsManager()
	if command == "accept" {
		ad, err := mm.FindAdByID(param)
		if err != nil || ad.AdActiveStatus != ads.AdActiveStatusYes || ad.AdAdminStatus != ads.AdAdminStatusPending || ad.AdPayStatus != ads.AdPayStatusYes {
			doMessage(bot, m.Chat.ID, "<b>Invalid ad<b>")
			return
		}
		ad.AdAdminStatus = ads.AdAdminStatusAccepted
		assert.Nil(mm.UpdateAd(ad))
		resp = fmt.Sprintf("Ad %s is accepted", ad.Name)
	} else if command == "reject" {
		ad, err := mm.FindAdByID(param)
		if err != nil || ad.AdActiveStatus != ads.AdActiveStatusYes || ad.AdAdminStatus != ads.AdAdminStatusPending || ad.AdPayStatus != ads.AdPayStatusYes {
			doMessage(bot, m.Chat.ID, "<b>Invalid ad<b>")
			return
		}
		ad.AdAdminStatus = ads.AdAdminStatusRejected
		assert.Nil(mm.UpdateAd(ad))
		resp = fmt.Sprintf("Ad %s is rejected", ad.Name)
	} else {
		ad, err := mm.LoadNextAd(param)
		if err != nil {
			doMessage(bot, m.Chat.ID, "<b>No ad available at this time<b>")
			return
		}

		RenderMessage(bot, m.Chat.ID, ad)
		doMessage(bot, m.Chat.ID, fmt.Sprintf("Accept /confirm_%s_%d", "accept", ad.ID))
		doMessage(bot, m.Chat.ID, fmt.Sprintf("Reject /confirm_%s_%d", "reject", ad.ID))
		doMessage(bot, m.Chat.ID, fmt.Sprintf("Next /confirm_%s_%d", "next", ad.ID))
	}
	doMessage(bot, m.Chat.ID, resp)
	doMessage(bot, m.Chat.ID, "<i>OK</i>")
}
