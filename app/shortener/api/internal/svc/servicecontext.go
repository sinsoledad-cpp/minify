// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.1

package svc

import (
	"minify/app/shortener/api/internal/config"
	"minify/app/shortener/data/model"
	"minify/app/shortener/domain/repository"
	"minify/app/shortener/domain/service"
	"minify/common/middleware"
	"minify/common/service/snowflake"

	"github.com/casbin/casbin/v2"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/rest"
)

type ServiceContext struct {
	Config          config.Config
	AuthzMiddleware rest.Middleware
	//RedisClient   *redis.Redis                   // ⭐ Redis 客户端 (主要供 links 缓存使用)
	LinkRepo      repository.LinkRepository      // ⭐ 注入 Link 仓储接口
	AnalyticsRepo repository.AnalyticsRepository // ⭐ 注入 Analytics 仓储接口
	IdGenerator   service.IdGenerator
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 1. 初始化数据库连接
	conn := sqlx.NewMysql(c.Database.DataSource)

	// 2. 初始化 goctl models (linksModel 需要 c.CacheRedis)
	linksModel := model.NewLinksModel(conn, c.CacheRedis)
	summaryModel := model.NewAnalyticsSummaryDailyModel(conn)

	// 3. 初始化 data.repository 实现
	linkRepo := repository.NewLinkRepoImpl(linksModel)
	analyticsRepo := repository.NewAnalyticsRepoImpl(summaryModel, linksModel)

	idGen, err := snowflake.NewGenerator(c.Snowflake.WorkerId)

	e, err := casbin.NewEnforcer(c.Casbin.ModelPath, c.Casbin.PolicyPath)
	if err != nil {
		logx.Must(err)
	}

	// 6. 从文件加载策略
	if err := e.LoadPolicy(); err != nil {
		logx.Must(err)
	}

	if err != nil {
		logx.Must(err) // 初始化失败，直接 panic
	}
	return &ServiceContext{
		Config:          c,
		AuthzMiddleware: middleware.NewAuthzMiddleware(e).Handle,
		LinkRepo:        linkRepo,
		AnalyticsRepo:   analyticsRepo,
		IdGenerator:     idGen,
	}
}
