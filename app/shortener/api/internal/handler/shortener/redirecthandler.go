// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.1

package shortener

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"lucid/app/shortener/api/internal/logic/shortener"
	"lucid/app/shortener/api/internal/svc"
	"lucid/app/shortener/api/internal/types"
)

// 短链接重定向
func RedirectHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RedirectReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		// 将请求和响应写入上下文，以便在逻辑层中使用
		ctx := context.WithValue(r.Context(), "request", r)
		ctx = context.WithValue(ctx, "response", w)

		l := shortener.NewRedirectLogic(ctx, svcCtx)
		err := l.Redirect(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		}
		// 注意：不调用 httpx.Ok(w)，因为重定向已经在逻辑层处理
	}
}
