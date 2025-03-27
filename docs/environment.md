# 环境配置

本项目使用 `.env` 文件来存储环境配置，特别是数据库连接信息。这样可以避免将敏感信息（如数据库密码）直接硬编码在代码中，提高了安全性。

## 快速设置

项目根目录提供了 `.env.example` 示例文件，您可以通过以下命令快速创建您自己的环境配置：

```bash
# 复制示例文件
cp .env.example .env

# 编辑文件填入您的实际配置
vim .env
```

## 配置项

在项目根目录创建 `.env` 文件，包含以下配置项：

```
# 数据库配置
DB_HOST=your_db_host        # 数据库主机地址
DB_USER=your_db_user        # 数据库用户名
DB_PASSWORD=your_db_password # 数据库密码
DB_PORT=3306                # 数据库端口
DB_NAME=your_db_name        # 数据库名称

# JWT配置
JWT_SECRET=your_jwt_secret  # JWT 密钥，用于生成和验证用户令牌

# 应用配置
APP_PORT=8080               # 应用监听端口

# 微信登录配置
WECHAT_APP_ID=your_wechat_app_id           # 微信开放平台 AppID
WECHAT_APP_SECRET=your_wechat_app_secret   # 微信开放平台 AppSecret

# 苹果登录配置
APPLE_TEAM_ID=your_apple_team_id           # 苹果开发者 Team ID
APPLE_KEY_ID=your_apple_key_id             # 苹果私钥 ID
APPLE_PRIVATE_KEY=path_or_content          # 苹果私钥文件路径或内容
APPLE_BUNDLE_ID=your_app_bundle_id         # 应用的 Bundle ID
```

## 注意事项

1. **安全性**：`.env` 文件包含敏感信息，不应该被提交到版本控制系统中。确保在 `.gitignore` 中包含了 `.env`。

2. **开发环境**：在开发环境中，可以使用以下示例配置：
   ```
   DB_HOST=localhost
   DB_USER=root
   DB_PASSWORD=your_password
   DB_PORT=3306
   DB_NAME=yuanqi_ios
   JWT_SECRET=dev_jwt_secret
   APP_PORT=8080
   
   # 开发环境微信配置
   WECHAT_APP_ID=your_dev_wechat_app_id
   WECHAT_APP_SECRET=your_dev_wechat_app_secret
   
   # 开发环境苹果配置
   APPLE_TEAM_ID=your_apple_team_id
   APPLE_KEY_ID=your_apple_key_id
   APPLE_PRIVATE_KEY=./certs/apple_key.p8
   APPLE_BUNDLE_ID=com.example.app.development
   ```

3. **测试环境**：测试环境应使用专门的测试数据库：
   ```
   DB_HOST=localhost
   DB_USER=root
   DB_PASSWORD=your_password
   DB_PORT=3306
   DB_NAME=yuanqi_ios_test
   JWT_SECRET=test_jwt_secret
   APP_PORT=8080
   
   # 测试微信配置（可使用开发账号）
   WECHAT_APP_ID=your_test_wechat_app_id
   WECHAT_APP_SECRET=your_test_wechat_app_secret
   
   # 测试苹果配置
   APPLE_TEAM_ID=your_apple_team_id
   APPLE_KEY_ID=your_apple_key_id
   APPLE_PRIVATE_KEY=./certs/apple_key.p8
   APPLE_BUNDLE_ID=com.example.app.testing
   ```

4. **生产环境**：在生产环境中，请使用强密码和安全的 JWT 密钥：
   ```
   DB_HOST=your_production_host
   DB_USER=production_user
   DB_PASSWORD=strong_password
   DB_PORT=3306
   DB_NAME=yuanqi_ios_prod
   JWT_SECRET=long_random_string
   APP_PORT=80
   
   # 生产环境微信配置
   WECHAT_APP_ID=your_prod_wechat_app_id
   WECHAT_APP_SECRET=your_prod_wechat_app_secret
   
   # 生产环境苹果配置
   APPLE_TEAM_ID=your_apple_team_id
   APPLE_KEY_ID=your_apple_key_id
   APPLE_PRIVATE_KEY=/secure/path/to/apple_key.p8
   APPLE_BUNDLE_ID=com.example.app
   ```

## 第三方登录配置说明

### 微信登录

微信登录需要在[微信开放平台](https://open.weixin.qq.com/)注册开发者账号，并创建移动应用。配置项说明：

- **WECHAT_APP_ID**: 在微信开放平台创建应用后获得的 AppID
- **WECHAT_APP_SECRET**: 对应的 AppSecret，用于服务端接口调用

### 苹果登录

苹果登录需要在[Apple Developer](https://developer.apple.com/)账号中配置 "Sign in with Apple" 功能。配置项说明：

- **APPLE_TEAM_ID**: 苹果开发者账号的 Team ID
- **APPLE_KEY_ID**: 用于签名 JWT 令牌的私钥 ID
- **APPLE_PRIVATE_KEY**: 私钥文件的路径或内容（P8 格式）
- **APPLE_BUNDLE_ID**: 应用的 Bundle Identifier

## 如何加载配置

项目使用 `github.com/joho/godotenv` 库从 `.env` 文件加载配置。配置逻辑在 `config/config.go` 文件中实现。

在代码中，可以通过以下方式获取和使用配置：

```go
import "ios-api/config"

func main() {
    // 加载配置
    cfg, err := config.LoadConfig()
    if err != nil {
        // 处理错误
    }
    
    // 使用配置
    dbHost := cfg.DBHost
    dbPassword := cfg.DBPassword
    // ...
}
``` 