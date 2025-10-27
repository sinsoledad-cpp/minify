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

// --- 辅助转换函数 ---

// toModel 将领域实体(Entity)转换为数据模型(PO)
// (这是仓储实现层的私有方法)
func toModel(u *entity.User) *model.Users {
	return &model.Users{
		Id:           uint64(u.ID), // 注意类型转换
		Username:     u.Username,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		Role:         u.Role,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}

// fromModel 将数据模型(PO)转换为领域实体(Entity)
// (这是仓储实现层的私有方法)
func fromModel(m *model.Users) *entity.User {
	if m == nil {
		return nil
	}
	return &entity.User{
		ID:           int64(m.Id), // 注意类型转换
		Username:     m.Username,
		Email:        m.Email,
		PasswordHash: m.PasswordHash,
		Role:         m.Role,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}

// Create 实现了接口，内部调用 goctl model
func (r *userRepoImpl) Create(ctx context.Context, user *entity.User) error {
	// 调用实体的 ToModel 方法转为 PO
	po := toModel(user)

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
	return fromModel(po), nil
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
	return fromModel(po), nil
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
	return fromModel(po), nil
}

// (新增) ListAll 实现
func (r *userRepoImpl) ListAll(ctx context.Context, page, pageSize int) ([]*entity.User, int64, error) {
	// 1. 获取总数
	total, err := r.userModel.CountAll(ctx)
	if err != nil {
		logx.WithContext(ctx).Errorf("userRepoImpl.ListAll CountAll error: %v", err)
		return nil, 0, err
	}
	if total == 0 {
		return []*entity.User{}, 0, nil
	}

	// 2. 计算 offset
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20 // 默认值
	}
	offset := (page - 1) * pageSize

	// 3. 查询分页数据
	pos, err := r.userModel.FindAll(ctx, offset, pageSize)
	if err != nil {
		logx.WithContext(ctx).Errorf("userRepoImpl.ListAll FindAll error: %v", err)
		return nil, 0, err
	}

	// 4. 转换 PO -> Entity
	entities := make([]*entity.User, len(pos))
	for i, po := range pos {
		entities[i] = fromModel(po)
	}

	return entities, total, nil
}
