// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.1

package admin

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"minify/app/shortener/api/internal/logic/admin"
	"minify/app/shortener/api/internal/svc"
	"minify/app/shortener/api/internal/types"

	"minify/common/utils/response"
)

// 获取任意链接的详细报表 (Admin)
func GetLinkAnalyticsAdminHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetAnalyticsRequest
		if err := httpx.Parse(r, &req); err != nil {
			response.ClientError(r.Context(), w, response.RequestError, err.Error())
			return
		}

		l := admin.NewGetLinkAnalyticsAdminLogic(r.Context(), svcCtx)
		resp, err := l.GetLinkAnalyticsAdmin(&req)
		if err != nil {
			response.LogicError(r.Context(), w, err)
		} else {
			response.Ok(r.Context(), w, resp)
		}
	}
}
