package mail

import (
	"common/config"
	"common/initializer"
	"time"

	"bytes"
	"common/assert"

	"html/template"

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

type MailTemplateData struct {
	content interface{}
	date    time.Time
}

func (t *MailTemplateData) newMail(data interface{}) {
	t.date = time.Now()
	t.content = data
}

// Initialize try to connect to mail server
func (mailInitializer) Initialize() {

}

func (mailInitializer) Finalize() {
	logrus.Debug("Mail is done")
}

func templateFiller(templ []byte, data interface{}) []byte {
	tmpl := template.Must(template.New("tmpl").Parse(string(templ)))
	buf := &bytes.Buffer{}
	assert.Nil(tmpl.Execute(buf, data))
	return buf.Bytes()
}

func masterTemplate(content []byte) []byte {
	src, err := Asset("resource/email-master.html")
	assert.Nil(err)
	ctn := template.HTML(content)
	return templateFiller(src, struct {
		Date    time.Time
		Content template.HTML
	}{
		time.Now(),
		ctn,
	})
}

// SendTemplate is a simple email sender with text/html
func SendByTemplateName(subject string, TemplateName string, data interface{}, from string, to ...string) error {
	src, err := Asset(TemplateName)
	assert.Nil(err)
	return SendByTemplate(subject, src, data, from, to...)
}

// SendTemplate is a simple email sender with text/html
func SendByTemplate(subject string, EmailTemplate []byte, data interface{}, from string, to ...string) error {
	content := templateFiller(EmailTemplate, data)
	body := masterTemplate(content)
	return Send(subject, string(body), from, to...)
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
