// common/middleware/authz.go
package middleware

import (
	"minify/common/utils/jwtx"
	"minify/common/utils/response" // 引用你的通用响应包
	"net/http"

	"github.com/casbin/casbin/v2"
	"github.com/zeromicro/go-zero/core/logx"
)

type AuthzMiddleware struct {
	Enforcer *casbin.Enforcer
}

// NewAuthzMiddleware 创建一个新的 AuthzMiddleware 实例
func NewAuthzMiddleware(e *casbin.Enforcer) *AuthzMiddleware {
	return &AuthzMiddleware{
		Enforcer: e,
	}
}

// Handle 是中间件的核心逻辑
func (m *AuthzMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. 从 context 获取 claims (必须在 jwt 中间件之后运行)
		claims, err := jwtx.GetClaimsFromCtx(r.Context())
		if err != nil {
			logx.WithContext(r.Context()).Errorf("AuthzMiddleware: GetClaimsFromCtx error: %v", err)
			response.Error(r.Context(), w, http.StatusUnauthorized, http.StatusUnauthorized, "Unauthorized")
			return
		}

		// 2. 准备 Casbin 参数
		sub := claims.Role // Subject: 使用 JWT 中的角色
		obj := r.URL.Path  // Object: 使用请求路径
		act := r.Method    // Action: 使用 HTTP 方法

		// 3. 鉴权
		ok, err := m.Enforcer.Enforce(sub, obj, act)
		if err != nil {
			logx.WithContext(r.Context()).Errorf("AuthzMiddleware: Enforce error: %v", err)
			response.Error(r.Context(), w, http.StatusInternalServerError, 500, "Permission check error")
			return
		}

		// 4. 处理结果
		if !ok {
			logx.WithContext(r.Context()).Infof("AuthzMiddleware: Forbidden - sub: %s, obj: %s, act: %s", sub, obj, act)
			response.Error(r.Context(), w, http.StatusForbidden, http.StatusForbidden, "Forbidden")
			return
		}

		// 5. 放行
		next(w, r)
	}
}
