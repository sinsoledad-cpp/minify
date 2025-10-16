// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.1

package shortener

import (
	"context"

	"lucid/app/shortener/api/internal/svc"
	"lucid/app/shortener/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateShortUrlLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建新的短链接
func NewCreateShortUrlLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateShortUrlLogic {
	return &CreateShortUrlLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateShortUrlLogic) CreateShortUrl(req *types.CreateShortUrlReq) (resp *types.CreateShortUrlResp, err error) {
	// todo: add your logic here and delete this line

	return
}
