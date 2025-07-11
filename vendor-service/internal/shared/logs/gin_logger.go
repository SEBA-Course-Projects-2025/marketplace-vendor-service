package logs

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

func GinLogger() gin.HandlerFunc {

	return func(c *gin.Context) {

		c.Next()

		latency := time.Since(time.Now())
		status := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.URL.Path

		logLine := fmt.Sprintf("%s \"%s %s\" %d %s", time.Now().Format(time.RFC3339), method, path, status, latency)

		labels := map[string]string{
			"job":    "vendor_service",
			"status": fmt.Sprintf("%d", status),
			"method": method,
			"path":   path,
		}

		go func() {
			if err := SendLogsToLoki(logLine, labels); err != nil {
				fmt.Println("Loki Gin logs error:", err)
			}
		}()

	}

}
