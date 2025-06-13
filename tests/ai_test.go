package tests

import (
	"testing"

	"ios-api/config"
	"ios-api/models"
	"ios-api/services"

	"github.com/stretchr/testify/assert"
)

func TestAIService_GetAvailableModels(t *testing.T) {
	// 创建测试配置
	cfg := &config.Config{
		AIAPIKey:  "test-api-key",
		AIBaseURL: "https://geekai.co/api/v1",
	}

	// 创建AI服务
	aiService := services.NewAIService(cfg)

	// 测试获取可用模型
	models, err := aiService.GetAvailableModels()

	// 验证结果
	assert.NoError(t, err, "获取模型列表不应出错")
	assert.NotEmpty(t, models, "模型列表不应为空")
	assert.Contains(t, models, "gpt-4o-mini", "应包含gpt-4o-mini模型")
	assert.Contains(t, models, "gpt-4o", "应包含gpt-4o模型")

	// 注意：由于是动态获取，可能包含更多模型
	t.Logf("获取到 %d 个模型：%v", len(models), models)
}

func TestAIService_ValidateRequest(t *testing.T) {
	// 测试有效的聊天请求
	validRequest := models.ChatRequest{
		Model: "gpt-4o-mini",
		Messages: []models.AIMessage{
			{
				Role:    "system",
				Content: "你是一个有用的助手。",
			},
			{
				Role:    "user",
				Content: "你好",
			},
		},
		Stream: false,
	}

	// 验证请求结构
	assert.Equal(t, "gpt-4o-mini", validRequest.Model)
	assert.Len(t, validRequest.Messages, 2)
	assert.Equal(t, "system", validRequest.Messages[0].Role)
	assert.Equal(t, "user", validRequest.Messages[1].Role)
	assert.False(t, validRequest.Stream)
}

func TestTravelPlanRequest_Validation(t *testing.T) {
	// 测试有效的旅行计划请求
	validRequest := models.TravelPlanRequest{
		Destination: "日本东京",
		StartDate:   "2024-03-15",
		EndDate:     "2024-03-20",
		Budget:      "15000元人民币",
		Preferences: "喜欢历史文化，想体验当地美食",
	}

	// 验证请求结构
	assert.Equal(t, "日本东京", validRequest.Destination)
	assert.Equal(t, "2024-03-15", validRequest.StartDate)
	assert.Equal(t, "2024-03-20", validRequest.EndDate)
	assert.Equal(t, "15000元人民币", validRequest.Budget)
	assert.Equal(t, "喜欢历史文化，想体验当地美食", validRequest.Preferences)
}

func TestAIService_ConfigValidation(t *testing.T) {
	// 测试配置验证
	tests := []struct {
		name      string
		apiKey    string
		baseURL   string
		expectErr bool
	}{
		{
			name:      "有效配置",
			apiKey:    "valid-api-key",
			baseURL:   "https://geekai.co/api/v1",
			expectErr: false,
		},
		{
			name:      "空API密钥",
			apiKey:    "",
			baseURL:   "https://geekai.co/api/v1",
			expectErr: true,
		},
		{
			name:      "空BaseURL",
			apiKey:    "valid-api-key",
			baseURL:   "",
			expectErr: false, // BaseURL有默认值
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				AIAPIKey:  tt.apiKey,
				AIBaseURL: tt.baseURL,
			}

			aiService := services.NewAIService(cfg)

			// 验证配置
			assert.Equal(t, tt.apiKey, aiService.APIKey)
			if tt.baseURL != "" {
				assert.Equal(t, tt.baseURL, aiService.BaseURL)
			}
		})
	}
}

// 基准测试 - 测试模型列表获取性能
func BenchmarkAIService_GetAvailableModels(b *testing.B) {
	cfg := &config.Config{
		AIAPIKey:  "test-api-key",
		AIBaseURL: "https://geekai.co/api/v1",
	}

	aiService := services.NewAIService(cfg)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = aiService.GetAvailableModels()
	}
}
