package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ LinkAccessLogsModel = (*customLinkAccessLogsModel)(nil)

type (
	// LinkAccessLogsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customLinkAccessLogsModel.
	LinkAccessLogsModel interface {
		linkAccessLogsModel
		withSession(session sqlx.Session) LinkAccessLogsModel
	}

	customLinkAccessLogsModel struct {
		*defaultLinkAccessLogsModel
	}
)

// NewLinkAccessLogsModel returns a model for the database table.
func NewLinkAccessLogsModel(conn sqlx.SqlConn) LinkAccessLogsModel {
	return &customLinkAccessLogsModel{
		defaultLinkAccessLogsModel: newLinkAccessLogsModel(conn),
	}
}

func (m *customLinkAccessLogsModel) withSession(session sqlx.Session) LinkAccessLogsModel {
	return NewLinkAccessLogsModel(sqlx.NewSqlConnFromSession(session))
}
