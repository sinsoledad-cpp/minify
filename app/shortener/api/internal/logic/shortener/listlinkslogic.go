// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.1

package shortener

import (
	"context"
	"errors"
	"minify/app/shortener/api/internal/logic"
	"minify/common/utils/jwtx"

	"minify/app/shortener/api/internal/svc"
	"minify/app/shortener/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListLinksLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取短链接列表 (分页)
func NewListLinksLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLinksLogic {
	return &ListLinksLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListLinksLogic) ListLinks(req *types.ListLinksRequest) (resp *types.ListLinksResponse, err error) {
	// 1. 从 JWT Context 获取用户 ID (身份认证)
	claims, err := jwtx.GetClaimsFromCtx(l.ctx)
	if err != nil {
		return nil, errors.New("invalid token")
	}
	userId := uint64(claims.UserID)

	// 2. 处理分页参数
	page := req.Page
	if page <= 0 {
		page = 1
	}

	pageSize := req.PageSize
	switch {
	case pageSize <= 0:
		pageSize = 20 // 默认值
	case pageSize > 100:
		pageSize = 100 // 设置一个合理的最大值
	}
	// 3. 调用仓储层(Repository)
	// 仓储层 link_repo_impl.go 已经实现了按 status 过滤的逻辑
	links, total, err := l.svcCtx.LinkRepo.ListByUser(l.ctx, userId, req.Status, page, pageSize)
	if err != nil {
		l.Logger.Errorf("LinkRepo.ListByUser error: %v", err)
		return nil, errors.New("failed to list links")
	}
	// 4. 将领域实体(Entity)列表转换为 DTO 列表
	// 我们使用你在 logic/converter.go 中定义的 ToTypesLink 转换器
	dtoLinks := make([]types.Link, len(links))
	for i, linkEntity := range links {
		// 注意：converter 在上一级 logic 包中
		dtoLink := logic.ToTypesLink(linkEntity)
		if dtoLink != nil {
			dtoLinks[i] = *dtoLink
		}
	}

	// 5. 返回响应
	return &types.ListLinksResponse{
		Links: dtoLinks,
		Total: total, // total 是符合过滤条件的总记录数，由仓储层 CountByUserIdAndStatus 返回
	}, nil
}
