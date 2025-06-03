package tests

import (
	"ios-api/middlewares"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCORSMiddleware(t *testing.T) {
	// 设置Gin为测试模式
	gin.SetMode(gin.TestMode)

	// 创建路由器
	r := gin.New()
	r.Use(middlewares.CORSMiddleware())

	// 添加测试路由
	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "test"})
	})

	// 测试OPTIONS预检请求
	t.Run("OPTIONS预检请求", func(t *testing.T) {
		req, _ := http.NewRequest("OPTIONS", "/test", nil)
		req.Header.Set("Origin", "http://localhost:3000")
		req.Header.Set("Access-Control-Request-Method", "GET")
		req.Header.Set("Access-Control-Request-Headers", "Content-Type,Authorization")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// 验证状态码
		assert.Equal(t, 204, w.Code)

		// 验证CORS头部
		assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
		assert.Equal(t, "GET, POST, PUT, DELETE, OPTIONS", w.Header().Get("Access-Control-Allow-Methods"))
		assert.Contains(t, w.Header().Get("Access-Control-Allow-Headers"), "Authorization")
		assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
	})

	// 测试实际GET请求
	t.Run("GET请求带CORS头", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Origin", "http://localhost:3000")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// 验证状态码
		assert.Equal(t, 200, w.Code)

		// 验证CORS头部
		assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
		assert.Equal(t, "GET, POST, PUT, DELETE, OPTIONS", w.Header().Get("Access-Control-Allow-Methods"))
		assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))

		// 验证响应内容
		assert.Contains(t, w.Body.String(), "test")
	})

	// 测试无Origin头的请求
	t.Run("无Origin头的请求", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/test", nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// 验证状态码
		assert.Equal(t, 200, w.Code)

		// 验证依然有CORS头部
		assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
	})
}
