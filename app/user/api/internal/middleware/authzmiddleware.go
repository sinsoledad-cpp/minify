// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.1

package middleware

import (
	"net/http"
)

type AuthzMiddleware struct {
}

func NewAuthzMiddleware() *AuthzMiddleware {
	return &AuthzMiddleware{}
}

func (m *AuthzMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO generate middleware implement function, delete after code implementation

		// Passthrough to next handler if need
		next(w, r)
	}
}
