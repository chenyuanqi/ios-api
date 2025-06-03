package main

import (
	"fmt"
	"log"
	"regexp"

	"ios-api/config"
	"ios-api/middlewares"
	"ios-api/routes"
	"ios-api/services"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 连接主数据库
	db, err := gorm.Open(mysql.Open(cfg.GetDSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("连接主数据库失败: %v", err)
	}

	// 连接通用数据库（用于settings表）
	generalDB, err := gorm.Open(mysql.Open(cfg.GetGeneralDSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("连接通用数据库失败: %v", err)
	}

	// 创建用户服务
	userService := &services.UserService{
		DB:        db,
		JWTSecret: cfg.JWTSecret,
		Config:    cfg,
	}

	// 创建设置服务
	settingService := &services.SettingService{
		DB:   generalDB,
		Salt: cfg.SettingSalt,
	}

	// 创建 Gin 实例
	r := gin.Default()

	pattern := `^https?://([a-z0-9-]+\.)?chenyuanqi\.com(:[0-9]+)?$`
	regex, err := regexp.Compile(pattern)
	if err != nil {
		panic("无法解析 CORS 正则：" + err.Error())
	}

	corsCfg := middlewares.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:8080", "http://127.0.0.1:5500"},
		RegexAllowOrigin: regex,
	}

	// 配置CORS中间件
	r.Use(middlewares.CORSMiddleware(corsCfg))

	// 设置路由
	routes.SetupRoutes(r, userService, settingService)

	// 启动服务器
	port := fmt.Sprintf(":%d", cfg.AppPort)
	log.Printf("服务器启动，监听端口: %s", port)
	if err := r.Run(port); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
