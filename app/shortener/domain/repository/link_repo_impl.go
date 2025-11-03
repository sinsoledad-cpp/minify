package repository

import (
	"context"
	"errors"
	"minify/app/shortener/data/model" // ⭐ 依赖 data/model
	"minify/app/shortener/domain/entity"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

// 确保 linkRepoImpl 实现了 LinkRepository 接口 (接口定义在 link_repository.go)
var _ LinkRepository = (*linkRepoImpl)(nil)

type linkRepoImpl struct {
	linksModel model.LinksModel // ⭐ 只依赖 model 接口
}

// NewLinkRepoImpl 创建 LinkRepository 的实现
func NewLinkRepoImpl(linksModel model.LinksModel) LinkRepository { // ⭐ 返回接口类型
	return &linkRepoImpl{
		linksModel: linksModel,
	}
}

// --- ⭐ 转换辅助函数 (放在实现文件中) ---
func toModel(e *entity.Link) *model.Links {
	isActiveInt := int64(0)
	if e.IsActive {
		isActiveInt = 1
	}
	poID := uint64(0)
	if e.ID > 0 {
		poID = uint64(e.ID)
	}
	return &model.Links{
		Id:             poID,
		UserId:         e.UserID,
		ShortCode:      e.ShortCode,
		OriginalUrl:    e.OriginalUrl,
		IsActive:       isActiveInt,
		ExpirationTime: e.ExpirationTime,
		CreatedAt:      e.CreatedAt,
		UpdatedAt:      e.UpdatedAt,
		DeletedAt:      e.DeletedAt,
	}
}

func fromModel(m *model.Links) *entity.Link {
	return &entity.Link{
		ID:             int64(m.Id),
		UserID:         m.UserId,
		ShortCode:      m.ShortCode,
		OriginalUrl:    m.OriginalUrl,
		IsActive:       m.IsActive == 1,
		ExpirationTime: m.ExpirationTime,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
		DeletedAt:      m.DeletedAt,
	}
}

// --- 接口实现 ---

func (r *linkRepoImpl) Create(ctx context.Context, link *entity.Link) error {
	po := toModel(link) // ⭐ 使用辅助函数
	po.Id = 0
	res, err := r.linksModel.Insert(ctx, po) // ⭐ 调用 model.Insert
	if err != nil {
		logx.WithContext(ctx).Errorf("linkRepoImpl.Create error: %v", err)
		return err
	}
	lastId, _ := res.LastInsertId()
	link.ID = lastId // 回填 Entity
	return nil
}

func (r *linkRepoImpl) FindByCode(ctx context.Context, code string) (*entity.Link, error) {
	po, err := r.linksModel.FindOneByShortCode(ctx, code) // ⭐ 调用 model.FindOneByShortCode
	if err != nil {
		if errors.Is(err, model.ErrNotFound) { // ⭐ 比较 model 错误
			return nil, entity.ErrLinkNotFound // ⭐ 返回 entity 错误
		}
		logx.WithContext(ctx).Errorf("linkRepoImpl.FindByCode error: %v", err)
		return nil, err
	}
	if po.DeletedAt.Valid {
		return nil, entity.ErrLinkNotFound // ⭐ 返回 entity 错误
	}
	return fromModel(po), nil // ⭐ 使用辅助函数
}

func (r *linkRepoImpl) FindByID(ctx context.Context, id int64) (*entity.Link, error) {
	po, err := r.linksModel.FindOne(ctx, uint64(id)) // ⭐ 调用 model.FindOne
	if err != nil {
		if errors.Is(err, model.ErrNotFound) { // ⭐ 比较 model 错误
			return nil, entity.ErrLinkNotFound // ⭐ 返回 entity 错误
		}
		logx.WithContext(ctx).Errorf("linkRepoImpl.FindByID error: %v", err)
		return nil, err
	}
	if po.DeletedAt.Valid {
		return nil, entity.ErrLinkNotFound // ⭐ 返回 entity 错误
	}
	return fromModel(po), nil // ⭐ 使用辅助函数
}

// ListByUser ⭐ 修改: ListByUser 签名变更
func (r *linkRepoImpl) ListByUser(ctx context.Context, userId uint64, status string, limit int, offset int, lastCreatedAt time.Time, lastId uint64) ([]*entity.Link, error) {

	// ⭐ 修改: 调用新的 model 方法 (传递所有参数)
	pos, err := r.linksModel.FindListByUserIdAndStatus(ctx, userId, status, limit, offset, lastCreatedAt, lastId)
	if err != nil {
		logx.WithContext(ctx).Errorf("linkRepoImpl.ListByUser FindList error: %v", err)
		if errors.Is(err, model.ErrNotFound) {
			return []*entity.Link{}, nil // ⭐ 返回空列表，不是错误
		}
		return nil, err // ⭐ 返回真实错误
	}

	links := make([]*entity.Link, len(pos))
	for i, po := range pos {
		links[i] = fromModel(po) // ⭐ 使用辅助函数
	}

	return links, nil
}

// CountByUser ⭐ 新增: CountByUser 实现
func (r *linkRepoImpl) CountByUser(ctx context.Context, userId uint64, status string) (int64, error) {
	total, err := r.linksModel.CountByUserIdAndStatus(ctx, userId, status)
	if err != nil {
		logx.WithContext(ctx).Errorf("linkRepoImpl.CountByUser error: %v", err)
		if errors.Is(err, model.ErrNotFound) {
			return 0, nil // 找不到算 0
		}
		return 0, err
	}
	return total, nil
}

func (r *linkRepoImpl) Update(ctx context.Context, link *entity.Link) error {
	po := toModel(link)                 // ⭐ 使用辅助函数
	err := r.linksModel.Update(ctx, po) // ⭐ 调用 model.Update
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return entity.ErrLinkNotFound
		} // ⭐ 返回 entity 错误
		logx.WithContext(ctx).Errorf("linkRepoImpl.Update error: %v", err)
		return err
	}
	return nil
}

func (r *linkRepoImpl) Delete(ctx context.Context, link *entity.Link) error {
	link.MarkDeleted()
	return r.Update(ctx, link) // 调用自身的 Update
}
func (r *linkRepoImpl) ListGlobal(ctx context.Context, userId *uint64, status string, page, pageSize int) ([]*entity.Link, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	// ⭐ 1. 调用 model.CountGlobal (我们刚在 linksmodel.go 中创建的)
	total, err := r.linksModel.CountGlobal(ctx, userId, status)
	if err != nil {
		logx.WithContext(ctx).Errorf("linkRepoImpl.ListGlobal Count error: %v", err)
		if errors.Is(err, model.ErrNotFound) {
			return []*entity.Link{}, 0, nil
		}
		return nil, 0, err
	}
	if total == 0 {
		return []*entity.Link{}, 0, nil
	}

	// ⭐ 2. 调用 model.FindListGlobal (我们刚在 linksmodel.go 中创建的)
	pos, err := r.linksModel.FindListGlobal(ctx, userId, status, pageSize, offset)
	if err != nil {
		logx.WithContext(ctx).Errorf("linkRepoImpl.ListGlobal FindList error: %v", err)
		if errors.Is(err, model.ErrNotFound) {
			return []*entity.Link{}, total, nil
		}
		return nil, 0, err
	}

	links := make([]*entity.Link, len(pos))
	for i, po := range pos {
		links[i] = fromModel(po) // 复用 fromModel 转换函数
	}
	return links, total, nil
}
