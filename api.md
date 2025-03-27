这是一个 golang 的项目，使用的是 gin 框架+gorm 框架，使用的是 mysql 数据库。
在完成功能的时候，需要迭代项目的创建 sql 语句的文件 migrate.sql，方便后期维护。

需要完成以下 api：
1.ios 用户注册，支持邮箱和密码注册，支持第三方登录（微信、苹果）
2.ios 用户登录，支持邮箱和密码登录，支持第三方登录（微信、苹果）
3.ios 用户退出登录
4.ios 用户信息获取，包括头像、昵称、注册日期、个性签名
5.ios 用户信息修改，包括头像、昵称、个性签名

以下是连接数据库的信息(保留到 .env 文件中，不参与版本管理)：
host: xxxxx
user: root
password: xxx
port: 3306
database: yuanqi_ios
