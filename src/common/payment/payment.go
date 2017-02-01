package payment

import (
	"bytes"
	"crypto/tls"
	"encoding/xml"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

// Request lint
type Request struct {
	XMLName xml.Name `xml:"http://zarinpal.com/ PaymentRequest"`

	MerchantID string `xml:"MerchantID,omitempty"`

	Amount int64 `xml:"Amount,omitempty"`

	Description string `xml:"Description,omitempty"`

	Email string `xml:"Email,omitempty"`

	Mobile string `xml:"Mobile,omitempty"`

	CallbackURL string `xml:"CallbackURL,omitempty"`
}

// RequestResponse lint
type RequestResponse struct {
	XMLName xml.Name `xml:"http://zarinpal.com/ PaymentRequestResponse"`

	Status int64 `xml:"Status,omitempty"`

	Authority string `xml:"Authority,omitempty"`
}

// paymentRequestWithExtra lint
type paymentRequestWithExtra struct {
	XMLName xml.Name `xml:"http://zarinpal.com/ PaymentRequestWithExtra"`

	MerchantID string `xml:"MerchantID,omitempty"`

	Amount int32 `xml:"Amount,omitempty"`

	Description string `xml:"Description,omitempty"`

	AdditionalData string `xml:"AdditionalData,omitempty"`

	Email string `xml:"Email,omitempty"`

	Mobile string `xml:"Mobile,omitempty"`

	CallbackURL string `xml:"CallbackURL,omitempty"`
}

// RequestWithExtraResponse lint
type RequestWithExtraResponse struct {
	XMLName xml.Name `xml:"http://zarinpal.com/ PaymentRequestWithExtraResponse"`

	Status int32 `xml:"Status,omitempty"`

	Authority string `xml:"Authority,omitempty"`
}

// Verification lint
type Verification struct {
	XMLName xml.Name `xml:"http://zarinpal.com/ PaymentVerification"`

	MerchantID string `xml:"MerchantID,omitempty"`

	Authority string `xml:"Authority,omitempty"`

	Amount int64 `xml:"Amount,omitempty"`
}

// VerificationResponse lint
type VerificationResponse struct {
	XMLName xml.Name `xml:"http://zarinpal.com/ PaymentVerificationResponse"`

	Status int64 `xml:"Status,omitempty"`

	RefID int64 `xml:"RefID,omitempty"`
}

// paymentVerificationWithExtra lint
type paymentVerificationWithExtra struct {
	XMLName xml.Name `xml:"http://zarinpal.com/ PaymentVerificationWithExtra"`

	MerchantID string `xml:"MerchantID,omitempty"`

	Authority string `xml:"Authority,omitempty"`

	Amount int32 `xml:"Amount,omitempty"`
}

// VerificationWithExtraResponse lint
type VerificationWithExtraResponse struct {
	XMLName xml.Name `xml:"http://zarinpal.com/ PaymentVerificationWithExtraResponse"`

	Status int32 `xml:"Status,omitempty"`

	RefID int64 `xml:"RefID,omitempty"`

	ExtraDetail string `xml:"ExtraDetail,omitempty"`
}

// GetUnverifiedTransactions lint
type GetUnverifiedTransactions struct {
	XMLName xml.Name `xml:"http://zarinpal.com/ GetUnverifiedTransactions"`

	MerchantID string `xml:"MerchantID,omitempty"`
}

// GetUnverifiedTransactionsResponse lint
type GetUnverifiedTransactionsResponse struct {
	XMLName xml.Name `xml:"http://zarinpal.com/ GetUnverifiedTransactionsResponse"`

	Status int32 `xml:"Status,omitempty"`

	Authorities string `xml:"Authorities,omitempty"`
}

// RefreshAuthority lint
type RefreshAuthority struct {
	XMLName xml.Name `xml:"http://zarinpal.com/ RefreshAuthority"`

	MerchantID string `xml:"MerchantID,omitempty"`

	Authority string `xml:"Authority,omitempty"`

	ExpireIn int32 `xml:"ExpireIn,omitempty"`
}

// RefreshAuthorityResponse lint
type RefreshAuthorityResponse struct {
	XMLName xml.Name `xml:"http://zarinpal.com/ RefreshAuthorityResponse"`

	Status int32 `xml:"Status,omitempty"`
}

// GatewayImplementationServicePortType lint
type GatewayImplementationServicePortType struct {
	client *SOAPClient
}

// NewPaymentGatewayImplementationServicePortType lint
func NewPaymentGatewayImplementationServicePortType(url string, tls bool, auth *BasicAuth) *GatewayImplementationServicePortType {
	if url == "" {
		url = "https://de.zarinpal.com/pg/services/WebGate/service"
	}
	client := NewSOAPClient(url, tls, auth)

	return &GatewayImplementationServicePortType{
		client: client,
	}
}

// Request lint
func (service *GatewayImplementationServicePortType) Request(request *Request) (*RequestResponse, error) {
	response := new(RequestResponse)
	err := service.client.Call("#PaymentRequest", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// RequestWithExtra lint
func (service *GatewayImplementationServicePortType) RequestWithExtra(request *paymentRequestWithExtra) (*RequestWithExtraResponse, error) {
	response := new(RequestWithExtraResponse)
	err := service.client.Call("#PaymentRequestWithExtra", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Verification lint
func (service *GatewayImplementationServicePortType) Verification(request *Verification) (*VerificationResponse, error) {
	response := new(VerificationResponse)
	err := service.client.Call("#PaymentVerification", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// VerificationWithExtra lint
func (service *GatewayImplementationServicePortType) VerificationWithExtra(request *paymentVerificationWithExtra) (*VerificationWithExtraResponse, error) {
	response := new(VerificationWithExtraResponse)
	err := service.client.Call("#PaymentVerificationWithExtra", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// GetUnverifiedTransactions lint
func (service *GatewayImplementationServicePortType) GetUnverifiedTransactions(request *GetUnverifiedTransactions) (*GetUnverifiedTransactionsResponse, error) {
	response := new(GetUnverifiedTransactionsResponse)
	err := service.client.Call("#GetUnverifiedTransactions", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// RefreshAuthority lint
func (service *GatewayImplementationServicePortType) RefreshAuthority(request *RefreshAuthority) (*RefreshAuthorityResponse, error) {
	response := new(RefreshAuthorityResponse)
	err := service.client.Call("#RefreshAuthority", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

var timeout time.Duration

func dialTimeout(network, addr string) (net.Conn, error) {
	return net.DialTimeout(network, addr, timeout)
}

// SOAPEnvelope lint
type SOAPEnvelope struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`

	Body SOAPBody
}

// SOAPHeader lint
type SOAPHeader struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Header"`

	Header interface{}
}

// SOAPBody lint
type SOAPBody struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Body"`

	Fault   *SOAPFault  `xml:",omitempty"`
	Content interface{} `xml:",omitempty"`
}

// SOAPFault lint
type SOAPFault struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Fault"`

	Code   string `xml:"faultcode,omitempty"`
	String string `xml:"faultstring,omitempty"`
	Actor  string `xml:"faultactor,omitempty"`
	Detail string `xml:"detail,omitempty"`
}

// BasicAuth lint
type BasicAuth struct {
	Login    string
	Password string
}

// SOAPClient lint
type SOAPClient struct {
	url  string
	tls  bool
	auth *BasicAuth
}

// UnmarshalXML lint
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

// NewSOAPClient lint
func NewSOAPClient(url string, tls bool, auth *BasicAuth) *SOAPClient {
	return &SOAPClient{
		url:  url,
		tls:  tls,
		auth: auth,
	}
}

// Call lint
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

	// log.Println(buffer.String())

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
	defer res.Body.Close()

	rawbody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if len(rawbody) == 0 {
		// log.Println("empty response")
		return nil
	}

	// log.Println(string(rawbody))
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

func init() {
	timeout = time.Duration(30 * time.Second)
}
