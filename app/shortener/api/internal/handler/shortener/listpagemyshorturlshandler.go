// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.1

package shortener

import (
	"net/http"

	"lucid/app/shortener/api/internal/logic/shortener"
	"lucid/app/shortener/api/internal/svc"
	"lucid/app/shortener/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// 获取当前用户的所有短链接
func ListPageMyShortUrlsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ListPageMyShortUrlsReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := shortener.NewListPageMyShortUrlsLogic(r.Context(), svcCtx)
		resp, err := l.ListPageMyShortUrls(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
