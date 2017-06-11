package bot

import (
	"common/assert"

	"modules/telegram/ad/ads"

	"regexp"

	"strconv"

	"fmt"

	"modules/misc/base"
	"modules/telegram/teleuser/tlu"

	"modules/misc/trans"

	"common/config"
	"common/mail"
	"modules/user/aaa"

	"gopkg.in/telegram-bot-api.v4"
)

var (
	conMsg = regexp.MustCompile("/confirm_([a-z]+)_([0-9]+)")
)

func doMessageString(bot *tgbotapi.BotAPI, chatID int64, m string) {
	msg := tgbotapi.NewMessage(chatID, m)
	msg.ParseMode = "HTML"
	_, err := bot.Send(msg)
	assert.Nil(err)
	return

}

func doMessage(bot *tgbotapi.BotAPI, chatID int64, m trans.T9String) {
	doMessageString(bot, chatID, m.Translate(trans.PersianLang))
}

func (bb *bot) confirm(bot *tgbotapi.BotAPI, m *tgbotapi.Message) {
	u, err := tlu.NewTluManager().GetUser(m.Chat.ID)
	if err != nil {
		// No action
		doMessage(bot, m.Chat.ID, trans.T("Not authorized"))
		return
	}
	if _, ok := u.HasPerm(base.ScopeGlobal, "confirm_ad"); !ok {
		// No action
		doMessage(bot, m.Chat.ID, trans.T("Not authorized"))
		return
	}
	var (
		param   int64
		command string
		resp    trans.T9String
	)
	cns := conMsg.FindStringSubmatch(m.Text)
	if len(cns) == 3 {
		command = cns[1]
		param, _ = strconv.ParseInt(cns[2], 10, 0)
	}
	mm := ads.NewAdsManager()
	if command == "accept" {
		ad, err := mm.FindAdByID(param)
		if err != nil || ad.AdActiveStatus != ads.ActiveStatusYes || ad.AdAdminStatus != ads.AdminStatusPending || ad.AdPayStatus != ads.ActiveStatusYes {
			doMessage(bot, m.Chat.ID, trans.T("Invalid ad"))
			return
		}
		ad.AdAdminStatus = ads.AdminStatusAccepted
		assert.Nil(mm.UpdateAd(ad))
		resp = trans.T("Ad %s is accepted", ad.Name)
	} else if command == "reject" {
		ad, err := mm.FindAdByID(param)
		if err != nil || ad.AdActiveStatus != ads.ActiveStatusYes || ad.AdAdminStatus != ads.AdminStatusPending || ad.AdPayStatus != ads.ActiveStatusYes {
			doMessage(bot, m.Chat.ID, trans.T("Invalid ad"))
			return
		}
		owner, err := aaa.NewAaaManager().FindUserByID(ad.UserID)
		if err != nil {
			return
		}
		ad.AdAdminStatus = ads.AdminStatusRejected
		assert.Nil(mm.UpdateAd(ad))
		//send email
		go func() {
			mail.SendByTemplateName(trans.T("ad rejected").Translate("fa_IR"), "reject-ad", struct {
				Name     string
				Campaign string
			}{
				Name:     owner.Email,
				Campaign: ad.Name,
			}, config.Config.Mail.From, owner.Email)
		}()
		resp = trans.T("Ad %s is rejected", ad.Name)
	} else {
		ad, err := mm.LoadNextAd(param)
		if err != nil {
			doMessage(bot, m.Chat.ID, trans.T("No ad available at this time"))
			return
		}

		RenderMessage(bot, m.Chat.ID, ad)
		doMessageString(bot, m.Chat.ID, fmt.Sprintf("%s /confirm_%s_%d", trans.T("Accept").String(), "accept", ad.ID))
		doMessageString(bot, m.Chat.ID, fmt.Sprintf("%s /confirm_%s_%d", trans.T("Reject").String(), "reject", ad.ID))
		doMessageString(bot, m.Chat.ID, fmt.Sprintf("%s /confirm_%s_%d", trans.T("Next").String(), "next", ad.ID))
		return
	}
	doMessage(bot, m.Chat.ID, resp)
	doMessage(bot, m.Chat.ID, trans.T("OK"))
}

func init() {
	base.RegisterPermission("confirm_ad", "confirm_ad")
}
