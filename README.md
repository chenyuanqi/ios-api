# iOS 用户系统 API

这是一个基于 Golang、Gin 框架和 GORM 的 iOS 用户系统 API，提供用户注册、登录、退出登录以及用户信息管理功能。

## 项目架构

项目采用典型的 MVC 架构：
- `config`：配置文件处理
- `controllers`：API 控制器，处理请求和响应
- `middlewares`：中间件，如身份验证、日志等
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
创建 `.env` 文件并设置数据库连接信息：
```
DB_HOST=your_db_host
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_PORT=your_db_port
DB_NAME=your_db_name
JWT_SECRET=your_jwt_secret
APP_PORT=8080
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

## API 文档

详细的 API 文档请参考 [API文档](./docs/api.md)

## 测试

运行单元测试：
```bash
go test ./... -v
```

## 安全说明

本项目使用环境变量来存储敏感信息（如数据库密码和 JWT 密钥），避免将这些信息硬编码在代码中。确保在生产环境中使用强密码和安全的密钥。

## 许可证

[MIT](LICENSE) 