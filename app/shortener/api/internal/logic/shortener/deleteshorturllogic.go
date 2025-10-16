// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.1

package shortener

import (
	"context"

	"lucid/app/shortener/api/internal/svc"
	"lucid/app/shortener/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteShortUrlLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除短链接
func NewDeleteShortUrlLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteShortUrlLogic {
	return &DeleteShortUrlLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteShortUrlLogic) DeleteShortUrl(req *types.DeleteShortUrlReq) (resp *types.DeleteShortUrlResp, err error) {
	// todo: add your logic here and delete this line

	return
}
