package controllers

import (
	"net/http"

	"ios-api/services"

	"github.com/gin-gonic/gin"
)

// OAuthController OAuth控制器
type OAuthController struct {
	UserService   *services.UserService
	WechatService *services.WechatService
	AppleService  *services.AppleService
}

// 微信授权请求参数
type WechatAuthRequest struct {
	RedirectURI string `json:"redirect_uri" binding:"required"`
	State       string `json:"state"`
}

// 微信授权回调请求参数
type WechatCallbackRequest struct {
	Code  string `form:"code" binding:"required"`
	State string `form:"state"`
}

// 苹果授权请求参数
type AppleAuthRequest struct {
	RedirectURI string `json:"redirect_uri" binding:"required"`
}

// 苹果授权回调请求参数
type AppleCallbackRequest struct {
	Code      string `json:"code"`
	IdToken   string `json:"id_token"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// WechatAuthURL 获取微信授权URL
func (c *OAuthController) WechatAuthURL(ctx *gin.Context) {
	var req WechatAuthRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	// 设置重定向URI
	c.WechatService.RedirectURI = req.RedirectURI

	// 获取授权URL
	authURL := c.WechatService.GetAuthURL(req.State)

	ctx.JSON(http.StatusOK, gin.H{
		"auth_url": authURL,
	})
}

// WechatCallback 处理微信授权回调
func (c *OAuthController) WechatCallback(ctx *gin.Context) {
	var req WechatCallbackRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	// 处理微信回调
	oauthParams, err := c.WechatService.HandleCallback(req.Code)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "微信授权处理失败: " + err.Error()})
		return
	}

	// 使用OAuth参数进行登录
	user, token, err := c.UserService.OAuthLogin(*oauthParams)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "登录失败: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "微信登录成功",
		"user":    user,
		"token":   token,
	})
}

// AppleAuth 苹果授权
func (c *OAuthController) AppleAuth(ctx *gin.Context) {
	var req AppleAuthRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	// 客户端直接处理苹果登录，后端不需要返回URL
	ctx.JSON(http.StatusOK, gin.H{
		"message": "苹果登录需要在客户端实现，请在客户端完成授权后，将授权结果发送到 /api/v1/oauth/apple/callback",
	})
}

// AppleCallback 处理苹果授权回调
func (c *OAuthController) AppleCallback(ctx *gin.Context) {
	var req AppleCallbackRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	// 合并名字字段
	name := req.Name
	if name == "" && (req.FirstName != "" || req.LastName != "") {
		name = req.FirstName + " " + req.LastName
	}

	// 处理苹果回调
	oauthParams, err := c.AppleService.HandleCallback(req.Code, req.IdToken, name, req.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "苹果授权处理失败: " + err.Error()})
		return
	}

	// 使用OAuth参数进行登录
	user, token, err := c.UserService.OAuthLogin(*oauthParams)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "登录失败: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "苹果登录成功",
		"user":    user,
		"token":   token,
	})
}
