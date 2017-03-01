package mail

import (
	"common/config"
	"common/initializer"

	"bytes"
	"common/assert"

	"html/template"

	"errors"

	"github.com/Sirupsen/logrus"
	"github.com/go-gomail/gomail"
)

var (
	// Client the object connect to mail server
	Client *gomail.Dialer
	//once   = &sync.Once{}
	mailTemplate = template.New("mail")
)

type mailInitializer struct {
}

// Initialize try to connect to mail server
func (mailInitializer) Initialize() {
	loadTemplates()
}

func (mailInitializer) Finalize() {
	logrus.Debug("Mail is done")
}

// SendByTemplateName is a simple email sender with text/html
func SendByTemplateName(subject string, TemplateName string, data interface{}, from string, to ...string) error {
	buff := &bytes.Buffer{}

	err := mailTemplate.Lookup(TemplateName).Execute(buff, data)
	assert.Nil(err)
	return Send(subject, buff.Bytes(), from, to...)
}

// Send is a simple email sender with text/html
func Send(subject string, body []byte, from string, to ...string) error {
	Client := gomail.NewDialer(config.Config.Mail.Host, config.Config.Mail.Port, config.Config.Mail.UserName, config.Config.Mail.Password)
	m := gomail.NewMessage()
	if from == "" {
		from = config.Config.Mail.From
	}
	m.SetHeader("From", from)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", string(body))
	//m.Attach("/home/Alex/lolcat.jpg")

	logrus.Infof("%+v", Client)
	// Send the email to Bob, Cora and Dan.
	err := Client.DialAndSend(m)
	return err
}

func loadTemplates() {
	if len(_bindata) < 1 {
		assert.Nil(errors.New("There is no template to compile"))
		return
	}

	for t := range _bindata {
		data, err := Asset(t)
		assert.Nil(err)
		mailTemplate.New(t).Parse(string(data))
	}
}

func init() {
	initializer.Register(&mailInitializer{})
}
