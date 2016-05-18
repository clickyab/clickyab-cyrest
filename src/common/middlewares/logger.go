package middlewares

import (
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

// Logger is the middleware for log system
func Logger(c *gin.Context) {
	// Start timer
	start := time.Now()
	path := c.Request.URL.Path

	// Process request
	c.Next()

	latency := time.Since(start)

	clientIP := c.ClientIP()
	method := c.Request.Method
	statusCode := c.Writer.Status()
	comment := c.Errors.ByType(gin.ErrorTypePrivate).String()
	logrus.WithFields(
		logrus.Fields{
			"Method":   method,
			"Path":     path,
			"Latency":  latency,
			"ClientIP": clientIP,
			"Status":   statusCode,
			"Comment":  comment,
		},
	).Info(http.StatusText(statusCode))
}
