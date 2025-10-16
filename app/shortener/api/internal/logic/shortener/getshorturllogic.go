// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.1

package shortener

import (
	"context"

	"lucid/app/shortener/api/internal/svc"
	"lucid/app/shortener/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetShortUrlLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取短链接详情
func NewGetShortUrlLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetShortUrlLogic {
	return &GetShortUrlLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetShortUrlLogic) GetShortUrl(req *types.GetShortUrlReq) (resp *types.GetShortUrlResp, err error) {
	// todo: add your logic here and delete this line

	return
}
