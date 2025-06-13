package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"ios-api/config"
	"ios-api/models"
)

// AIService AI服务
type AIService struct {
	APIKey  string
	BaseURL string
	Client  *http.Client
}

// NewAIService 创建新的AI服务实例
func NewAIService(cfg *config.Config) *AIService {
	return &AIService{
		APIKey:  cfg.AIAPIKey,
		BaseURL: cfg.AIBaseURL,
		Client: &http.Client{
			Timeout: 60 * time.Second, // AI请求可能需要更长时间
		},
	}
}

// ChatCompletion 通用聊天完成接口
func (s *AIService) ChatCompletion(request models.ChatRequest) (*models.AIResponse, error) {
	// 验证API密钥
	if s.APIKey == "" {
		return nil, fmt.Errorf("AI API密钥未配置")
	}

	// 构建API请求
	apiRequest := models.AIRequest{
		Model:       request.Model,
		Messages:    request.Messages,
		Stream:      request.Stream,
		Temperature: request.Temperature,
		MaxTokens:   request.MaxTokens,
	}

	// 序列化请求
	jsonData, err := json.Marshal(apiRequest)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %v", err)
	}

	// 创建HTTP请求
	url := fmt.Sprintf("%s/chat/completions", s.BaseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("创建HTTP请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Authorization", "Bearer "+s.APIKey)
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("AI API请求失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	// 解析响应
	var apiResponse models.AIResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	// 检查响应是否有效
	if len(apiResponse.Choices) == 0 {
		return nil, fmt.Errorf("AI API响应中没有选择项")
	}

	return &apiResponse, nil
}

// GenerateTravelPlan 生成旅行计划
func (s *AIService) GenerateTravelPlan(request models.TravelPlanRequest) (string, error) {
	// 构建系统提示词
	systemPrompt := `你是一个专业的旅行规划师。请根据用户提供的信息，生成详细的旅行计划。
旅行计划应该包括：
1. 行程概览
2. 每日详细安排
3. 推荐景点和活动
4. 住宿建议
5. 交通安排
6. 预算分配
7. 注意事项和建议

请以结构化的方式输出，便于阅读和理解。`

	// 构建用户请求
	userPrompt := fmt.Sprintf(`请为我制定一个旅行计划：
目的地：%s
出发日期：%s
返回日期：%s
预算：%s
偏好和特殊要求：%s

请生成详细的旅行计划。`,
		request.Destination,
		request.StartDate,
		request.EndDate,
		request.Budget,
		request.Preferences)

	// 构建聊天请求
	chatRequest := models.ChatRequest{
		Model: "gpt-4o-mini",
		Messages: []models.AIMessage{
			{
				Content: systemPrompt,
				Role:    "system",
			},
			{
				Content: userPrompt,
				Role:    "user",
			},
		},
		Stream: false,
	}

	// 调用聊天完成接口
	response, err := s.ChatCompletion(chatRequest)
	if err != nil {
		return "", err
	}

	return response.Choices[0].Message.Content, nil
}

// GetAvailableModels 获取可用的AI模型列表
func (s *AIService) GetAvailableModels() ([]string, error) {
	// 验证API密钥
	if s.APIKey == "" {
		return nil, fmt.Errorf("AI API密钥未配置")
	}

	// 创建HTTP请求获取模型列表
	url := fmt.Sprintf("%s/models", s.BaseURL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建HTTP请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Authorization", "Bearer "+s.APIKey)
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		// 如果API不支持获取模型列表，返回默认模型列表
		return s.getDefaultModels(), nil
	}

	// 解析响应
	var modelsResponse struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}

	err = json.Unmarshal(body, &modelsResponse)
	if err != nil {
		// 如果解析失败，返回默认模型列表
		return s.getDefaultModels(), nil
	}

	// 提取模型ID
	models := make([]string, 0, len(modelsResponse.Data))
	for _, model := range modelsResponse.Data {
		if model.ID != "" {
			models = append(models, model.ID)
		}
	}

	// 如果没有获取到模型，返回默认列表
	if len(models) == 0 {
		return s.getDefaultModels(), nil
	}

	return models, nil
}

// getDefaultModels 获取默认的模型列表（作为备选）
func (s *AIService) getDefaultModels() []string {
	return []string{
		"gpt-4o-mini",
		"gpt-4o",
		"gpt-4-turbo",
		"gpt-3.5-turbo",
		"claude-3-haiku",
		"claude-3-sonnet",
		"claude-3-opus",
		"claude-3.5-sonnet",
		"gemini-1.5-flash",
		"gemini-1.5-pro",
	}
}
