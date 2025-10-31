package model

import (
	"context"
	"errors"
	"fmt"
	"minify/app/shortener/domain/entity"
	"time"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ LinksModel = (*customLinksModel)(nil)

type (
	// LinksModel is an interface to be customized, add more methods here,
	// and implement the added methods in customLinksModel.
	LinksModel interface {
		linksModel
		FindListByUserIdAndStatus(ctx context.Context, userId uint64, status string, limit, offset int) ([]*Links, error)
		CountByUserIdAndStatus(ctx context.Context, userId uint64, status string) (int64, error)
		RawConn() (sqlx.SqlConn, error)
		TableName() string // ⭐ 接口中定义大写 T
		// --- (新增) ---
		// FindListGlobal 按可选的 userId 和 status 分页查询 (nil userId = 查询所有)
		FindListGlobal(ctx context.Context, userId *uint64, status string, limit, offset int) ([]*Links, error)
		// CountGlobal 按可选的 userId 和 status 统计 (nil userId = 查询所有)
		CountGlobal(ctx context.Context, userId *uint64, status string) (int64, error)
	}

	customLinksModel struct {
		*defaultLinksModel
		rawConn sqlx.SqlConn
	}
)

// NewLinksModel returns a model for the database table.
func NewLinksModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) LinksModel {
	return &customLinksModel{
		defaultLinksModel: newLinksModel(conn, c, opts...),
		rawConn:           conn,
	}
}

// FindListByUserIdAndStatus: 列表查询不走缓存，直接查库
func (m *customLinksModel) FindListByUserIdAndStatus(ctx context.Context, userId uint64, status string, limit, offset int) ([]*Links, error) {
	// ... (构建 WHERE 条件和 args 的逻辑不变)
	baseWhere := "user_id = ? AND deleted_at IS NULL"
	args := []interface{}{userId}
	now := time.Now()
	switch status {
	case entity.StatusActive, "":
		baseWhere += " AND is_active = 1 AND (expiration_time IS NULL OR expiration_time > ?)"
		args = append(args, now)
	case entity.StatusExpired:
		baseWhere += " AND expiration_time IS NOT NULL AND expiration_time <= ?"
		args = append(args, now)
	case entity.StatusInactive:
		baseWhere += " AND is_active = 0"
	case entity.StatusAll:
	default:
		baseWhere += " AND is_active = 1 AND (expiration_time IS NULL OR expiration_time > ?)"
		args = append(args, now)
	}

	query := fmt.Sprintf("select %s from %s where %s order by created_at desc limit ? offset ?", linksRows, m.table, baseWhere) // m.table 来自 embedded defaultLinksModel
	args = append(args, limit, offset)

	var resp []*Links
	// ⭐ 使用 m.QueryRowsNoCacheCtx 直接查询数据库
	// defaultLinksModel 嵌入了 CachedConn，它提供了 QueryRowsNoCacheCtx 方法
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, args...)
	return resp, err
}

// CountByUserIdAndStatus: 计数查询不走缓存，直接查库
func (m *customLinksModel) CountByUserIdAndStatus(ctx context.Context, userId uint64, status string) (int64, error) {
	// ... (构建 WHERE 条件和 args 的逻辑不变)
	baseWhere := "user_id = ? AND deleted_at IS NULL"
	args := []interface{}{userId}
	now := time.Now()
	switch status {
	case entity.StatusActive, "":
		baseWhere += " AND is_active = 1 AND (expiration_time IS NULL OR expiration_time > ?)"
		args = append(args, now)
	case entity.StatusExpired:
		baseWhere += " AND expiration_time IS NOT NULL AND expiration_time <= ?"
		args = append(args, now)
	case entity.StatusInactive:
		baseWhere += " AND is_active = 0"
	case entity.StatusAll:
	default:
		baseWhere += " AND is_active = 1 AND (expiration_time IS NULL OR expiration_time > ?)"
		args = append(args, now)
	}

	query := fmt.Sprintf("select count(*) from %s where %s", m.table, baseWhere)

	var count int64
	// ⭐ 使用 m.QueryRowNoCacheCtx 直接查询数据库
	err := m.QueryRowNoCacheCtx(ctx, &count, query, args...)
	return count, err
}

// ⭐ 实现 RawConn (通过嵌入的 defaultLinksModel)
func (m *customLinksModel) RawConn() (sqlx.SqlConn, error) {
	if m.rawConn == nil {
		return nil, errors.New("raw connection is not available")
	}
	return m.rawConn, nil
}

// ⭐ 实现大写的 TableName() 方法 它内部调用嵌入的、小写的 tableName() 方法
func (m *customLinksModel) TableName() string {
	return m.defaultLinksModel.tableName() // 调用嵌入的未导出方法
}

// (新增) buildGlobalWhere 构建全局查询的 WHERE 子句
func (m *customLinksModel) buildGlobalWhere(ctx context.Context, userId *uint64, status string) (string, []interface{}) {
	// 基础条件：未软删除
	baseWhere := "deleted_at IS NULL"
	args := []interface{}{}

	// 1. (关键) 如果 userId 不是 nil，才添加 user_id 过滤
	if userId != nil {
		baseWhere += " AND user_id = ?"
		args = append(args, *userId)
	}

	// 2. 状态过滤 (逻辑与 ListByUser 保持一致)
	now := time.Now()
	switch status {
	case entity.StatusActive, "":
		baseWhere += " AND is_active = 1 AND (expiration_time IS NULL OR expiration_time > ?)"
		args = append(args, now)
	case entity.StatusExpired:
		baseWhere += " AND expiration_time IS NOT NULL AND expiration_time <= ?"
		args = append(args, now)
	case entity.StatusInactive:
		baseWhere += " AND is_active = 0"
	case entity.StatusAll:
		// all = 只过滤 deleted_at，所以不需要额外条件
	default:
		// 默认等同于 active
		baseWhere += " AND is_active = 1 AND (expiration_time IS NULL OR expiration_time > ?)"
		args = append(args, now)
	}

	return baseWhere, args
}

// (新增) FindListGlobal: 列表查询不走缓存，直接查库
func (m *customLinksModel) FindListGlobal(ctx context.Context, userId *uint64, status string, limit, offset int) ([]*Links, error) {
	// 调用辅助函数构建查询
	whereClause, args := m.buildGlobalWhere(ctx, userId, status)

	query := fmt.Sprintf("select %s from %s where %s order by created_at desc limit ? offset ?", linksRows, m.table, whereClause)
	args = append(args, limit, offset)

	var resp []*Links
	// 使用 m.QueryRowsNoCacheCtx 直接查询数据库
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, args...)
	return resp, err
}

// (新增) CountGlobal: 计数查询不走缓存，直接查库
func (m *customLinksModel) CountGlobal(ctx context.Context, userId *uint64, status string) (int64, error) {
	// 调用辅助函数构建查询
	whereClause, args := m.buildGlobalWhere(ctx, userId, status)

	query := fmt.Sprintf("select count(*) from %s where %s", m.table, whereClause)

	var count int64
	// 使用 m.QueryRowNoCacheCtx 直接查询数据库
	err := m.QueryRowNoCacheCtx(ctx, &count, query, args...)
	return count, err
}
