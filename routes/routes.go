package routes

import (
	"ios-api/controllers"
	"ios-api/middlewares"
	"ios-api/services"

	"github.com/gin-gonic/gin"
)

// SetupRoutes 设置路由
func SetupRoutes(r *gin.Engine, userService *services.UserService, settingService *services.SettingService, aiService *services.AIService) {
	// 创建微信服务
	wechatService := &services.WechatService{
		AppID:     userService.Config.WechatAppID,
		AppSecret: userService.Config.WechatAppSecret,
	}

	// 创建苹果服务
	appleService := &services.AppleService{
		TeamID:     userService.Config.AppleTeamID,
		KeyID:      userService.Config.AppleKeyID,
		PrivateKey: userService.Config.ApplePrivateKey,
		BundleID:   userService.Config.AppleBundleID,
	}

	// 创建控制器
	userController := &controllers.UserController{
		UserService: userService,
	}

	// 创建OAuth控制器
	oauthController := &controllers.OAuthController{
		UserService:   userService,
		WechatService: wechatService,
		AppleService:  appleService,
	}

	// 创建设置控制器
	settingController := &controllers.SettingController{
		SettingService: settingService,
	}

	// 创建AI控制器
	aiController := controllers.NewAIController(aiService)

	// 无需认证的路由
	v1 := r.Group("/api/v1")
	{
		// 用户注册
		v1.POST("/register", userController.Register)
		// 用户登录
		v1.POST("/login", userController.Login)
		// 第三方登录
		v1.POST("/oauth/login", userController.OAuthLogin)

		// 微信授权相关
		v1.POST("/oauth/wechat/auth", oauthController.WechatAuthURL)
		v1.GET("/oauth/wechat/callback", oauthController.WechatCallback)

		// 苹果授权相关
		v1.POST("/oauth/apple/auth", oauthController.AppleAuth)
		v1.POST("/oauth/apple/callback", oauthController.AppleCallback)

		// 设置相关API（不需要认证）
		v1.GET("/settings/:key", settingController.GetSetting)
		v1.PUT("/settings/:key", settingController.SetSetting)

		// 缓存管理API（不需要认证，但建议在生产环境中添加认证）
		v1.DELETE("/settings/:key/cache", settingController.ClearCache)  // 清除指定key的缓存
		v1.DELETE("/settings/cache", settingController.ClearAllCache)    // 清除所有缓存
		v1.GET("/settings/cache/stats", settingController.GetCacheStats) // 获取缓存统计

		// AI相关API（不需要认证）
		v1.POST("/ai/chat/completions", aiController.ChatCompletion) // 通用AI聊天
		v1.POST("/ai/travel/plan", aiController.GenerateTravelPlan)  // 生成旅行计划
		v1.GET("/ai/models", aiController.GetAvailableModels)        // 获取可用模型
		v1.GET("/ai/status", aiController.GetAIStatus)               // 获取AI服务状态
	}

	// 需要认证的路由
	auth := r.Group("/api/v1")
	auth.Use(middlewares.AuthMiddleware(userService))
	{
		// 退出登录
		auth.POST("/logout", userController.Logout)
		// 获取用户信息
		auth.GET("/user", userController.GetUserInfo)
		// 更新用户信息
		auth.PUT("/user", userController.UpdateUserInfo)
	}
}
