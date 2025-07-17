package middleware

import (
	"blog/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

// Auth 认证中间件
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取Authorization头
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.Warn("未提供认证令牌")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供认证令牌"})
			c.Abort()
			return
		}

		// 检查Bearer前缀
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.Warn("无效的认证令牌格式")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的认证令牌格式"})
			c.Abort()
			return
		}

		// 解析JWT令牌
		token, err := jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
			return []byte(utils.JWTSecret), nil
		})

		if err != nil {
			utils.Error("解析JWT令牌失败", zap.Error(err))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的认证令牌"})
			c.Abort()
			return
		}

		// 验证令牌
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// 将用户ID存储在上下文中
			c.Set("userID", uint(claims["user_id"].(float64)))
			c.Next()
		} else {
			utils.Warn("无效的认证令牌")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的认证令牌"})
			c.Abort()
			return
		}
	}
}
