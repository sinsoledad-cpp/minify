package repository

import (
	"context"
	"minify/app/shortener/domain/entity"
)

// LinkAccessLogsRepository 是访问日志仓储的接口
type LinkAccessLogsRepository interface {
	// Create 保存一个新的访问日志
	Create(ctx context.Context, log *entity.LinkAccessLog) error
}
