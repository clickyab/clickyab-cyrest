package ucfg

import (
	"common/config"
	"time"

	"gopkg.in/fzerorubigd/onion.v2"
)

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

	return []onion.Layer{def}
}

// Loaded inform the modules that all layer are ready
func (c *configLoader) Loaded() {
	Cfg.TokenTimeout = c.o.GetDurationDefault("user.token_timeout", 48*time.Hour)
	c.o.GetStruct("sms", &Cfg.SMS)
	c.o.GetStruct("smtp", &Cfg.SMTP)
}

func init() {
	config.Register(&configLoader{})
}
