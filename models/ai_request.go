package models

// AIMessage AI消息结构
type AIMessage struct {
	Content string `json:"content" binding:"required"`
	Role    string `json:"role" binding:"required,oneof=system user assistant"`
}

// AIRequest AI请求结构
type AIRequest struct {
	Model       string      `json:"model" binding:"required"`
	Messages    []AIMessage `json:"messages" binding:"required,min=1"`
	Stream      bool        `json:"stream"`
	Temperature *float64    `json:"temperature,omitempty"`
	MaxTokens   *int        `json:"max_tokens,omitempty"`
}

// AIChoice AI响应选择结构
type AIChoice struct {
	Message AIMessage `json:"message"`
	Index   int       `json:"index"`
}

// AIResponse AI响应结构
type AIResponse struct {
	ID      string     `json:"id"`
	Object  string     `json:"object"`
	Created int64      `json:"created"`
	Model   string     `json:"model"`
	Choices []AIChoice `json:"choices"`
	Usage   struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// TravelPlanRequest 旅行计划请求结构
type TravelPlanRequest struct {
	Destination string `json:"destination" binding:"required"`
	StartDate   string `json:"start_date" binding:"required"`
	EndDate     string `json:"end_date" binding:"required"`
	Budget      string `json:"budget" binding:"required"`
	Preferences string `json:"preferences"`
}

// ChatRequest 通用聊天请求结构
type ChatRequest struct {
	Model       string      `json:"model" binding:"required"`
	Messages    []AIMessage `json:"messages" binding:"required,min=1"`
	Stream      bool        `json:"stream"`
	Temperature *float64    `json:"temperature,omitempty"`
	MaxTokens   *int        `json:"max_tokens,omitempty"`
} 