package repository

import (
	"context"
	"minify/app/shortener/domain/entity" // 引用 Link 实体
)

// LinkRepository 是短链接仓储的接口
type LinkRepository interface {
	Create(ctx context.Context, link *entity.Link) error
	// FindByCode 根据 ShortCode 查找 Link (应包含缓存逻辑)
	FindByCode(ctx context.Context, code string) (*entity.Link, error)
	FindByID(ctx context.Context, id int64) (*entity.Link, error)
	// ListByUser 根据用户ID和状态分页查询 Link 列表
	// 返回: 链接列表, 符合条件的总数, 错误
	ListByUser(ctx context.Context, userId uint64, status string, page, pageSize int) ([]*entity.Link, int64, error)
	Update(ctx context.Context, link *entity.Link) error
	// Delete 执行软删除
	Delete(ctx context.Context, link *entity.Link) error
	// ListGlobal 按可选的用户ID和状态分页查询 (nil userId = 查询所有)
	// 返回: 链接列表, 符合条件的总数, 错误
	ListGlobal(ctx context.Context, userId *uint64, status string, page, pageSize int) ([]*entity.Link, int64, error)
}
