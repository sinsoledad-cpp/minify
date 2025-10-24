package response

import (
	"context"
	"errors"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// Body 是您 API 统一的响应结构
type Body struct {
	Code int         `json:"code"`           // 自定义业务码
	Msg  string      `json:"msg"`            // 提示信息
	Data interface{} `json:"data,omitempty"` // 成功时的数据
}

// BizError 是一个通用的业务错误结构体
type BizError struct {
	Code int    // 业务码
	Msg  string // 错误信息
}

// 实现 error 接口
func (e *BizError) Error() string {
	return e.Msg
}

// NewBizError 是一个创建 BizError 的辅助函数
// Logic 层将调用这个函数来创建动态错误
func NewBizError(code int, msg string) error {
	return &BizError{Code: code, Msg: msg}
}

// --------------------------------------------------------------------
// 全局错误码（非业务相关）
// --------------------------------------------------------------------
const (
	// OK 成功
	OK = 0
	// InternalError 内部错误
	InternalError = 500
	// RequestError 客户端请求错误
	RequestError = 400
)

// --------------------------------------------------------------------
// Handler 调用的函数
// --------------------------------------------------------------------

// Ok 用于 handler 成功返回
// 用法: response.Ok(r.Context(), w, resp)
func Ok(ctx context.Context, w http.ResponseWriter, resp interface{}) {
	body := Body{
		Code: OK,
		Msg:  "OK",
		Data: resp,
	}
	// 始终返回 200 OK
	httpx.OkJsonCtx(ctx, w, body)
}

// ClientError 用于 handler 处理已知的客户端错误 (如 httpx.Parse 失败)
// 用法: response.ClientError(r.Context(), w, response.RequestError, err.Error())
func ClientError(ctx context.Context, w http.ResponseWriter, bizCode int, errMsg string) {
	body := Body{
		Code: bizCode,
		Msg:  errMsg,
	}
	// 始终返回 400 Bad Request
	httpx.WriteJsonCtx(ctx, w, http.StatusBadRequest, body)
}

// LogicError 用于 handler 处理来自 logic 层的错误
// 它会自动区分 BizError (400) 和 未知Error (500)
// 用法: response.LogicError(r.Context(), w, err)
func LogicError(ctx context.Context, w http.ResponseWriter, err error) {
	var httpStatus int
	var bizCode int
	var bizMsg string

	var bizErr *BizError
	if errors.As(err, &bizErr) {
		// 1. 是我们自定义的 BizError
		httpStatus = http.StatusBadRequest // 业务错误默认为 400
		bizCode = bizErr.Code
		bizMsg = bizErr.Msg
		// 增加 Warn 日志，用于追踪业务错误
		logx.WithContext(ctx).Errorf("Business warning: code=%d, msg=%s", bizCode, bizMsg)
	} else {
		// 2. 是一个未知的系统错误
		logx.WithContext(ctx).Errorf("Internal server error: %v", err)
		httpStatus = http.StatusInternalServerError
		bizCode = InternalError
		bizMsg = "Internal Server Error" // 不暴露详细错误
	}

	body := Body{
		Code: bizCode,
		Msg:  bizMsg,
	}
	// 按提取的 httpStatus 写入响应
	httpx.WriteJsonCtx(ctx, w, httpStatus, body)
}
