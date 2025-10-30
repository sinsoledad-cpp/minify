package logic

import (
	"minify/app/shortener/api/internal/types"
	"minify/app/shortener/domain/entity"
	"strings"
	"time"
)

// ⭐ 2. 定义 Converter 结构体，它持有依赖
type Converter struct {
	ShortDomain string
}

// ⭐ 3. 提供一个 NewConverter 的工厂函数
func NewConverter(shortDomain string) *Converter {
	// 确保域名最后没有 /
	domain := strings.TrimSuffix(shortDomain, "/")
	return &Converter{
		ShortDomain: domain,
	}
}

// ToTypesLink 将 Link 实体转换为 Link DTO
func (c *Converter) ToTypesLink(e *entity.Link) *types.Link {
	if e == nil {
		return nil
	}

	expTime := ""
	if e.ExpirationTime.Valid {
		// 格式化为 ISO 8601，与 API 定义一致
		expTime = e.ExpirationTime.Time.Format(time.RFC3339)
	}

	return &types.Link{
		Id:             e.ID,
		ShortCode:      c.ShortDomain + "/" + e.ShortCode,
		OriginalUrl:    e.OriginalUrl,
		IsActive:       e.IsActive,
		ExpirationTime: expTime,
		CreatedAt:      e.CreatedAt.Format(time.RFC3339),
	}
}

// ParseAnalyticsDates 解析 API 请求中的日期字符串，并提供默认值
// (此函数设为导出，供 shortener 包中的 logic 调用)
func ParseAnalyticsDates(reqStartDate, reqEndDate string) (time.Time, time.Time, error) {
	var startDate, endDate time.Time
	var err error
	// 1. 获取时区 (基于 shortener.go 的设置, time.Local 已被设置为 "Asia/Shanghai")
	loc := time.Local
	// 2. 解析结束日期
	if reqEndDate == "" {
		// 默认到今天 (的 00:00:00)
		endDate = time.Now().In(loc)
	} else {
		// API 定义为 "ISO 8601 日期"，我们按 "YYYY-MM-DD" 格式解析
		endDate, err = time.ParseInLocation(time.DateOnly, reqEndDate, loc)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
	}
	// 3. 解析开始日期
	if reqStartDate == "" {
		// 默认 30 天前
		startDate = endDate.AddDate(0, 0, -30)
	} else {
		startDate, err = time.ParseInLocation(time.DateOnly, reqStartDate, loc)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
	}
	// 4. 确保 startDate <= endDate
	if startDate.After(endDate) {
		startDate = endDate.AddDate(0, 0, -30) // 如果开始日期晚于结束日期，重置为 30 天前
	}
	return startDate, endDate, nil
}

// ToTypesAnalyticsResponse 将 LinkAnalytics 实体转换为 GetAnalyticsResponse DTO
// (此函数设为导出，供 shortener 包中的 logic 调用)
func ToTypesAnalyticsResponse(e *entity.LinkAnalytics) *types.GetAnalyticsResponse {
	if e == nil {
		return &types.GetAnalyticsResponse{} // 返回一个空响应，而不是 nil
	}

	return &types.GetAnalyticsResponse{
		ShortCode:    e.ShortCode,
		TotalClicks:  e.TotalClicks,
		TimeSeries:   toTypesClickPointSlice(e.TimeSeries), // 调用本包内的 unexported 函数
		TopReferers:  toTypesStatItemSlice(e.TopReferers),  // 调用本包内的 unexported 函数
		TopCountries: toTypesStatItemSlice(e.TopCountries), // ...
		TopDevices:   toTypesStatItemSlice(e.TopDevices),   // ...
		TopBrowsers:  toTypesStatItemSlice(e.TopBrowsers),  // ...
		TopOS:        toTypesStatItemSlice(e.TopOS),        // ...
	}
}

// toTypesClickPointSlice 转换时间序列 (unexported, 仅供本包使用)
func toTypesClickPointSlice(items []entity.AnalyticsClickPoint) []types.ClickPoint {
	if items == nil {
		return []types.ClickPoint{} // 返回空切片
	}
	dtos := make([]types.ClickPoint, len(items))
	for i, item := range items {
		dtos[i] = types.ClickPoint{
			// 仓储层返回的 Time 已经是 time.Time，我们将其格式化为 YYYY-MM-DD
			Time:  item.Time.Format(time.DateOnly),
			Value: item.Value,
		}
	}
	return dtos
}

// toTypesStatItemSlice 转换统计项 (unexported, 仅供本包使用)
func toTypesStatItemSlice(items []entity.AnalyticsStatItem) []types.StatItem {
	if items == nil {
		return []types.StatItem{} // 返回空切片
	}
	dtos := make([]types.StatItem, len(items))
	for i, item := range items {
		dtos[i] = types.StatItem{
			Key:   item.Key,
			Value: item.Value,
		}
	}
	return dtos
}
