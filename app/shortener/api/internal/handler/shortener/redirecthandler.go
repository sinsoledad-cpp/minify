// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.1

package shortener

import (
	"errors"
	"minify/app/shortener/api/internal/logic/errcode"
	"net/http"
	"strings"

	"minify/app/shortener/api/internal/logic/shortener"
	"minify/app/shortener/api/internal/svc"
	"minify/app/shortener/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"

	"minify/common/utils/response"
)

// 短链接重定向 (301/302)
func RedirectHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RedirectRequest
		if err := httpx.Parse(r, &req); err != nil {
			response.ClientError(r.Context(), w, response.RequestError, err.Error())
			return
		}

		ip := getClientIp(r)
		ua := r.UserAgent()
		referer := r.Referer()
		l := shortener.NewRedirectLogic(r.Context(), svcCtx)
		// ⭐ 1. 接收 url 和 err
		url, err := l.Redirect(&req, ip, ua, referer)
		if err != nil {
			// ⭐ 2. 如果是链接不存在，返回 404
			if errors.Is(err, errcode.ErrLinkNotFound) {
				http.NotFound(w, r)
			} else if errors.Is(err, errcode.ErrLinkExpiredOrInactive) {
				// 如果是链接过期或禁用，可以返回 400
				//http.Error(w, err.Error(), http.StatusBadRequest)
				response.LogicError(r.Context(), w, err)
			} else {
				// 其他错误返回 500
				response.LogicError(r.Context(), w, err)
			}
		} else {
			// ⭐ 3. 执行重定向，不再调用 httpx.Ok(w)
			http.Redirect(w, r, url, http.StatusFound) // 302 临时重定向
		}
	}
}

// ⭐ 4. (辅助函数) 添加一个辅助函数来获取真实客户端 IP
// 优先从 X-Forwarded-For (Nginx等代理) 获取, 其次 X-Real-IP, 最后才是 RemoteAddr
func getClientIp(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip != "" {
		// X-Forwarded-For 可能是 "client, proxy1, proxy2"
		if parts := strings.Split(ip, ","); len(parts) > 0 {
			return strings.TrimSpace(parts[0])
		}
	}
	ip = r.Header.Get("X-Real-IP")
	if ip != "" {
		return ip
	}
	// r.RemoteAddr 格式可能是 "ip:port", 尝试分割
	if parts := strings.Split(r.RemoteAddr, ":"); len(parts) > 0 {
		if parts[0] != "" {
			// 返回 IP 部分
			return parts[0]
		}
	}
	return r.RemoteAddr
}
