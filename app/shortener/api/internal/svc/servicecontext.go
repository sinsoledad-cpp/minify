// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.1

package svc

import (
	"lucid/app/shortener/api/internal/config"
	"lucid/data/model/shortener"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config               config.Config
	ShortUrlsModel       shortener.ShortUrlsModel
	AggDailySummaryModel shortener.AggDailySummaryModel
	UrlAnalyticsModel    shortener.UrlAnalyticsModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.DB.DataSource)
	return &ServiceContext{
		Config:               c,
		ShortUrlsModel:       shortener.NewShortUrlsModel(conn),
		AggDailySummaryModel: shortener.NewAggDailySummaryModel(conn),
		UrlAnalyticsModel:    shortener.NewUrlAnalyticsModel(conn),
	}
}
