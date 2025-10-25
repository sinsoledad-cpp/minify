package mq

import (
	"context"
	"encoding/json"
	"minify/app/shortener/api/internal/svc"
	"minify/app/shortener/domain/entity"

	"github.com/zeromicro/go-zero/core/logx"
)

// LinkAccessEventConsumer 是处理链接访问事件的消费者
type LinkAccessEventConsumer struct {
	ctx    context.Context // 这个 ctx 是服务启动时的全局 context
	svcCtx *svc.ServiceContext
}

// NewLinkAccessEventConsumer 创建一个新的消费者实例
func NewLinkAccessEventConsumer(ctx context.Context, svcCtx *svc.ServiceContext) *LinkAccessEventConsumer {
	return &LinkAccessEventConsumer{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Consume 是 go-queue 调用的核心方法
func (l *LinkAccessEventConsumer) Consume(ctx context.Context, key, val string) error {
	logx.WithContext(ctx).Infof("Kafka Consumer: Received message key: %s, val: %s", key, val)

	var event LinkAccessEvent
	if err := json.Unmarshal([]byte(val), &event); err != nil {
		logx.WithContext(ctx).Errorf("Failed to unmarshal LinkAccessEvent: %v", err)
		return err // 反序列化失败，返回错误，消息可能会被重试
	}

	// 1. 异步更新 links 表的 visit_count (原子操作，处理缓存)
	// (这部分已经在使用 LinkRepo 接口，是正确的)
	if err := l.svcCtx.LinkRepo.IncrementVisitCount(ctx, event.LinkID, 1); err != nil {
		// 即使这里失败，我们仍然尝试记录原始日志
		logx.WithContext(ctx).Errorf("Failed to increment visit count for linkID %d: %v", event.LinkID, err)
	}

	// 2. ⭐ (已修改) 创建领域实体(Entity)
	logEntry := entity.NewLinkAccessLog(
		event.LinkID,
		event.ShortCode,
		event.AccessedAt,
		event.IpAddress,
		event.UserAgent,
		event.Referer,
	)

	// 3. ⭐ (已修改) 调用仓储接口
	if err := l.svcCtx.LinkAccessLogsRepo.Create(ctx, logEntry); err != nil {
		logx.WithContext(ctx).Errorf("Failed to insert link access log for linkID %d: %v", event.LinkID, err)
		return err // 插入日志失败，返回错误
	}

	return nil
}
