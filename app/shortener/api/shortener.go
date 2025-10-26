// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.1

package main

import (
	"context"
	"flag"
	"fmt"
	"minify/app/shortener/api/internal/mq"
	"time"

	"minify/app/shortener/api/internal/config"
	"minify/app/shortener/api/internal/handler"
	"minify/app/shortener/api/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/shortener.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	// 加载东八时区（Asia/Shanghai）
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		panic(err)
	}

	// 设置全局时区，确保每次时间解析都使用东八时区
	time.Local = loc

	// ⭐ 4. (修改) 创建服务组，它将管理所有服务
	group := service.NewServiceGroup()
	defer group.Stop() // 确保服务组停止时，所有子服务都会被关闭

	// ⭐ 5. (修改) 创建 API 服务，但不再 defer server.Stop()
	server := rest.MustNewServer(c.RestConf)
	// defer server.Stop() // (移除) - group.Stop() 会接管

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	group.Add(server) // ⭐ 6. (新增) 将 API 服务添加到服务组

	// ⭐ 7. (新增) 初始化并添加 Kafka 消费者
	// 我们需要一个后台 context 来运行消费者
	consumerCtx := context.Background()
	// mq.Consumers 会返回一个 []service.Service 列表
	consumerServices := mq.Consumers(c, consumerCtx, ctx)
	for _, srv := range consumerServices {
		group.Add(srv) // ⭐ 8. (新增) 将每个消费者服务添加到服务组
	}

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	// server.Start() // (移除)
	group.Start() // ⭐ 9. (修改) 启动服务组!
}

//func main() {
//	flag.Parse()
//
//	var c config.Config
//	conf.MustLoad(*configFile, &c)
//
//	// 加载东八时区（Asia/Shanghai）
//	loc, err := time.LoadLocation("Asia/Shanghai")
//	if err != nil {
//		panic(err)
//	}
//
//	// 设置全局时区，确保每次时间解析都使用东八时区
//	time.Local = loc
//
//	server := rest.MustNewServer(c.RestConf)
//	defer server.Stop()
//
//	ctx := svc.NewServiceContext(c)
//	handler.RegisterHandlers(server, ctx)
//
//	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
//	server.Start()
//}
