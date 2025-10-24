package response

import (
	"context"
	"errors"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

type Body struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}
type BizError struct {
	Code int
	Msg  string
}

func (e *BizError) Error() string {
	return e.Msg
}
func NewBizError(code int, msg string) error {
	return &BizError{Code: code, Msg: msg}
}

const (
	OK    = 0
	OKMsg = "OK"
)
const (
	InternalError = 500
	RequestError  = 400
)

func Ok(ctx context.Context, w http.ResponseWriter, resp interface{}) {
	Success(ctx, w, http.StatusOK, resp, OKMsg)
}

func OkMsg(ctx context.Context, w http.ResponseWriter, resp interface{}, msg string) {
	Success(ctx, w, http.StatusOK, resp, msg)
}

func Success(ctx context.Context, w http.ResponseWriter, httpStatus int, resp interface{}, msg string) {
	if msg == "" {
		msg = OKMsg
	}
	body := Body{
		Code: OK,
		Msg:  msg,
		Data: resp,
	}
	httpx.WriteJsonCtx(ctx, w, httpStatus, body)
}

func Error(ctx context.Context, w http.ResponseWriter, httpStatus int, bizCode int, errMsg string) {
	body := Body{
		Code: bizCode,
		Msg:  errMsg,
	}
	// 使用你传入的自定义 httpStatus
	httpx.WriteJsonCtx(ctx, w, httpStatus, body)
}

// ClientError 用于 handler 处理已知的客户端错误 (如 httpx.Parse 失败)
// 它封装了 Error(ctx, w, http.StatusBadRequest, bizCode, errMsg)
// (这是修改后的快捷方式)
func ClientError(ctx context.Context, w http.ResponseWriter, bizCode int, errMsg string) {
	// 调用 Error，并硬编码 http.StatusBadRequest
	Error(ctx, w, http.StatusBadRequest, bizCode, errMsg)
}

// LogicError 用于 handler 处理来自 logic 层的错误
// 它会自动区分 BizError (400) 和 未知Error (500)
// (这是修改后，调用 Error 的快捷方式)
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
		// (建议：使用 Infof)
		logx.WithContext(ctx).Infof("Business warning: code=%d, msg=%s", bizCode, bizMsg)
	} else {
		// 2. 是一个未知的系统错误
		logx.WithContext(ctx).Infof("Internal server error: %v", err)
		httpStatus = http.StatusInternalServerError
		bizCode = InternalError
		bizMsg = "Internal Server Error" // 不暴露详细错误
	}

	// 统一调用新的 Error 函数
	Error(ctx, w, httpStatus, bizCode, bizMsg)
}
