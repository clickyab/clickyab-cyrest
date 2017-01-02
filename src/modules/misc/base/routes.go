package base

import (
	"common/config"
	"common/utils"
	"os"
	"path/filepath"
	"sync"

	"modules/misc/middlewares"

	"modules/misc/trans"
	"net/http"
	"strings"

	"fmt"

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
		fmt.Println(config.Config.FrontMountPoint+"/", c.Request().URL.Path)
		if strings.HasPrefix(c.Request().URL.Path, config.Config.FrontMountPoint+"/") {
			// this request is from the front so return the index.html
			f := filepath.Join(config.Config.FrontPath, "index.html")
			return c.File(f)
		}

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
		engine.Logger = NewLogger()
		echo.NotFoundHandler = notFoundHandler
		echo.MethodNotAllowedHandler = methodNotAllowedHandler
	})

	return engine
}
