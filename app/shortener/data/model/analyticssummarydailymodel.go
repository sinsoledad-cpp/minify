package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ AnalyticsSummaryDailyModel = (*customAnalyticsSummaryDailyModel)(nil)

type (
	// AnalyticsSummaryDailyModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAnalyticsSummaryDailyModel.
	AnalyticsSummaryDailyModel interface {
		analyticsSummaryDailyModel
		withSession(session sqlx.Session) AnalyticsSummaryDailyModel
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
