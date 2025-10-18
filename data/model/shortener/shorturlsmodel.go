package shortener

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ShortUrlsModel = (*customShortUrlsModel)(nil)

type (
	// ShortUrlsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customShortUrlsModel.
	ShortUrlsModel interface {
		shortUrlsModel
		withSession(session sqlx.Session) ShortUrlsModel
		FindAllByUserId(ctx context.Context, userId uint64) ([]*ShortUrls, error)
		FindPagedByUserId(ctx context.Context, userId uint64, page, pageSize int) ([]*ShortUrls, int64, error)
	}

	customShortUrlsModel struct {
		*defaultShortUrlsModel
	}
)

// NewShortUrlsModel returns a model for the database table.
func NewShortUrlsModel(conn sqlx.SqlConn) ShortUrlsModel {
	return &customShortUrlsModel{
		defaultShortUrlsModel: newShortUrlsModel(conn),
	}
}

func (m *customShortUrlsModel) FindAllByUserId(ctx context.Context, userId uint64) ([]*ShortUrls, error) {
	var resp []*ShortUrls
	query := fmt.Sprintf("select %s from %s where `user_id` = ? and `deleted_at` IS NULL order by `created_at` desc", shortUrlsRows, m.table)
	err := m.conn.QueryRowsCtx(ctx, &resp, query, userId)
	switch err {
	case nil:
		return resp, nil
	case sqlx.ErrNotFound:
		return nil, nil // 没有找到，返回空列表而不是错误
	default:
		return nil, err
	}
}

func (m *customShortUrlsModel) FindPagedByUserId(ctx context.Context, userId uint64, page, pageSize int) ([]*ShortUrls, int64, error) {
	// 1. 查询总数
	var total int64
	countQuery := fmt.Sprintf("select count(*) from %s where `user_id` = ? and `deleted_at` IS NULL", m.table)
	err := m.conn.QueryRowCtx(ctx, &total, countQuery, userId)
	if err != nil {
		return nil, 0, err
	}

	// 如果总数为0，直接返回空列表，避免无效的列表查询
	if total == 0 {
		return []*ShortUrls{}, 0, nil
	}

	// 2. 查询分页数据
	var resp []*ShortUrls
	offset := (page - 1) * pageSize
	query := fmt.Sprintf("select %s from %s where `user_id` = ? and `deleted_at` IS NULL order by `created_at` desc LIMIT ? OFFSET ?", shortUrlsRows, m.table)
	err = m.conn.QueryRowsCtx(ctx, &resp, query, userId, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}

	return resp, total, nil
}

func (m *customShortUrlsModel) withSession(session sqlx.Session) ShortUrlsModel {
	return NewShortUrlsModel(sqlx.NewSqlConnFromSession(session))
}
