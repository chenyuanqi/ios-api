# iOS 用户系统 API

这是一个基于 Golang、Gin 框架和 GORM 的 iOS 用户系统 API，提供用户注册、登录、退出登录以及用户信息管理功能。

## 项目架构

项目采用典型的 MVC 架构：
- `config`：配置文件处理
- `controllers`：API 控制器，处理请求和响应
- `middlewares`：中间件，如身份验证、日志、CORS等
- `models`：数据模型，对应数据库表结构
- `repositories`：数据访问层，封装数据库操作
- `routes`：路由定义，API 路径配置
- `services`：业务逻辑层，处理业务规则
- `utils`：工具函数，如加密、JWT 等

## 功能列表

1. 用户注册
   - 支持邮箱和密码注册
   - 支持第三方登录（微信、苹果）

2. 用户登录
   - 支持邮箱和密码登录
   - 支持第三方登录（微信、苹果）

3. 用户退出登录

4. 用户信息获取
   - 获取头像、昵称、注册日期、个性签名

5. 用户信息修改
   - 修改头像、昵称、个性签名

6. 设置管理
   - 获取指定key的设置值
   - 设置/更新指定key的值（需要可配置盐值的MD5校验）

7. **跨域访问支持（CORS）**
   - 支持所有来源的跨域请求（开发环境）
   - 支持常用的HTTP方法（GET、POST、PUT、DELETE、OPTIONS）
   - 支持认证头部（Authorization）
   - 自动处理预检请求（OPTIONS）
   - 可配置允许的域名列表（生产环境推荐）

## 安装和运行

### 前置条件
- Go 1.16+
- MySQL 5.7+

### 安装步骤
1. 克隆仓库
```bash
git clone [仓库地址]
cd ios-api
```

2. 安装依赖
```bash
go mod download
```

3. 配置环境变量
从示例文件创建 .env 文件：
```bash
cp .env.example .env
```
然后编辑 .env 文件，设置您的实际配置：
```
# 主数据库配置（用户系统）
DB_HOST=your_db_host
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_PORT=your_db_port
DB_NAME=your_db_name

# 通用数据库配置（设置系统）
GENERAL_DB_HOST=your_general_db_host
GENERAL_DB_USER=your_general_db_user
GENERAL_DB_PASSWORD=your_general_db_password
GENERAL_DB_PORT=your_general_db_port
GENERAL_DB_NAME=yuanqi_general

JWT_SECRET=your_jwt_secret
APP_PORT=8080

# 设置管理配置
SETTING_SALT=chenyuanqi2025CYQ

# 微信登录配置
WECHAT_APP_ID=your_wechat_app_id
WECHAT_APP_SECRET=your_wechat_app_secret

# 苹果登录配置
APPLE_TEAM_ID=your_apple_team_id
APPLE_KEY_ID=your_apple_key_id
APPLE_PRIVATE_KEY=path_to_your_private_key_or_key_string
APPLE_BUNDLE_ID=your_app_bundle_id
```

详细的环境配置说明请参考 [环境配置文档](./docs/environment.md)。

4. 初始化数据库
```bash
mysql -u root -p < migrate.sql
```

5. 运行项目
```bash
go run main.go
```

## 第三方登录配置

### 微信登录配置

要启用微信登录功能，需要完成以下步骤：

1. 在[微信开放平台](https://open.weixin.qq.com/)注册开发者账号
2. 创建移动应用并获取 AppID 和 AppSecret
3. 在`.env`文件中配置以下参数：
   ```
   WECHAT_APP_ID=your_wechat_app_id
   WECHAT_APP_SECRET=your_wechat_app_secret
   ```
4. 在 iOS 客户端集成微信 SDK：
   - 在 [微信开放平台](https://developers.weixin.qq.com/doc/oplatform/Mobile_App/Access_Guide/iOS.html) 下载 iOS SDK
   - 按照文档配置 URL Scheme 为 `wx` + AppID
   - 在 `Info.plist` 中添加 LSApplicationQueriesSchemes 字段，包含 `weixin` 值
   - 实现微信授权并获取授权码 (code)
   - 将授权码发送到后端 API 获取用户信息

### 苹果登录配置

要启用苹果登录功能，需要完成以下步骤：

1. 在 [Apple Developer](https://developer.apple.com/) 注册开发者账号
2. 配置 "Sign in with Apple" 功能：
   - 在 Certificates, Identifiers & Profiles 中启用 "Sign in with Apple" capability
   - 创建 Services ID 并配置域名和重定向 URL
   - 获取 Team ID、Key ID 和私钥文件
3. 在`.env`文件中配置以下参数：
   ```
   APPLE_TEAM_ID=your_apple_team_id
   APPLE_KEY_ID=your_apple_key_id
   APPLE_PRIVATE_KEY=path_to_your_private_key_or_key_string
   APPLE_BUNDLE_ID=your_app_bundle_id
   ```
4. 在 iOS 客户端集成 Apple 登录：
   - 启用 "Sign in with Apple" capability
   - 使用 AuthenticationServices 框架实现苹果登录
   - 获取用户标识符和授权码
   - 将授权信息发送到后端 API 验证身份

### API 调用示例

使用微信登录：
```json
POST /api/v1/oauth/login
{
  "provider": "wechat",
  "provider_user_id": "微信用户唯一标识",
  "nickname": "微信昵称",
  "avatar": "头像URL"
}
```

使用苹果登录：
```json
POST /api/v1/oauth/login
{
  "provider": "apple",
  "provider_user_id": "苹果用户唯一标识",
  "nickname": "用户昵称",
  "avatar": "头像URL"
}
```

## API 文档

详细的 API 文档请参考 [API文档](./docs/api.md)

## 跨域访问配置

本项目已集成完整的CORS跨域支持，详细配置说明请参考 [CORS配置文档](./docs/cors.md)

## 部署指南

如果您需要将此项目部署到生产环境，请参考 [部署指南](./deployment.md)，其中包含了在 Ubuntu 16.04 系统上部署的详细步骤。

## 测试

运行单元测试：
```bash
go test ./... -v
```

## 安全说明

本项目使用环境变量来存储敏感信息（如数据库密码和 JWT 密钥），避免将这些信息硬编码在代码中。确保在生产环境中使用强密码和安全的密钥。

## 项目状态

目前项目已经完成了以下工作：

1. **基础架构**：完成了基于Gin和GORM的项目结构搭建
2. **用户系统**：实现了用户注册、登录、退出、信息查询和修改功能
3. **第三方登录**：支持微信和苹果第三方授权登录
4. **设置管理**：支持键值对设置的获取和更新，使用可配置盐值的MD5校验保证安全性
5. **API统一规范**：所有API返回统一的响应格式，便于前端处理
6. **配置管理**：使用环境变量管理敏感配置，支持多数据库连接
7. **部署文档**：完整的部署指南，支持在Ubuntu 16.04上部署

项目代码结构清晰，使用了MVC架构模式，具有良好的可测试性和可维护性。已经完成了基本的单元测试，确保代码的健壮性。

## 后续工作

以下是未来可能的改进方向：

1. 增加更多第三方登录支持（如Google、Facebook等）
2. 实现用户权限管理系统
3. 添加用户密码重置功能
4. 实现完整的API访问日志
5. 优化数据库查询性能
6. 添加更多的单元测试和集成测试
7. 实现API限流和防刷机制

## 贡献指南

欢迎贡献代码，请按以下步骤：

1. Fork 本仓库
2. 创建你的特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交你的改动 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 开启一个 Pull Request

## 许可证

本项目使用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情 