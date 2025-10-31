// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.1

package admin

import (
	"context"
	"errors"
	"minify/app/shortener/api/internal/logic"
	"minify/app/shortener/api/internal/logic/errcode"
	"minify/app/shortener/domain/entity"

	"minify/app/shortener/api/internal/svc"
	"minify/app/shortener/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetLinkAnalyticsAdminLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取任意链接的详细报表 (Admin)
func NewGetLinkAnalyticsAdminLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLinkAnalyticsAdminLogic {
	return &GetLinkAnalyticsAdminLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetLinkAnalyticsAdminLogic) GetLinkAnalyticsAdmin(req *types.GetAnalyticsRequest) (resp *types.GetAnalyticsResponse, err error) {
	// 1. (鉴权) Casbin 中间件已确保是 admin。

	// 2. 查找链接实体 (获取 link.ID)
	// (关键: 跳过 GetAnalyticsLogic 中的所有者检查)
	link, err := l.svcCtx.LinkRepo.FindByCode(l.ctx, req.Code)
	if err != nil {
		if errors.Is(err, entity.ErrLinkNotFound) {
			// 链接不存在
			return nil, errcode.ErrLinkNotFound //
		}
		// 其他数据库错误
		l.Logger.Errorf("FindByCode error: %v", err)
		return nil, errcode.ErrInternalError //
	}

	// 3. 解析日期范围 (复用 logic/converter.go 中的辅助函数)
	startDate, endDate, err := logic.ParseAnalyticsDates(req.StartDate, req.EndDate)
	if err != nil {
		l.Logger.Infof("parseAnalyticsDates error: %v", err)
		return nil, errcode.ErrInvalidParams //
	}

	// 4. 调用仓储(Repository)获取报表数据
	analyticsEntity, err := l.svcCtx.AnalyticsRepo.GetLinkAnalytics(l.ctx, link.ID, startDate, endDate) //
	if err != nil {
		l.Logger.Errorf("AnalyticsRepo.GetLinkAnalytics error: %v", err)
		return nil, errcode.ErrInternalError
	}

	// 5. 将领域实体(Entity)转换为 DTO (复用 logic/converter.go 中的辅助函数)
	return logic.ToTypesAnalyticsResponse(analyticsEntity), nil
}
