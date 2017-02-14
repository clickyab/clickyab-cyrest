package bot

import (
	"common/assert"

	"modules/telegram/ad/ads"

	"regexp"

	"strconv"

	"modules/misc/base"
	"modules/telegram/teleuser/tlu"

	"modules/misc/trans"

	"gopkg.in/telegram-bot-api.v4"
)

var (
	conMsg = regexp.MustCompile("/confirm_([a-z]+)_([0-9]+)")
)

func doMessage(bot *tgbotapi.BotAPI, chatID int64, m string) {
	msg := tgbotapi.NewMessage(chatID, m)
	msg.ParseMode = "HTML"
	_, err := bot.Send(msg)
	assert.Nil(err)
	return
}

func (bb *bot) confirm(bot *tgbotapi.BotAPI, m *tgbotapi.Message) {
	u, err := tlu.NewTluManager().GetUser(m.Chat.ID)
	if err != nil {
		// No action
		doMessage(bot, m.Chat.ID, trans.T("<b>Not authorized</b> #1001").Translate())
		return
	}
	if _, ok := u.HasPerm(base.ScopeGlobal, "confirm_ad"); !ok {
		// No action
		doMessage(bot, m.Chat.ID, trans.T("<b>Not authorized</b> #1002").Translate())
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
			doMessage(bot, m.Chat.ID, trans.T("<b>Invalid ad</b>").Translate())
			return
		}
		ad.AdAdminStatus = ads.AdAdminStatusAccepted
		assert.Nil(mm.UpdateAd(ad))
		resp = trans.T("Ad %s is accepted", ad.Name).Translate()
	} else if command == "reject" {
		ad, err := mm.FindAdByID(param)
		if err != nil || ad.AdActiveStatus != ads.AdActiveStatusYes || ad.AdAdminStatus != ads.AdAdminStatusPending || ad.AdPayStatus != ads.AdPayStatusYes {
			doMessage(bot, m.Chat.ID, trans.T("<b>Invalid ad</b>").Translate())
			return
		}
		ad.AdAdminStatus = ads.AdAdminStatusRejected
		assert.Nil(mm.UpdateAd(ad))
		resp = trans.T("Ad %s is rejected", ad.Name).Translate()
	} else {
		ad, err := mm.LoadNextAd(param)
		if err != nil {
			doMessage(bot, m.Chat.ID, trans.T("<b>No ad available at this time</b>").Translate())
			return
		}

		RenderMessage(bot, m.Chat.ID, ad)
		doMessage(bot, m.Chat.ID, trans.T("Accept /confirm_%s_%d", "accept", ad.ID).Translate())
		doMessage(bot, m.Chat.ID, trans.T("Reject /confirm_%s_%d", "reject", ad.ID).Translate())
		doMessage(bot, m.Chat.ID, trans.T("Next /confirm_%s_%d", "next", ad.ID).Translate())
	}
	doMessage(bot, m.Chat.ID, resp)
	doMessage(bot, m.Chat.ID, trans.T("<i>OK</i>").Translate())
}

func init() {
	base.RegisterPermission("confirm_ad", "confirm_ad")
}
