# CORS 跨域访问配置

本项目已经集成了完整的CORS（Cross-Origin Resource Sharing）跨域访问支持，允许前端应用从不同的域名访问API。

## 配置说明

### 当前CORS配置

项目使用了自定义的CORS中间件（`middlewares/cors.go`），配置如下：

- **Access-Control-Allow-Origin**: `*` - 允许所有来源的跨域请求
- **Access-Control-Allow-Methods**: `GET, POST, PUT, DELETE, OPTIONS` - 允许的HTTP方法
- **Access-Control-Allow-Headers**: `Origin, Content-Type, Accept, Authorization, User-Agent, Content-Length, X-Requested-With` - 允许的请求头
- **Access-Control-Expose-Headers**: `Content-Length, Content-Type` - 暴露给客户端的响应头
- **Access-Control-Allow-Credentials**: `true` - 允许发送凭据（如cookies、Authorization头）

### 预检请求处理

中间件自动处理OPTIONS预检请求，返回204状态码，确保复杂跨域请求能够正常工作。

## 使用方式

### 前端JavaScript调用示例

```javascript
// 使用fetch发起跨域请求
fetch('http://your-api-domain.com/api/v1/user', {
  method: 'GET',
  headers: {
    'Content-Type': 'application/json',
    'Authorization': 'Bearer your-jwt-token'
  },
  credentials: 'include' // 如果需要发送cookies
})
.then(response => response.json())
.then(data => console.log(data));

// 使用axios发起跨域请求
axios.defaults.withCredentials = true;
axios.get('http://your-api-domain.com/api/v1/user', {
  headers: {
    'Authorization': 'Bearer your-jwt-token'
  }
})
.then(response => console.log(response.data));
```

### 前端Vue.js示例

```javascript
// main.js
import axios from 'axios'

// 配置axios默认值
axios.defaults.baseURL = 'http://your-api-domain.com/api/v1'
axios.defaults.withCredentials = true

// 请求拦截器，自动添加Authorization头
axios.interceptors.request.use(config => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

Vue.prototype.$http = axios
```

### 前端React示例

```javascript
// api.js
const API_BASE_URL = 'http://your-api-domain.com/api/v1';

class API {
  static async request(endpoint, options = {}) {
    const url = `${API_BASE_URL}${endpoint}`;
    
    const config = {
      headers: {
        'Content-Type': 'application/json',
        ...options.headers
      },
      credentials: 'include',
      ...options
    };

    // 添加认证头
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }

    const response = await fetch(url, config);
    return response.json();
  }

  static async get(endpoint) {
    return this.request(endpoint, { method: 'GET' });
  }

  static async post(endpoint, data) {
    return this.request(endpoint, {
      method: 'POST',
      body: JSON.stringify(data)
    });
  }
}

// 使用示例
API.get('/user').then(data => console.log(data));
```

## 生产环境配置建议

### 限制允许的域名

在生产环境中，建议限制允许的来源域名以提高安全性。修改 `middlewares/cors.go`：

```go
// 生产环境配置示例
func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        origin := c.Request.Header.Get("Origin")
        
        // 允许的域名列表
        allowedOrigins := []string{
            "https://your-frontend-domain.com",
            "https://your-mobile-app-domain.com",
            "https://your-admin-panel.com",
        }
        
        // 检查来源是否在允许列表中
        allowed := false
        for _, allowedOrigin := range allowedOrigins {
            if origin == allowedOrigin {
                allowed = true
                break
            }
        }
        
        if allowed {
            c.Header("Access-Control-Allow-Origin", origin)
        }
        
        // 其他CORS头设置...
    }
}
```

### 环境变量配置

可以通过环境变量来配置允许的域名：

```bash
# .env
CORS_ALLOWED_ORIGINS=https://example.com,https://app.example.com
```

然后在代码中读取：

```go
allowedOrigins := strings.Split(os.Getenv("CORS_ALLOWED_ORIGINS"), ",")
```

## 使用测试页面验证CORS

项目提供了一个完整的CORS测试页面，位于 `examples/cors_test.html`，可以帮助您验证CORS配置是否正常工作。

### 使用方法

1. **启动API服务器**：
   ```bash
   go run main.go
   ```

2. **打开测试页面**：
   - 使用浏览器直接打开 `examples/cors_test.html` 文件
   - 或者通过本地HTTP服务器访问（推荐）：
     ```bash
     # 使用Python启动简单HTTP服务器
     cd examples
     python3 -m http.server 3000
     # 然后访问 http://localhost:3000/cors_test.html
     ```

3. **进行测试**：
   - 页面提供了5种不同的测试场景
   - 每个测试都会显示详细的请求和响应信息
   - 绿色表示成功，红色表示失败

### 测试场景说明

1. **简单GET请求**：测试基本的跨域GET请求
2. **POST请求**：测试发送JSON数据的POST请求
3. **带认证的请求**：测试携带Authorization头的请求
4. **OPTIONS预检请求**：测试浏览器的预检请求处理
5. **PUT请求**：测试更新数据的PUT请求

### 预期结果

正确配置的CORS应该显示：
- 所有请求的响应头中包含 `Access-Control-Allow-Origin: *`
- OPTIONS请求返回204状态码
- 所有请求都能正常处理，不会被浏览器的CORS策略阻止

### 常见错误

如果测试失败，可能的原因包括：
- API服务器未启动
- CORS中间件未正确配置
- 浏览器缓存了旧的CORS策略（尝试硬刷新）
- 防火墙或网络代理阻止了请求

## 常见问题

### 1. 预检请求失败

**问题**: 浏览器发起OPTIONS请求时返回错误
**解决**: 确保服务器正确处理OPTIONS请求，返回适当的CORS头

### 2. 凭据发送失败

**问题**: 无法发送cookies或Authorization头
**解决**: 
- 前端设置 `credentials: 'include'` 或 `withCredentials: true`
- 后端设置 `Access-Control-Allow-Credentials: true`
- 当允许凭据时，`Access-Control-Allow-Origin` 不能为 `*`

### 3. 自定义头部被拒绝

**问题**: 发送自定义请求头时被CORS阻止
**解决**: 将自定义头部添加到 `Access-Control-Allow-Headers` 中

## 测试CORS配置

### 使用curl测试

```bash
# 测试预检请求
curl -X OPTIONS \
  -H "Origin: http://localhost:3000" \
  -H "Access-Control-Request-Method: POST" \
  -H "Access-Control-Request-Headers: Content-Type, Authorization" \
  http://your-api-domain.com/api/v1/user

# 测试实际请求
curl -X GET \
  -H "Origin: http://localhost:3000" \
  -H "Authorization: Bearer your-token" \
  http://your-api-domain.com/api/v1/user
```

### 使用浏览器开发者工具

1. 打开浏览器开发者工具
2. 在Network选项卡中观察请求
3. 检查响应头中是否包含正确的CORS头部
4. 确认预检请求（OPTIONS）得到正确响应

## 安全考虑

1. **不要在生产环境使用 `*` 通配符**：限制具体的域名
2. **谨慎设置 `Access-Control-Allow-Credentials`**：只在需要时启用
3. **定期审查允许的域名列表**：移除不再使用的域名
4. **使用HTTPS**：在生产环境中始终使用HTTPS
5. **监控跨域请求**：记录和监控跨域访问模式

通过以上配置，您的API已经支持完整的跨域访问功能，可以安全地为各种前端应用提供服务。 