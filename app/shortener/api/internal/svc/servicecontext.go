// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.1

package svc

import (
	"lucid/app/shortener/api/internal/config"
	"lucid/data/model/shortener"
)

type ServiceContext struct {
	Config config.Config
	ShortUrlsModel shortener.ShortUrlsModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
	}
}
