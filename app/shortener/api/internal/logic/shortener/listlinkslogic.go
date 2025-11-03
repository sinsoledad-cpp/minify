// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.1

package shortener

import (
	"context"
	"errors"
	"minify/common/utils/jwtx"
	"time"

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

	// 2. ⭐ 修改: 处理分页参数
	pageSize := req.PageSize
	switch {
	case pageSize <= 0:
		pageSize = 20 // 默认值
	case pageSize > 100:
		pageSize = 100 // 设置一个合理的最大值
	}
	// 我们总是请求 N+1 条数据，用于判断 "HasMore"
	limit := pageSize + 1

	// 3. ⭐ 关键修改：始终查询总数 (Total)
	// 这样无论前端使用 page 还是 cursor，total 字段始终有值。
	total, err := l.svcCtx.LinkRepo.CountByUser(l.ctx, userId, req.Status)
	if err != nil {
		l.Logger.Errorf("LinkRepo.CountByUser error: %v", err)
		return nil, errors.New("failed to count links")
	}
	if total == 0 {
		// 如果总数为0，提前返回空结果
		return &types.ListLinksResponse{
			Links: []types.Link{},
			Total: 0,
			// HasMore 等默认为 false
		}, nil
	}

	// 4. ⭐ 修改: 处理游标或 Offset
	var lastCreatedAt time.Time
	var lastId uint64
	var offset int

	// 核心决策：判断使用游标还是OFFSET
	useCursor := req.LastId > 0 && req.LastCreatedAt > 0

	if useCursor {
		// --- A. 游标分页路径 (高性能 "下一页") ---
		lastCreatedAt = time.UnixMilli(req.LastCreatedAt).In(time.Local)
		lastId = uint64(req.LastId)
		offset = 0 // 游标模式不使用 offset

	} else {
		// --- B. 传统分页路径 ("第1页" 或 "跳页") ---
		page := req.Page
		if page <= 0 {
			page = 1
		}
		offset = (page - 1) * pageSize
		// (total 已在步骤 3 中查询)
	}

	// 5. 调用仓储层(Repository)
	links, err := l.svcCtx.LinkRepo.ListByUser(l.ctx, userId, req.Status, limit, offset, lastCreatedAt, lastId)
	if err != nil {
		l.Logger.Errorf("LinkRepo.ListByUser error: %v", err)
		return nil, errors.New("failed to list links")
	}

	// 6. 处理返回结果，判断 HasMore
	hasMore := false
	if len(links) == limit {
		// 请求 N+1 条，返回了 N+1 条，说明有更多
		hasMore = true
		// 移除多出来的那一条 (只返回 N 条)
		links = links[:pageSize]
	}

	var nextLastCreatedAt int64
	var nextLastId int64

	// 7. 将领域实体(Entity)列表转换为 DTO 列表
	dtoLinks := make([]types.Link, len(links))
	for i, linkEntity := range links {
		dtoLink := l.svcCtx.Converter.ToTypesLink(linkEntity)
		if dtoLink != nil {
			dtoLinks[i] = *dtoLink
		}
	}

	// 8. 计算下一个游标 (无论哪种模式都计算)
	if len(links) > 0 {
		lastEntity := links[len(links)-1]
		nextLastCreatedAt = lastEntity.CreatedAt.UnixMilli()
		nextLastId = lastEntity.ID
	}

	// 9. ⭐ 修改: 返回响应
	return &types.ListLinksResponse{
		Links:             dtoLinks,
		Total:             total, // ⭐ 解决: total 现在总是有值
		HasMore:           hasMore,
		NextLastCreatedAt: nextLastCreatedAt,
		NextLastId:        nextLastId,
	}, nil
}
