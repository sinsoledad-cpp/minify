package shortener

import (
	"context"
	"fmt"

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

func (m *customUrlAnalyticsModel) withSession(session sqlx.Session) UrlAnalyticsModel {
	return NewUrlAnalyticsModel(sqlx.NewSqlConnFromSession(session))
}
