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
	// 1. 检查用户名是否已存在
	_, err = l.svcCtx.UserRepo.FindByUsername(l.ctx, req.Username)
	if err == nil {
		// 3. 翻译错误：找到了用户，说明已存在
		return nil, response.NewBizError(errcode.ErrUsernameExists, "username already exists")
	}
	if !errors.Is(err, entity.ErrUserNotFound) {
		// 4. 是数据库查询错误，返回 500
		l.Logger.Errorf("FindByUsername error: %v", err)
		return nil, err
	}
	// 走到这里说明 err == entity.ErrUserNotFound, 用户名可用, 继续

	// 2. 检查邮箱是否已存在
	_, err = l.svcCtx.UserRepo.FindByEmail(l.ctx, req.Email)
	if err == nil {
		// 3. 翻译错误：找到了邮箱，说明已存在
		return nil, response.NewBizError(errcode.ErrEmailExists, "email already exists")
	}
	if !errors.Is(err, entity.ErrUserNotFound) {
		// 4. 是数据库查询错误，返回 500
		l.Logger.Errorf("FindByEmail error: %v", err)
		return nil, err
	}
	// 走到这里说明 err == entity.ErrUserNotFound, 邮箱可用, 继续

	// 3. 创建 User 实体
	user, err := entity.NewUser(req.Username, req.Email, req.Password)
	if err != nil {
		// 3. 翻译错误：来自实体的校验错误 (如密码太短)
		l.Logger.Infof("NewUser validation error: %v", err)
		// 假设这是客户端请求参数错误
		return nil, response.NewBizError(response.RequestError, err.Error())
	}

	// 4. 持久化
	if err := l.svcCtx.UserRepo.Create(l.ctx, user); err != nil {
		// 4. 数据库插入错误，返回 500
		l.Logger.Errorf("UserRepo.Create error: %v", err)
		return nil, err
	}

	// 5. 注册成功，自动登录 (生成 JWT)
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
		// 4. JWT 生成失败，返回 500
		l.Logger.Errorf("GenerateToken error: %v", err)
		return nil, err
	}

	return &types.LoginResponse{
		AccessToken:  accessToken,
		AccessExpire: now + accessExpire,
	}, nil
}
