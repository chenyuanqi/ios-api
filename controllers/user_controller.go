package controllers

import (
	"ios-api/services"
	"ios-api/utils"

	"github.com/gin-gonic/gin"
)

// UserController 用户控制器
type UserController struct {
	UserService *services.UserService
}

// Register 注册用户
func (c *UserController) Register(ctx *gin.Context) {
	var params services.RegisterParams
	if err := ctx.ShouldBindJSON(&params); err != nil {
		utils.ParamError(ctx, "请求参数错误: "+err.Error())
		return
	}

	user, token, err := c.UserService.Register(params)
	if err != nil {
		if err == services.ErrEmailExists {
			utils.Conflict(ctx, err.Error())
		} else {
			utils.ServerError(ctx, err.Error())
		}
		return
	}

	utils.Created(ctx, "注册成功", gin.H{
		"user":  user,
		"token": token,
	})
}

// Login 用户登录
func (c *UserController) Login(ctx *gin.Context) {
	var params services.LoginParams
	if err := ctx.ShouldBindJSON(&params); err != nil {
		utils.ParamError(ctx, "请求参数错误: "+err.Error())
		return
	}

	user, token, err := c.UserService.Login(params)
	if err != nil {
		if err == services.ErrUserNotFound || err == services.ErrInvalidPassword {
			utils.Unauthorized(ctx, err.Error())
		} else {
			utils.ServerError(ctx, err.Error())
		}
		return
	}

	utils.Success(ctx, "登录成功", gin.H{
		"user":  user,
		"token": token,
	})
}

// OAuthLogin 第三方登录
func (c *UserController) OAuthLogin(ctx *gin.Context) {
	var params services.OAuthLoginParams
	if err := ctx.ShouldBindJSON(&params); err != nil {
		utils.ParamError(ctx, "请求参数错误: "+err.Error())
		return
	}

	// 校验 provider 参数值
	if params.Provider != "wechat" && params.Provider != "apple" {
		utils.ParamError(ctx, "不支持的登录方式，仅支持 wechat 和 apple")
		return
	}

	user, token, err := c.UserService.OAuthLogin(params)
	if err != nil {
		utils.ServerError(ctx, err.Error())
		return
	}

	utils.Success(ctx, "登录成功", gin.H{
		"user":  user,
		"token": token,
	})
}

// Logout 用户退出登录
func (c *UserController) Logout(ctx *gin.Context) {
	token, _ := ctx.Get("token")
	tokenStr, ok := token.(string)
	if !ok {
		utils.ServerError(ctx, "获取会话信息失败")
		return
	}

	if err := c.UserService.Logout(tokenStr); err != nil {
		utils.ServerError(ctx, err.Error())
		return
	}

	utils.Success(ctx, "退出登录成功", nil)
}

// GetUserInfo 获取用户信息
func (c *UserController) GetUserInfo(ctx *gin.Context) {
	userID, _ := ctx.Get("userID")
	userIDUint, ok := userID.(uint)
	if !ok {
		utils.ServerError(ctx, "获取用户信息失败")
		return
	}

	user, err := c.UserService.GetUserByID(userIDUint)
	if err != nil {
		if err == services.ErrUserNotFound {
			utils.NotFound(ctx, err.Error())
		} else {
			utils.ServerError(ctx, err.Error())
		}
		return
	}

	utils.Success(ctx, "获取用户信息成功", gin.H{
		"user": user,
	})
}

// UpdateUserInfo 更新用户信息
func (c *UserController) UpdateUserInfo(ctx *gin.Context) {
	userID, _ := ctx.Get("userID")
	userIDUint, ok := userID.(uint)
	if !ok {
		utils.ServerError(ctx, "获取用户信息失败")
		return
	}

	var params services.UpdateUserParams
	if err := ctx.ShouldBindJSON(&params); err != nil {
		utils.ParamError(ctx, "请求参数错误: "+err.Error())
		return
	}

	user, err := c.UserService.UpdateUser(userIDUint, params)
	if err != nil {
		if err == services.ErrUserNotFound {
			utils.NotFound(ctx, err.Error())
		} else {
			utils.ServerError(ctx, err.Error())
		}
		return
	}

	utils.Success(ctx, "更新用户信息成功", gin.H{
		"user": user,
	})
}
