# 拾序｜Seqlite

slogan：**拾字成篇，有序留存**

## Sqlite Server

Seqlite 的 Go 后端服务，使用 Gin 处理 HTTP 请求，PostgreSQL 存储业务数据，Redis 存储签名密钥和刷新会话。

项目按业务划分 Controller、Service、Repository，目前以认证模块为主要实现。其他业务接口沿用相同的分层方式扩展。

## 项目结构

```text
.
├── cmd/
│   └── main.go                      # 程序入口：加载配置、初始化日志、启动与优雅退出
│
├── migrations/                      # Goose SQL 迁移，负责表结构、索引与外键
│   ├── 00001_create_users.sql       # 用户表与有效邮箱唯一索引
│   ├── 00002_create_spaces_and_documents.sql # 空间、文档及所属关系
│   ├── 00003_create_comments.sql    # 评论及同一文档内的父子评论约束
│   ├── 00004_create_roles_and_permissions.sql # 角色、权限与关联表
│   └── 00005_create_grants.sql      # 空间授权、文档授权及具体资源外键
│
├── internal/
│   ├── config/                      # 应用配置
│   │   ├── config.go                # HTTP、PostgreSQL、Redis、JWT 配置读取与校验
│   │   └── env.go                   # 按 APP_ENV 加载环境文件
│   │
│   ├── controller/                  # HTTP 请求处理
│   │   ├── auth.go                  # 登录参数绑定、Token 响应与 Cookie；其他认证处理预留
│   │   └── …
│   │
│   ├── service/                     # 业务流程与规则
│   │   ├── auth.go                  # 登录编排：用户查询、密码校验、Token 签发、会话创建
│   │   └── …
│   │
│   ├── repository/                  # 数据存取，封装 PostgreSQL / Redis 操作
│   │   ├── user_postgres.go         # 用户创建、邮箱存在检查与用户查询
│   │   ├── refresh_session_redis.go # 刷新会话写入与过期时间计算
│   │   ├── sign_key_redis.go        # 签名密钥存取、按需轮换与保留期限管理
│   │   └── …
│   │
│   ├── model/                       # GORM 数据模型
│   │   ├── user.go                  # 用户模型、邮箱唯一索引与软删除
│   │   ├── document.go              # 文档模型、所属空间与所有者
│   │   ├── grant.go                 # SpaceGrants、DocumentGrants 两类资源授权
│   │   └── …
│   │
│   ├── dto/                         # HTTP 请求、响应的数据结构
│   │   ├── auth.go                  # 登录请求、登录响应与注册请求
│   │   └── …
│   │
│   ├── router/                      # 路由与中间件注册
│   │   ├── router.go                # 创建 Gin 引擎、注册全局中间件与 /api 路由组
│   │   ├── middleware.go            # 全局中间件注册辅助函数
│   │   ├── auth.go                  # 登录、注册、刷新、退出的路由绑定
│   │   └── …
│   │
│   ├── middleware/                  # 请求处理中的通用逻辑
│   │   ├── auth.go                  # Bearer Access Token 验证与用户、会话上下文
│   │   ├── logger.go                # 基于 slog 的请求日志与错误记录
│   │   └── request_id.go            # 读取或生成请求 ID，设置 X-Request-ID
│   │
│   ├── response/                    # 统一 HTTP JSON 响应
│   │   ├── response.go              # 响应结构、成功响应与错误响应输出
│   │   ├── code.go                  # 响应业务码
│   │   ├── write_error.go           # 应用错误到 HTTP 状态、业务码和消息的映射
│   │   └── bind_error.go            # JSON 绑定和参数校验错误格式化
│   │
│   ├── errs/                        # 统一错误定义、分类与包装
│   │   ├── error.go                 # 应用错误类型、预定义错误及错误链
│   │   └── kind.go                  # 错误种类及对应 HTTP 状态
│   │
│   ├── security/                    # 认证所需的基础能力
│   │   ├── jwt.go                   # RS256 JWT 签发、验签及密钥保留期限计算
│   │   ├── password.go              # Argon2id 密码哈希生成与校验
│   │   ├── rsa.go                   # RSA 密钥对生成与 PEM 编码
│   │   ├── session_id.go            # 安全随机会话 ID 生成
│   │   └── sign_key_store.go        # TokenManager 使用的签名密钥存取接口
│   │
│   ├── database/                    # 数据库连接初始化与管理
│   │   ├── postgres.go              # GORM 连接、连接池、健康检查与关闭
│   │   └── redis.go                 # Redis 客户端、TLS、连接检查与关闭
│   │
│   └── server/                      # HTTP 服务构建与依赖组装
│       ├── server.go                # 初始化数据库、构建 HTTP Server 与超时配置
│       └── wire.go                  # 组装密钥仓储、TokenManager 和认证模块依赖
│
├── .env.example                     # 应用及 Goose 数据库连接配置示例
├── docker-compose.yml               # 容器编排文件，当前为空
├── Makefile                         # 应用构建与数据库迁移命令
├── go.mod                           # Go 版本与依赖声明
├── go.sum                           # 依赖校验信息
└── README.md
```

## 分层职责

请求主要沿以下路径处理：

```text
HTTP 请求 → Router / Middleware → Controller → Service → Repository → PostgreSQL / Redis
```

- **Controller**：绑定和校验请求参数、调用业务逻辑、处理 Cookie，并通过
  `response` 输出结果。
- **Service**：编排业务流程和业务规则，调用 Repository 与 Security。
- **Repository**：封装查询、写入、Redis Lua 脚本和数据保留期限等存储逻辑。
- **Model / DTO**：Model 描述数据库模型；DTO 描述 HTTP 请求和响应。
- **Security**：提供密码、JWT、RSA 和会话 ID 等能力；通过 `SigningKeyStore`
  接口访问签名密钥。
- **Errs /
  Response**：Errs 提供错误类型和分类；Response 将错误映射为 HTTP 状态、业务码及 JSON 响应。
- **Database / Server**：Database 初始化连接；Server 组装依赖并构建 HTTP 服务。

新增业务时，可参照认证模块，在 `controller`、`service`、`repository`、`dto` 和
`router` 中增加对应文件；涉及持久化实体时，在 `model` 中增加模型。
