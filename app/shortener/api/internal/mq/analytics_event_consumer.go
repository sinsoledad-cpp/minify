package mq

import (
	"context"
	"encoding/json"
	"minify/app/shortener/api/internal/svc"
	"minify/app/shortener/domain/entity"

	"github.com/zeromicro/go-zero/core/logx"
	// ⭐ (推荐) 导入 UA 解析库: "github.com/ua-parser/uap-go/uaparser"
	// ⭐ (推荐) 导入 IP 解析库: "github.com/oschwald/geoip2-golang"
)

// AnalyticsEventConsumer 是处理聚合分析的消费者
type AnalyticsEventConsumer struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	// ⭐ (推荐) uaParser *uaparser.Parser
	// ⭐ (推荐) geoIPDB  *geoip2.Reader
}

// NewAnalyticsEventConsumer 创建一个新的消费者实例
func NewAnalyticsEventConsumer(ctx context.Context, svcCtx *svc.ServiceContext) *AnalyticsEventConsumer {
	// ⭐ (推荐) 在这里初始化解析器
	// parser, _ := uaparser.New("path/to/regexes.yaml")
	// db, _ := geoip2.Open("path/to/GeoLite2-City.mmdb")

	return &AnalyticsEventConsumer{
		ctx:    ctx,
		svcCtx: svcCtx,
		// ⭐ uaParser: parser,
		// ⭐ geoIPDB:  db,
	}
}

// Consume 是 go-queue 调用的核心方法
func (l *AnalyticsEventConsumer) Consume(ctx context.Context, key, val string) error {
	logx.WithContext(ctx).Infof("Analytics Consumer: Received message key: %s, val: %s", key, val)

	var event LinkAccessEvent
	if err := json.Unmarshal([]byte(val), &event); err != nil {
		logx.WithContext(ctx).Errorf("Failed to unmarshal LinkAccessEvent: %v", err)
		return err
	}

	// --- 1. (核心) 解析 IP 和 UserAgent ---
	// (你需要自己实现这部分)

	// TODO: 解析 IP 地址
	// ip := net.ParseIP(event.IpAddress)
	// record, err := l.geoIPDB.City(ip)
	// country := record.Country.Names["en"]
	country := "unknown" // Stub

	// TODO: 解析 User Agent
	// client := l.uaParser.Parse(event.UserAgent)
	// browser := client.UserAgent.Family
	// os := client.Os.Family
	// device := client.Device.Family
	browser := "unknown" // Stub
	os := "unknown"      // Stub
	device := "unknown"  // Stub

	// --- 2. 构造维度实体 ---
	dims := &entity.AnalyticsDimensions{
		Referer: event.Referer, // Referer 不需要解析
		Country: country,
		Browser: browser,
		OS:      os,
		Device:  device,
	}

	// --- 3. 调用仓储执行聚合 ---
	// (仓储层会处理事务)
	// 我们使用 event.AccessedAt 作为聚合的 'date'
	if err := l.svcCtx.AnalyticsRepo.IncrementDimensions(ctx, event.LinkID, event.AccessedAt, dims); err != nil {
		logx.WithContext(ctx).Errorf("Analytics Consumer: Failed to increment dimensions for linkID %d: %v", event.LinkID, err)
		return err // 返回错误，消息将被重试
	}

	return nil
}
