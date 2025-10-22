package repository

import (
	"context"
	"errors"
	"minify/app/user/data/model" // 引用 goctl model
	"minify/app/user/domain/entity"

	"github.com/zeromicro/go-zero/core/logx"
)

// 确保 userRepoImpl 实现了 UserRepository 接口
var _ UserRepository = (*userRepoImpl)(nil)

type userRepoImpl struct {
	userModel model.UsersModel // 依赖 goctl model
}

func NewUserRepoImpl(userModel model.UsersModel) UserRepository {
	return &userRepoImpl{
		userModel: userModel,
	}
}

// Create 实现了接口，内部调用 goctl model
func (r *userRepoImpl) Create(ctx context.Context, user *entity.User) error {
	// 调用实体的 ToModel 方法转为 PO
	po := user.ToModel()

	// goctl model 的 Insert 不返回 ID, 并且 CreatedAt/UpdatedAt 由数据库生成
	// 我们需要调整一下 PO，让 goctl model 能够正确插入

	//poToInsert := &model.Users{
	//	Username:     po.Username,
	//	Email:        po.Email,
	//	PasswordHash: po.PasswordHash,
	//	Role:         po.Role,
	//}

	res, err := r.userModel.Insert(ctx, po)
	if err != nil {
		logx.WithContext(ctx).Errorf("userRepoImpl.Create error: %v", err)
		return err
	}

	// 回填 ID
	lastId, _ := res.LastInsertId()
	user.ID = lastId
	return nil
}

// FindByUsername 实现了接口，内部调用 goctl model
func (r *userRepoImpl) FindByUsername(ctx context.Context, username string) (*entity.User, error) {
	// 调用 goctl model
	po, err := r.userModel.FindOneByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, entity.ErrUserNotFound
		}
		logx.WithContext(ctx).Errorf("userRepoImpl.FindByUsername error: %v", err)
		return nil, err
	}

	// 将 PO 转换为 Entity 返回
	return entity.FromModel(po), nil
}

func (r *userRepoImpl) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	po, err := r.userModel.FindOneByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, entity.ErrUserNotFound
		}
		logx.WithContext(ctx).Errorf("userRepoImpl.FindByEmail error: %v", err)
		return nil, err
	}
	return entity.FromModel(po), nil
}

func (r *userRepoImpl) FindByID(ctx context.Context, id int64) (*entity.User, error) {
	po, err := r.userModel.FindOne(ctx, uint64(id)) // 注意类型转换
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, entity.ErrUserNotFound
		}
		logx.WithContext(ctx).Errorf("userRepoImpl.FindByID error: %v", err)
		return nil, err
	}
	return entity.FromModel(po), nil
}
