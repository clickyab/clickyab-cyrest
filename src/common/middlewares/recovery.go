package middlewares

import (
	"common/config"
	"common/utils"
	"fmt"
	"net/http"
	"net/http/httputil"
	"runtime/debug"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

// Recovery is the middleware to prevent the panic to crash the app
func Recovery(ctx *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			ctx.JSON(
				http.StatusInternalServerError,
				struct {
					Error string `json:"error"`
				}{
					Error: http.StatusText(http.StatusInternalServerError),
				},
			)
			ctx.Abort()
			stack := debug.Stack()
			dump, _ := httputil.DumpRequest(ctx.Request, true)
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

	ctx.Next()
}
