package errcode

// 业务错误码
const (
	// --- 用户认证/注册 (10001+) ---

	// ErrUserPassword 用户名或密码错误
	// (用于登录时，FindByID 或 CheckPassword 失败)
	ErrUserPassword = 10001

	// ErrUserNotFound 用户未找到
	// (用于 GetUserInfo 时 FindByID 失败)
	ErrUserNotFound = 10002

	// ErrUsernameExists 用户名已存在
	// (用于注册时)
	ErrUsernameExists = 10003

	// ErrEmailExists 邮箱已存在
	// (用于注册时)
	ErrEmailExists = 10004

	// ErrTokenInvalid Token 无效或已过期
	// (用于 GetClaimsFromCtx 失败)
	ErrTokenInvalid = 10005
)
