package bcfg

import (
	"common/config"

	"gopkg.in/fzerorubigd/onion.v2"
)

// Bcfg is the file module settings
var Bcfg Config

// Config is the file module config
type Config struct {
	Gate struct {
		MerchantID       string `onion:"merchant_id"`
		MerchantOkStatus int64  `onion:"merchant_ok_status"`
		CallbackURL      string `onion:"callback_url"`
		FrontCallbackURL string `onion:"front_callback_url"`
		ZarinURL         string `onion:"zarin_url"`
		APIURL           string `onion:"api_url"`
		Mobile           string `onion:"mobile"`
		Email            string `onion:"email"`
		Description      string `onion:"description"`
	} `onion:"gate"`
	Withdrawal struct {
		MinWithdrawal int64 `onion:min_withdrawal`
	} `onion:"withdrawal"`
}

type configLoader struct {
	o *onion.Onion
}

// Initialize is called when the module is going to add its layer
func (c *configLoader) Initialize(o *onion.Onion) []onion.Layer {
	c.o = o
	def := onion.NewDefaultLayer()
	_ = def.SetDefault("gate.merchant_id", "52efa4ef-9074-4300-afff-01315ee8a9d4")
	_ = def.SetDefault("gate.merchant_ok_status", 100)
	_ = def.SetDefault("gate.callback_url", "/api/campaign/verify/")
	_ = def.SetDefault("gate.front_callback_url", "/v1/verify/")
	_ = def.SetDefault("gate.zarin_url", "https://www.zarinpal.com/pg/StartPay/")
	_ = def.SetDefault("gate.api_url", "https://de.zarinpal.com/pg/services/WebGate/wsdl")
	_ = def.SetDefault("gate.mobile", "09375722346")
	_ = def.SetDefault("gate.email", "dara51php@gmail.com")
	_ = def.SetDefault("gate.description", "Plan Requested")

	_ = def.SetDefault("withdrawal.min_withdrawal", 50000)
	return []onion.Layer{def}
}

// Loaded inform the modules that all layer are ready
func (c *configLoader) Loaded() {
	c.o.GetStruct("gate", &Bcfg.Gate)
}

func init() {
	config.Register(&configLoader{})
}
