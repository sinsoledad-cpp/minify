// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.1

package svc

import (
	"minify/app/user/api/internal/config"
	"minify/app/user/data/model"
	"minify/app/user/domain/repository"
	datarepo "minify/app/user/domain/repository"
	"minify/common/middleware"

	"github.com/casbin/casbin/v2"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/rest"
)

type ServiceContext struct {
	Config          config.Config
	AuthzMiddleware rest.Middleware
	UserRepo        repository.UserRepository
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 1. 创建数据库连接
	conn := sqlx.NewMysql(c.Database.DataSource)

	// 2. 初始化 goctl model
	userModel := model.NewUsersModel(conn)

	// 3. 初始化 data.repository 实现
	// 注意：这里的 NewUserRepoImpl 是你手写的
	userRepo := datarepo.NewUserRepoImpl(userModel)
	e, err := casbin.NewEnforcer(c.Casbin.ModelPath, c.Casbin.PolicyPath)
	if err != nil {
		logx.Must(err)
	}

	// 6. 从文件加载策略
	if err := e.LoadPolicy(); err != nil {
		logx.Must(err)
	}
	return &ServiceContext{
		Config:          c,
		UserRepo:        userRepo,
		AuthzMiddleware: middleware.NewAuthzMiddleware(e).Handle,
	}
}
