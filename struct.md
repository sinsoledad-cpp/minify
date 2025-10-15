/your-project-monorepo
├── app/
│   ├── user/                   // 微服务: 用户服务
│   │   ├── api/
│   │   ├── rpc/
│   │   │   └── etc/
│   │   │       └── user-rpc.yaml   // 配置文件: 连接到独立的 user_db
│   │   └── domain/
│   │
│   └── order/                  // 微服务: 订单服务
│       ├── api/
│       ├── rpc/
│       │   └── etc/
│       │       └── order-rpc.yaml  // 配置文件: 连接到独立的 order_db
│       └── domain/
│
├── common/                     // 通用共享库
│
├── data/                       // ⭐【共享的基础设施内核】
│   │                           // 这里定义了“如何”与外部系统交互的统一模式和实现。
│   │                           // 即使连接到不同的数据库，连接、查询、事务的“模式”是共享的。
│   ├── model/                  // 所有 goctl model 生成的数据库模型 (PO)
│   │   ├── user/               // ✅ user 服务的数据库模型 (user_db)
│   │   └── order/              // ✅ order 服务的数据库模型 (order_db)
│   └── repository_impl/        // 所有 domain.repository 接口的具体实现
│       ├── user_repo_impl.go   // 此实现会被注入 user_db 的连接
│       └── order_repo_impl.go  // 此实现会被注入 order_db 的连接
│
├── gen/                        // 生成的客户端SDK (独立的 Go Module)
│
├── protos/                     // 服务契约 (.proto 文件)
│
├── schema/                     // ⭐【数据库定义蓝图 (按服务隔离)】
│   └── sql/
│       ├── user/               // ✅ user 服务的数据库蓝图 (user_db)
│       │   ├── 001_create_users_table.sql
│       │   └── 002_add_extra_field.sql
│       │
│       └── order/              // ✅ order 服务的数据库蓝图 (order_db)
│           ├── 001_create_orders_table.sql
│           └── 002_create_order_items_table.sql
│
├── deployments/                // 部署与配置
├── scripts/                    // 自动化脚本
├── go.mod                      // 主项目的 Go Module
└── Makefile