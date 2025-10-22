package repository

import (
	"context"
	"minify/app/user/domain/entity"
)

// UserRepository 是用户仓储的接口
// 它定义在 Domain 层，将在 Data 层被实现
type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	FindByUsername(ctx context.Context, username string) (*entity.User, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	FindByID(ctx context.Context, id int64) (*entity.User, error)
	// Update(ctx context.Context, user *entity.User) error
}
