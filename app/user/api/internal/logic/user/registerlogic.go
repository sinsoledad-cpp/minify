// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.1

package user

import (
	"context"
	"errors"
	"minify/app/user/domain/entity"
	"minify/common/utils/jwtx"
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
	// 1. 检查用户名是否已存在
	// ⭐ 调用仓储接口
	if _, err := l.svcCtx.UserRepo.FindByUsername(l.ctx, req.Username); err == nil {
		return nil, errors.New("username already exists")
	}

	// 2. 检查邮箱是否已存在
	// ⭐ 调用仓储接口
	if _, err := l.svcCtx.UserRepo.FindByEmail(l.ctx, req.Email); err == nil {
		return nil, errors.New("email already exists")
	}

	// 3. ⭐ 调用领域实体(Entity)的工厂函数创建 User
	user, err := entity.NewUser(req.Username, req.Email, req.Password)
	if err != nil {
		l.Logger.Errorf("NewUser error: %v", err)
		return nil, errors.New("failed to create user")
	}

	// 4. ⭐ 调用仓储接口(Repository)持久化
	if err := l.svcCtx.UserRepo.Create(l.ctx, user); err != nil {
		l.Logger.Errorf("UserRepo.Create error: %v", err)
		return nil, errors.New("failed to save user")
	}

	// 5. 注册成功，自动登录 (生成 JWT)
	// (user.ID 已在 Create 方法中被回填)
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
		return nil, err
	}

	return &types.LoginResponse{
		AccessToken:  accessToken,
		AccessExpire: now + accessExpire,
	}, nil
}
