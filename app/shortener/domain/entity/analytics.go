package entity

import (
	"time"
)

// AnalyticsStatItem 表示报表中的一个统计项 (e.g., "Chrome", 1024)
type AnalyticsStatItem struct {
	Key   string
	Value int64
}

// AnalyticsClickPoint 表示时间序列数据点 (e.g., "2025-10-21", 150)
type AnalyticsClickPoint struct {
	Time  time.Time // 使用 time.Time 类型更符合领域语义
	Value int64
}

// LinkAnalytics 封装单个链接的详细分析数据 (领域对象)
type LinkAnalytics struct {
	LinkID       int64  // 聚合表基于 link_id
	ShortCode    string // 需要额外查询 link 表获取
	TotalClicks  int64
	TimeSeries   []AnalyticsClickPoint // 使用领域对象
	TopReferers  []AnalyticsStatItem   // 使用领域对象
	TopCountries []AnalyticsStatItem   // 使用领域对象
	TopDevices   []AnalyticsStatItem   // 使用领域对象
	TopBrowsers  []AnalyticsStatItem   // 使用领域对象
	TopOS        []AnalyticsStatItem   // 使用领域对象
}

// DashboardSummary 封装仪表盘总览数据 (领域对象)
type DashboardSummary struct {
	TotalLinks  int64
	TotalClicks int64
	TopLink     *Link // 引用 Link 实体
}

// ⭐ (新增) AnalyticsDimensions 封装了从事件中解析出的维度
type AnalyticsDimensions struct {
	Referer string
	Country string
	Browser string
	OS      string
	Device  string
	// 还可以添加 'hourly' (小时)
	// Hour string
}

// (新增) 定义维度类型常量，与数据库 `dimension_type` 字段对应
const (
	DimTotal   = "total"
	DimReferer = "referer"
	DimCountry = "country"
	DimBrowser = "browser"
	DimOS      = "os"
	DimDevice  = "device"
)
