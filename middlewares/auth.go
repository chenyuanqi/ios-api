package middlewares

import (
	"net/http"
	"strings"

	"ios-api/services"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware 认证中间件
func AuthMiddleware(userService *services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Authorization header 获取 token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供认证信息"})
			c.Abort()
			return
		}

		// Bearer Token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "认证格式错误"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 验证 token
		userID, err := userService.VerifyToken(tokenString)
		if err != nil {
			errMsg := "认证失败"
			if err == services.ErrTokenExpired {
				errMsg = "登录已过期，请重新登录"
			} else if err == services.ErrInvalidToken {
				errMsg = "无效的认证信息"
			}
			c.JSON(http.StatusUnauthorized, gin.H{"error": errMsg})
			c.Abort()
			return
		}

		// 将用户ID存入上下文
		c.Set("userID", userID)
		c.Set("token", tokenString)

		c.Next()
	}
}
