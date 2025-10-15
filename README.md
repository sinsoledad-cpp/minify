lucid

中文释义: 清晰的，易懂的，明亮的。

涵义: 微服务治理的一大挑战就是复杂性。lucid 这个名字直接回应了这一痛点，寓意此框架能让原本混乱、模糊的系统架构变得清晰透明、易于理解和维护。它代表了一种追求“大道至简”的哲学，让开发者能保持头脑清醒，洞悉系统的每一处细节。

# template

1. etc.tpl

goctl rpc protoc user.proto --go_out="./internal/user" --go-grpc_out="./internal/user" --zrpc_out="./internal/user/zrpc" --style=go_zero --home=./template -c


```text
/your-project-monorepo
├── app/
│   └── user/                   // 微服务: 用户服务 (一个限界上下文)
│       ├── api/                // 【表示层】: 对外的 HTTP 接口, 服务的门面
│       │   ├── etc/
│       │   │   └── user-api.yaml
│       │   ├── internal/
│       │   │   ├── config/
│       │   │   ├── handler/
│       │   │   ├── logic/      // 应用层: 编排用例, 依赖 domain.repository 接口
│       │   │   ├── svc/
│       │   │   └── types/
│       │   ├── user.api        // API 接口定义文件
│       │   └── user.go         // API 服务启动入口
│       │
│       ├── rpc/                // 【应用层】: 对内的 gRPC 接口, 承载核心业务
│       │   ├── etc/
│       │   │   └── user-rpc.yaml
│       │   ├── internal/
│       │   │   ├── config/
│       │   │   ├── logic/      // 应用层: 实现 proto 接口, 依赖 domain.repository 接口
│       │   │   ├── server/     // 🔒 RPC 服务端(Server)的具体实现
│       │   │   └── svc/
│       │   └── user.go         // RPC 服务启动入口
│       │
│       └── domain/             // ⭐【领域层】: 业务核心, 框架无关, 纯 Go
│           ├── aggregate/      // 聚合
│           ├── entity/         // 实体 (包含业务行为)
│           ├── valueobject/    // 值对象
│           └── repository/     // 仓库接口 (定义数据持久化的抽象)
│
├── common/                     // 【通用共享库】: 手写的、稳定的、跨服务共享的通用工具
│   ├── constants/              // 全局常量
│   ├── utils/                  // 通用工具函数
│   └── xerr/                   // 统一错误处理
│
├── data/                       // ⭐【基础设施层】: 所有技术细节的具体实现
│   ├── model/                  // ✅ goctl model 生成的数据库模型 (PO) 和 DAO
│   └── repository_impl/        // ✅ domain.repository 接口的具体实现 (在此使用 model)
│
├── gen/                        // ⭐【生成代码/客户端SDK】: 所有服务生成的客户端代码
│   └── go/                     // Go 语言的生成代码
│       ├── user/
│       │   └── v1/
│       │       ├── user.pb.go          // ✅ 消息体定义
│       │       ├── user_grpc.pb.go     // ✅ gRPC 接口定义
│       │       └── user/               // 文件夹
│       │           └── user.go         // ✅ zRPC 客户端封装
│       └── go.mod              // ⭐ 关键: 这是一个独立的 Go Module, 用于解耦
│
├── protos/                     // 【服务契约】: 所有 .proto 文件, 服务间通信的蓝图
│   └── user/
│       └── v1/
│           └── user.proto
│
├── schema/                     // 【数据库定义】: 数据库结构的唯一事实来源
│   └── sql/
│       └── user.sql
│
├── deployments/                // 【部署与配置】: 与部署相关的所有文件
│   ├── docker/                 // Dockerfile
│   ├── docker-compose/         // 本地开发环境编排
│   └── kubernetes/             // Kubernetes manifests (yaml)
│
├── scripts/                    // 【自动化脚本】: 提升开发效率
│   ├── gen_model.sh            // 一键生成数据库模型
│   ├── gen_proto.sh            // 一键生成所有 proto 客户端
│   └── build.sh                // 编译脚本
│
├── .gitignore
├── go.mod                      // 主项目的 Go Module
├── go.sum
└── Makefile                    // 推荐使用 Makefile 简化常用命令


```