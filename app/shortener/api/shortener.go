// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.1

package main

import (
	"flag"
	"fmt"
	"time"

	"minify/app/shortener/api/internal/config"
	"minify/app/shortener/api/internal/handler"
	"minify/app/shortener/api/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
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

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
