package repository

import (
	"context"
	"database/sql"
	"minify/app/shortener/data/model"
	"minify/app/shortener/domain/entity"

	"github.com/zeromicro/go-zero/core/logx"
)

// 确保 impl 实现了接口
var _ LinkAccessLogsRepository = (*linkAccessLogsRepoImpl)(nil)

type linkAccessLogsRepoImpl struct {
	model model.LinkAccessLogsModel // 依赖 goctl model
}

// NewLinkAccessLogsRepoImpl 创建仓储实现
func NewLinkAccessLogsRepoImpl(model model.LinkAccessLogsModel) LinkAccessLogsRepository {
	return &linkAccessLogsRepoImpl{
		model: model,
	}
}

// toModelLog 将日志实体(Entity)转换为数据模型(PO)
// (仓储实现层的私有方法)
func toModelLog(e *entity.LinkAccessLog) *model.LinkAccessLogs {
	return &model.LinkAccessLogs{
		// Id 由数据库自增，不需要传
		LinkId:     uint64(e.LinkID),
		ShortCode:  e.ShortCode,
		AccessedAt: e.AccessedAt,
		IpAddress:  e.IpAddress,
		UserAgent:  sql.NullString{String: e.UserAgent, Valid: e.UserAgent != ""},
		Referer:    sql.NullString{String: e.Referer, Valid: e.Referer != ""},
		// Geo/UA 字段
		GeoCountry:  sql.NullString{String: e.GeoCountry, Valid: e.GeoCountry != ""},
		GeoCity:     sql.NullString{String: e.GeoCity, Valid: e.GeoCity != ""},
		DeviceType:  sql.NullString{String: e.DeviceType, Valid: e.DeviceType != ""},
		BrowserName: sql.NullString{String: e.BrowserName, Valid: e.BrowserName != ""},
		OsName:      sql.NullString{String: e.OsName, Valid: e.OsName != ""},
	}
}

// Create 实现了接口
func (r *linkAccessLogsRepoImpl) Create(ctx context.Context, log *entity.LinkAccessLog) error {
	po := toModelLog(log)

	// 调用 goctl model 的 Insert
	if _, err := r.model.Insert(ctx, po); err != nil {
		logx.WithContext(ctx).Errorf("linkAccessLogsRepoImpl.Create error: %v", err)
		return err
	}
	return nil
}
