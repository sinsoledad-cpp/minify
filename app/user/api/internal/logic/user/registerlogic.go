// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.1

package user

import (
	"context"
	"fmt"

	"lucid/app/user/api/internal/svc"
	"lucid/app/user/api/internal/types"
	"lucid/data/model/user"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 用户注册
func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterReq) (resp *types.RegisterResp, err error) {
	// 检查用户名是否已存在
	_, err = l.svcCtx.UsersModel.FindOneByUsername(l.ctx, req.Username)
	if err == nil {
		return nil, fmt.Errorf("用户名已存在")
	} else if err != sqlx.ErrNotFound {
		logx.Errorf("查询用户失败: %v", err)
		return nil, fmt.Errorf("注册失败，请稍后重试")
	}

	// 对密码进行加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logx.Errorf("密码加密失败: %v", err)
		return nil, fmt.Errorf("注册失败，请稍后重试")
	}

	// 创建用户
	result, err := l.svcCtx.UsersModel.Insert(l.ctx, &user.Users{
		Username: req.Username,
		Password: string(hashedPassword),
	})
	if err != nil {
		logx.Errorf("创建用户失败: %v", err)
		return nil, fmt.Errorf("注册失败，请稍后重试")
	}

	// 获取用户ID
	userId, err := result.LastInsertId()
	if err != nil {
		logx.Errorf("获取用户ID失败: %v", err)
		return nil, fmt.Errorf("注册失败，请稍后重试")
	}

	return &types.RegisterResp{
		UserID:   userId,
		Username: req.Username,
	}, nil
}
