// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.1

package shortener

import (
	"context"
	"errors"
	"minify/app/shortener/domain/entity"

	"minify/app/shortener/api/internal/svc"
	"minify/app/shortener/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RedirectLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 短链接重定向 (301/302)
func NewRedirectLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RedirectLogic {
	return &RedirectLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RedirectLogic) Redirect(req *types.RedirectRequest) (string, error) {
	// 1. 调用 Repo
	link, err := l.svcCtx.LinkRepo.FindByCode(l.ctx, req.Code)
	if err != nil {
		// 如果找不到，返回给 handler
		if errors.Is(err, entity.ErrLinkNotFound) {
			return "", entity.ErrLinkNotFound // entity.ErrLinkNotFound 是一个标准 error
		}
		l.Logger.Errorf("FindByCode error: %v", err)
		return "", err
	}

	// 2. 检查链接是否可用
	if err := link.CanRedirect(); err != nil {
		// 例如返回 ErrLinkExpired
		return "", err
	}

	// 3. (TODO) 在这里异步发送日志...

	// 4. ⭐ 返回目标 URL
	return link.OriginalUrl, nil
}
