// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.1

package shortener

import (
	"context"
	"errors"
	"minify/app/shortener/api/internal/logic"
	"minify/app/shortener/domain/entity"
	"minify/common/utils/jwtx"

	"minify/app/shortener/api/internal/svc"
	"minify/app/shortener/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateLinkLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新短链接
func NewUpdateLinkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLinkLogic {
	return &UpdateLinkLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateLinkLogic) UpdateLink(req *types.UpdateLinkRequest) (resp *types.UpdateLinkResponse, err error) {
	// 1. 从 JWT Context 获取用户 ID (身份认证)
	claims, err := jwtx.GetClaimsFromCtx(l.ctx)
	if err != nil {
		return nil, errors.New("invalid token")
	}
	userId := uint64(claims.UserID)

	// 2. 查找链接实体 (获取聚合根)
	link, err := l.svcCtx.LinkRepo.FindByCode(l.ctx, req.Code)
	if err != nil {
		if errors.Is(err, entity.ErrLinkNotFound) {
			// 要更新的链接不存在
			return nil, entity.ErrLinkNotFoundOrForbidden
		}
		// 其他数据库错误
		l.Logger.Errorf("FindByCode error: %v", err)
		return nil, err
	}

	// 3. 检查所有权 (DDD 核心：应用层执行授权策略)
	if link.UserID != userId {
		// 链接存在，但不属于你
		return nil, entity.ErrLinkNotFoundOrForbidden
	}

	// 4. 检查是否为 "No-Op" (空操作)
	// API 定义 中 OriginalUrl 是 string,
	// 而 IsActive 和 ExpirationTime 是指针。
	// 我们假设 OriginalUrl 为空字符串时代表 "不更新"。
	isNoOp := req.OriginalUrl == "" && req.IsActive == nil && req.ExpirationTime == nil
	if isNoOp {
		// 如果没有提供任何可更新的字段，直接返回当前链接信息
		return &types.UpdateLinkResponse{
			Link: *logic.ToTypesLink(link), // ⭐ 使用 converter 转换
		}, nil
	}

	// 5. 委托领域实体执行业务规则 (Domain Logic)
	//

	// 处理 OriginalUrl：API 类型为 string，
	// 实体方法 UpdateDetails 期望 *string
	var urlToUpdate *string
	if req.OriginalUrl != "" {
		urlToUpdate = &req.OriginalUrl
	}

	err = link.UpdateDetails(urlToUpdate, req.IsActive, req.ExpirationTime)
	if err != nil {
		// 捕获实体返回的业务验证错误 (例如：无效的 URL 格式或过期时间格式)
		//
		l.Logger.Infof("UpdateDetails validation error: %v", err)
		// 返回一个对用户友好的错误，可以考虑使用 httpx.NewCodeError
		return nil, err
	}

	// 6. 调用仓储持久化 (Infrastructure)
	// link 实体现在是 "脏" 的 (已修改)，将其保存到 repo
	// Update 方法会处理数据库更新和缓存失效
	//
	if err := l.svcCtx.LinkRepo.Update(l.ctx, link); err != nil {
		l.Logger.Errorf("LinkRepo.Update error: %v", err)
		return nil, errors.New("failed to update link")
	}

	// 7. 返回更新后的 DTO
	return &types.UpdateLinkResponse{
		Link: *logic.ToTypesLink(link), // ⭐ 使用 converter 转换
	}, nil
}
