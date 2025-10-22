// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.1

package shortener

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"lucid/app/shortener/api/internal/logic/shortener"
	"lucid/app/shortener/api/internal/svc"
	"lucid/app/shortener/api/internal/types"
)

// 删除短链接 (软删除)
func DeleteLinkHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DeleteLinkRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := shortener.NewDeleteLinkLogic(r.Context(), svcCtx)
		err := l.DeleteLink(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
