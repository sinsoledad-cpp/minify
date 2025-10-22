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

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewLoginLogic 用户登录
func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginRequest) (resp *types.LoginResponse, err error) {
	// 1. ⭐ 调用仓储接口(Repository)查找用户
	// (根据 API 定义，允许用 username 或 email 登录)
	user, err := l.svcCtx.UserRepo.FindByUsername(l.ctx, req.Username)
	if err != nil {
		// 如果用 Username 找不到，尝试用 Email 找
		if errors.Is(err, entity.ErrUserNotFound) { //
			user, err = l.svcCtx.UserRepo.FindByEmail(l.ctx, req.Username)
			if err != nil {
				// 无论是没找到还是数据库错误，都返回密码错误，防止信息泄露
				return nil, entity.ErrPasswordMismatch //
			}
		} else {
			// 其他数据库错误
			l.Logger.Errorf("FindByUsername error: %v", err)
			return nil, err
		}
	}

	// 2. ⭐ 调用领域实体(Entity)的业务方法校验密码
	if !user.CheckPassword(req.Password) {
		return nil, entity.ErrPasswordMismatch //
	}

	// 3. 密码正确，生成 JWT
	now := time.Now().Unix()
	// 从 svcCtx 获取 Auth 配置
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

	// 4. 返回响应
	return &types.LoginResponse{
		AccessToken:  accessToken,
		AccessExpire: now + accessExpire,
	}, nil
}
