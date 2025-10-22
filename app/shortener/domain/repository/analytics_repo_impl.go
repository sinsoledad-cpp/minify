package repository // ⭐ 包名是 repository

import (
	"context"
	"database/sql" // ⭐ 导入 sql
	"errors"
	// "fmt" // 如果自定义SQL移到model层，不再需要
	"minify/app/shortener/data/model" // ⭐ 依赖 data/model
	"minify/app/shortener/domain/entity"
	// "minify/app/shortener/domain/repository" // 不再需要导入自己 (因为在同一个包)
	"sort"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	// "github.com/zeromicro/go-zero/core/stores/sqlx" // 如果自定义SQL移到model层，不再需要
)

// 确保 analyticsRepoImpl 实现了 AnalyticsRepository 接口 (接口定义在 analytics_repository.go)
var _ AnalyticsRepository = (*analyticsRepoImpl)(nil)

type analyticsRepoImpl struct {
	summaryModel model.AnalyticsSummaryDailyModel // ⭐ 只依赖 model 接口
	linksModel   model.LinksModel                 // ⭐ 只依赖 model 接口
}

// NewAnalyticsRepoImpl 创建 AnalyticsRepository 的实现
// 注入 model 接口
func NewAnalyticsRepoImpl(summaryModel model.AnalyticsSummaryDailyModel, linksModel model.LinksModel) AnalyticsRepository { // ⭐ 返回接口类型
	return &analyticsRepoImpl{
		summaryModel: summaryModel,
		linksModel:   linksModel,
	}
}

// --- 接口实现 ---

// GetDashboardData 获取仪表盘总览数据
func (r *analyticsRepoImpl) GetDashboardData(ctx context.Context, userId *uint64, startDate, endDate time.Time) (*entity.DashboardSummary, error) {
	summary := &entity.DashboardSummary{}
	var err error

	// 1. 获取 TotalLinks
	// 调用 linksModel 的自定义 Count 方法
	linkStatus := entity.StatusAll // 查所有未删除的
	if userId != nil {
		summary.TotalLinks, err = r.linksModel.CountByUserIdAndStatus(ctx, *userId, linkStatus)
	} else {
		// 假设 admin 的 dashboard 也是查自己的 total links (如前讨论)
		// 如果需要全局 count，linksModel 必须添加 CountAllByStatus 方法
		return nil, errors.New("cannot get dashboard total links without user id (even for admin)")
	}
	// 处理 Count 可能返回的 ErrNotFound
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			summary.TotalLinks = 0 // 没有找到记录，总数为 0
		} else {
			logx.WithContext(ctx).Errorf("GetDashboardData TotalLinks error: %v", err)
			return nil, err
		}
	}

	// 2. 获取 TotalClicks
	// 调用 summaryModel 的自定义 FindTotalClicks 方法
	linksTableName := r.linksModel.TableName() // 获取表名给 summaryModel 使用
	summary.TotalClicks, err = r.summaryModel.FindTotalClicks(ctx, userId, startDate, endDate, linksTableName)
	if err != nil {
		logx.WithContext(ctx).Errorf("GetDashboardData TotalClicks error: %v", err)
		return nil, err
	}

	// 3. 获取 TopLink
	// 调用自身的 GetTopLink 方法
	summary.TopLink, err = r.GetTopLink(ctx, userId, startDate, endDate)
	// GetTopLink 内部会处理 ErrNotFound，这里只需检查其他错误
	if err != nil && !errors.Is(err, entity.ErrLinkNotFound) {
		logx.WithContext(ctx).Errorf("GetDashboardData GetTopLink error: %v", err)
		// 忽略错误，允许 TopLink 为空
	}

	return summary, nil
}

// GetLinkAnalytics 获取单个链接的详细分析数据
func (r *analyticsRepoImpl) GetLinkAnalytics(ctx context.Context, linkId int64, startDate, endDate time.Time) (*entity.LinkAnalytics, error) {
	analytics := &entity.LinkAnalytics{
		LinkID: linkId, // 先设置 LinkID
	}

	// 调用 summaryModel 的自定义 FindSummariesByLinkID 方法
	results, err := r.summaryModel.FindSummariesByLinkID(ctx, uint64(linkId), startDate, endDate)
	if err != nil {
		// 如果是 Not Found，说明这个时间段内没有数据
		if errors.Is(err, model.ErrNotFound) {
			// 仍然需要获取 ShortCode
			linkPO, findErr := r.linksModel.FindOne(ctx, uint64(linkId))
			if findErr == nil {
				analytics.ShortCode = linkPO.ShortCode
			} else if !errors.Is(findErr, model.ErrNotFound) {
				logx.WithContext(ctx).Infof("GetLinkAnalytics (no summary data) failed to get shortcode for linkId %d: %v", linkId, findErr)
			}
			return analytics, nil // 返回空的 analytics 对象 (TotalClicks=0, 列表为空)
		}
		logx.WithContext(ctx).Errorf("GetLinkAnalytics FindSummariesByLinkID error: %v", err)
		return nil, err
	}

	// --- 处理查询结果 ---
	tmpTimeSeries := make(map[time.Time]int64)

	for _, row := range results {
		// 聚合 TotalClicks
		if row.DimensionType == "total" {
			analytics.TotalClicks += int64(row.ClickCount) // 类型转换
		}

		// 聚合 TimeSeries (按天)
		// 假设 'total' 类型的数据按天存储
		if row.DimensionType == "total" {
			day := row.Date // Date 已经是 time.Time
			tmpTimeSeries[day] += int64(row.ClickCount)
		}

		// 分类填充 TopN 列表
		item := entity.AnalyticsStatItem{Key: row.DimensionValue, Value: int64(row.ClickCount)} // 类型转换
		switch row.DimensionType {
		case "referer":
			analytics.TopReferers = append(analytics.TopReferers, item)
		case "country":
			analytics.TopCountries = append(analytics.TopCountries, item)
		case "device":
			analytics.TopDevices = append(analytics.TopDevices, item)
		case "browser":
			analytics.TopBrowsers = append(analytics.TopBrowsers, item)
		case "os":
			analytics.TopOS = append(analytics.TopOS, item)
			// case "timeseries_hourly": // 如果有按小时的数据，在这里处理
		}
	}

	// --- 后处理 ---
	// 转换并排序 TimeSeries
	analytics.TimeSeries = make([]entity.AnalyticsClickPoint, 0, len(tmpTimeSeries))
	for t, v := range tmpTimeSeries {
		analytics.TimeSeries = append(analytics.TimeSeries, entity.AnalyticsClickPoint{Time: t, Value: v})
	}
	sort.Slice(analytics.TimeSeries, func(i, j int) bool {
		return analytics.TimeSeries[i].Time.Before(analytics.TimeSeries[j].Time)
	})

	// 对 TopN 列表进行排序和截断 (取 Top 10)
	limitTopN := func(items []entity.AnalyticsStatItem) []entity.AnalyticsStatItem {
		sort.Slice(items, func(i, j int) bool {
			return items[i].Value > items[j].Value // 按 Value 降序
		})
		if len(items) > 10 {
			return items[:10]
		}
		return items
	}
	analytics.TopReferers = limitTopN(analytics.TopReferers)
	analytics.TopCountries = limitTopN(analytics.TopCountries)
	analytics.TopDevices = limitTopN(analytics.TopDevices)
	analytics.TopBrowsers = limitTopN(analytics.TopBrowsers)
	analytics.TopOS = limitTopN(analytics.TopOS)

	// 获取 ShortCode (需要额外查一次 links 表)
	// 使用 linksModel (它有缓存)
	linkPO, err := r.linksModel.FindOne(ctx, uint64(linkId))
	if err == nil {
		analytics.ShortCode = linkPO.ShortCode
	} else if !errors.Is(err, model.ErrNotFound) { // 找不到不算致命错误
		logx.WithContext(ctx).Infof("GetLinkAnalytics failed to get shortcode for linkId %d: %v", linkId, err)
	}

	return analytics, nil
}

// GetTopLink 获取点击量最高的链接
func (r *analyticsRepoImpl) GetTopLink(ctx context.Context, userId *uint64, startDate, endDate time.Time) (*entity.Link, error) {
	// 调用 summaryModel 的自定义 FindTopClickedLinkID 方法
	linksTableName := r.linksModel.TableName()
	topLinkId, err := r.summaryModel.FindTopClickedLinkID(ctx, userId, startDate, endDate, linksTableName)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) { // model 返回 ErrNotFound
			return nil, entity.ErrLinkNotFound // repo 转换为 entity 的错误
		}
		logx.WithContext(ctx).Errorf("GetTopLink FindTopClickedLinkID error: %v", err)
		return nil, err
	}

	// 调用 linksModel.FindOne 获取 Link 详情 (会走缓存)
	linkPO, err := r.linksModel.FindOne(ctx, topLinkId)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			logx.WithContext(ctx).Errorf("GetTopLink inconsistency: link ID %d found in summary but not in links table", topLinkId)
			return nil, entity.ErrLinkNotFound // 仍然返回 Not Found
		}
		logx.WithContext(ctx).Errorf("GetTopLink FindOne error: %v", err)
		return nil, err
	}

	// ⭐ 使用 link_repo_impl.go 中定义的 fromModel 函数 (需要确保它在该包可见)
	// 或者在这里重新定义一个私有的转换函数
	return fromModelLinks(linkPO), nil
}

// --- 辅助函数 ---

// fromModelLinks 将 links PO 转换为 Link Entity
// (这个函数应该与 link_repo_impl.go 中的 fromModel 保持一致)
func fromModelLinks(m *model.Links) *entity.Link {
	if m == nil {
		return nil
	}
	// 注意：这里直接调用了 fromModel 函数，它定义在同一个包的 link_repo_impl.go 文件中
	// 这是将实现放在同一个包的好处之一 (但也耦合了实现细节)
	return fromModel(m) // 调用 link_repo_impl.go 中的 fromModel
}

// (如果 formatNullTime 需要，也放在这里或公共包)
func formatNullTime(nt sql.NullTime) string {
	if nt.Valid {
		return nt.Time.Format(time.RFC3339)
	}
	return ""
}
