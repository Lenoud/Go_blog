package middleware

import (
	"blog/utils"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Logger 日志中间件
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// 处理请求
		c.Next()

		// 记录请求日志
		cost := time.Since(start)
		utils.Info("Request",
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("method", c.Request.Method),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("cost", cost),
		)

		// 如果发生错误，记录错误日志
		if len(c.Errors) > 0 {
			utils.Error("Request Error",
				zap.String("path", path),
				zap.String("query", query),
				zap.String("ip", c.ClientIP()),
				zap.String("method", c.Request.Method),
				zap.Int("status", c.Writer.Status()),
				zap.Duration("cost", cost),
				zap.String("error", c.Errors.String()),
			)
		}
	}
}
