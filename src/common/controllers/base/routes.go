package base

import (
	"common/config"
	"common/middlewares"
	"sync"

	"gopkg.in/labstack/echo.v3"
)

// Routes the base rote structure
type Routes interface {
	// Routes is for adding new controller
	Routes(r *echo.Echo, mountPoint string)
}

var (
	engine *echo.Echo
	all    []Routes
	once   = &sync.Once{}
)

// Register a new controller class
func Register(c ...Routes) {
	all = append(all, c...)
}

// Initialize the controller
func Initialize(mountPoint string) *echo.Echo {
	once.Do(func() {
		engine = echo.New()
		mid := []echo.MiddlewareFunc{middlewares.Recovery, middlewares.Logger}
		if config.Config.CORS {
			mid = append(mid, middlewares.CORS())
		}
		engine.Use(mid...)
		for i := range all {
			all[i].Routes(engine, mountPoint)
		}
	})
	//engine.SetLogLevel(log.DEBUG)
	if config.Config.DevelMode {
		engine.Static("/swagger", config.Config.SwaggerRoot)
	}
	engine.Logger = NewLogger()
	return engine
}
