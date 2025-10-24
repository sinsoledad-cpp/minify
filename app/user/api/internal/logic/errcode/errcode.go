package errcode

import "minify/common/utils/response"

// 业务错误实例
// Logic 层将直接返回这些预先定义好的错误
var (
	// --- 用户认证/注册 (10001+) ---

	// ErrUserPassword 用户名或密码错误
	ErrUserPassword = response.NewBizError(10001, "invalid username or password")

	// ErrUserNotFound 用户未找到
	ErrUserNotFound = response.NewBizError(10002, "user not found")

	// ErrUsernameExists 用户名已存在
	ErrUsernameExists = response.NewBizError(10003, "username already exists")

	// ErrEmailExists 邮箱已存在
	ErrEmailExists = response.NewBizError(10004, "email already exists")

	// ErrTokenInvalid Token 无效或已过期
	ErrTokenInvalid = response.NewBizError(10005, "invalid token")

	ErrInternalError = response.NewBizError(11500, "服务器内部错误")
)
