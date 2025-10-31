// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.1

package admin

import (
	"context"
	"minify/app/shortener/api/internal/logic/errcode"

	"minify/app/shortener/api/internal/svc"
	"minify/app/shortener/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListAllLinksLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取全站短链接列表 (Admin)
func NewListAllLinksLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListAllLinksLogic {
	return &ListAllLinksLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListAllLinksLogic) ListAllLinks(req *types.ListAllLinksRequest) (resp *types.ListLinksResponse, err error) {
	// 1. (鉴权) Casbin 中间件已经确保了只有 admin 才能访问
	//    我们不需要在这里检查 JWT 的 role，可以直接执行 admin 逻辑

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

	// 3. (关键) 处理 userId 过滤器
	var userIdFilter *uint64
	if req.UserId > 0 {
		// 如果请求中指定了 UserId (不为0)，则将其指针赋给过滤器
		userIdFilter = &req.UserId
	}
	// 如果 req.UserId 为 0 (或未提供)，userIdFilter 保持 nil，仓储层将查询所有用户

	// 4. 调用仓储层(Repository)的新方法 ListGlobal
	links, total, err := l.svcCtx.LinkRepo.ListGlobal(l.ctx, userIdFilter, req.Status, page, pageSize)
	if err != nil {
		l.Logger.Errorf("LinkRepo.ListGlobal error: %v", err)
		return nil, errcode.ErrInternalError // ⭐ 使用 errcode
	}

	// 5. 将领域实体(Entity)列表转换为 DTO 列表
	//    复用 shortener-api 的 Converter
	dtoLinks := make([]types.Link, len(links))
	for i, linkEntity := range links {
		// svcCtx.Converter 是在 logic/converter.go 中定义的
		dtoLink := l.svcCtx.Converter.ToTypesLink(linkEntity)
		if dtoLink != nil {
			dtoLinks[i] = *dtoLink
		}
	}

	// 6. 返回响应 (复用 types.ListLinksResponse)
	return &types.ListLinksResponse{
		Links: dtoLinks,
		Total: total,
	}, nil
}
