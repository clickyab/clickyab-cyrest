package base

import (
	"common/config"
	"common/utils"
	"modules/misc/middlewares"
	"modules/misc/trans"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/Sirupsen/logrus"
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

func notFoundHandler(c echo.Context) error {
	fP := filepath.Join(config.Config.FrontPath, c.Request().URL.Path)
	if _, err := os.Stat(fP); os.IsNotExist(err) {
		return c.JSON(http.StatusNotFound, ErrorResponseSimple{
			Error: trans.E(http.StatusText(http.StatusNotFound)),
		})
	}

	return c.File(fP)
}

func methodNotAllowedHandler(c echo.Context) error {
	return c.JSON(http.StatusMethodNotAllowed, ErrorResponseSimple{
		Error: trans.E(http.StatusText(http.StatusMethodNotAllowed)),
	})
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
		engine.Any(config.Config.FrontMountPoint+"/*", func(ctx echo.Context) error {
			// this request is from the front so return the index.html
			f := filepath.Join(config.Config.FrontPath, "index.html")
			return ctx.File(f)
		})
		if config.Config.DevelMode {
			err := utils.ChangeInFile(filepath.Join(config.Config.SwaggerRoot, "cyrest.yaml"), "swaggerbase", config.Config.Site)
			if err != nil {
				logrus.Warnf(
					"can not change the %s file and set the site %s",
					filepath.Join(config.Config.SwaggerRoot, "cyrest.yaml"),
					config.Config.Site,
				)
			}

			engine.Static("/swagger", config.Config.SwaggerRoot)
		}
		engine.Static("/statics", config.Config.StaticRoot)

		engine.Logger = NewLogger()
		echo.NotFoundHandler = notFoundHandler
		echo.MethodNotAllowedHandler = methodNotAllowedHandler
	})

	return engine
}
