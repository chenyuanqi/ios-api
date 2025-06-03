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
# 主数据库配置（用户系统）
DB_HOST=your_db_host        # 数据库主机地址
DB_USER=your_db_user        # 数据库用户名
DB_PASSWORD=your_db_password # 数据库密码
DB_PORT=3306                # 数据库端口
DB_NAME=your_db_name        # 数据库名称

# 通用数据库配置（设置系统）
GENERAL_DB_HOST=your_general_db_host        # 通用数据库主机地址
GENERAL_DB_USER=your_general_db_user        # 通用数据库用户名
GENERAL_DB_PASSWORD=your_general_db_password # 通用数据库密码
GENERAL_DB_PORT=3306                        # 通用数据库端口
GENERAL_DB_NAME=yuanqi_general              # 通用数据库名称

# JWT配置
JWT_SECRET=your_jwt_secret  # JWT 密钥，用于生成和验证用户令牌

# 应用配置
APP_PORT=8080               # 应用监听端口

# 设置管理配置
SETTING_SALT=your_setting_salt  # 设置管理MD5校验盐值，增强安全性

# 缓存配置
CACHE_DIR=./cache              # LevelDB缓存目录，用于设置数据缓存

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
   
   # 设置管理配置
   SETTING_SALT=dev_setting_salt_2024
   
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
   
   # 测试环境设置管理配置
   SETTING_SALT=test_setting_salt_2024
   
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
   
   # 生产环境设置管理配置（使用强盐值）
   SETTING_SALT=prod_strong_salt_abc123xyz789
   
   # 生产环境缓存配置
   CACHE_DIR=/var/cache/ios-api
   
   # 生产环境微信配置
   WECHAT_APP_ID=your_prod_wechat_app_id
   WECHAT_APP_SECRET=your_prod_wechat_app_secret
   
   # 生产环境苹果配置
   APPLE_TEAM_ID=your_apple_team_id
   APPLE_KEY_ID=your_apple_key_id
   APPLE_PRIVATE_KEY=/secure/path/to/apple_key.p8
   APPLE_BUNDLE_ID=com.example.app
   ```

## 缓存配置说明

### LevelDB缓存

项目使用LevelDB作为设置数据的本地缓存，以提高读取性能：

- **CACHE_DIR**: 缓存文件存储目录
  - 开发环境：`./cache`（项目根目录下的cache文件夹）
  - 生产环境：推荐使用绝对路径，如 `/var/cache/ios-api`

### 缓存功能

1. **自动缓存**：首次读取设置时自动缓存到LevelDB
2. **缓存失效**：更新设置时自动删除对应缓存
3. **缓存管理**：提供API接口进行缓存管理

### 缓存管理API

- `GET /api/v1/settings/cache/stats` - 获取缓存统计信息
- `DELETE /api/v1/settings/:key/cache` - 清除指定key的缓存
- `DELETE /api/v1/settings/cache` - 清除所有设置缓存

### 注意事项

1. **磁盘空间**：确保缓存目录有足够的磁盘空间
2. **权限**：确保应用有读写缓存目录的权限
3. **备份**：缓存数据会在重启后保持，如需清理可删除缓存目录
4. **性能**：缓存可显著提高设置读取性能，特别是频繁访问的设置

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