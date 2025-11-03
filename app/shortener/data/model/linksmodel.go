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
		FindListByUserIdAndStatus(ctx context.Context, userId uint64, status string, limit int, offset int, lastCreatedAt time.Time, lastId uint64) ([]*Links, error)
		//FindListByUserIdAndStatus(ctx context.Context, userId uint64, status string, limit, offset int) ([]*Links, error)
		CountByUserIdAndStatus(ctx context.Context, userId uint64, status string) (int64, error)
		RawConn() (sqlx.SqlConn, error)
		TableName() string // ⭐ 接口中定义大写 T
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

// FindListByUserIdAndStatus : 列表查询不走缓存，直接查库 ⭐ 修改: 签名变更，并实现混合逻辑
func (m *customLinksModel) FindListByUserIdAndStatus(ctx context.Context, userId uint64, status string, limit int, offset int, lastCreatedAt time.Time, lastId uint64) ([]*Links, error) {

	// 1. 构建基础 WHERE (不变)
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

	var query string

	// 2. ⭐ 核心：选择分页策略
	// 如果前端提供了游标 (lastId > 0)，则优先使用游标分页
	if lastId > 0 && !lastCreatedAt.IsZero() {
		// 策略一：游标分页 (高性能)
		baseWhere += " AND (created_at < ? OR (created_at = ? AND id < ?))"
		args = append(args, lastCreatedAt, lastCreatedAt, lastId)

		query = fmt.Sprintf("select %s from %s where %s order by created_at DESC, id DESC limit ?", linksRows, m.table, baseWhere)
		args = append(args, limit)

	} else {
		// 策略二：传统 OFFSET 分页 (用于跳页)
		query = fmt.Sprintf("select %s from %s where %s order by created_at DESC, id DESC limit ? offset ?", linksRows, m.table, baseWhere)
		args = append(args, limit, offset)
	}

	// 3. 执行查询
	var resp []*Links
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, args...)
	return resp, err
}

// CountByUserIdAndStatus 计数查询不走缓存，直接查库
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

// 实现 RawConn (通过嵌入的 defaultLinksModel)
func (m *customLinksModel) RawConn() (sqlx.SqlConn, error) {
	if m.rawConn == nil {
		return nil, errors.New("raw connection is not available")
	}
	return m.rawConn, nil
}

// 实现大写的 TableName() 方法 它内部调用嵌入的、小写的 tableName() 方法
func (m *customLinksModel) TableName() string {
	return m.defaultLinksModel.tableName() // 调用嵌入的未导出方法
}

// buildGlobalWhere 构建全局查询的 WHERE 子句
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

// FindListGlobal  列表查询不走缓存，直接查库 这里是方案二（延迟关联）
func (m *customLinksModel) FindListGlobal(ctx context.Context, userId *uint64, status string, limit, offset int) ([]*Links, error) {
	// 调用辅助函数构建查询
	whereClause, args := m.buildGlobalWhere(ctx, userId, status)

	// --- 方案二：延迟关联 ---

	// 1. 子查询：仅使用覆盖索引 `idx_user_status_cursor` 快速定位 ID。
	//    我们使用 `created_at DESC, id DESC` 来 100% 匹配索引顺序。
	subQuery := fmt.Sprintf("SELECT id FROM %s WHERE %s ORDER BY created_at DESC, id DESC LIMIT ? OFFSET ?", m.table, whereClause)

	// 2. 主查询：JOIN 子查询的结果，只拉取目标页的完整数据。
	//    我们必须在外部再次 ORDER BY，以保证最终结果的顺序。
	query := fmt.Sprintf("SELECT t1.* FROM %s AS t1 JOIN (%s) AS t2 ON t1.id = t2.id ORDER BY t1.created_at DESC, t1.id DESC", m.table, subQuery)

	// 3. 将 limit 和 offset 添加到 args (它们是给子查询用的)
	args = append(args, limit, offset)
	// --- 结束 ---

	var resp []*Links
	// 使用 m.QueryRowsNoCacheCtx 直接查询数据库
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, args...)
	return resp, err
}

// CountGlobal  计数查询不走缓存，直接查库
func (m *customLinksModel) CountGlobal(ctx context.Context, userId *uint64, status string) (int64, error) {
	// 调用辅助函数构建查询
	whereClause, args := m.buildGlobalWhere(ctx, userId, status)

	query := fmt.Sprintf("select count(*) from %s where %s", m.table, whereClause)

	var count int64
	// 使用 m.QueryRowNoCacheCtx 直接查询数据库
	err := m.QueryRowNoCacheCtx(ctx, &count, query, args...)
	return count, err
}
