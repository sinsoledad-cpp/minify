// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.1

package shortener

import (
	"errors"
	"minify/app/shortener/domain/entity"
	"net/http"

	"minify/app/shortener/api/internal/logic/shortener"
	"minify/app/shortener/api/internal/svc"
	"minify/app/shortener/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// 短链接重定向 (301/302)
func RedirectHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RedirectRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := shortener.NewRedirectLogic(r.Context(), svcCtx)
		// ⭐ 1. 接收 url 和 err
		url, err := l.Redirect(&req)
		if err != nil {
			// ⭐ 2. 如果是链接不存在，返回 404
			if errors.Is(err, entity.ErrLinkNotFound) {
				http.NotFound(w, r)
			} else if errors.Is(err, entity.ErrLinkExpired) || errors.Is(err, entity.ErrLinkInactive) {
				// 如果是链接过期或禁用，可以返回 400
				http.Error(w, err.Error(), http.StatusBadRequest)
			} else {
				// 其他错误返回 500
				httpx.ErrorCtx(r.Context(), w, err)
			}
		} else {
			// ⭐ 3. 执行重定向，不再调用 httpx.Ok(w)
			http.Redirect(w, r, url, http.StatusFound) // 302 临时重定向
		}
	}
}
