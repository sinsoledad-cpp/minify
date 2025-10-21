// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.1

package user

import (
	"context"

	"lucid/app/user/api/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
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
	// 对于无状态 JWT 认证，服务器端 "登出"
	// 实际上并不需要做任何事情。
	//
	// 1. 你的路由配置 确保了 @jwt(Auth)
	//    中间件已经运行，它验证了这是一个合法的、未过期的 Token。
	//
	// 2. 服务器是无状态的，它不存储 Token。
	//
	// 3. 返回 nil (即 HTTP 200 OK)
	//    就是给客户端的信号：“我已收到你的登出请求，请你（客户端）自行删除本地 Token”。
	//
	// 客户端（前端）在收到这个 200 响应后，
	// 必须负责从 localStorage / cookie 中删除该 Token。

	return nil
}
