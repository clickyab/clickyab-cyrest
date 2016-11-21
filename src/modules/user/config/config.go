package ucfg

import (
	"common/config"
	"time"

	"gopkg.in/fzerorubigd/onion.v2"
)

// Cfg is the user module settings
var Cfg Config

// Config is the user module config
type Config struct {
	TokenTimeout time.Duration

	SMTP struct {
		Host     string
		Port     int
		User     string
		Password string
		Sender   string
		Name     string
	} `onion:"smtp"`

	SMS struct {
		Sender   string
		User     string
		Password string
		Send     bool
	} `onion:"sms"`

	OAuth struct {
		ClientID         string `onion:"client_id"`
		ClientSecret     string `onion:"client_secret"`
		RedirectURI      string `onion:"redirect_uri"`
		LoginRedirect    string `onion:"login_redirect"`
		RegisterRedirect string `onion:"register_redirect"`
	} `onion:"oauth"`
}

type configLoader struct {
	o *onion.Onion
}

// Initialize is called when the module is going to add its layer
func (c *configLoader) Initialize(o *onion.Onion) []onion.Layer {
	c.o = o
	def := onion.NewDefaultLayer()
	_ = def.SetDefault("user.token_timeout", "48h")
	_ = def.SetDefault("smtp.host", "127.0.0.1")
	_ = def.SetDefault("smtp.port", 1025)
	_ = def.SetDefault("smtp.user", "")
	_ = def.SetDefault("smtp.password", "")
	_ = def.SetDefault("smtp.sender", "app@azmoona.com")
	_ = def.SetDefault("smtp.name", "Azmoona")

	_ = def.SetDefault("sms.sender", "30007957954893")
	_ = def.SetDefault("sms.user", "dev@azmoona.com")
	_ = def.SetDefault("sms.password", "bita123*")

	_ = def.SetDefault("oauth.client_id", "975262007411-1nk67l3s49ua2lt41pr805flr5a8c1n5.apps.googleusercontent.com")
	_ = def.SetDefault("oauth.client_secret", "kxddsIpmuSWJ3iAo1Ghs_uR6")
	_ = def.SetDefault("oauth.redirect_uri", "http://home.rubi.gd/api/user/oauth/callback")
	_ = def.SetDefault("oauth.login_redirect", "/login")
	_ = def.SetDefault("oauth.register_redirect", "/register")

	return []onion.Layer{def}
}

// Loaded inform the modules that all layer are ready
func (c *configLoader) Loaded() {
	Cfg.TokenTimeout = c.o.GetDurationDefault("user.token_timeout", 48*time.Hour)
	c.o.GetStruct("sms", &Cfg.SMS)
	c.o.GetStruct("smtp", &Cfg.SMTP)
	c.o.GetStruct("oauth", &Cfg.OAuth)
}

func init() {
	config.Register(&configLoader{})
}
