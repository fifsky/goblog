package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ilibs/logger"
)

func Ginrus() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		// some evil middlewares modify this values
		path := c.Request.URL.Path
		c.Next()

		end := time.Now()
		latency := end.Sub(start)
		end = end.Local()

		info := map[string]interface{}{
			"status":     c.Writer.Status(),
			"method":     c.Request.Method,
			"path":       path,
			"ip":         c.ClientIP(),
			"latency":    latency,
			"user-agent": c.Request.UserAgent(),
			"time":       end.Format("2006-01-02 15:04:05"),
		}

		log := logger.NewLogger(func(c *logger.Config) {
			c.LogName = "access_log"
			c.LogDetail = false
		})

		if len(c.Errors) > 0 {
			// Append error field if this is an erroneous request.
			log.Errorf(c.Errors.String(),info)
		} else {
			log.Info(info)
		}
	}
}
