// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.1

package user

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"lucid/app/user/api/internal/svc"
	"lucid/app/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取当前用户信息
func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserInfoLogic) GetUserInfo() (resp *types.UserInfoResp, err error) {
	// 从上下文中获取用户ID
	userId, err := l.ctx.Value("userId").(json.Number).Int64()
	if err != nil {
		return nil, fmt.Errorf("获取用户信息失败")
	}

	// 查询用户信息
	user, err := l.svcCtx.UsersModel.FindOne(l.ctx, uint64(userId))
	if err != nil {
		return nil, fmt.Errorf("获取用户信息失败")
	}

	return &types.UserInfoResp{
		UserID:    int64(user.Id),
		Username:  user.Username,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
	}, nil
}
