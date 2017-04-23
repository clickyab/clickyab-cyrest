package tcfg

import (
	"common/assert"
	"common/config"
	"strconv"
	"strings"

	"time"

	onion "gopkg.in/fzerorubigd/onion.v2"
)

// Cfg is the telegram module config
var Cfg Config

// Config is the telegram config
type Config struct {
	o *onion.Onion
	// Chat id of admins
	admins []int64
	verify []int64

	Telegram struct {
		APIKey    string `onion:"api_key"`
		BotID     string `onion:"bot_id"`
		BotName   string `onion:"bot_name"`
		MsgCount  int    `onion:"message_count"`
		MsgOffset int    `onion:"message_offset"`

		CLIAddress        string        `onion:"cli_host"`
		CLIPort           int           `onion:"cli_port"`
		LastPostChannel   int           `onion:"last_post_channel"`
		LimitCountWarning int64         `onion:"limit_count_warning"`
		TimeReQueUe       time.Duration `onion:"time_requeue"`
		PositionAdDefault int64         `onion:"position_ad_default"`
		SendDelay         time.Duration `onion:"send_delay"`
	}
}

// Initialize is called when the module is going to add its layer
func (c *Config) Initialize(o *onion.Onion) []onion.Layer {
	c.o = o
	d := onion.NewDefaultLayer()
	assert.Nil(d.SetDefault("telegram.admins", "70018667"))
	assert.Nil(d.SetDefault("telegram.verify", "70018667"))

	assert.Nil(d.SetDefault("telegram.message_count", "50"))
	assert.Nil(d.SetDefault("telegram.message_offset", "0"))
	assert.Nil(d.SetDefault("telegram.api_key", "347601159:AAEangmt4d67iRwd3-afAaKINzQJKA6q6G4"))
	assert.Nil(d.SetDefault("telegram.bot_name", "rubikaddemobot"))
	assert.Nil(d.SetDefault("telegram.cli_host", "localhost"))
	assert.Nil(d.SetDefault("telegram.cli_port", 9999))
	assert.Nil(d.SetDefault("telegram.last_post_channel", 20))
	assert.Nil(d.SetDefault("telegram.limit_count_warning", 5))
	assert.Nil(d.SetDefault("telegram.time_requeue", 2*time.Minute))
	assert.Nil(d.SetDefault("telegram.position_ad_default", 10))
	assert.Nil(d.SetDefault("telegram.send_delay", 2*time.Second))
	return []onion.Layer{d}
}

// Loaded inform the modules that all layer are ready
func (c *Config) Loaded() {
	admins := c.o.GetString("telegram.admins")
	adminsArray := strings.Split(admins, ",")
	for i := range adminsArray {
		x, err := strconv.ParseInt(adminsArray[i], 10, 0)
		assert.Nil(err)
		c.admins = append(c.admins, x)
	}

	verify := c.o.GetString("telegram.verify")
	verifyArray := strings.Split(verify, ",")
	for i := range verifyArray {
		x, err := strconv.ParseInt(verifyArray[i], 10, 0)
		assert.Nil(err)
		c.verify = append(c.verify, x)
	}

	c.o.GetStruct("", c)
	c.Telegram.TimeReQueUe = c.o.GetDuration("telegram.time_requeue")
	c.Telegram.SendDelay = c.o.GetDuration("telegram.send_delay")
}

// IsAdmin check if the current user is admin
func (c *Config) IsAdmin(chID int64) bool {
	for i := range c.admins {
		if c.admins[i] == chID {
			return true
		}
	}

	return false
}

// IsVerify check if the current user is admin
func (c *Config) IsVerify(chID int64) bool {
	for i := range c.verify {
		if c.verify[i] == chID {
			return true
		}
	}

	return false
}

func init() {
	config.Register(&Cfg)
}
