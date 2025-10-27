package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UsersModel = (*customUsersModel)(nil)

type (
	// UsersModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUsersModel.
	UsersModel interface {
		usersModel
		withSession(session sqlx.Session) UsersModel
		FindAll(ctx context.Context, offset, limit int) ([]*Users, error)
		CountAll(ctx context.Context) (int64, error)
	}

	customUsersModel struct {
		*defaultUsersModel
	}
)

// NewUsersModel returns a model for the database table.
func NewUsersModel(conn sqlx.SqlConn) UsersModel {
	return &customUsersModel{
		defaultUsersModel: newUsersModel(conn),
	}
}

func (m *customUsersModel) withSession(session sqlx.Session) UsersModel {
	return NewUsersModel(sqlx.NewSqlConnFromSession(session))
}

// (新增) FindAll 分页查询所有用户
func (m *customUsersModel) FindAll(ctx context.Context, offset, limit int) ([]*Users, error) {
	query := fmt.Sprintf("select %s from %s order by `id` asc limit ? offset ?", usersRows, m.table)
	var resp []*Users
	err := m.conn.QueryRowsCtx(ctx, &resp, query, limit, offset)
	switch err {
	case nil:
		return resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// (新增) CountAll 查询用户总数
func (m *customUsersModel) CountAll(ctx context.Context) (int64, error) {
	query := fmt.Sprintf("select count(*) from %s", m.table)
	var count int64
	err := m.conn.QueryRowCtx(ctx, &count, query)
	switch err {
	case nil:
		return count, nil
	case sqlx.ErrNotFound:
		return 0, ErrNotFound
	default:
		return 0, err
	}
}
