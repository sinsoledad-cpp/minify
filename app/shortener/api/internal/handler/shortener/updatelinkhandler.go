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

// 更新短链接
func UpdateLinkHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateLinkRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := shortener.NewUpdateLinkLogic(r.Context(), svcCtx)
		resp, err := l.UpdateLink(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
