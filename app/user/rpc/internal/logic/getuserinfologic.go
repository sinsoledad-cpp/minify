package logic

import (
	"context"

	"lucid/app/user/rpc/internal/svc"
	"lucid/gen/go/user/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 根据用户ID获取用户信息
func (l *GetUserInfoLogic) GetUserInfo(in *v1.GetUserInfoRequest) (*v1.GetUserInfoResponse, error) {
	// todo: add your logic here and delete this line

	return &v1.GetUserInfoResponse{}, nil
}
