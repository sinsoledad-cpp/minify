package shortener

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ AggDailySummaryModel = (*customAggDailySummaryModel)(nil)

type (
	// AggDailySummaryModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAggDailySummaryModel.
	AggDailySummaryModel interface {
		aggDailySummaryModel
		withSession(session sqlx.Session) AggDailySummaryModel
		FindAllByShortUrlId(ctx context.Context, shortUrlId uint64) ([]*AggDailySummary, error)
	}

	customAggDailySummaryModel struct {
		*defaultAggDailySummaryModel
	}
)

// NewAggDailySummaryModel returns a model for the database table.
func NewAggDailySummaryModel(conn sqlx.SqlConn) AggDailySummaryModel {
	return &customAggDailySummaryModel{
		defaultAggDailySummaryModel: newAggDailySummaryModel(conn),
	}
}

// FindAllByShortUrlId 添加这个新方法
func (m *customAggDailySummaryModel) FindAllByShortUrlId(ctx context.Context, shortUrlId uint64) ([]*AggDailySummary, error) {
	var resp []*AggDailySummary
	// 按日期升序排序
	query := fmt.Sprintf("select %s from %s where `short_url_id` = ? order by `summary_date` asc", aggDailySummaryRows, m.table)

	err := m.conn.QueryRowsCtx(ctx, &resp, query, shortUrlId)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *customAggDailySummaryModel) withSession(session sqlx.Session) AggDailySummaryModel {
	return NewAggDailySummaryModel(sqlx.NewSqlConnFromSession(session))
}
