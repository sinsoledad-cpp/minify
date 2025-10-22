// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.1

package shortener

import (
	"context"

	"minify/app/shortener/api/internal/svc"
	"minify/app/shortener/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteLinkLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除短链接 (软删除)
func NewDeleteLinkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLinkLogic {
	return &DeleteLinkLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteLinkLogic) DeleteLink(req *types.DeleteLinkRequest) error {
	// todo: add your logic here and delete this line

	return nil
}
