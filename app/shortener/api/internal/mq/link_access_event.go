package mq

import (
	"context"
	"encoding/json"
	"minify/app/shortener/api/internal/svc"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

// LinkAccessEvent 对应生产者 (RedirectLogic) 发送的消息结构
type LinkAccessEvent struct {
	ShortCode  string    `json:"shortCode"`
	LinkID     int64     `json:"linkId"`
	AccessedAt time.Time `json:"accessedAt"`
	IpAddress  string    `json:"ipAddress"`
	UserAgent  string    `json:"userAgent"`
	Referer    string    `json:"referer"`
}

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
// ⭐ 2. (已修正) 签名现在包含 ctx context.Context
func (l *LinkAccessEventConsumer) Consume(ctx context.Context, key, val string) error {
	// ⭐ 3. (已修正) 使用传入的 ctx 进行日志记录
	logx.WithContext(ctx).Infof("Kafka Consumer: Received message key: %s, val: %s", key, val)

	var event LinkAccessEvent
	if err := json.Unmarshal([]byte(val), &event); err != nil {
		// ⭐ 3. (已修正) 使用传入的 ctx 进行日志记录
		logx.WithContext(ctx).Errorf("Failed to unmarshal LinkAccessEvent: %v", err)
		return err // 反序列化失败，返回错误，消息可能会被重试
	}

	// 1. 异步更新 links 表的 visit_count (原子操作，处理缓存)
	// ⭐ 3. (已修正) 使用传入的 ctx 调用仓储
	if err := l.svcCtx.LinkRepo.IncrementVisitCount(ctx, event.LinkID, 1); err != nil {
		// 即使这里失败，我们仍然尝试记录原始日志
		logx.WithContext(ctx).Errorf("Failed to increment visit count for linkID %d: %v", event.LinkID, err)
	}

	//// 2. 写入原始访问日志 (link_access_logs)
	//logEntry := &model.LinkAccessLogs{
	//	LinkId:     uint64(event.LinkID),
	//	ShortCode:  event.ShortCode,
	//	AccessedAt: event.AccessedAt,
	//	IpAddress:  event.IpAddress,
	//	UserAgent:  sql.NullString{String: event.UserAgent, Valid: event.UserAgent != ""},
	//	Referer:    sql.NullString{String: event.Referer, Valid: event.Referer != ""},
	//	// Geo/UA 解析字段 (GeoCountry, GeoCity, etc.) 留空，等待 ETL 任务处理
	//}
	//
	//// ⭐ 3. (已修正) 使用传入的 ctx 插入数据库
	//if _, err := l.svcCtx.LinkAccessLogsModel.Insert(ctx, logEntry); err != nil {
	//	logx.WithContext(ctx).Errorf("Failed to insert link access log for linkID %d: %v", event.LinkID, err)
	//	return err // 插入日志失败，返回错误
	//}

	return nil
}
