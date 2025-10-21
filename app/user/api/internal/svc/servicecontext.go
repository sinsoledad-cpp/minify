// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.1

package svc

import (
	"lucid/app/user/api/internal/config"
	"lucid/app/user/data/model"
	"lucid/app/user/domain/repository"
	datarepo "lucid/app/user/domain/repository"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config   config.Config
	UserRepo repository.UserRepository
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 1. 创建数据库连接
	conn := sqlx.NewMysql(c.Database.DataSource)

	// 2. 初始化 goctl model
	userModel := model.NewUsersModel(conn)

	// 3. 初始化 data.repository 实现
	// 注意：这里的 NewUserRepoImpl 是你手写的
	userRepo := datarepo.NewUserRepoImpl(userModel)
	return &ServiceContext{
		Config:   c,
		UserRepo: userRepo,
	}
}
