# AI 模块使用示例

## 环境配置

在 `.env` 文件中添加以下配置：

```bash
# AI服务配置
AI_API_KEY=your_geekai_api_key_here
AI_BASE_URL=https://geekai.co/api/v1
```

## API 接口说明

### 1. 获取AI服务状态

**请求：**
```bash
GET /api/v1/ai/status
```

**响应：**
```json
{
  "code": 0,
  "message": "AI服务状态正常",
  "data": {
    "service_name": "GeekAI",
    "base_url": "https://geekai.co/api/v1",
    "api_key_set": true,
    "timeout": "60s"
  }
}
```

### 2. 获取可用模型列表

**说明**：动态获取GeekAI平台支持的AI模型列表，如果API获取失败则返回默认模型列表

**请求：**
```bash
GET /api/v1/ai/models
```

**响应：**
```json
{
  "code": 0,
  "message": "获取模型列表成功",
  "data": {
    "models": [
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
      "grok-2-beta",
      "deepseek-v3"
    ],
    "count": 12,
    "source": "dynamic"
  }
}
```

**注意**：
- `source: "dynamic"` 表示模型列表是从GeekAI API动态获取的
- 如果API获取失败，会自动返回内置的默认模型列表
- 模型列表会根据GeekAI平台的实际可用模型实时更新

### 3. 通用AI聊天接口

**请求：**
```bash
POST /api/v1/ai/chat/completions
Content-Type: application/json

{
  "model": "gpt-4o-mini",
  "messages": [
    {
      "role": "system",
      "content": "你是一个有用的助手。"
    },
    {
      "role": "user",
      "content": "请介绍一下人工智能的发展历史。"
    }
  ],
  "stream": false,
  "temperature": 0.7,
  "max_tokens": 1000
}
```

**响应：**
```json
{
  "code": 0,
  "message": "AI聊天完成成功",
  "data": {
    "id": "chatcmpl-xxx",
    "object": "chat.completion",
    "created": 1703980800,
    "model": "gpt-4o-mini",
    "choices": [
      {
        "index": 0,
        "message": {
          "role": "assistant",
          "content": "人工智能的发展历史可以追溯到20世纪40年代..."
        }
      }
    ],
    "usage": {
      "prompt_tokens": 50,
      "completion_tokens": 200,
      "total_tokens": 250
    }
  }
}
```

### 4. 生成旅行计划

**请求：**
```bash
POST /api/v1/ai/travel/plan
Content-Type: application/json

{
  "destination": "日本东京",
  "start_date": "2024-03-15",
  "end_date": "2024-03-20",
  "budget": "15000元人民币",
  "preferences": "喜欢历史文化，想体验当地美食，对购物也有兴趣"
}
```

**响应：**
```json
{
  "code": 0,
  "message": "旅行计划生成成功",
  "data": {
    "plan": "# 东京5日游详细旅行计划\n\n## 行程概览\n...",
    "destination": "日本东京",
    "start_date": "2024-03-15",
    "end_date": "2024-03-20",
    "budget": "15000元人民币"
  }
}
```

## 使用示例

### JavaScript/前端调用示例

```javascript
// 获取AI服务状态
async function getAIStatus() {
  const response = await fetch('/api/v1/ai/status');
  const data = await response.json();
  console.log('AI服务状态:', data);
}

// 通用AI聊天
async function chatWithAI(messages, model = 'gpt-4o-mini') {
  const response = await fetch('/api/v1/ai/chat/completions', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      model: model,
      messages: messages,
      stream: false,
      temperature: 0.7
    })
  });
  
  const data = await response.json();
  return data.data.choices[0].message.content;
}

// 生成旅行计划
async function generateTravelPlan(destination, startDate, endDate, budget, preferences) {
  const response = await fetch('/api/v1/ai/travel/plan', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      destination,
      start_date: startDate,
      end_date: endDate,
      budget,
      preferences
    })
  });
  
  const data = await response.json();
  return data.data.plan;
}

// 使用示例
(async () => {
  // 检查AI服务状态
  await getAIStatus();
  
  // 进行AI对话
  const messages = [
    { role: 'system', content: '你是一个有用的助手。' },
    { role: 'user', content: '请推荐一些学习编程的方法。' }
  ];
  const aiResponse = await chatWithAI(messages);
  console.log('AI回复:', aiResponse);
  
  // 生成旅行计划
  const travelPlan = await generateTravelPlan(
    '日本东京',
    '2024-03-15',
    '2024-03-20',
    '15000元人民币',
    '喜欢历史文化，想体验当地美食'
  );
  console.log('旅行计划:', travelPlan);
})();
```

### cURL 调用示例

```bash
# 获取AI服务状态
curl -X GET "http://localhost:8080/api/v1/ai/status"

# 获取可用模型
curl -X GET "http://localhost:8080/api/v1/ai/models"

# 通用AI聊天
curl -X POST "http://localhost:8080/api/v1/ai/chat/completions" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "gpt-4o-mini",
    "messages": [
      {
        "role": "system",
        "content": "你是一个有用的助手。"
      },
      {
        "role": "user",
        "content": "请介绍一下Go语言的特点。"
      }
    ],
    "stream": false,
    "temperature": 0.7
  }'

# 生成旅行计划
curl -X POST "http://localhost:8080/api/v1/ai/travel/plan" \
  -H "Content-Type: application/json" \
  -d '{
    "destination": "日本东京",
    "start_date": "2024-03-15",
    "end_date": "2024-03-20",
    "budget": "15000元人民币",
    "preferences": "喜欢历史文化，想体验当地美食，对购物也有兴趣"
  }'
```

## 错误处理

### 常见错误响应

1. **API密钥未配置**
```json
{
  "code": 2000,
  "message": "AI服务未配置API密钥",
  "data": null
}
```

2. **参数错误**
```json
{
  "code": 1001,
  "message": "请求参数格式错误: Key: 'ChatRequest.Model' Error:Field validation for 'Model' failed on the 'required' tag",
  "data": null
}
```

3. **AI API调用失败**
```json
{
  "code": 2000,
  "message": "AI请求失败: AI API请求失败，状态码: 401, 响应: {\"error\":\"Invalid API key\"}",
  "data": null
}
```

## 注意事项

1. **API密钥安全**：请妥善保管您的GeekAI API密钥，不要在代码中硬编码
2. **请求频率**：注意API调用频率限制，避免过于频繁的请求
3. **超时设置**：AI请求可能需要较长时间，当前设置为60秒超时
4. **模型选择**：不同模型有不同的性能和成本特点，请根据需求选择合适的模型
5. **内容过滤**：请确保输入内容符合AI服务提供商的使用政策 