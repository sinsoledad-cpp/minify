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

type GetDashboardLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取仪表盘总览数据
func NewGetDashboardLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDashboardLogic {
	return &GetDashboardLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDashboardLogic) GetDashboard(req *types.GetDashboardRequest) (resp *types.GetDashboardResponse, err error) {
	// 1. 从 JWT Context 获取用户 ID
	claims, err := jwtx.GetClaimsFromCtx(l.ctx)
	if err != nil {
		return nil, errcode.ErrInvalidToken // ⭐ 翻译错误
	}
	userId := uint64(claims.UserID)

	// 2. 解析日期范围 (调用 logic/converter.go 中的辅助函数)
	startDate, endDate, err := logic.ParseAnalyticsDates(req.StartDate, req.EndDate)
	if err != nil {
		l.Logger.Infof("ParseAnalyticsDates error: %v", err)
		return nil, errcode.ErrInvalidParams // ⭐ 翻译错误
	}

	// 3. 调用仓储(Repository)获取报表数据
	summary, err := l.svcCtx.AnalyticsRepo.GetDashboardData(l.ctx, &userId, startDate, endDate)
	if err != nil {
		// GetTopLink 返回 entity.ErrLinkNotFound 是正常的 (TopLink 为 nil)，不需要报错
		if !errors.Is(err, entity.ErrLinkNotFound) {
			l.Logger.Errorf("AnalyticsRepo.GetDashboardData error: %v", err)
			return nil, errcode.ErrInternalError // ⭐ 翻译错误
		}
	}

	// 4. 转换 DTO
	return &types.GetDashboardResponse{
		TotalLinks:  summary.TotalLinks,
		TotalClicks: summary.TotalClicks,
		TopLink:     logic.ToTypesLink(summary.TopLink), // ⭐ 使用 converter
	}, nil
}
