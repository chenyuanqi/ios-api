# Ubuntu 16.04 系统部署指南

本文档提供在 Ubuntu 16.04 系统上部署 iOS 用户系统 API 的详细步骤。

## 一、环境准备

### 1. 更新系统

```bash
sudo apt update
sudo apt upgrade -y
```

### 2. 安装 Go 语言环境

由于 Ubuntu 16.04 仓库中的 Go 版本较旧，我们需要手动安装新版本：

```bash
# 下载 Go 1.19 (项目所需版本)
wget https://golang.org/dl/go1.19.linux-amd64.tar.gz

# 解压到 /usr/local
sudo tar -C /usr/local -xzf go1.19.linux-amd64.tar.gz

# 设置环境变量
echo 'export PATH=$PATH:/usr/local/go/bin' | sudo tee -a /etc/profile
echo 'export GOPATH=$HOME/go' | tee -a ~/.bashrc
echo 'export PATH=$PATH:$GOPATH/bin' | tee -a ~/.bashrc

# 应用环境变量
source /etc/profile
source ~/.bashrc

# 验证安装
go version
```

### 3. 安装 MySQL 5.7

```bash
sudo apt install mysql-server-5.7 -y

# 启动 MySQL 服务
sudo systemctl start mysql
sudo systemctl enable mysql

# 设置 MySQL 安全配置
sudo mysql_secure_installation
```

### 4. 安装 Nginx

```bash
sudo apt install nginx -y
sudo systemctl start nginx
sudo systemctl enable nginx
```

## 二、部署项目

### 1. 创建项目目录

```bash
sudo mkdir -p /var/www/ios-api
sudo chown -R $USER:$USER /var/www/ios-api
```

### 2. 获取项目代码

```bash
# 如果使用 git
git clone [仓库地址] /var/www/ios-api

# 或者手动上传项目文件
# 使用 scp 或 sftp 工具
```

### 3. 安装项目依赖

```bash
cd /var/www/ios-api
go mod download
```

### 4. 配置环境变量

创建 `.env` 文件：

```bash
cat > /var/www/ios-api/.env << EOF
# 数据库配置
DB_HOST=127.0.0.1
DB_USER=你的MySQL用户名
DB_PASSWORD=你的MySQL密码
DB_PORT=3306
DB_NAME=yuanqi_ios

# JWT配置
JWT_SECRET=长而随机的字符串

# 应用配置
APP_PORT=8080
EOF
```

### 5. 初始化数据库

```bash
# 创建数据库
sudo mysql -e "CREATE DATABASE IF NOT EXISTS yuanqi_ios CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

# 创建专用数据库用户并授权
sudo mysql -e "CREATE USER 'ios_api_user'@'localhost' IDENTIFIED BY '密码';"
sudo mysql -e "GRANT ALL PRIVILEGES ON yuanqi_ios.* TO 'ios_api_user'@'localhost';"
sudo mysql -e "FLUSH PRIVILEGES;"

# 导入数据表结构
sudo mysql yuanqi_ios < /var/www/ios-api/migrate.sql
```

### 6. 编译项目

```bash
cd /var/www/ios-api
go build -o ios-api
```

## 三、配置服务

### 1. 创建 Systemd 服务文件

```bash
sudo cat > /etc/systemd/system/ios-api.service << EOF
[Unit]
Description=iOS API Service
After=network.target mysql.service

[Service]
Type=simple
User=www-data
Group=www-data
WorkingDirectory=/data/www/ios-api
ExecStart=/data/www/ios-api/ios-api
Restart=on-failure
RestartSec=5
StandardOutput=syslog
StandardError=syslog
SyslogIdentifier=ios-api

[Install]
WantedBy=multi-user.target
EOF
```

### 2. 设置文件权限

```bash
sudo chown -R www-data:www-data /data/www/ios-api
sudo chmod +x /data/www/ios-api/ios-api
```

### 3. 启动服务

```bash
sudo systemctl daemon-reload
sudo systemctl start ios-api
sudo systemctl enable ios-api
sudo systemctl status ios-api
```

## 四、配置 Nginx 反向代理

### 1. 创建 Nginx 配置文件

```bash
sudo cat > /etc/nginx/sites-available/ios-api << EOF
server {
    listen 80;
    server_name your-domain.com;  # 替换为你的域名

    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade \$http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host \$host;
        proxy_cache_bypass \$http_upgrade;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
    }
}
EOF
```

### 2. 启用站点配置并重启 Nginx

```bash
sudo ln -s /etc/nginx/sites-available/ios-api /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl restart nginx
```

## 五、HTTPS 配置（推荐）

使用 Let's Encrypt 获取免费 SSL 证书：

```bash
# 安装 Certbot
sudo apt install certbot python-certbot-nginx -y

# 获取并安装证书
sudo certbot --nginx -d your-domain.com

# 自动续期
sudo certbot renew --dry-run
```

## 六、防火墙配置

```bash
# 安装 UFW
sudo apt install ufw -y

# 配置防火墙规则
sudo ufw allow ssh
sudo ufw allow http
sudo ufw allow https

# 启用防火墙
sudo ufw enable
```

## 七、测试验证

1. 检查服务状态：
```bash
sudo systemctl status ios-api
```

2. 测试 API 接口：
```bash
curl -X POST http://your-domain.com/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123",
    "nickname": "测试用户"
  }'
```

## 八、日志和监控

### 1. 查看日志

```bash
# 查看应用日志
sudo journalctl -u ios-api -f

# 查看 Nginx 访问日志
sudo tail -f /var/log/nginx/access.log

# 查看 Nginx 错误日志
sudo tail -f /var/log/nginx/error.log
```

### 2. 设置简单监控

```bash
# 安装监控工具
sudo apt install htop -y
```

## 九、备份策略

### 1. 数据库备份

创建定时备份脚本：

```bash
# 创建备份目录
sudo mkdir -p /var/backups/ios-api

# 创建备份脚本
cat > /var/www/ios-api/backup.sh << EOF
#!/bin/bash
DATE=\$(date +%Y-%m-%d)
BACKUP_DIR=/var/backups/ios-api
mysqldump -u root -p密码 yuanqi_ios > \$BACKUP_DIR/yuanqi_ios_\$DATE.sql
find \$BACKUP_DIR -type f -name "yuanqi_ios_*.sql" -mtime +7 -delete
EOF

# 设置权限
chmod +x /var/www/ios-api/backup.sh

# 设置定时任务
(crontab -l ; echo "0 2 * * * /var/www/ios-api/backup.sh") | crontab -
```

## 十、常见问题排查

1. 服务无法启动
   - 检查环境变量配置: `cat /var/www/ios-api/.env`
   - 检查 MySQL 连接: `mysql -u ios_api_user -p -h localhost yuanqi_ios`
   - 检查服务日志: `sudo journalctl -u ios-api -n 50`

2. Nginx 连接问题
   - 检查 Nginx 配置: `sudo nginx -t`
   - 检查网络连接: `curl -v http://localhost:8080`

3. 数据库问题
   - 检查 MySQL 服务状态: `sudo systemctl status mysql`
   - 检查数据库连接: `mysql -u ios_api_user -p -h localhost yuanqi_ios -e "SELECT 1;"` 