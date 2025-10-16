// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.1

package shortener

import (
	"context"

	"lucid/app/shortener/api/internal/svc"
	"lucid/app/shortener/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListMyShortUrlsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取当前用户的所有短链接
func NewListMyShortUrlsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListMyShortUrlsLogic {
	return &ListMyShortUrlsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListMyShortUrlsLogic) ListMyShortUrls() (resp *types.ListShortUrlsResp, err error) {
	// todo: add your logic here and delete this line

	return
}
