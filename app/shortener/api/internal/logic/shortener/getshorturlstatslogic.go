// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.1

package shortener

import (
	"context"

	"lucid/app/shortener/api/internal/svc"
	"lucid/app/shortener/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetShortUrlStatsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取短链接的访问统计
func NewGetShortUrlStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetShortUrlStatsLogic {
	return &GetShortUrlStatsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetShortUrlStatsLogic) GetShortUrlStats(req *types.GetShortUrlStatsReq) (resp *types.GetShortUrlStatsResp, err error) {
	// todo: add your logic here and delete this line

	return
}
