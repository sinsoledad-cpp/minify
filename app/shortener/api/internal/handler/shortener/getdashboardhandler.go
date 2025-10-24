// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.1

package shortener

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"minify/app/shortener/api/internal/logic/shortener"
	"minify/app/shortener/api/internal/svc"
	"minify/app/shortener/api/internal/types"

	"minify/common/utils/response"
)

// 获取仪表盘总览数据
func GetDashboardHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetDashboardRequest
		if err := httpx.Parse(r, &req); err != nil {
			response.ClientError(r.Context(), w, response.RequestError, err.Error())
			return
		}

		l := shortener.NewGetDashboardLogic(r.Context(), svcCtx)
		resp, err := l.GetDashboard(&req)
		if err != nil {
			response.LogicError(r.Context(), w, err)
		} else {
			response.Ok(r.Context(), w, resp)
		}
	}
}
