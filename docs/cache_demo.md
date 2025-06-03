# LevelDB 缓存功能演示

本文档演示如何使用项目中的LevelDB缓存功能来优化设置数据的读取性能。

## 功能概述

项目使用LevelDB作为设置数据的本地缓存，提供以下功能：

1. **自动缓存**：首次读取设置时自动缓存
2. **缓存失效**：更新设置时自动删除对应缓存
3. **手动缓存管理**：提供API接口手动管理缓存
4. **缓存统计**：查看缓存使用情况

## 性能优势

- **减少数据库查询**：频繁访问的设置从缓存读取
- **提高响应速度**：本地文件访问比网络数据库访问更快
- **降低数据库负载**：减少对数据库的压力

## 使用演示

### 1. 启动服务器

```bash
go run main.go
```

服务器启动时会显示缓存目录信息：
```
服务器启动，监听端口: :8080
缓存目录: ./cache
```

### 2. 设置数据

首先，我们需要设置一些数据：

```bash
# 计算MD5值（假设SETTING_SALT为"default_setting_salt"）
# MD5("app.theme" + "default_setting_salt") = 某个MD5值

curl -X PUT http://localhost:8080/api/v1/settings/app.theme \
  -H "Content-Type: application/json" \
  -d '{
    "value": "dark",
    "key_md5": "您计算的MD5值"
  }'
```

### 3. 读取数据（第一次 - 从数据库）

```bash
curl -X GET http://localhost:8080/api/v1/settings/app.theme
```

响应：
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

### 4. 读取数据（第二次 - 从缓存）

再次执行相同的GET请求，这次数据将从缓存读取，速度更快：

```bash
curl -X GET http://localhost:8080/api/v1/settings/app.theme
```

### 5. 查看缓存统计

```bash
curl -X GET http://localhost:8080/api/v1/settings/cache/stats
```

响应：
```json
{
  "code": 0,
  "message": "获取缓存统计成功",
  "data": {
    "cache_enabled": true,
    "cached_settings_count": 1
  }
}
```

### 6. 更新数据（自动清除缓存）

```bash
curl -X PUT http://localhost:8080/api/v1/settings/app.theme \
  -H "Content-Type: application/json" \
  -d '{
    "value": "light",
    "key_md5": "您计算的MD5值"
  }'
```

更新操作会自动清除对应的缓存，下次读取时会从数据库获取最新数据。

### 7. 手动清除指定缓存

```bash
curl -X DELETE http://localhost:8080/api/v1/settings/app.theme/cache
```

响应：
```json
{
  "code": 0,
  "message": "缓存清除成功",
  "data": null
}
```

### 8. 清除所有缓存

```bash
curl -X DELETE http://localhost:8080/api/v1/settings/cache
```

响应：
```json
{
  "code": 0,
  "message": "所有缓存清除成功",
  "data": null
}
```

## 性能测试

可以使用以下脚本进行简单的性能测试：

```bash
#!/bin/bash

# 性能测试脚本
KEY="performance.test"
VALUE="performance test value"
MD5="您的MD5值"  # 需要计算 MD5(performance.test + SETTING_SALT)

echo "设置测试数据..."
curl -s -X PUT http://localhost:8080/api/v1/settings/$KEY \
  -H "Content-Type: application/json" \
  -d "{\"value\":\"$VALUE\",\"key_md5\":\"$MD5\"}"

echo -e "\n\n测试读取性能（10次）..."

for i in {1..10}; do
  echo "第 $i 次读取:"
  time curl -s http://localhost:8080/api/v1/settings/$KEY > /dev/null
done
```

预期结果：
- 第1次读取：从数据库，响应时间较长
- 第2-10次读取：从缓存，响应时间明显缩短

## 缓存文件结构

缓存文件存储在配置的目录中（默认为 `./cache`）：

```
cache/
└── settings_cache/
    ├── 000001.log
    ├── CURRENT
    ├── LOCK
    ├── LOG
    └── MANIFEST-000000
```

这些是LevelDB的内部文件，不建议手动修改。

## 注意事项

1. **权限**：确保应用有读写缓存目录的权限
2. **磁盘空间**：缓存会占用磁盘空间，需要定期监控
3. **数据一致性**：缓存与数据库可能存在短暂的不一致
4. **备份**：缓存数据在重启后依然保持，可以加快启动后的首次访问
5. **清理**：如果需要完全重置缓存，可以删除整个缓存目录

## 最佳实践

1. **定期清理**：在系统维护期间清理不必要的缓存
2. **监控统计**：定期查看缓存统计，了解使用情况
3. **合理配置**：根据实际需求配置缓存目录位置
4. **错误处理**：缓存出错时应优雅降级到数据库访问

通过以上功能，您可以显著提高设置数据的读取性能，特别是对于频繁访问的配置项。 