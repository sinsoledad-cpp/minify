// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.1

package user

import (
	"context"
	"minify/app/user/api/internal/logic/errcode"
	"minify/common/utils/jwtx"
	"time"

	"minify/app/user/api/internal/svc"
	"minify/app/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGetUserInfoLogic 获取当前登录用户信息
func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserInfoLogic) GetUserInfo() (resp *types.UserInfoResponse, err error) {
	claims, err := jwtx.GetClaimsFromCtx(l.ctx)
	if err != nil {
		return nil, errcode.ErrTokenInvalid
	}

	// 3. ⭐ 调用仓储接口
	user, err := l.svcCtx.UserRepo.FindByID(l.ctx, claims.UserID)
	if err != nil {
		l.Logger.Errorf("FindByID error: %v", err)
		return nil, errcode.ErrUserNotFound
	}

	// 4. ⭐ 转换为 API 响应 (DTO)
	return &types.UserInfoResponse{
		Id:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
	}, nil
}
