// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.1

package svc

import (
	"lucid/app/shortener/api/internal/config"
	"lucid/data/model/shortener"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config         config.Config
	ShortUrlsModel shortener.ShortUrlsModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.DB.DataSource)
	return &ServiceContext{
		Config:         c,
		ShortUrlsModel: shortener.NewShortUrlsModel(conn),
	}
}
