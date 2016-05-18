package base

import (
	"common/middlewares"
	"sync"

	"common/config"

	"github.com/gin-gonic/gin"
)

// Routes the base rote structure
type Routes interface {
	// Routes is for adding new controller
	Routes(r *gin.Engine, mountPoint string)
}

var (
	engine *gin.Engine
	all    []Routes
	once   = &sync.Once{}
)

// Register a new controller class
func Register(c ...Routes) {
	all = append(all, c...)
}

// Initialize the controller
func Initialize(mountPoint string) *gin.Engine {
	once.Do(func() {
		engine = gin.New()
		mid := []gin.HandlerFunc{middlewares.Recovery, middlewares.Logger}
		if config.Config.CORS {
			mid = append(mid, middlewares.CORSMiddlewareGenerator())
		}
		engine.Use(mid...)
		for i := range all {
			all[i].Routes(engine, mountPoint)
		}
	})

	return engine
}
