package mailer

import (
	baseconfig "common/config"
	"common/utils"
	"modules/user/config"
	"runtime/debug"

	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/aymerick/douceur/inliner"
	"github.com/go-gomail/gomail"
)

var (
	mailer *gomail.Dialer
	lock   = &sync.Mutex{}
)

// SendMail try to send an email base on template. mostly called as gorutine, so I log the error if any
func SendMail(to, toName, text, subject string) (err error) {
	lock.Lock()
	defer lock.Unlock()
	defer func() {
		if err != nil {
			logrus.Warn(err)
			stack := debug.Stack()

			if baseconfig.Config.Redmine.Active {
				go utils.RedmineDoError(err, stack)
			}

			if baseconfig.Config.Slack.Active {
				go utils.SlackDoMessage(err, ":shit:", utils.SlackAttachment{Text: string(stack), Color: "#AA3939"})
			}

		}
	}()

	if mailer == nil {
		mailer = gomail.NewPlainDialer(
			ucfg.Cfg.SMTP.Host,
			ucfg.Cfg.SMTP.Port,
			ucfg.Cfg.SMTP.User,
			ucfg.Cfg.SMTP.Password,
		)
	}

	// Its time to send email
	msg := gomail.NewMessage()
	msg.SetAddressHeader("From", ucfg.Cfg.SMTP.Sender, ucfg.Cfg.SMTP.Name)
	msg.SetAddressHeader("To", to, toName)
	msg.SetHeader("Subject", subject)

	inlined, err := inliner.Inline(text)
	if err != nil {
		// This is not critical, but let us know
		inlined = text
		logrus.Warn(err)
	}
	msg.SetBody("text/html", inlined)

	if err := mailer.DialAndSend(msg); err != nil {
		return err
	}

	return nil
}
