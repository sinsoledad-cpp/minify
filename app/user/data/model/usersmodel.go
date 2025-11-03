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
// ⭐⭐⭐ 这里是方案二（延迟关联）的实现 ⭐⭐⭐
func (m *customUsersModel) FindAll(ctx context.Context, offset, limit int) ([]*Users, error) {

	// --- 方案二：延迟关联 ---

	// 1. 子查询：仅使用主键索引（PRIMARY）快速定位 ID。
	subQuery := fmt.Sprintf("SELECT `id` FROM %s ORDER BY `id` ASC LIMIT ? OFFSET ?", m.table)

	// 2. 主查询：JOIN 子查询的结果，拉取完整数据。
	//    `usersRows` (来自 usersmodel_gen.go) 包含所有字段。
	//    我们必须在外部再次 ORDER BY。
	query := fmt.Sprintf("SELECT t1.* FROM %s AS t1 JOIN (%s) AS t2 ON t1.id = t2.id ORDER BY t1.id ASC", m.table, subQuery)

	// 3. args (limit, offset) 是给子查询用的。
	// --- 结束 ---

	var resp []*Users
	// ⭐ 注意：这里使用的是 m.conn (来自 defaultUsersModel)，
	// goctl 1.9.1 生成的不带缓存的 model 默认使用 m.conn。
	err := m.conn.QueryRowsCtx(ctx, &resp, query, limit, offset)
	switch err {
	case nil:
		return resp, nil
	case sqlx.ErrNotFound:
		// QueryRowsCtx 在找不到时返回空切片和 nil err，
		// 但以防万一，我们保留这个检查。
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
