package shortener

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ UrlAnalyticsModel = (*customUrlAnalyticsModel)(nil)

type (
	// UrlAnalyticsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUrlAnalyticsModel.
	UrlAnalyticsModel interface {
		urlAnalyticsModel
		withSession(session sqlx.Session) UrlAnalyticsModel
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

func (m *customUrlAnalyticsModel) withSession(session sqlx.Session) UrlAnalyticsModel {
	return NewUrlAnalyticsModel(sqlx.NewSqlConnFromSession(session))
}
