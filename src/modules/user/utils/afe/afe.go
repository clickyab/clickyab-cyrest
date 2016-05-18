package afe

import (
	"bytes"
	"crypto/tls"
	"encoding/xml"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"
)

// against "unused imports"
var _ time.Time
var _ xml.Name

type SendMessage struct {
	XMLName xml.Name `xml:"http://www.afe.ir/ SendMessage"`

	Username string         `xml:"Username,omitempty"`
	Password string         `xml:"Password,omitempty"`
	Number   string         `xml:"Number,omitempty"`
	Mobile   *ArrayOfString `xml:"Mobile,omitempty"`
	Message  string         `xml:"Message,omitempty"`
	Type     string         `xml:"Type,omitempty"`
}

type SendMessageResponse struct {
	XMLName xml.Name `xml:"http://www.afe.ir/ SendMessageResponse"`

	SendMessageResult *ArrayOfString `xml:"SendMessageResult,omitempty"`
}

type SendBusinessCard struct {
	XMLName xml.Name `xml:"http://www.afe.ir/ SendBusinessCard"`

	Username    string         `xml:"Username,omitempty"`
	Password    string         `xml:"Password,omitempty"`
	Number      string         `xml:"Number,omitempty"`
	Mobile      *ArrayOfString `xml:"Mobile,omitempty"`
	ContactName string         `xml:"ContactName,omitempty"`
	PhoneNumber string         `xml:"PhoneNumber,omitempty"`
	Type        string         `xml:"Type,omitempty"`
}

type SendBusinessCardResponse struct {
	XMLName xml.Name `xml:"http://www.afe.ir/ SendBusinessCardResponse"`

	SendBusinessCardResult *ArrayOfString `xml:"SendBusinessCardResult,omitempty"`
}

type SendWappush struct {
	XMLName xml.Name `xml:"http://www.afe.ir/ SendWappush"`

	Username    string         `xml:"Username,omitempty"`
	Password    string         `xml:"Password,omitempty"`
	Number      string         `xml:"Number,omitempty"`
	Mobile      *ArrayOfString `xml:"Mobile,omitempty"`
	Url         string         `xml:"Url,omitempty"`
	Description string         `xml:"Description,omitempty"`
	Type        string         `xml:"Type,omitempty"`
}

type SendWappushResponse struct {
	XMLName xml.Name `xml:"http://www.afe.ir/ SendWappushResponse"`

	SendWappushResult *ArrayOfString `xml:"SendWappushResult,omitempty"`
}

type ArrayOfString struct {
	XMLName xml.Name `xml:"http://www.afe.ir/ ArrayOfString"`

	String []string `xml:"string,omitempty"`
}

type BoxServiceSoap struct {
	client *SOAPClient
}

func NewBoxServiceSoap(url string, tls bool, auth *BasicAuth) *BoxServiceSoap {
	if url == "" {
		url = "http://www.afe.ir/WebService/V4/BoxService.asmx"
	}
	client := NewSOAPClient(url, tls, auth)

	return &BoxServiceSoap{
		client: client,
	}
}

func (service *BoxServiceSoap) SendMessage(request *SendMessage) (*SendMessageResponse, error) {
	response := new(SendMessageResponse)
	err := service.client.Call("http://www.afe.ir/SendMessage", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *BoxServiceSoap) SendBusinessCard(request *SendBusinessCard) (*SendBusinessCardResponse, error) {
	response := new(SendBusinessCardResponse)
	err := service.client.Call("http://www.afe.ir/SendBusinessCard", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *BoxServiceSoap) SendWappush(request *SendWappush) (*SendWappushResponse, error) {
	response := new(SendWappushResponse)
	err := service.client.Call("http://www.afe.ir/SendWappush", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

var timeout = time.Duration(30 * time.Second)

func dialTimeout(network, addr string) (net.Conn, error) {
	return net.DialTimeout(network, addr, timeout)
}

type SOAPEnvelope struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`

	Body SOAPBody
}

type SOAPHeader struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Header"`

	Header interface{}
}

type SOAPBody struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Body"`

	Fault   *SOAPFault  `xml:",omitempty"`
	Content interface{} `xml:",omitempty"`
}

type SOAPFault struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Fault"`

	Code   string `xml:"faultcode,omitempty"`
	String string `xml:"faultstring,omitempty"`
	Actor  string `xml:"faultactor,omitempty"`
	Detail string `xml:"detail,omitempty"`
}

type BasicAuth struct {
	Login    string
	Password string
}

type SOAPClient struct {
	url  string
	tls  bool
	auth *BasicAuth
}

func (b *SOAPBody) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	if b.Content == nil {
		return xml.UnmarshalError("Content must be a pointer to a struct")
	}

	var (
		token    xml.Token
		err      error
		consumed bool
	)

Loop:
	for {
		if token, err = d.Token(); err != nil {
			return err
		}

		if token == nil {
			break
		}

		switch se := token.(type) {
		case xml.StartElement:
			if consumed {
				return xml.UnmarshalError("Found multiple elements inside SOAP body; not wrapped-document/literal WS-I compliant")
			} else if se.Name.Space == "http://schemas.xmlsoap.org/soap/envelope/" && se.Name.Local == "Fault" {
				b.Fault = &SOAPFault{}
				b.Content = nil

				err = d.DecodeElement(b.Fault, &se)
				if err != nil {
					return err
				}

				consumed = true
			} else {
				if err = d.DecodeElement(b.Content, &se); err != nil {
					return err
				}

				consumed = true
			}
		case xml.EndElement:
			break Loop
		}
	}

	return nil
}

func (f *SOAPFault) Error() string {
	return f.String
}

func NewSOAPClient(url string, tls bool, auth *BasicAuth) *SOAPClient {
	return &SOAPClient{
		url:  url,
		tls:  tls,
		auth: auth,
	}
}

func (s *SOAPClient) Call(soapAction string, request, response interface{}) error {
	envelope := SOAPEnvelope{
	//Header:        SoapHeader{},
	}

	envelope.Body.Content = request
	buffer := new(bytes.Buffer)

	encoder := xml.NewEncoder(buffer)
	//encoder.Indent("  ", "    ")

	if err := encoder.Encode(envelope); err != nil {
		return err
	}

	if err := encoder.Flush(); err != nil {
		return err
	}

	log.Println(buffer.String())

	req, err := http.NewRequest("POST", s.url, buffer)
	if err != nil {
		return err
	}
	if s.auth != nil {
		req.SetBasicAuth(s.auth.Login, s.auth.Password)
	}

	req.Header.Add("Content-Type", "text/xml; charset=\"utf-8\"")
	if soapAction != "" {
		req.Header.Add("SOAPAction", soapAction)
	}

	req.Header.Set("User-Agent", "gowsdl/0.1")
	req.Close = true

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: s.tls,
		},
		Dial: dialTimeout,
	}

	client := &http.Client{Transport: tr}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = res.Body.Close() }()

	rawbody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if len(rawbody) == 0 {
		log.Println("empty response")
		return nil
	}

	log.Println(string(rawbody))
	respEnvelope := new(SOAPEnvelope)
	respEnvelope.Body = SOAPBody{Content: response}
	err = xml.Unmarshal(rawbody, respEnvelope)
	if err != nil {
		return err
	}

	fault := respEnvelope.Body.Fault
	if fault != nil {
		return fault
	}

	return nil
}
