// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.1

package user

import (
	"context"
	"errors"
	"minify/app/user/api/internal/logic/errcode"
	"minify/app/user/domain/entity"
	"minify/common/utils/jwtx"
	"minify/common/utils/response"
	"time"

	"minify/app/user/api/internal/svc"
	"minify/app/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewRegisterLogic 用户注册
func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}
func (l *RegisterLogic) Register(req *types.RegisterRequest) (resp *types.LoginResponse, err error) {
	// 1. 检查用户名
	_, err = l.svcCtx.UserRepo.FindByUsername(l.ctx, req.Username)
	if err == nil {
		return nil, errcode.ErrUsernameExists // <-- 3. 直接返回错误实例
	}
	if !errors.Is(err, entity.ErrUserNotFound) {
		l.Logger.Errorf("FindByUsername error: %v", err)
		return nil, err // 500 错误
	}

	// 2. 检查邮箱
	_, err = l.svcCtx.UserRepo.FindByEmail(l.ctx, req.Email)
	if err == nil {
		return nil, errcode.ErrEmailExists // <-- 3. 直接返回错误实例
	}
	if !errors.Is(err, entity.ErrUserNotFound) {
		l.Logger.Errorf("FindByEmail error: %v", err)
		return nil, err // 500 错误
	}

	// 3. 创建 User 实体
	user, err := entity.NewUser(req.Username, req.Email, req.Password)
	if err != nil {
		// 4. 对于动态错误消息 (如参数校验)，仍然使用 NewBizError
		l.Logger.Infof("NewUser validation error: %v", err)
		return nil, response.NewBizError(response.RequestError, err.Error())
	}

	// 4. 持久化
	if err := l.svcCtx.UserRepo.Create(l.ctx, user); err != nil {
		l.Logger.Errorf("UserRepo.Create error: %v", err)
		return nil, err // 500 错误
	}

	// 5. 生成 JWT
	now := time.Now().Unix()
	accessExpire := l.svcCtx.Config.Auth.AccessExpire
	accessToken, err := jwtx.GenerateToken(
		l.svcCtx.Config.Auth.AccessSecret,
		now,
		accessExpire,
		user.ID,
		user.Role,
	)
	if err != nil {
		l.Logger.Errorf("GenerateToken error: %v", err)
		return nil, err // 500 错误
	}

	return &types.LoginResponse{
		AccessToken:  accessToken,
		AccessExpire: now + accessExpire,
	}, nil
}
