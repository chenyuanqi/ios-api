# iOS 用户系统 API 文档

## 基本信息

- 基础 URL: `http://your-domain.com/api/v1`
- 所有包含请求体的请求必须指定内容类型为 `Content-Type: application/json`
- 所有需要认证的请求都需要在请求头中包含 `Authorization: Bearer {token}`

## 响应格式

所有 API 请求都将返回统一的响应格式：

```json
{
  "code": 0,        // 响应码，0表示成功，非0表示错误
  "message": "成功", // 响应消息
  "data": {}        // 响应数据，具体内容根据接口不同而不同
}
```

### 响应码说明
- `0`: 成功
- `1001`: 参数错误
- `1002`: 未授权
- `1004`: 资源不存在
- `1009`: 资源冲突（如邮箱已注册）
- `2000`: 服务器内部错误

### HTTP 状态码

API 同时使用标准的 HTTP 状态码表示请求结果：
- 200: 请求成功
- 201: 创建成功
- 400: 请求参数错误
- 401: 未授权或授权失败
- 404: 资源不存在
- 409: 冲突（例如邮箱已注册）
- 500: 服务器内部错误

## API 列表

### 1. 用户注册

**POST /register**

通过邮箱和密码注册新用户。

请求参数：

```json
{
  "email": "user@example.com",
  "password": "password123",
  "nickname": "用户昵称"
}
```

成功响应 (201)：

```json
{
  "code": 0,
  "message": "注册成功",
  "data": {
    "user": {
      "id": 1,
      "email": "user@example.com",
      "nickname": "用户昵称",
      "avatar": null,
      "signature": null,
      "created_at": "2023-03-27T08:00:00Z",
      "updated_at": "2023-03-27T08:00:00Z"
    },
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

### 2. 用户登录

**POST /login**

通过邮箱和密码登录。

请求参数：

```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

成功响应 (200)：

```json
{
  "code": 0,
  "message": "登录成功",
  "data": {
    "user": {
      "id": 1,
      "email": "user@example.com",
      "nickname": "用户昵称",
      "avatar": "https://example.com/avatar.jpg",
      "signature": "个性签名",
      "created_at": "2023-03-27T08:00:00Z",
      "updated_at": "2023-03-27T08:00:00Z"
    },
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

### 3. 第三方登录

**POST /oauth/login**

通过第三方授权登录（微信、苹果）。

请求参数：

```json
{
  "provider": "wechat", // 或 "apple"
  "provider_user_id": "第三方用户ID",
  "nickname": "用户昵称",
  "avatar": "头像URL",
  "email": "用户邮箱"  // 可选
}
```

成功响应 (200)：

```json
{
  "code": 0,
  "message": "登录成功",
  "data": {
    "user": {
      "id": 1,
      "email": null,
      "nickname": "用户昵称",
      "avatar": "https://example.com/avatar.jpg",
      "signature": null,
      "created_at": "2023-03-27T08:00:00Z",
      "updated_at": "2023-03-27T08:00:00Z"
    },
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

### 4. 微信授权 URL

**POST /oauth/wechat/auth**

获取微信授权链接。

请求参数：

```json
{
  "redirect_uri": "https://your-domain.com/callback",
  "state": "随机字符串，防止CSRF攻击"  // 可选
}
```

成功响应 (200)：

```json
{
  "code": 0,
  "message": "获取微信授权链接成功",
  "data": {
    "auth_url": "https://open.weixin.qq.com/connect/oauth2/authorize?appid=..."
  }
}
```

### 5. 微信授权回调

**GET /oauth/wechat/callback?code=授权码&state=状态**

处理微信授权回调，通常不需要前端直接调用，应该配置为微信开放平台的回调地址。

查询参数：
- `code`: 微信授权码
- `state`: 与授权请求时传入的state参数相同

成功响应 (200)：

```json
{
  "code": 0,
  "message": "微信登录成功",
  "data": {
    "user": {
      "id": 1,
      "email": null,
      "nickname": "微信用户昵称",
      "avatar": "https://微信头像URL",
      "signature": null,
      "created_at": "2023-03-27T08:00:00Z",
      "updated_at": "2023-03-27T08:00:00Z"
    },
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

### 6. 苹果授权

**POST /oauth/apple/auth**

获取苹果授权信息。

请求参数：

```json
{
  "redirect_uri": "https://your-domain.com/callback"
}
```

成功响应 (200)：

```json
{
  "code": 0,
  "message": "苹果登录说明",
  "data": {
    "instruction": "苹果登录需要在客户端实现，请在客户端完成授权后，将授权结果发送到 /api/v1/oauth/apple/callback"
  }
}
```

### 7. 苹果授权回调

**POST /oauth/apple/callback**

处理苹果授权回调。

请求参数：

```json
{
  "code": "授权码", // 授权码和ID令牌至少需要提供一个
  "id_token": "ID令牌",
  "name": "用户姓名", // 可选
  "email": "用户邮箱", // 可选
  "first_name": "名", // 可选，如果提供了name则忽略
  "last_name": "姓"  // 可选，如果提供了name则忽略
}
```

成功响应 (200)：

```json
{
  "code": 0,
  "message": "苹果登录成功",
  "data": {
    "user": {
      "id": 1,
      "email": "user@example.com",
      "nickname": "Apple User",
      "avatar": null,
      "signature": null,
      "created_at": "2023-03-27T08:00:00Z",
      "updated_at": "2023-03-27T08:00:00Z"
    },
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

### 8. 退出登录

**POST /logout**

需要认证。

请求头：
```
Authorization: Bearer {token}
```

成功响应 (200)：

```json
{
  "code": 0,
  "message": "退出登录成功",
  "data": null
}
```

### 9. 获取用户信息

**GET /user**

获取当前登录用户的信息。需要认证。

请求头：
```
Authorization: Bearer {token}
```

成功响应 (200)：

```json
{
  "code": 0,
  "message": "获取用户信息成功",
  "data": {
    "user": {
      "id": 1,
      "email": "user@example.com",
      "nickname": "用户昵称",
      "avatar": "https://example.com/avatar.jpg",
      "signature": "个性签名",
      "created_at": "2023-03-27T08:00:00Z",
      "updated_at": "2023-03-27T08:00:00Z"
    }
  }
}
```

### 10. 修改用户信息

**PUT /user**

修改当前登录用户的信息。需要认证。

请求头：
```
Authorization: Bearer {token}
```

请求参数（可选字段）：

```json
{
  "nickname": "新昵称",
  "avatar": "新头像URL",
  "signature": "新个性签名"
}
```

成功响应 (200)：

```json
{
  "code": 0,
  "message": "更新用户信息成功",
  "data": {
    "user": {
      "id": 1,
      "email": "user@example.com",
      "nickname": "新昵称",
      "avatar": "https://example.com/new-avatar.jpg",
      "signature": "新个性签名",
      "created_at": "2023-03-27T08:00:00Z",
      "updated_at": "2023-03-27T09:00:00Z"
    }
  }
}
```

### 11. 获取设置

**GET /settings/{key}**

获取指定key的设置值。

路径参数：
- `key`: 设置的键名

成功响应 (200)：

```json
{
  "code": 0,
  "message": "获取设置成功",
  "data": {
    "key": "app.theme",
    "value": "dark",
    "created_at": "2023-03-27T08:00:00Z",
    "updated_at": "2023-03-27T09:00:00Z"
  }
}
```

### 12. 设置/更新设置

**PUT /settings/{key}**

设置或更新指定key的值。需要提供key的MD5值进行安全校验。

路径参数：
- `key`: 设置的键名

请求参数：

```json
{
  "value": "设置的值",
  "key_md5": "key的MD5值"
}
```

**MD5校验说明**：
- 需要计算key+盐值的MD5值并在请求中提供
- 盐值通过环境变量 `SETTING_SALT` 配置，请联系系统管理员获取
- 计算方式：MD5(key + SETTING_SALT)
- 例如，key为"app.theme"，盐值为"your_salt"，则计算MD5("app.theme" + "your_salt")
- 可以使用在线MD5工具或编程方式计算

成功响应 (200)：

```json
{
  "code": 0,
  "message": "设置保存成功",
  "data": {
    "key": "app.theme",
    "value": "dark",
    "created_at": "2023-03-27T08:00:00Z",
    "updated_at": "2023-03-27T09:00:00Z"
  }
}
```

**设置管理错误响应示例**

### 设置不存在 (404)

```json
{
  "code": 1004,
  "message": "设置不存在",
  "data": null
}
```

### MD5校验失败 (401)

```json
{
  "code": 1002,
  "message": "MD5校验失败，无权限操作此设置",
  "data": null
}
```

### Key格式错误 (400)

```json
{
  "code": 1001,
  "message": "key格式不正确，只允许字母、数字、下划线、点号",
  "data": null
}
```

### Key长度错误 (400)

```json
{
  "code": 1001,
  "message": "key长度必须在1-64字符之间",
  "data": null
}
```

### Value长度错误 (400)

```json
{
  "code": 1001,
  "message": "value长度不能超过10KB",
  "data": null
}
```

### 13. 清除指定设置的缓存

**DELETE /settings/{key}/cache**

清除指定key的缓存数据。

路径参数：
- `key`: 设置的键名

成功响应 (200)：

```json
{
  "code": 0,
  "message": "缓存清除成功",
  "data": null
}
```

### 14. 清除所有设置缓存

**DELETE /settings/cache**

清除所有设置的缓存数据。

成功响应 (200)：

```json
{
  "code": 0,
  "message": "所有缓存清除成功",
  "data": null
}
```

### 15. 获取缓存统计信息

**GET /settings/cache/stats**

获取缓存系统的统计信息。

成功响应 (200)：

```json
{
  "code": 0,
  "message": "获取缓存统计成功",
  "data": {
    "cache_enabled": true,
    "cached_settings_count": 15
  }
}
```

响应字段说明：
- `cache_enabled`: 缓存是否启用
- `cached_settings_count`: 当前缓存的设置数量

## 错误响应示例

### 参数错误 (400)

```json
{
  "code": 1001,
  "message": "请求参数错误: 邮箱格式不正确",
  "data": null
}
```

### 未授权 (401)

```json
{
  "code": 1002,
  "message": "登录已过期，请重新登录",
  "data": null
}
```

### 资源不存在 (404)

```json
{
  "code": 1004,
  "message": "用户不存在",
  "data": null
}
```

### 资源冲突 (409)

```json
{
  "code": 1009,
  "message": "邮箱已被注册",
  "data": null
}
```

### 服务器错误 (500)

```json
{
  "code": 2000,
  "message": "服务器内部错误",
  "data": null
}
``` 