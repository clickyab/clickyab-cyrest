package models

import (
	"common/assert"
	"common/config"
	"common/models/common"
	"common/utils"
	"database/sql"
	"errors"
	"sync"

	"common/try"

	"strings"

	"common/initializer"

	"github.com/Sirupsen/logrus"
	"gopkg.in/gorp.v1"
)

var (
	dbmap *gorp.DbMap
	db    *sql.DB
	once  = sync.Once{}
	all   []utils.Initializer
)

// Hooker interface :))))) You have a dirty mind.
type Hooker interface {
	// AddHook is called after initialize only if the manager implement it
	AddHook()
}

type gorpLogger struct {
}

type modelsInitializer struct {
}

func (g gorpLogger) Printf(format string, v ...interface{}) {
	logrus.Infof(format, v...)
}

// Initialize the modules, its safe to call this as many time as you want.
func (modelsInitializer) Initialize() {
	once.Do(func() {
		var err error
		db, err = sql.Open("mysql", config.Config.Mysql.DSN+config.Config.Mysql.DataBase+"?parseTime=true&charset=utf8")
		assert.Nil(err)

		db.SetMaxIdleConns(config.Config.Mysql.MaxIdleConnection)
		db.SetMaxOpenConns(config.Config.Mysql.MaxConnection)
		err = db.Ping()
		assert.Nil(err)

		dbmap = &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{}}

		if config.Config.DevelMode {
			logger := gorpLogger{}
			dbmap.TraceOn("[db]", logger)
		} else {
			dbmap.TraceOff()
		}
		common.Initialize(db, dbmap)
		for i := range all {
			all[i].Initialize()

		}
		// If they are hooker call them.
		for i := range all {
			if h, ok := all[i].(Hooker); ok {
				h.AddHook()
			}
		}

		// Add the hook to error
		try.CatchHook(func(err error) error {
			if err == sql.ErrNoRows {
				return errors.New("not found")
			}
			if strings.HasPrefix(err.Error(), "gorp:") {
				// this is not correct
				logrus.Panic(err)
			}
			return err
		})
	})
}
func (modelsInitializer) Finalize() {
	logrus.Debug("models are done")
}

// Register a new initializer module
func Register(m ...utils.Initializer) {
	all = append(all, m...)
}

func init() {
	initializer.Register(&modelsInitializer{})
}
