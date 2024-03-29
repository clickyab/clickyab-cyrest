package middlewares

import (
	"common/config"
	"common/utils"
	"fmt"
	"net/http"
	"runtime/debug"

	"net/http/httputil"

	"common/assert"

	"github.com/Sirupsen/logrus"
	"gopkg.in/labstack/echo.v3"
)

// Recovery is the middleware to prevent the panic to crash the app
func Recovery(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		defer func() {
			if err := recover(); err != nil {
				err := ctx.JSON(
					http.StatusInternalServerError,
					struct {
						Error string `json:"error"`
					}{
						Error: http.StatusText(http.StatusInternalServerError),
					},
				)
				assert.Nil(err)

				stack := debug.Stack()
				dump, _ := httputil.DumpRequest(ctx.Request(), true)
				data := fmt.Sprintf("Request : \n %s \n\nStack : \n %s", dump, stack)
				logrus.WithField("error", err).Warn(err, data)
				if config.Config.Redmine.Active {
					go utils.RedmineDoError(err, []byte(data))
				}

				if config.Config.Slack.Active {
					go utils.SlackDoMessage(err, ":shit:", utils.SlackAttachment{Text: data, Color: "#AA3939"})
				}
			}
		}()

		return next(ctx)
	}
}
