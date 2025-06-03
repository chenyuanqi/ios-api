# 设置管理API使用示例

本文档展示如何使用设置管理API进行键值对的存储和获取。

## 安全机制

设置管理API使用MD5校验来确保安全性。在设置或更新任何key时，都需要提供该key加上配置盐值的MD5值。

### 计算MD5值

**重要**：MD5计算需要加上配置的盐值（通过环境变量 `SETTING_SALT` 设置）

计算公式：`MD5(key + SETTING_SALT)`

**注意**：请联系系统管理员获取当前环境使用的盐值，或查看服务器的环境配置。

可以使用以下方式计算：

**在线工具**：
- 访问 https://www.md5hashgenerator.com/
- 输入 `你的key值` + `配置的盐值`（连接后的字符串）
- 获取MD5哈希值

**命令行**：
```bash
echo -n "your_key_here配置的盐值" | md5sum
```

**JavaScript**：
```javascript
const crypto = require('crypto');
const key = 'app.theme';
const salt = '配置的盐值'; // 从环境配置获取
const md5 = crypto.createHash('md5').update(key + salt).digest('hex');
console.log(md5); // 输出MD5值
```

**Python**：
```python
import hashlib
key = 'app.theme'
salt = '配置的盐值'  # 从环境配置获取
md5 = hashlib.md5((key + salt).encode()).hexdigest()
print(md5)
```

## API使用示例

### 1. 设置应用主题

假设我们要设置应用的主题为深色模式：

**Key**: `app.theme`
**Value**: `dark`
**MD5**: `2ba4eb1c5fbd10524258c4382d7c47e6` (这是"app.theme" + "配置的盐值"的MD5值)

```bash
curl -X PUT http://localhost:8080/api/v1/settings/app.theme \
  -H "Content-Type: application/json" \
  -d '{
    "value": "dark",
    "key_md5": "2ba4eb1c5fbd10524258c4382d7c47e6"
  }'
```

**响应**：
```json
{
  "code": 0,
  "message": "设置保存成功",
  "data": {
    "key": "app.theme",
    "value": "dark",
    "created_at": "2023-03-27T08:00:00Z",
    "updated_at": "2023-03-27T08:00:00Z"
  }
}
```

### 2. 获取应用主题

```bash
curl -X GET http://localhost:8080/api/v1/settings/app.theme
```

**响应**：
```json
{
  "code": 0,
  "message": "获取设置成功",
  "data": {
    "key": "app.theme",
    "value": "dark",
    "created_at": "2023-03-27T08:00:00Z",
    "updated_at": "2023-03-27T08:00:00Z"
  }
}
```

### 3. 设置用户默认头像

**Key**: `user.default_avatar`
**Value**: `https://example.com/default-avatar.png`

首先计算MD5：
```bash
echo -n "user.default_avatar配置的盐值" | md5sum
# 输出: 根据具体盐值而定
```

然后设置：
```bash
curl -X PUT http://localhost:8080/api/v1/settings/user.default_avatar \
  -H "Content-Type: application/json" \
  -d '{
    "value": "https://example.com/default-avatar.png",
    "key_md5": "根据具体盐值计算的MD5值"
  }'
```

### 4. 设置系统配置

**Key**: `system.maintenance_mode`
**Value**: `false`

```bash
# 计算MD5
echo -n "system.maintenance_mode配置的盐值" | md5sum

# 设置值
curl -X PUT http://localhost:8080/api/v1/settings/system.maintenance_mode \
  -H "Content-Type: application/json" \
  -d '{
    "value": "false",
    "key_md5": "根据具体盐值计算的MD5值"
  }'
```

## 常见错误处理

### 1. MD5校验失败

**错误响应**：
```json
{
  "code": 1002,
  "message": "MD5校验失败，无权限操作此设置",
  "data": null
}
```

**解决方案**：确保提供的MD5值是key的正确MD5哈希值。

### 2. Key格式错误

**错误响应**：
```json
{
  "code": 1001,
  "message": "key格式不正确，只允许字母、数字、下划线、点号",
  "data": null
}
```

**解决方案**：确保key只包含字母、数字、下划线(_)和点号(.)。

### 3. 设置不存在

**错误响应**：
```json
{
  "code": 1004,
  "message": "设置不存在",
  "data": null
}
```

**解决方案**：检查key是否正确，或者先创建该设置。

## 最佳实践

### 1. Key命名规范

建议使用分层的命名方式：
- `app.theme` - 应用主题
- `user.default_avatar` - 用户默认头像
- `system.maintenance_mode` - 系统维护模式
- `notification.email_enabled` - 邮件通知开关

### 2. 值的格式

- 简单值：直接存储字符串
- 复杂值：使用JSON格式存储

```bash
# 存储JSON配置
curl -X PUT http://localhost:8080/api/v1/settings/app.config \
  -H "Content-Type: application/json" \
  -d '{
    "value": "{\"theme\":\"dark\",\"language\":\"zh-CN\",\"notifications\":true}",
    "key_md5": "app.config的MD5值"
  }'
```

### 3. 安全考虑

- 不要在客户端硬编码敏感的key
- 定期更换重要设置的key
- 对敏感值进行加密后再存储

## 前端集成示例

### JavaScript/TypeScript

```javascript
class SettingsAPI {
  constructor(baseURL) {
    this.baseURL = baseURL;
  }

  // 计算MD5（包含固定盐值）
  calculateMD5(str) {
    // 使用crypto-js库
    const salt = '配置的盐值'; // 从环境配置获取
    return CryptoJS.MD5(str + salt).toString();
  }

  // 获取设置
  async getSetting(key) {
    const response = await fetch(`${this.baseURL}/api/v1/settings/${key}`);
    return response.json();
  }

  // 设置值
  async setSetting(key, value) {
    const keyMD5 = this.calculateMD5(key);
    const response = await fetch(`${this.baseURL}/api/v1/settings/${key}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        value: value,
        key_md5: keyMD5
      })
    });
    return response.json();
  }
}

// 使用示例
const settings = new SettingsAPI('http://localhost:8080');

// 设置主题
settings.setSetting('app.theme', 'dark').then(result => {
  console.log('设置成功:', result);
});

// 获取主题
settings.getSetting('app.theme').then(result => {
  console.log('当前主题:', result.data.value);
});
```

### Swift (iOS)

```swift
import Foundation
import CryptoKit

class SettingsAPI {
    let baseURL: String
    
    init(baseURL: String) {
        self.baseURL = baseURL
    }
    
    // 计算MD5（包含固定盐值）
    func calculateMD5(_ string: String) -> String {
        let salt = "配置的盐值" // 从环境配置获取
        let dataToHash = string + salt
        let digest = Insecure.MD5.hash(data: dataToHash.data(using: .utf8) ?? Data())
        return digest.map { String(format: "%02hhx", $0) }.joined()
    }
    
    // 获取设置
    func getSetting(key: String, completion: @escaping (Result<SettingResponse, Error>) -> Void) {
        guard let url = URL(string: "\(baseURL)/api/v1/settings/\(key)") else { return }
        
        URLSession.shared.dataTask(with: url) { data, response, error in
            if let error = error {
                completion(.failure(error))
                return
            }
            
            guard let data = data else { return }
            
            do {
                let result = try JSONDecoder().decode(SettingResponse.self, from: data)
                completion(.success(result))
            } catch {
                completion(.failure(error))
            }
        }.resume()
    }
    
    // 设置值
    func setSetting(key: String, value: String, completion: @escaping (Result<SettingResponse, Error>) -> Void) {
        guard let url = URL(string: "\(baseURL)/api/v1/settings/\(key)") else { return }
        
        let keyMD5 = calculateMD5(key)
        let requestBody = SetSettingRequest(value: value, keyMD5: keyMD5)
        
        var request = URLRequest(url: url)
        request.httpMethod = "PUT"
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
        
        do {
            request.httpBody = try JSONEncoder().encode(requestBody)
        } catch {
            completion(.failure(error))
            return
        }
        
        URLSession.shared.dataTask(with: request) { data, response, error in
            if let error = error {
                completion(.failure(error))
                return
            }
            
            guard let data = data else { return }
            
            do {
                let result = try JSONDecoder().decode(SettingResponse.self, from: data)
                completion(.success(result))
            } catch {
                completion(.failure(error))
            }
        }.resume()
    }
}

// 数据模型
struct SetSettingRequest: Codable {
    let value: String
    let keyMD5: String
    
    enum CodingKeys: String, CodingKey {
        case value
        case keyMD5 = "key_md5"
    }
}

struct SettingResponse: Codable {
    let code: Int
    let message: String
    let data: Setting?
}

struct Setting: Codable {
    let key: String
    let value: String
    let createdAt: String
    let updatedAt: String
    
    enum CodingKeys: String, CodingKey {
        case key, value
        case createdAt = "created_at"
        case updatedAt = "updated_at"
    }
}
```

这个设置管理系统提供了简单而安全的键值对存储功能，适合存储应用配置、用户偏好设置等数据。
