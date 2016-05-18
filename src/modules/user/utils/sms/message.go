package sms

import (
	"modules/user/config"
	"modules/user/utils/afe"
	"modules/user/utils/mailer"
	"runtime/debug"

	"common/config"
	"common/utils"

	"github.com/Sirupsen/logrus"
)

// SendSMS use the web service to send a sms
func SendSMS(phone, text string) (err error) {
	defer func() {
		if err != nil {
			logrus.Warn(err)
			stack := debug.Stack()

			if config.Config.Redmine.Active {
				go utils.RedmineDoError(err, stack)
			}

			if config.Config.Slack.Active {
				go utils.SlackDoMessage(err, ":shit:", utils.SlackAttachment{Text: string(stack), Color: "#AA3939"})
			}

		} else {
		}
	}()

	if config.Config.DevelMode {
		if !ucfg.Cfg.SMS.Send {
			logrus.Warn(mailer.SendMail("temp@azmoona.com", phone, text, "SMS Simulation for number : "+phone))
			return
		}
	}
	bx := afe.NewBoxServiceSoap("", false, nil)

	req := afe.SendMessage{}
	req.Message = text
	req.Number = ucfg.Cfg.SMS.Sender
	req.Username = ucfg.Cfg.SMS.User
	req.Password = ucfg.Cfg.SMS.Password
	req.Mobile.String = []string{"98" + phone}
	req.Type = "1"

	_, err = bx.SendMessage(&req)
	return err
}
