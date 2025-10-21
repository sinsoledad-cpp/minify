package jwtx

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/golang-jwt/jwt/v4"
)

// Claims 是我们放入 Token 的自定义数据
type Claims struct {
	UserID int64  `json:"userId"`
	Role   string `json:"role"`
}

// GenerateToken 生成一个新的 JWT token (供 user-api 使用)
func GenerateToken(secretKey string, iat, seconds, userId int64, role string) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	claims["userId"] = userId // 使用 Claims 结构体中的键
	claims["role"] = role     // 使用 Claims 结构体中的键
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secretKey))
}

// GetClaimsFromCtx 从 go-zero 中间件注入的 context 中安全地提取 Claims
// (供所有微服务使用)
func GetClaimsFromCtx(ctx context.Context) (*Claims, error) {
	var claims Claims

	// 1. 提取 UserID
	val := ctx.Value("userId")
	if val == nil {
		return nil, errors.New("missing userId in token context")
	}

	var userId int64
	if idFloat, ok := val.(float64); ok {
		userId = int64(idFloat)
	} else if idJson, ok := val.(json.Number); ok {
		userId, _ = idJson.Int64()
	} else {
		return nil, errors.New("invalid userId type in token context")
	}
	if userId == 0 {
		return nil, errors.New("invalid user id in token")
	}
	claims.UserID = userId

	// 2. 提取 Role
	val = ctx.Value("role")
	if val == nil {
		return nil, errors.New("missing role in token context")
	}

	if role, ok := val.(string); ok {
		claims.Role = role
	} else {
		return nil, errors.New("invalid role type in token context")
	}

	return &claims, nil
}
