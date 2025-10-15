// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.1

package svc

import (
	"lucid/app/user/api/internal/config"
	"lucid/data/model/user"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config    config.Config
	UsersModel user.UsersModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.DB.DataSource)
	return &ServiceContext{
		Config:    c,
		UsersModel: user.NewUsersModel(conn),
	}
}
