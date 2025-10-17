package shortener

import (
	"context"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UrlAnalyticsModel = (*customUrlAnalyticsModel)(nil)

type (
	// UrlAnalyticsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUrlAnalyticsModel.
	UrlAnalyticsModel interface {
		urlAnalyticsModel
		withSession(session sqlx.Session) UrlAnalyticsModel
		GetTotalClicksAndUniqueVisitorsByShortUrlId(ctx context.Context, shortUrlId uint64) (totalClicks, uniqueVisitors int64, err error)
		GetAnalyticsStatsByShortUrlIds(ctx context.Context, shortUrlIds []uint64) (map[uint64]AnalyticsStats, error)
	}

	customUrlAnalyticsModel struct {
		*defaultUrlAnalyticsModel
	}
)

// NewUrlAnalyticsModel returns a model for the database table.
func NewUrlAnalyticsModel(conn sqlx.SqlConn) UrlAnalyticsModel {
	return &customUrlAnalyticsModel{
		defaultUrlAnalyticsModel: newUrlAnalyticsModel(conn),
	}
}

// clicksAndVisitors 用于存储查询结果
type clicksAndVisitors struct {
	TotalClicks    int64 `db:"total_clicks"`
	UniqueVisitors int64 `db:"unique_visitors"`
}

func (m *customUrlAnalyticsModel) GetTotalClicksAndUniqueVisitorsByShortUrlId(ctx context.Context, shortUrlId uint64) (totalClicks, uniqueVisitors int64, err error) {
	query := fmt.Sprintf("SELECT COUNT(*) AS total_clicks, COUNT(DISTINCT ip_address) AS unique_visitors FROM %s WHERE short_url_id = ?", m.table)

	var result clicksAndVisitors
	err = m.conn.QueryRowCtx(ctx, &result, query, shortUrlId)
	if err != nil {
		return 0, 0, err
	}
	return result.TotalClicks, result.UniqueVisitors, nil
}

type AnalyticsStats struct {
	ShortUrlId     uint64 `db:"short_url_id"`
	TotalClicks    int64  `db:"total_clicks"`
	UniqueVisitors int64  `db:"unique_visitors"`
}

func (m *customUrlAnalyticsModel) GetAnalyticsStatsByShortUrlIds(ctx context.Context, shortUrlIds []uint64) (map[uint64]AnalyticsStats, error) {
	// 如果没有ID，直接返回空map，避免查询数据库
	if len(shortUrlIds) == 0 {
		return make(map[uint64]AnalyticsStats), nil
	}

	// 1. 手动构建 IN 子句的占位符 (?,?,?)
	placeholders := strings.Repeat("?,", len(shortUrlIds)-1) + "?"

	// 2. 构建完整的 SQL 查询语句
	query := fmt.Sprintf(
		"SELECT short_url_id, COUNT(*) AS total_clicks, COUNT(DISTINCT ip_address) AS unique_visitors FROM %s WHERE short_url_id IN (%s) GROUP BY short_url_id",
		m.table,
		placeholders,
	)

	// 3. 将 []uint64 转换为 []interface{} 以便作为可变参数传递
	args := make([]interface{}, len(shortUrlIds))
	for i, id := range shortUrlIds {
		args[i] = id
	}

	var stats []AnalyticsStats
	// 4. 执行查询
	err := m.conn.QueryRowsCtx(ctx, &stats, query, args...)
	if err != nil {
		return nil, err
	}

	// 将查询结果从切片转换为map，方便后续查找
	statsMap := make(map[uint64]AnalyticsStats, len(stats))
	for _, s := range stats {
		statsMap[s.ShortUrlId] = s
	}

	return statsMap, nil
}

func (m *customUrlAnalyticsModel) withSession(session sqlx.Session) UrlAnalyticsModel {
	return NewUrlAnalyticsModel(sqlx.NewSqlConnFromSession(session))
}
