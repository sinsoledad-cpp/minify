package repository

import (
	"context"
	"minify/app/shortener/domain/entity"

	// ⭐ 不再导入 api/internal/types
	"time"
)

// AnalyticsRepository 是报表仓储的接口
// ⭐ 返回类型已修改为 domain 对象
type AnalyticsRepository interface {
	// GetDashboardData 获取仪表盘总览数据
	GetDashboardData(ctx context.Context, userId *uint64, startDate, endDate time.Time) (*entity.DashboardSummary, error) // ⭐ 返回 *entity.DashboardSummary

	// GetLinkAnalytics 获取单个链接的详细分析数据
	GetLinkAnalytics(ctx context.Context, linkId int64, startDate, endDate time.Time) (*entity.LinkAnalytics, error) // ⭐ 返回 *entity.LinkAnalytics

	// GetTopLink 获取点击量最高的链接 (供 Dashboard 使用)
	GetTopLink(ctx context.Context, userId *uint64, startDate, endDate time.Time) (*entity.Link, error) // ⭐ 返回 *entity.Link
}
