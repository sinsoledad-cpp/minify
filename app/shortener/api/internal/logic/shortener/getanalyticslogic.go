// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.1

package shortener

import (
	"context"
	"errors"
	"minify/app/shortener/api/internal/logic"
	"minify/app/shortener/api/internal/logic/errcode"
	"minify/app/shortener/domain/entity"
	"minify/common/utils/jwtx"

	"minify/app/shortener/api/internal/svc"
	"minify/app/shortener/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAnalyticsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取单个链接的详细报表
func NewGetAnalyticsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAnalyticsLogic {
	return &GetAnalyticsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAnalyticsLogic) GetAnalytics(req *types.GetAnalyticsRequest) (resp *types.GetAnalyticsResponse, err error) {
	// 1. 从 JWT Context 获取用户 ID (身份认证)
	claims, err := jwtx.GetClaimsFromCtx(l.ctx)
	if err != nil {
		return nil, errcode.ErrInvalidToken
	}
	userId := uint64(claims.UserID)

	// 2. 查找链接实体 (获取聚合根)
	link, err := l.svcCtx.LinkRepo.FindByCode(l.ctx, req.Code)
	if err != nil {
		if errors.Is(err, entity.ErrLinkNotFound) {
			// 链接不存在
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

	// 4. 解析日期范围 (⭐ 调用 logic 包中的辅助函数)
	startDate, endDate, err := logic.ParseAnalyticsDates(req.StartDate, req.EndDate)
	if err != nil {
		l.Logger.Infof("parseAnalyticsDates error: %v", err)
		return nil, errcode.ErrInternalError
	}

	// 5. 调用仓储(Repository)获取报表数据
	analyticsEntity, err := l.svcCtx.AnalyticsRepo.GetLinkAnalytics(l.ctx, link.ID, startDate, endDate)
	if err != nil {
		l.Logger.Errorf("AnalyticsRepo.GetLinkAnalytics error: %v", err)
		return nil, errcode.ErrInternalError
	}

	// 6. 将领域实体(Entity)转换为 DTO (⭐ 调用 logic 包中的辅助函数)
	return logic.ToTypesAnalyticsResponse(analyticsEntity), nil
}
