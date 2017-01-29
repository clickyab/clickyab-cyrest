package mail

import (
	"common/config"
	"common/initializer"

	"github.com/Sirupsen/logrus"
	"github.com/go-gomail/gomail"
)

var (
	// Client the object connect to mail server
	Client *gomail.Dialer
	//once   = &sync.Once{}
)

type mailInitializer struct {
}

// Initialize try to connect to mail server
func (mailInitializer) Initialize() {

}

func (mailInitializer) Finalize() {
	logrus.Debug("Mail is done")
}

// Send is a simple email sender with text/html
func Send(subject string, body string, from string, to ...string) error {
	Client := gomail.NewDialer(config.Config.Mail.Host, config.Config.Mail.Port, config.Config.Mail.UserName, config.Config.Mail.Password)
	m := gomail.NewMessage()
	if from == "" {
		from = config.Config.Mail.From
	}
	m.SetHeader("From", from)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	//m.Attach("/home/Alex/lolcat.jpg")

	logrus.Infof("%+v", Client)
	// Send the email to Bob, Cora and Dan.
	err := Client.DialAndSend(m)
	return err
}

func init() {
	initializer.Register(&mailInitializer{})
}
