# iOS 用户系统 API 文档

## 基本信息

- 基础 URL: `http://your-domain.com/api/v1`
- 所有包含请求体的请求必须指定内容类型为 `Content-Type: application/json`
- 所有需要认证的请求都需要在请求头中包含 `Authorization: Bearer {token}`

## 错误处理

所有 API 请求如果发生错误，将返回相应的 HTTP 状态码和错误信息，格式如下：

```json
{
  "error": "错误信息"
}
```

常见 HTTP 状态码：
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
  "message": "注册成功",
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
  "message": "登录成功",
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
  "avatar": "头像URL"
}
```

成功响应 (200)：

```json
{
  "message": "登录成功",
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
```

### 4. 退出登录

**POST /logout**

需要认证。

请求头：
```
Authorization: Bearer {token}
```

成功响应 (200)：

```json
{
  "message": "退出登录成功"
}
```

### 5. 获取用户信息

**GET /user**

获取当前登录用户的信息。需要认证。

请求头：
```
Authorization: Bearer {token}
```

成功响应 (200)：

```json
{
  "message": "获取用户信息成功",
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
```

### 6. 修改用户信息

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
  "message": "更新用户信息成功",
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
``` 