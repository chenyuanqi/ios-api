package routes

import (
	"ios-api/controllers"
	"ios-api/middlewares"
	"ios-api/services"

	"github.com/gin-gonic/gin"
)

// SetupRoutes 设置路由
func SetupRoutes(r *gin.Engine, userService *services.UserService) {
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
