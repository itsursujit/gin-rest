package Middleware

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var timeFormat = "02/Jan/2006:15:04:05 -0700"

// Logger is the logrus logger handler
func HttpLogger(logger logrus.FieldLogger) gin.HandlerFunc {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}
	return func(c *gin.Context) {
		// other handler can change c.Path so:
		path := c.Request.URL.Path
		query := c.Request.URL.Query()
		start := time.Now()
		c.Next()
		latency := fmt.Sprintf("%s", time.Since(start))
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		clientUserAgent := c.Request.UserAgent()
		referer := c.Request.Referer()
		dataLength := c.Writer.Size()
		requestId, _ := c.Get("RequestId")
		if dataLength < 0 {
			dataLength = 0
		}

		entry := logger.WithFields(logrus.Fields{
			"hostname":   hostname,
			"statusCode": statusCode,
			"latency":    latency, // time to process
			"clientIP":   clientIP,
			"method":     c.Request.Method,
			"path":       path,
			"query":      query,
			"referer":    referer,
			"dataLength": dataLength,
			"userAgent":  clientUserAgent,
			"requestId":  requestId,
		})

		if len(c.Errors) > 0 {
			entry.Error(c.Errors.ByType(gin.ErrorTypePrivate).String())
		} else {
			if statusCode > 499 {
				entry.Error()
			} else if statusCode > 399 {
				entry.Warn()
			} else {
				entry.Info()
			}
		}
	}
}
