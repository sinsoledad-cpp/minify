// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.1

package shortener

import (
	"net/http"

	"minify/app/shortener/api/internal/logic/shortener"
	"minify/app/shortener/api/internal/svc"
	"minify/app/shortener/api/internal/types"

	"minify/common/utils/response"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// 获取短链接列表 (分页)
func ListLinksHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ListLinksRequest
		if err := httpx.Parse(r, &req); err != nil {
			response.ClientError(r.Context(), w, response.RequestError, err.Error())
			return
		}

		l := shortener.NewListLinksLogic(r.Context(), svcCtx)
		resp, err := l.ListLinks(&req)
		if err != nil {
			response.LogicError(r.Context(), w, err)
		} else {
			response.Ok(r.Context(), w, resp)
		}
	}
}
