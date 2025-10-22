// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.1

package shortener

import (
	"context"

	"lucid/app/shortener/api/internal/svc"
	"lucid/app/shortener/api/internal/types"

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
	// todo: add your logic here and delete this line

	return
}
