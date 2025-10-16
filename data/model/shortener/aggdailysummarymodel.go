package shortener

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ AggDailySummaryModel = (*customAggDailySummaryModel)(nil)

type (
	// AggDailySummaryModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAggDailySummaryModel.
	AggDailySummaryModel interface {
		aggDailySummaryModel
		withSession(session sqlx.Session) AggDailySummaryModel
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

func (m *customAggDailySummaryModel) withSession(session sqlx.Session) AggDailySummaryModel {
	return NewAggDailySummaryModel(sqlx.NewSqlConnFromSession(session))
}
