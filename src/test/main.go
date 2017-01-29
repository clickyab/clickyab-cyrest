package main

import (
	"common/config"
	"common/initializer"

	"github.com/Sirupsen/logrus"
	"github.com/go-gomail/gomail"
)

func main() {
	config.Initialize()
	config.InitApplication()

	defer initializer.Initialize().Finalize()
	/*from := mail.NewEmail("Example User", "mahmud.tehrani+from@gmail.com")
	subject := "Sending with SendGrid is Fun"
	to := mail.NewEmail("Example User", "mahmud.tehrani+to@gmail.com")
	content := mail.NewContent("text/plain", "and easy to do anywhere, even with Go")
	m := mail.NewV3MailInit(from, subject, to, content)
	logrus.Infof("%+v", m)

	request := sendgrid.GetRequest("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", "/v3/mail/send", "https://api.sendgrid.com")
	logrus.Infof("%+v", request)
	request.Method = "POST"
	request.Body = mail.GetRequestBody(m)
	response, err := sendgrid.API(request)
	logrus.Infof("%+v", response)
	logrus.Warn(err)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}*/
	/*apiKey := "pJYAZl_1QiaS81UwGOnjrQ"
		request := sendgrid.GetRequest(apiKey, "/v3/mail/send", "https://api.sendgrid.com")
		request.Method = "POST"
		request.Body = []byte(` {
	    "personalizations": [
	        {
	            "to": [
	                {
	                    "email": "mahmud.tehrani+to@gmail.com"
	                }
	            ],
	            "subject": "I'm replacing the subject tag",
	                        "substitutions": {
	                            "-name-": "Example User",
	                            "-city-": "Denver"
	                        },
	        }
	    ],
	    "from": {
	        "email": "mahmud.tehrani+from@gmail.com"
	    },
	    "content": [
	        {
	            "type": "text/html",
	            "value": "I'm replacing the <strong>body tag</strong>"
	        }
	    ],
	    "template_id": "13b8f94f-bcae-4ec6-b752-70d6cb59f932"
	}`)
		response, err := sendgrid.API(request)
		if err != nil {
			log.Println(err)
		} else {
			fmt.Println(response.StatusCode)
			fmt.Println(response.Body)
			fmt.Println(response.Headers)
		}*/
	/*request := sendgrid.GetRequest(os.Getenv("SENDGRID_API_KEY"), "/v3/api_keys", "https://api.sendgrid.com")
	request.Method = "GET"

	response, err := sendgrid.API(request)
	logrus.Infof("%+v", request.BaseURL, request.Body, request.Headers, request.Method, request.QueryParams)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}
	*/
	/*sg := sendgrid.NewSendGridClient("sendgrid_user", "sendgrid_key")
	message := sendgrid.NewMail()
	message.AddTo("yamil@sendgrid.com")
	message.AddSubject("My first email!")
	message.AddText("Sending Email from Go using SendGrid")
	message.AddFrom("yamil.asusta@sendgrid.com")
	if r := sg.Send(message); r == nil {
		fmt.Println("Email sent!")
	} else {
		fmt.Println(r)
	}*/
	/*from := mail.NewEmail("Example User", "test@example.com")
	subject := "Sending with SendGrid is Fun"
	to := mail.NewEmail("Example User", "test@example.com")
	content := mail.NewContent("text/plain", "and easy to do anywhere, even with Go")
	m := mail.NewV3MailInit(from, subject, to, content)

	request := sendgrid.GetRequest(os.Getenv("SENDGRID_API_KEY"), "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	request.Body = mail.GetRequestBody(m)
	response, err := sendgrid.API(request)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}*/
	// Build the URL
	/*const host = "https://api.sendgrid.com"
	endpoint := "/v3/api_keys"
	baseURL := host + endpoint

	// Build the request headers
	key := os.Getenv("SENDGRID_API_KEY")
	Headers := make(map[string]string)
	Headers["Authorization"] = "Bearer " + key

	// GET Collection
	method := rest.Get

	// Build the query parameters
	queryParams := make(map[string]string)
	queryParams["limit"] = "100"
	queryParams["offset"] = "0"

	// Make the API call
	request := rest.Request{
		Method:      method,
		BaseURL:     baseURL,
		Headers:     Headers,
		QueryParams: queryParams,
	}
	response, err := rest.API(request)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}*/
	m := gomail.NewMessage()
	m.SetHeader("From", config.Config.Mail.From)
	m.SetHeader("To", "mahmud.tehrani@gmail.com", "cora@example.com")
	m.SetAddressHeader("Cc", "dan@example.com", "Dan")
	m.SetHeader("Subject", "Hello!")
	m.SetBody("text/html", "Hello <b>Bob</b> and <i>Cora</i>!")
	//m.Attach("/home/Alex/lolcat.jpg")

	d := gomail.NewDialer(config.Config.Mail.Host, config.Config.Mail.Port, config.Config.Mail.UserName, config.Config.Mail.Password)
	logrus.Infof("%+v", d)
	// Send the email to Bob, Cora and Dan.
	err := d.DialAndSend(m)
	logrus.Warn(err)
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}
