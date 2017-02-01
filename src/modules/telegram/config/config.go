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
		ChannelID int64  `onion:"channel_id"`

		CLIAddress        string        `onion:"cli_host"`
		CLIPort           int           `onion:"cli_port"`
		LastPostChannel   int           `onion:"last_post_channel"`
		LimitCountWarning int64         `onion:"limit_count_warning"`
		TimeReQueUe       time.Duration `onion:"time_requeue"`
		PositionAdDefault int64         `onion:"position_ad_default"`
	}
}

// Initialize is called when the module is going to add its layer
func (c *Config) Initialize(o *onion.Onion) []onion.Layer {
	c.o = o
	d := onion.NewDefaultLayer()
	assert.Nil(d.SetDefault("telegram.admins", "70018667"))
	assert.Nil(d.SetDefault("telegram.verify", "70018667"))

	assert.Nil(d.SetDefault("telegram.channel_id", ""))
	assert.Nil(d.SetDefault("telegram.api_key", "232630313:AAHRVcaQxFvs3u2-VGAAlsD3Xe1TIUr5rhk"))
	assert.Nil(d.SetDefault("telegram.bot_id", "$0100000068c34a10ed72226be64e8d4d"))
	assert.Nil(d.SetDefault("telegram.cli_host", "localhost"))
	assert.Nil(d.SetDefault("telegram.cli_port", 9999))
	assert.Nil(d.SetDefault("telegram.last_post_channel", 20))
	assert.Nil(d.SetDefault("telegram.limit_count_warning", 5))
	assert.Nil(d.SetDefault("telegram.time_requeue", 5*time.Minute))
	assert.Nil(d.SetDefault("telegram.position_ad_default", 10))
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
