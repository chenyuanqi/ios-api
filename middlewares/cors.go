package middlewares

import (
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORSConfig 用于配置允许的域名白名单或正则规则
type CORSConfig struct {
	// 允许的精确白名单域名（例如 "http://localhost:3000", "http://localhost:8080" 等）
	AllowOrigins []string
	// 允许的正则表达式规则，用来匹配子域名或更复杂的情况
	// 例如：`^https?://([a-z0-9-]+\.)?chenyuanqi\.com(:[0-9]+)?$`
	RegexAllowOrigin *regexp.Regexp
	// 如果需要，还可以添加：AllowMethods、AllowHeaders、ExposeHeaders、AllowCredentials 等
}

// CORSMiddleware 跨域中间件
func CORSMiddleware(config CORSConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin == "" {
			// 如果请求没有 Origin 头，就直接放行（或者你可以选择禁止）
			c.Next()
			return
		}

		allowed := false

		// 1. 先检查白名单数组
		for _, ao := range config.AllowOrigins {
			if strings.EqualFold(ao, origin) {
				allowed = true
				break
			}
		}

		// 2. 如果白名单里没找到，再用正则去匹配
		if !allowed && config.RegexAllowOrigin != nil {
			if config.RegexAllowOrigin.MatchString(origin) {
				allowed = true
			}
		}

		if allowed {
			// 设置允许跨域的Headers - 返回具体的origin而不是*
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, User-Agent, Content-Length, X-Requested-With")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		method := c.Request.Method

		// 处理预检请求
		if method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		// 处理跨域
		c.Next()
	}
}
