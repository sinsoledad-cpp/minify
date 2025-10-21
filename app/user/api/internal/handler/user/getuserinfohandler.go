// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.1

package user

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"lucid/app/user/api/internal/logic/user"
	"lucid/app/user/api/internal/svc"
)

// 获取当前登录用户信息
func GetUserInfoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewGetUserInfoLogic(r.Context(), svcCtx)
		resp, err := l.GetUserInfo()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
