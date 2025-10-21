// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.1

package user

import (
	"net/http"

	"lucid/app/user/api/internal/logic/user"
	"lucid/app/user/api/internal/svc"
	"lucid/app/user/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// RegisterHandler 用户注册
func RegisterHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RegisterRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := user.NewRegisterLogic(r.Context(), svcCtx)
		resp, err := l.Register(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
