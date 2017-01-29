package ad

import (
	"bytes"
	"common/assert"
	"common/models/common"
	"crypto/tls"
	"modules/file/fila"
	"modules/misc/middlewares"
	"modules/misc/trans"
	"modules/telegram/ad/ads"
	"modules/user/aaa"
	"modules/user/middlewares"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"common/rabbit"
	"modules/telegram/cyborg/commands"

	"fmt"

	"encoding/xml"
	"io/ioutil"

	"modules/billing/config"

	echo "gopkg.in/labstack/echo.v3"
)

// @Validate {
// }
type adPayload struct {
	Name string `json:"name" validate:"required" error:"name is required"`
}

// @Validate {
// }
type adPromotePayload struct {
	CliMessageID string `json:"cli_message_id" validate:"required" error:"cli_message_id is required"`
}

// @Validate {
// }
type adUploadPayload struct {
	URL string `json:"url" validate:"required" error:"url is required"`
}

// @Validate {
// }
type adAdminStatusPayload struct {
	AdAdminStatus ads.AdAdminStatus `json:"admin_status" validate:"required" error:"status is required"`
}

// @Validate {
// }
type adDescriptionPayLoad struct {
	Body string `json:"body" validate:"required" error:"body is required"`
}

// Validate custom validation for user scope
func (lp *adAdminStatusPayload) ValidateExtra(ctx echo.Context) error {
	if !lp.AdAdminStatus.IsValid() {
		return middlewares.GroupError{
			"status": trans.E("status is invalid"),
		}
	}
	return nil
}

//	create create ad
//	@Route	{
//		url	=	/
//		method	= post
//		payload	= adPayload
//		resource = create_ad:self
//		middleware = authz.Authenticate
//		200 = ads.Ad
//		400 = base.ErrorResponseSimple
//	}
func (u *Controller) create(ctx echo.Context) error {
	pl := u.MustGetPayload(ctx).(*adPayload)
	m := ads.NewAdsManager()
	currentUser := authz.MustGetUser(ctx)
	newAd := &ads.Ad{
		Name:            pl.Name,
		AdArchiveStatus: ads.AdArchiveStatusNo,
		AdPayStatus:     ads.AdPayStatusNo,
		AdAdminStatus:   ads.AdAdminStatusPending,
		AdActiveStatus:  ads.AdActiveStatusNo,
		UserID:          currentUser.ID,
	}
	assert.Nil(m.CreateAd(newAd))
	return u.OKResponse(ctx, newAd)
}

//	changeAdminStatus change admin status for ad
//	@Route	{
//		url	=	/change-admin/:id
//		method	= put
//		payload	= adAdminStatusPayload
//		resource = change_admin_ad:parent
//		middleware = authz.Authenticate
//		200 = ads.Ad
//		400 = base.ErrorResponseSimple
//	}
func (u *Controller) changeAdminStatus(ctx echo.Context) error {
	pl := u.MustGetPayload(ctx).(*adAdminStatusPayload)
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 0)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	m := ads.NewAdsManager()
	currentAd, err := m.FindAdByID(id)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	currentUser := authz.MustGetUser(ctx)
	owner, err := aaa.NewAaaManager().FindUserByID(currentAd.UserID)
	assert.Nil(err)

	_, b := currentUser.HasPermOn("change_admin_ad", owner.ID, owner.DBParentID.Int64)
	if !b {
		return ctx.JSON(http.StatusForbidden, trans.E("user can't access"))
	}
	currentAd.AdAdminStatus = pl.AdAdminStatus
	assert.Nil(m.UpdateAd(currentAd))

	return u.OKResponse(ctx, currentAd)
}

//	changeArchiveStatus change archive status for ad
//	@Route	{
//		url = /change-archive/:id
//		method = put
//		resource = change_archive_ad:self
//		middleware = authz.Authenticate
//		200 = ads.Ad
//		400 = base.ErrorResponseSimple
//	}
func (u *Controller) changeArchiveStatus(ctx echo.Context) error {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 0)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	m := ads.NewAdsManager()
	currentAd, err := m.FindAdByID(id)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	currentUser := authz.MustGetUser(ctx)
	owner, err := aaa.NewAaaManager().FindUserByID(currentAd.UserID)
	assert.Nil(err)
	_, b := currentUser.HasPermOn("change_admin_ad", owner.ID, owner.DBParentID.Int64)
	if !b {
		return ctx.JSON(http.StatusForbidden, trans.E("user can't access"))
	}
	if currentAd.AdArchiveStatus == ads.AdArchiveStatusYes {
		currentAd.AdArchiveStatus = ads.AdArchiveStatusNo
	} else {
		currentAd.AdArchiveStatus = ads.AdArchiveStatusYes
	}
	assert.Nil(m.UpdateAd(currentAd))
	return u.OKResponse(ctx, currentAd)
}

//	addDescription add description to ad
//	@Route	{
//		url	=	/desc/:id
//		method	= put
//		payload	= adDescriptionPayLoad
//		resource = add_description_ad:self
//		middleware = authz.Authenticate
//		200 = ads.Ad
//		400 = base.ErrorResponseSimple
//	}
func (u *Controller) addDescription(ctx echo.Context) error {
	pl := u.MustGetPayload(ctx).(*adDescriptionPayLoad)
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 0)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	m := ads.NewAdsManager()
	currentAd, err := m.FindAdByID(id)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	currentUser := authz.MustGetUser(ctx)
	owner, err := aaa.NewAaaManager().FindUserByID(currentAd.UserID)
	assert.Nil(err)
	_, b := currentUser.HasPermOn("add_description_ad", owner.ID, owner.DBParentID.Int64)
	if !b {
		return ctx.JSON(http.StatusForbidden, trans.E("user can't access"))
	}
	currentAd.Description = common.MakeNullString(pl.Body)
	currentAd.CliMessageID = common.MakeNullString("")
	assert.Nil(m.UpdateAd(currentAd))
	return u.OKResponse(ctx, currentAd)
}

//	uploadBanner uploadBanner for ad
//	@Route	{
//		url	=	/upload/:id
//		method	= put
//		payload	= adUploadPayload
//		resource = upload_ad:self
//		middleware = authz.Authenticate
//		200 = ads.Ad
//		400 = base.ErrorResponseSimple
//	}
func (u *Controller) uploadBanner(ctx echo.Context) error {
	pl := u.MustGetPayload(ctx).(*adUploadPayload)
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 0)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	dURL := pl.URL
	_, err = url.Parse(dURL)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	m := ads.NewAdsManager()
	currentAd, err := m.FindAdByID(id)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	currentUser := authz.MustGetUser(ctx)
	owner, err := aaa.NewAaaManager().FindUserByID(currentAd.UserID)
	assert.Nil(err)
	_, b := currentUser.HasPermOn("upload_ad", owner.ID, owner.DBParentID.Int64)
	if !b {
		return ctx.JSON(http.StatusForbidden, trans.E("user can't access"))
	}

	//upload
	file, err := fila.CheckUpload(dURL, currentUser.ID)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	currentAd.Src = common.MakeNullString(file)
	currentAd.CliMessageID = common.MakeNullString("")
	assert.Nil(m.UpdateAd(currentAd))
	return u.OKResponse(ctx, currentAd)
}

//	promoteAd promoteAd for ad
//	@Route	{
//		url	=	/promote/:id
//		method	= put
//		payload	= adPromotePayload
//		resource = promote_ad:self
//		middleware = authz.Authenticate
//		200 = ads.Ad
//		400 = base.ErrorResponseSimple
//	}
func (u *Controller) promoteAd(ctx echo.Context) error {
	pl := u.MustGetPayload(ctx).(*adPromotePayload)
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 0)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	m := ads.NewAdsManager()
	currentAd, err := m.FindAdByID(id)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	currentUser := authz.MustGetUser(ctx)
	owner, err := aaa.NewAaaManager().FindUserByID(currentAd.UserID)
	assert.Nil(err)
	_, b := currentUser.HasPermOn("promote_ad", owner.ID, owner.DBParentID.Int64)
	if !b {
		return ctx.JSON(http.StatusForbidden, trans.E("user can't access"))
	}
	currentAd.CliMessageID = common.MakeNullString(pl.CliMessageID)
	currentAd.Src = common.MakeNullString("")
	currentAd.BotChatID = common.NullInt64{}
	currentAd.BotMessageID = common.NullInt64{}
	currentAd.Description = common.MakeNullString("")
	assert.Nil(m.UpdateAd(currentAd))
	return u.OKResponse(ctx, currentAd)
}

//	changeActiveStatus change active status for ad
//	@Route	{
//		url = /change-active/:id
//		method = put
//		resource = change_active_ad:self
//		middleware = authz.Authenticate
//		200 = ads.Ad
//		400 = base.ErrorResponseSimple
//	}
func (u *Controller) changeActiveStatus(ctx echo.Context) error {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 0)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	m := ads.NewAdsManager()
	currentAd, err := m.FindAdByID(id)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	currentUser := authz.MustGetUser(ctx)
	owner, err := aaa.NewAaaManager().FindUserByID(currentAd.UserID)
	assert.Nil(err)
	_, b := currentUser.HasPermOn("change_active_ad", owner.ID, owner.DBParentID.Int64)
	if !b {
		return ctx.JSON(http.StatusForbidden, trans.E("user can't access"))
	}
	//check everything is good TODO: check pay status later
	if currentAd.AdAdminStatus == "accepted" && currentAd.AdArchiveStatus == "no" && currentAd.Name != "" && currentAd.AdActiveStatus == ads.AdActiveStatusNo {
		currentAd.AdActiveStatus = ads.AdActiveStatusYes
		//check to add job
		if currentAd.CliMessageID.Valid {
			rabbit.MustPublish(commands.IdentifyAD{AdID: currentAd.ID})
		}
	} else {
		currentAd.AdActiveStatus = ads.AdActiveStatusNo
	}
	assert.Nil(m.UpdateAd(currentAd))
	return u.OKResponse(ctx, currentAd)
}

//	edit edit ad
//	@Route	{
//		url	=	/:id
//		method	= put
//		payload	= adPayload
//		resource = edit_ad:self
//		middleware = authz.Authenticate
//		200 = ads.Ad
//		400 = base.ErrorResponseSimple
//	}
func (u *Controller) edit(ctx echo.Context) error {
	pl := u.MustGetPayload(ctx).(*adPayload)
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 0)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	m := ads.NewAdsManager()
	currentAd, err := m.FindAdByID(id)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	currentUser := authz.MustGetUser(ctx)
	owner, err := aaa.NewAaaManager().FindUserByID(currentAd.UserID)
	assert.Nil(err)
	_, b := currentUser.HasPermOn("edit_ad", owner.ID, owner.DBParentID.Int64)
	if !b {
		return ctx.JSON(http.StatusForbidden, trans.E("user can't access"))
	}
	//check if it can be edited
	if currentAd.AdActiveStatus == ads.AdActiveStatusYes || currentAd.AdAdminStatus == ads.AdAdminStatusAccepted {
		return ctx.JSON(http.StatusForbidden, trans.E("user can't access"))
	}
	currentAd.Name = pl.Name
	assert.Nil(m.UpdateAd(currentAd))
	return u.OKResponse(ctx, currentAd)
}

//	getAd getAd for ad
//	@Route	{
//		url	=	/get/:id
//		method	= get
//		resource = get_ad:self
//		middleware = authz.Authenticate
//		200 = ads.Ad
//		400 = base.ErrorResponseSimple
//	}
func (u *Controller) getAd(ctx echo.Context) error {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 0)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	m := ads.NewAdsManager()
	currentAd, err := m.FindAdByID(id)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	currentUser := authz.MustGetUser(ctx)
	owner, err := aaa.NewAaaManager().FindUserByID(currentAd.UserID)
	assert.Nil(err)
	_, b := currentUser.HasPermOn("promote_ad", owner.ID, owner.DBParentID.Int64)
	if !b {
		return ctx.JSON(http.StatusForbidden, trans.E("user can't access"))
	}
	return u.OKResponse(ctx, currentAd)
}

//	charge charge user plan
//	@Route	{
//		url	=	/pay/:id
//		method	= get
//		resource = pay_ad:self
//		middleware = authz.Authenticate
//		200 = ads.Ad
//		400 = base.ErrorResponseSimple
//	}
func (u *Controller) charge(ctx echo.Context) error {
	//find ad
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 0)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	adManager := ads.NewAdsManager()
	currentAd, err := adManager.FindAdByID(id)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	//find plan
	plan, err := adManager.FindPlanByID(currentAd.PlanID.Int64)
	if err != nil {
		return u.NotFoundResponse(ctx, nil)
	}
	price := plan.Price
	client := NewPaymentGatewayImplementationServicePortType("", false, nil)
	// Create a new payment request to Zarinpal
	resp, err := client.PaymentRequest(&PaymentRequest{
		MerchantID:  bcfg.Bcfg.Gate.MerchantID,
		Amount:      price,
		Description: bcfg.Bcfg.Gate.Description,
		Email:       bcfg.Bcfg.Gate.Email,
		Mobile:      bcfg.Bcfg.Gate.Mobile,
		CallbackURL: bcfg.Bcfg.Gate.CallbackURL,
	})
	if err != nil {
		return u.BadResponse(ctx, nil)
	}
	if resp.Status == bcfg.Bcfg.Gate.MerchantOKStatus {
		//@TODO insert into payment table

		return u.OKResponse(ctx, fmt.Sprintf("%s%s", bcfg.Bcfg.Gate.ZarinURL, resp.Authority))
	}
	return u.BadResponse(ctx, nil)
}

//	charge charge user plan
//	@Route	{
//		url	=	/verify
//		method	= get
//		200 = ads.Ad
//		400 = base.ErrorResponseSimple
//	}
func (u *Controller) verify(ctx echo.Context) error {
	if ctx.FormValue("Status") != "ok" {
		return ctx.Redirect(http.StatusMovedPermanently, "")
	}
	client := NewPaymentGatewayImplementationServicePortType("", false, nil)
	// Create a new payment request to Zarinpal
	resp, err := client.PaymentVerification(&PaymentVerification{
		MerchantID: bcfg.Bcfg.Gate.MerchantID,
		Amount:     100,
		Authority:  ctx.FormValue("Authority"),
	})
	// Check if response is error free
	if err != nil {
		return u.BadResponse(ctx, nil)
	}
	if resp.Status == bcfg.Bcfg.Gate.MerchantOKStatus {
		//@TODO
		return u.OKResponse(ctx, nil)
	} else {
		//@TODO
		return u.BadResponse(ctx, nil)
	}

}

type PaymentRequest struct {
	XMLName xml.Name `xml:"http://zarinpal.com/ PaymentRequest"`

	MerchantID string `xml:"MerchantID,omitempty"`

	Amount int64 `xml:"Amount,omitempty"`

	Description string `xml:"Description,omitempty"`

	Email string `xml:"Email,omitempty"`

	Mobile string `xml:"Mobile,omitempty"`

	CallbackURL string `xml:"CallbackURL,omitempty"`
}

type PaymentRequestResponse struct {
	XMLName xml.Name `xml:"http://zarinpal.com/ PaymentRequestResponse"`

	Status int32 `xml:"Status,omitempty"`

	Authority string `xml:"Authority,omitempty"`
}

type PaymentRequestWithExtra struct {
	XMLName xml.Name `xml:"http://zarinpal.com/ PaymentRequestWithExtra"`

	MerchantID string `xml:"MerchantID,omitempty"`

	Amount int32 `xml:"Amount,omitempty"`

	Description string `xml:"Description,omitempty"`

	AdditionalData string `xml:"AdditionalData,omitempty"`

	Email string `xml:"Email,omitempty"`

	Mobile string `xml:"Mobile,omitempty"`

	CallbackURL string `xml:"CallbackURL,omitempty"`
}

type PaymentRequestWithExtraResponse struct {
	XMLName xml.Name `xml:"http://zarinpal.com/ PaymentRequestWithExtraResponse"`

	Status int32 `xml:"Status,omitempty"`

	Authority string `xml:"Authority,omitempty"`
}

type PaymentVerification struct {
	XMLName xml.Name `xml:"http://zarinpal.com/ PaymentVerification"`

	MerchantID string `xml:"MerchantID,omitempty"`

	Authority string `xml:"Authority,omitempty"`

	Amount int32 `xml:"Amount,omitempty"`
}

type PaymentVerificationResponse struct {
	XMLName xml.Name `xml:"http://zarinpal.com/ PaymentVerificationResponse"`

	Status int32 `xml:"Status,omitempty"`

	RefID int64 `xml:"RefID,omitempty"`
}

type PaymentVerificationWithExtra struct {
	XMLName xml.Name `xml:"http://zarinpal.com/ PaymentVerificationWithExtra"`

	MerchantID string `xml:"MerchantID,omitempty"`

	Authority string `xml:"Authority,omitempty"`

	Amount int32 `xml:"Amount,omitempty"`
}

type PaymentVerificationWithExtraResponse struct {
	XMLName xml.Name `xml:"http://zarinpal.com/ PaymentVerificationWithExtraResponse"`

	Status int32 `xml:"Status,omitempty"`

	RefID int64 `xml:"RefID,omitempty"`

	ExtraDetail string `xml:"ExtraDetail,omitempty"`
}

type GetUnverifiedTransactions struct {
	XMLName xml.Name `xml:"http://zarinpal.com/ GetUnverifiedTransactions"`

	MerchantID string `xml:"MerchantID,omitempty"`
}

type GetUnverifiedTransactionsResponse struct {
	XMLName xml.Name `xml:"http://zarinpal.com/ GetUnverifiedTransactionsResponse"`

	Status int32 `xml:"Status,omitempty"`

	Authorities string `xml:"Authorities,omitempty"`
}

type RefreshAuthority struct {
	XMLName xml.Name `xml:"http://zarinpal.com/ RefreshAuthority"`

	MerchantID string `xml:"MerchantID,omitempty"`

	Authority string `xml:"Authority,omitempty"`

	ExpireIn int32 `xml:"ExpireIn,omitempty"`
}

type RefreshAuthorityResponse struct {
	XMLName xml.Name `xml:"http://zarinpal.com/ RefreshAuthorityResponse"`

	Status int32 `xml:"Status,omitempty"`
}

type PaymentGatewayImplementationServicePortType struct {
	client *SOAPClient
}

func NewPaymentGatewayImplementationServicePortType(url string, tls bool, auth *BasicAuth) *PaymentGatewayImplementationServicePortType {
	if url == "" {
		url = "https://de.zarinpal.com/pg/services/WebGate/service"
	}
	client := NewSOAPClient(url, tls, auth)

	return &PaymentGatewayImplementationServicePortType{
		client: client,
	}
}

func (service *PaymentGatewayImplementationServicePortType) PaymentRequest(request *PaymentRequest) (*PaymentRequestResponse, error) {
	response := new(PaymentRequestResponse)
	err := service.client.Call("#PaymentRequest", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *PaymentGatewayImplementationServicePortType) PaymentRequestWithExtra(request *PaymentRequestWithExtra) (*PaymentRequestWithExtraResponse, error) {
	response := new(PaymentRequestWithExtraResponse)
	err := service.client.Call("#PaymentRequestWithExtra", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *PaymentGatewayImplementationServicePortType) PaymentVerification(request *PaymentVerification) (*PaymentVerificationResponse, error) {
	response := new(PaymentVerificationResponse)
	err := service.client.Call("#PaymentVerification", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *PaymentGatewayImplementationServicePortType) PaymentVerificationWithExtra(request *PaymentVerificationWithExtra) (*PaymentVerificationWithExtraResponse, error) {
	response := new(PaymentVerificationWithExtraResponse)
	err := service.client.Call("#PaymentVerificationWithExtra", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *PaymentGatewayImplementationServicePortType) GetUnverifiedTransactions(request *GetUnverifiedTransactions) (*GetUnverifiedTransactionsResponse, error) {
	response := new(GetUnverifiedTransactionsResponse)
	err := service.client.Call("#GetUnverifiedTransactions", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (service *PaymentGatewayImplementationServicePortType) RefreshAuthority(request *RefreshAuthority) (*RefreshAuthorityResponse, error) {
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
