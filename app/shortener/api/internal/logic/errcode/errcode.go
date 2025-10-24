package errcode

import "minify/common/utils/response"

// 业务错误实例
// Logic 层将直接返回这些预先定义好的错误 (shortener-api 使用 11000+)
var (
	// --- 通用/认证 (11001+) ---
	ErrInvalidToken            = response.NewBizError(11001, "无效的Token")
	ErrLinkNotFoundOrForbidden = response.NewBizError(11002, "链接不存在或无权访问")
	ErrInvalidParams           = response.NewBizError(11003, "请求参数无效") // 占位，通常配合 NewBizError(code, msg) 使用

	// --- 链接管理 (11101+) ---
	ErrCustomCodeExists = response.NewBizError(11101, "自定义短码已存在")
	ErrLinkExpired      = response.NewBizError(11102, "链接已过期")

	// --- 内部错误 (11500+) ---
	ErrInternalError   = response.NewBizError(11500, "服务器内部错误")
	ErrIdGenerateError = response.NewBizError(11501, "ID生成失败")
)
