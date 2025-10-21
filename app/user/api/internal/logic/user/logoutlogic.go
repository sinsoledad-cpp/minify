// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.1

package user

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"lucid/app/user/api/internal/svc"
)

type LogoutLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 用户登出
func NewLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogoutLogic {
	return &LogoutLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LogoutLogic) Logout() error {
	// todo: add your logic here and delete this line

	return nil
}
