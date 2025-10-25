// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.1

package shortener

import (
	"context"
	"encoding/json"
	"errors"
	"minify/app/shortener/api/internal/logic/errcode"
	"minify/app/shortener/api/internal/svc"
	"minify/app/shortener/api/internal/types"
	"minify/app/shortener/domain/entity"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/threading"
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

func (l *RedirectLogic) Redirect(req *types.RedirectRequest, ip, ua, referer string) (string, error) {
	// 1. 调用 Repo
	link, err := l.svcCtx.LinkRepo.FindByCode(l.ctx, req.Code)
	if err != nil {
		// 如果找不到，返回给 handler
		if errors.Is(err, entity.ErrLinkNotFound) {
			return "", errcode.ErrLinkNotFound // entity.ErrLinkNotFound 是一个标准 error
		}
		l.Logger.Errorf("FindByCode error: %v", err)
		return "", errcode.ErrInternalError
	}

	// 2. 检查链接是否可用
	if err := link.CanRedirect(); err != nil {
		// 例如返回 ErrLinkExpired
		return "", errcode.ErrLinkExpired
	}

	// 3. (TODO) ⭐ 使用 GoSafeCtx 异步发送 Kafka 日志
	// 这会安全地启动一个 goroutine，并自动处理 panic
	threading.GoSafeCtx(l.ctx, func() {
		l.logAccessEvent(l.ctx, link, ip, ua, referer)
	})

	// 4. ⭐ 返回目标 URL
	return link.OriginalUrl, nil
}

type LinkAccessEvent struct {
	ShortCode  string    `json:"shortCode"`
	LinkID     int64     `json:"linkId"`
	AccessedAt time.Time `json:"accessedAt"`
	IpAddress  string    `json:"ipAddress"`
	UserAgent  string    `json:"userAgent"`
	Referer    string    `json:"referer"`
}

// (辅助函数) 添加一个私有方法来处理 Kafka 发送
func (l *RedirectLogic) logAccessEvent(ctx context.Context, link *entity.Link, ip, ua, referer string) {
	// 构造消息
	event := LinkAccessEvent{
		ShortCode:  link.ShortCode,
		LinkID:     link.ID,
		AccessedAt: time.Now(),
		IpAddress:  ip,
		UserAgent:  ua,
		Referer:    referer,
	}

	// 序列化
	msgBody, err := json.Marshal(event)
	if err != nil {
		// ⭐ 7. (已修正) 使用传入的 ctx 记录日志, 确保 trace_id 等信息被正确传递
		logx.WithContext(ctx).Errorf("logAccessEvent: failed to marshal event for link %s: %v", link.ShortCode, err)
		return
	}

	// 8. ⭐ (已确认) 发送 Kafka 消息
	// kq.Pusher.Push(key string, value string) error PushWithKey
	err = l.svcCtx.LinkEventProducer.KPush(ctx, link.ShortCode, string(msgBody))
	if err != nil {
		logx.WithContext(ctx).Errorf("logAccessEvent: failed to push event to kafka for link %s: %v", link.ShortCode, err)
	}
}
