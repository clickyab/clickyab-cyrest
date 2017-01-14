package config

import (
	"runtime"

	"common/assert"
	"path/filepath"

	"github.com/fzerorubigd/expand"
	"gopkg.in/fzerorubigd/onion.v2"
	_ "gopkg.in/fzerorubigd/onion.v2/yamlloader" // for loading yaml files
)

const appName = "cyrest"

//Config is the global application config instance
var Config AppConfig

var o = onion.New()

// AppConfig is the application config
type AppConfig struct {
	DevelMode       bool   `onion:"devel_mode"`
	CORS            bool   `onion:"cors"`
	MaxCPUAvailable int    `onion:"max_cpu_available"`
	MountPoint      string `onion:"mount_point"`
	FrontMountPoint string `onion:"front_mount_point"`
	Profile         string

	Site      string
	Proto     string
	FrontPath string `onion:"front_path"`

	Port        string
	StaticRoot  string `onion:"static_root"`
	SwaggerRoot string `onion:"swagger_root"`
	ProfileRoot string `onion:"profile_root"`
	TimeZone    string `onion:"time_zone"`

	Redis struct {
		Size     int
		Network  string
		Address  string
		Password string
		Database int
	}

	Mysql struct {
		DSN               string `onion:"dsn"`
		DataBase          string `onion:"database"`
		MaxConnection     int    `onion:"max_connection"`
		MaxIdleConnection int    `onion:"max_idle_connection"`
	}

	AMQP struct {
		DSN        string
		Exchange   string
		Publisher  int
		ConfirmLen int
	}

	Page struct {
		PerPage    int `onion:"per_page"`
		MaxPerPage int `onion:"max_per_page"`
		MinPerPage int `onion:"min_per_page"`
	}

	Telegram struct {
		APIKey string `onion:"api_key"`
	}

	Redmine struct {
		APIKey         string
		URL            string
		ProjectID      int `onion:"project_id"`
		Active         bool
		NewIssueTypeID int `onion:"new_issue_type_id"`
	}

	Slack struct {
		Channel    string
		Username   string
		WebHookURL string
		Active     bool
	}
}

func defaultLayer() onion.DefaultLayer {
	res := onion.NewDefaultLayer()
	assert.Nil(res.SetDefault("site", "localhost"))
	assert.Nil(res.SetDefault("mount_point", "/api"))
	assert.Nil(res.SetDefault("front_mount_point", "/v1"))
	assert.Nil(res.SetDefault("devel_mode", true))
	assert.Nil(res.SetDefault("cors", true))
	assert.Nil(res.SetDefault("profile", "cpu"))
	assert.Nil(res.SetDefault("max_cpu_available", runtime.NumCPU()))
	assert.Nil(res.SetDefault("proto", "http"))
	assert.Nil(res.SetDefault("port", ":80"))

	path, err := expand.Path("$PWD/../statics/")
	assert.Nil(err)
	assert.Nil(res.SetDefault("static_root", path))
	path, err = expand.Path("$PWD/../3rd/swagger/")
	assert.Nil(err)
	path, err = filepath.Abs(path)
	assert.Nil(err)
	assert.Nil(res.SetDefault("swagger_root", path))

	path, err = expand.Path("$PWD/../front/public/")
	assert.Nil(err)
	path, err = filepath.Abs(path)
	assert.Nil(err)
	assert.Nil(res.SetDefault("front_path", path))

	path, err = expand.Path("$PWD/../tmp/profiles/")
	assert.Nil(err)
	path, err = filepath.Abs(path)
	assert.Nil(err)
	assert.Nil(res.SetDefault("profile_root", path))

	assert.Nil(res.SetDefault("redis.size", 10))
	assert.Nil(res.SetDefault("redis.network", "tcp"))
	assert.Nil(res.SetDefault("redis.address", ":6379"))
	assert.Nil(res.SetDefault("redis.password", ""))
	assert.Nil(res.SetDefault("redis.database", 0))

	assert.Nil(res.SetDefault("mysql.dsn", "root:bita123@/"))
	assert.Nil(res.SetDefault("mysql.database", "cyrest"))
	assert.Nil(res.SetDefault("mysql.max_connection", 100))
	assert.Nil(res.SetDefault("mysql.max_idle_connection", 10))

	assert.Nil(res.SetDefault("amqp.publisher", 30))
	assert.Nil(res.SetDefault("amqp.exchange", "cy"))
	assert.Nil(res.SetDefault("amqp.dsn", "amqp://cyrest:bita123@127.0.0.1:5672/"))
	assert.Nil(res.SetDefault("amqp.confirmlen", 50))

	assert.Nil(res.SetDefault("page.per_Page", 10))
	assert.Nil(res.SetDefault("page.max_per_page", 100))
	assert.Nil(res.SetDefault("page.min_per_Page", 1))

	assert.Nil(res.SetDefault("telegram.api_key", "273335144:AAEv4uPeo68X7Scc3MLKxwMO1YI3JFkWiJM"))

	assert.Nil(res.SetDefault("time_zone", "Asia/Tehran"))

	assert.Nil(res.SetDefault("slack.channel", "notifications"))
	assert.Nil(res.SetDefault("slack.username", "BigBrother"))
	assert.Nil(res.SetDefault("slack.webhookurl", "https://hooks.slack.com/services/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"))
	assert.Nil(res.SetDefault("slack.active", false))

	return res
}
