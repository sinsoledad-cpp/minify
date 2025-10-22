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

// 获取单个链接的详细报表
func GetAnalyticsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetAnalyticsRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := shortener.NewGetAnalyticsLogic(r.Context(), svcCtx)
		resp, err := l.GetAnalytics(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
