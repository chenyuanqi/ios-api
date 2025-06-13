package controllers

import (
	"net/http"

	"ios-api/models"
	"ios-api/services"
	"ios-api/utils"

	"github.com/gin-gonic/gin"
)

// AIController AI控制器
type AIController struct {
	AIService *services.AIService
}

// NewAIController 创建新的AI控制器
func NewAIController(aiService *services.AIService) *AIController {
	return &AIController{
		AIService: aiService,
	}
}

// ChatCompletion 通用聊天完成接口
// @Summary 通用AI聊天完成
// @Description 调用AI模型进行聊天对话
// @Tags AI
// @Accept json
// @Produce json
// @Param request body models.ChatRequest true "聊天请求参数"
// @Success 200 {object} utils.Response{data=models.AIResponse} "成功"
// @Failure 400 {object} utils.Response "参数错误"
// @Failure 500 {object} utils.Response "服务器内部错误"
// @Router /api/v1/ai/chat/completions [post]
func (ctrl *AIController) ChatCompletion(c *gin.Context) {
	var request models.ChatRequest

	// 绑定请求参数
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ParamError(c, "请求参数格式错误: "+err.Error())
		return
	}

	// 调用AI服务
	response, err := ctrl.AIService.ChatCompletion(request)
	if err != nil {
		utils.ServerError(c, "AI请求失败: "+err.Error())
		return
	}

	utils.Success(c, "AI聊天完成成功", response)
}

// GenerateTravelPlan 生成旅行计划
// @Summary 生成旅行计划
// @Description 根据用户输入生成详细的旅行计划
// @Tags AI
// @Accept json
// @Produce json
// @Param request body models.TravelPlanRequest true "旅行计划请求参数"
// @Success 200 {object} utils.Response{data=string} "成功"
// @Failure 400 {object} utils.Response "参数错误"
// @Failure 500 {object} utils.Response "服务器内部错误"
// @Router /api/v1/ai/travel/plan [post]
func (ctrl *AIController) GenerateTravelPlan(c *gin.Context) {
	var request models.TravelPlanRequest

	// 绑定请求参数
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ParamError(c, "请求参数格式错误: "+err.Error())
		return
	}

	// 调用AI服务生成旅行计划
	plan, err := ctrl.AIService.GenerateTravelPlan(request)
	if err != nil {
		utils.ServerError(c, "生成旅行计划失败: "+err.Error())
		return
	}

	// 返回旅行计划
	utils.Success(c, "旅行计划生成成功", map[string]interface{}{
		"plan":        plan,
		"destination": request.Destination,
		"start_date":  request.StartDate,
		"end_date":    request.EndDate,
		"budget":      request.Budget,
	})
}

// GetAvailableModels 获取可用的AI模型列表
// @Summary 获取可用AI模型
// @Description 获取系统支持的AI模型列表，优先从GeekAI API动态获取
// @Tags AI
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response{data=[]string} "成功"
// @Failure 500 {object} utils.Response "服务器内部错误"
// @Router /api/v1/ai/models [get]
func (ctrl *AIController) GetAvailableModels(c *gin.Context) {
	models, err := ctrl.AIService.GetAvailableModels()
	if err != nil {
		utils.ServerError(c, "获取模型列表失败: "+err.Error())
		return
	}

	utils.Success(c, "获取模型列表成功", map[string]interface{}{
		"models": models,
		"count":  len(models),
		"source": "dynamic", // 标识这是动态获取的模型列表
	})
}

// GetAIStatus 获取AI服务状态
// @Summary 获取AI服务状态
// @Description 检查AI服务的配置和连接状态
// @Tags AI
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response "成功"
// @Router /api/v1/ai/status [get]
func (ctrl *AIController) GetAIStatus(c *gin.Context) {
	status := map[string]interface{}{
		"service_name": "GeekAI",
		"base_url":     ctrl.AIService.BaseURL,
		"api_key_set":  ctrl.AIService.APIKey != "",
		"timeout":      "60s",
	}

	if ctrl.AIService.APIKey == "" {
		utils.Error(c, http.StatusServiceUnavailable, utils.CodeServerError, "AI服务未配置API密钥")
		return
	}

	utils.Success(c, "AI服务状态正常", status)
}
