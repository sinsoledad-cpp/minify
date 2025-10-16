// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.1

package shortener

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	"lucid/app/shortener/api/internal/svc"
	"lucid/app/shortener/api/internal/types"
	"lucid/data/model/shortener"

	"github.com/zeromicro/go-zero/core/logx"
)

type RedirectLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 短链接重定向
func NewRedirectLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RedirectLogic {
	return &RedirectLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RedirectLogic) Redirect(req *types.RedirectReq) error {
	// 1. 查询短链接是否存在
	shortUrl, err := l.svcCtx.ShortUrlsModel.FindOneByShortKey(l.ctx, req.ShortKey)
	if err != nil {
		if err == shortener.ErrNotFound {
			return fmt.Errorf("短链接不存在")
		}
		logx.Errorf("查询短链接失败: %v", err)
		return fmt.Errorf("系统错误，请稍后重试")
	}

	// 2. 检查短链接是否已删除
	if shortUrl.DeletedAt.Valid {
		return fmt.Errorf("短链接已失效")
	}

	// 3. 检查短链接是否过期
	if shortUrl.ExpiresAt.Valid && shortUrl.ExpiresAt.Time.Before(time.Now()) {
		return fmt.Errorf("短链接已过期")
	}

	// 4. 获取请求信息用于记录访问数据
	r := l.ctx.Value("request").(*http.Request)
	ipAddress := r.RemoteAddr
	// 如果有代理，尝试获取真实IP
	if forwardedFor := r.Header.Get("X-Forwarded-For"); forwardedFor != "" {
		ipAddress = strings.Split(forwardedFor, ",")[0]
	}

	// 5. 记录访问数据
	userAgent := r.UserAgent()
	referer := r.Referer()
	analytics := &shortener.UrlAnalytics{
		ShortUrlId: shortUrl.Id,
		IpAddress:  ipAddress,
		UserAgent: sql.NullString{
			String: userAgent,
			Valid:  userAgent != "",
		},
		Referer: sql.NullString{
			String: referer,
			Valid:  referer != "",
		},
	}

	// 异步记录访问数据，不阻塞重定向
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		_, err := l.svcCtx.UrlAnalyticsModel.Insert(ctx, analytics)
		if err != nil {
			logx.Errorf("记录访问数据失败: %v", err)
		}
	}()

	// 6. 执行重定向
	w := l.ctx.Value("response").(http.ResponseWriter)
	http.Redirect(w, r, shortUrl.OriginalUrl, http.StatusFound)
	return nil
}
