# 环境配置

本项目使用 `.env` 文件来存储环境配置，特别是数据库连接信息。这样可以避免将敏感信息（如数据库密码）直接硬编码在代码中，提高了安全性。

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
   ```

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