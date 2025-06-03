package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"regexp"
	"syscall"

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

	// 创建设置服务（带缓存）
	settingService, err := services.NewSettingService(generalDB, cfg.SettingSalt, cfg.CacheDir)
	if err != nil {
		log.Fatalf("创建设置服务失败: %v", err)
	}

	// 设置优雅关闭
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("正在关闭服务器...")

		// 关闭缓存连接
		if err := settingService.Close(); err != nil {
			log.Printf("关闭缓存失败: %v", err)
		}

		os.Exit(0)
	}()

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
	log.Printf("缓存目录: %s", cfg.CacheDir)
	if err := r.Run(port); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
