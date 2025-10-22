package model

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ AnalyticsSummaryDailyModel = (*customAnalyticsSummaryDailyModel)(nil)

type (
	// AnalyticsSummaryDailyModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAnalyticsSummaryDailyModel.
	AnalyticsSummaryDailyModel interface {
		analyticsSummaryDailyModel
		withSession(session sqlx.Session) AnalyticsSummaryDailyModel
		FindSummariesByLinkID(ctx context.Context, linkId uint64, startDate, endDate time.Time) ([]*AnalyticsSummaryDaily, error)
		FindTotalClicks(ctx context.Context, userId *uint64, startDate, endDate time.Time, linksTable string) (int64, error)
		FindTopClickedLinkID(ctx context.Context, userId *uint64, startDate, endDate time.Time, linksTable string) (uint64, error)
		TableName() string // ⭐ 暴露 TableName 接口
	}

	customAnalyticsSummaryDailyModel struct {
		*defaultAnalyticsSummaryDailyModel
	}
)

// NewAnalyticsSummaryDailyModel returns a model for the database table.
func NewAnalyticsSummaryDailyModel(conn sqlx.SqlConn) AnalyticsSummaryDailyModel {
	return &customAnalyticsSummaryDailyModel{
		defaultAnalyticsSummaryDailyModel: newAnalyticsSummaryDailyModel(conn),
	}
}

func (m *customAnalyticsSummaryDailyModel) withSession(session sqlx.Session) AnalyticsSummaryDailyModel {
	return NewAnalyticsSummaryDailyModel(sqlx.NewSqlConnFromSession(session))
}

// ⭐ 实现自定义方法：获取指定链接的聚合数据
func (m *customAnalyticsSummaryDailyModel) FindSummariesByLinkID(ctx context.Context, linkId uint64, startDate, endDate time.Time) ([]*AnalyticsSummaryDaily, error) {
	// 查询聚合表 analytics_summaries_daily
	// 使用 goctl 生成的字段列表 analyticsSummaryDailyRows
	query := fmt.Sprintf(`SELECT %s
                 FROM %s
                 WHERE link_id = ? AND date BETWEEN ? AND ?`, analyticsSummaryDailyRows, m.table) // m.table 来自嵌入的 default model

	var resp []*AnalyticsSummaryDaily
	// 使用 m.conn (来自嵌入的 default model) 执行查询
	err := m.conn.QueryRowsCtx(ctx, &resp, query, linkId, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	// QueryRowsCtx 在没有找到行时不返回 ErrNotFound，而是返回一个空切片和 nil error
	// 所以我们不需要特殊处理 ErrNotFound
	return resp, err
}

// ⭐ 实现自定义方法：获取总点击量 (Dashboard)
func (m *customAnalyticsSummaryDailyModel) FindTotalClicks(ctx context.Context, userId *uint64, startDate, endDate time.Time, linksTable string) (int64, error) {
	var totalClicks int64
	var err error
	// 使用 DATE_FORMAT 确保比较的是日期部分
	clickWhere := "`date` BETWEEN DATE(?) AND DATE(?)"
	clickArgs := []interface{}{startDate.Format("2006-01-02"), endDate.Format("2006-01-02")}

	if userId != nil {
		// 需要 JOIN links 表
		clickQuery := fmt.Sprintf(`SELECT COALESCE(SUM(asd.click_count), 0)
                         FROM %s asd JOIN %s l ON asd.link_id = l.id
                         WHERE l.user_id = ? AND asd.dimension_type = 'total' AND asd.%s`,
			m.table, linksTable, clickWhere)
		clickArgs = append([]interface{}{*userId}, clickArgs...)
		err = m.conn.QueryRowCtx(ctx, &totalClicks, clickQuery, clickArgs...)
	} else {
		// 查全局
		clickQuery := fmt.Sprintf("SELECT COALESCE(SUM(click_count), 0) FROM %s WHERE dimension_type = 'total' AND %s",
			m.table, clickWhere)
		err = m.conn.QueryRowCtx(ctx, &totalClicks, clickQuery, clickArgs...)
	}

	// QueryRowCtx 在 COUNT/SUM 为 0 时不返回 ErrNotFound
	if err != nil && !errors.Is(err, sql.ErrNoRows) && !errors.Is(err, sqlx.ErrNotFound) { // 但以防万一还是检查下
		return 0, err
	}
	// 如果确实没找到行 (理论上 SUM 不会，但 COUNT(*) 会)，返回 0
	if errors.Is(err, sql.ErrNoRows) || errors.Is(err, sqlx.ErrNotFound) {
		return 0, nil
	}
	return totalClicks, nil
}

// ⭐ 实现自定义方法：获取 Top Link ID (Dashboard)
func (m *customAnalyticsSummaryDailyModel) FindTopClickedLinkID(ctx context.Context, userId *uint64, startDate, endDate time.Time, linksTable string) (uint64, error) {
	var topLinkId uint64
	// 使用 DATE_FORMAT
	dateWhere := "`date` BETWEEN DATE(?) AND DATE(?)"
	args := []interface{}{startDate.Format("2006-01-02"), endDate.Format("2006-01-02")}

	query := fmt.Sprintf(`SELECT link_id FROM %s
                 WHERE dimension_type = 'total' AND %s
                 GROUP BY link_id
                 ORDER BY SUM(click_count) DESC
                 LIMIT 1`, m.table, dateWhere)

	if userId != nil {
		query = fmt.Sprintf(`SELECT asd.link_id FROM %s asd JOIN %s l ON asd.link_id = l.id
                    WHERE l.user_id = ? AND asd.dimension_type = 'total' AND asd.%s
                    GROUP BY asd.link_id
                    ORDER BY SUM(asd.click_count) DESC
                    LIMIT 1`, m.table, linksTable, dateWhere)
		args = append([]interface{}{*userId}, args...)
	}

	err := m.conn.QueryRowCtx(ctx, &topLinkId, query, args...)
	// QueryRowCtx 在 LIMIT 1 且没有找到行时会返回 ErrNotFound
	// 我们直接将这个错误（包括 ErrNotFound）返回给调用者 (repository) 处理
	return topLinkId, err
}

// ⭐ 暴露 TableName (goctl 已在 default model 中生成了小写的 tableName 方法)
func (m *customAnalyticsSummaryDailyModel) TableName() string {
	return m.table // 直接返回嵌入的 default model 的 table 字段
}
