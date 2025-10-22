// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.1

package shortener

import (
	"context"

	"minify/app/shortener/api/internal/svc"
	"minify/app/shortener/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateLinkLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新短链接
func NewUpdateLinkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLinkLogic {
	return &UpdateLinkLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateLinkLogic) UpdateLink(req *types.UpdateLinkRequest) (resp *types.UpdateLinkResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
