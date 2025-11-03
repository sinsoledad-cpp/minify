package repository

import (
	"context"
	"minify/app/shortener/domain/entity" // 引用 Link 实体
	"time"
)

// LinkRepository 是短链接仓储的接口
type LinkRepository interface {
	Create(ctx context.Context, link *entity.Link) error
	// FindByCode 根据 ShortCode 查找 Link (应包含缓存逻辑)
	FindByCode(ctx context.Context, code string) (*entity.Link, error)
	FindByID(ctx context.Context, id int64) (*entity.Link, error)
	//// ListByUser 根据用户ID和状态分页查询 Link 列表
	//// 返回: 链接列表, 符合条件的总数, 错误
	//ListByUser(ctx context.Context, userId uint64, status string, page, pageSize int) ([]*entity.Link, int64, error)

	// ListByUser 根据用户ID和状态分页查询 Link 列表
	// ⭐ 修改: 签名变更
	// 同时接收 offset (用于 page) 和 cursor (lastCreatedAt/lastId)
	// 返回: 链接列表, 错误 (不再返回 total)
	ListByUser(ctx context.Context, userId uint64, status string, limit int, offset int, lastCreatedAt time.Time, lastId uint64) ([]*entity.Link, error)

	// ⭐ 新增: CountByUser (从 ListByUser 中拆分出来，供 Logic 单独调用)
	CountByUser(ctx context.Context, userId uint64, status string) (int64, error)

	Update(ctx context.Context, link *entity.Link) error
	// Delete 执行软删除
	Delete(ctx context.Context, link *entity.Link) error
	// ListGlobal 按可选的用户ID和状态分页查询 (nil userId = 查询所有)
	// 返回: 链接列表, 符合条件的总数, 错误
	ListGlobal(ctx context.Context, userId *uint64, status string, page, pageSize int) ([]*entity.Link, int64, error)
}
