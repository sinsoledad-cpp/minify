// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.1

package user

import (
	"context"
	"fmt"
	"time"

	"lucid/app/user/api/internal/svc"
	"lucid/app/user/api/internal/types"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 用户登录
func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	// 查询用户是否存在
	user, err := l.svcCtx.UsersModel.FindOneByUsername(l.ctx, req.Username)
	if err != nil {
		return nil, fmt.Errorf("用户名或密码错误")
	}

	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, fmt.Errorf("用户名或密码错误")
	}

	// 生成JWT令牌
	now := time.Now().Unix()
	accessExpire := l.svcCtx.Config.Auth.AccessExpire
	accessToken, err := l.genToken(user.Id, user.Username, now, accessExpire)
	if err != nil {
		return nil, fmt.Errorf("登录失败，请稍后重试")
	}

	return &types.LoginResp{
		UserID:      int64(user.Id),
		Username:    user.Username,
		AccessToken: accessToken,
		ExpiresIn:   accessExpire,
	}, nil
}

// 生成JWT令牌
func (l *LoginLogic) genToken(userId uint64, username string, iat, seconds int64) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	claims["userId"] = userId
	claims["username"] = username
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(l.svcCtx.Config.Auth.AccessSecret))
}
