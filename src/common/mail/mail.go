package mail

import (
	"common/config"
	"common/initializer"

	"bytes"
	"common/assert"
	"text/template"

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

// SendTemplate is a simple email sender with text/html
func SendByTemplateName(subject string, TemplateName string, data interface{}, from string, to ...string) error {
	src, err := Asset(TemplateName)
	assert.Nil(err)
	return SendByTemplate(subject, src, data, from, to...)
}

// SendTemplate is a simple email sender with text/html
func SendByTemplate(subject string, EmailTemplate []byte, data interface{}, from string, to ...string) error {
	tmpl := template.Must(template.New("email").Parse(string(EmailTemplate)))
	buf := &bytes.Buffer{}
	assert.Nil(tmpl.Execute(buf, data))
	return Send(subject, buf.String(), from, to...)
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

func CreateBody() {

}

func init() {
	initializer.Register(&mailInitializer{})
}
