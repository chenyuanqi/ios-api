package main

import (
	"fmt"
	"log"

	"ios-api/config"
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

	// 连接数据库
	db, err := gorm.Open(mysql.Open(cfg.GetDSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}

	// 创建用户服务
	userService := &services.UserService{
		DB:        db,
		JWTSecret: cfg.JWTSecret,
	}

	// 创建 Gin 实例
	r := gin.Default()

	// 设置路由
	routes.SetupRoutes(r, userService)

	// 启动服务器
	port := fmt.Sprintf(":%d", cfg.AppPort)
	log.Printf("服务器启动，监听端口: %s", port)
	if err := r.Run(port); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
