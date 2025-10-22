// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.1

package shortener

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"minify/app/shortener/api/internal/logic/shortener"
	"minify/app/shortener/api/internal/svc"
	"minify/app/shortener/api/internal/types"
)

// 创建短链接
func CreateLinkHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateLinkRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := shortener.NewCreateLinkLogic(r.Context(), svcCtx)
		resp, err := l.CreateLink(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
