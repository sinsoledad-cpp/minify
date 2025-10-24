minify/ (项目根目录)
├── app/
│   ├── user/                   // ⭐ 微服务: 用户服务 (一个限界上下文)
│   │   │
│   │   ├── api/                // 1. 应用层 (Application): HTTP/API 接入
│   │   │   ├── etc/
│   │   │   │   └── user.yaml
│   │   │   ├── internal/
│   │   │   │   ├── config/
│   │   │   │   │   └── config.go
│   │   │   │   ├── handler/
│   │   │   │   │   └── ... (go-zero 生成)
│   │   │   │   ├── logic/
│   │   │   │   │   ├── errcode/
│   │   │   │   │   │   └── errcode.go          // (可选) 逻辑层共享的业务错误码
│   │   │   │   │   ├── user/
│   │   │   │   │   │   └── loginlogic.go       // (go-zero 生成)
│   │   │   │   │   │   └── ...
│   │   │   │   │   └── converter.go        // ⭐ 你的实践: 负责 DTO (types.go) 与 Entity (domain/entity) 互转
│   │   │   │   ├── svc/
│   │   │   │   │   └── servicecontext.go   // (go-zero 生成) 依赖注入 domain/repository 接口
│   │   │   │   └── types/
│   │   │   │       └── types.go            // (go-zero 生成) API 的 DTO
│   │   │   ├── user.api
│   │   │   └── user.go
│   │   │
│   │   ├── rpc/                // 1. 应用层 (Application): gRPC/RPC 接入
│   │   │   ├── etc/
│   │   │   │   └── user.yaml
│   │   │   ├── internal/
│   │   │   │   ├── logic/      // (go-zero 生成)
│   │   │   │   ├── server/     // (go-zero 生成)
│   │   │   │   └── svc/        // (go-zero 生成)
│   │   │   └── user.go
│   │   │
│   │   ├── domain/             // 2. 领域层 (Domain): 核心业务 (手写)
│   │   │   ├── entity/
│   │   │   │   └── user.go     //    - 领域实体 (e.g., User struct, 包含业务方法 CheckPassword)
│   │   │   ├── repository/
│   │   │   │   └── user_repository.go //    - 仓储 *接口* (Interface)
│   │   │   └── service/
│   │   │       └── user_service.go    // ⭐ 你的实践: 领域服务 (e.g., 复杂的密码策略、用户激活)
│   │   │
│   │   ├── data/               // 3. 基础设施层 (Infrastructure): 领域层的具体实现
│   │   │   ├── model/          //    - goctl 生成的 DB 模型 (PO, e.g., UsersModel)
│   │   │   │   └── ...
│   │   │   └── repository/
│   │   │       └── user_repo_impl.go //    - 仓储 *实现* (Implementation), 依赖 model
│   │   │
│   │   └── schema/             // 4. 数据库定义 (Migrations)
│   │       └── sql/
│   │           └── 000001_users.up.sql
│   │
│   └── shortener/              // (⭐ 另一个微服务: shortener, 结构同上)
│       ├── api/
│       ├── rpc/
│       ├── domain/
│       ├── data/
│       └── schema/
│
├── common/                     // ⭐ 通用共享库 (你的规划实现了)
│   ├── service/
│   │   └── snowflake/
│   │       └── node.go
│   └── utils/
│       ├── codec/
│       │   └── base62.go
│       ├── jwtx/
│       │   └── jwt.go
│       └── response/
│           └── response.go     (你的统一响应体)
│
├── protos/                     // ⭐ 服务契约 (.proto 文件)
│   └── user/
│       └── v1/
│           └── user.proto
│
├── scripts/                    // 自动化脚本
│   └── docker-compose.yml
│
├── template/                   // ⭐ 你的 goctl 自定义模板
│   ├── api/
│   │   ├── handler.tpl         (你刚刚修改的)
│   │   └── ...
│   ├── model/
│   └── ... (其他模板)
│
├── go.mod
├── go.sum
├── Makefile
├── README.md
└── struct.md                   (你最初的规划文件)