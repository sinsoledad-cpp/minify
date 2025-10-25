package mq

import (
	"context"
	"encoding/json"
	"minify/app/shortener/api/internal/svc"
	"minify/app/shortener/domain/entity"
	"net" // ⭐ 1. 导入 net 包

	"github.com/oschwald/geoip2-golang"    // ⭐ 2. 导入 geoip2
	"github.com/ua-parser/uap-go/uaparser" // ⭐ 3. 导入 uaparser
	"github.com/zeromicro/go-zero/core/logx"
)

// AnalyticsEventConsumer 是处理聚合分析的消费者
type AnalyticsEventConsumer struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	// ⭐ 4. 添加解析器实例
	uaParser *uaparser.Parser
	geoIPDB  *geoip2.Reader
}

// NewAnalyticsEventConsumer 创建一个新的消费者实例
func NewAnalyticsEventConsumer(ctx context.Context, svcCtx *svc.ServiceContext) *AnalyticsEventConsumer {
	// ⭐ 5. 在构造函数中初始化解析器
	// 我们假设 goctl 启动时的工作目录是项目根目录
	// 并且你已经将所需文件放在了 'etc/' 目录下

	// 初始化 UserAgent 解析器
	// 你需要从 https://github.com/ua-parser/uap-core/blob/master/regexes.yaml 下载
	ua, err := uaparser.New("etc/regexes.yaml")
	if err != nil {
		// 在服务启动时失败，直接 panic，类似 Casbin
		logx.Must(err)
	}

	// 初始化 GeoIP 解析器
	// 你需要从 MaxMind 下载 GeoLite2-City.mmdb
	db, err := geoip2.Open("etc/GeoLite2-City.mmdb")
	if err != nil {
		logx.Must(err)
	}

	return &AnalyticsEventConsumer{
		ctx:      ctx,
		svcCtx:   svcCtx,
		uaParser: ua, // ⭐ 6. 注入实例
		geoIPDB:  db, // ⭐ 6. 注入实例
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
	// (已移除 city 变量)
	country := "unknown"
	browser := "unknown"
	os := "unknown"
	device := "unknown"
	referer := "unknown"

	// a. 解析 IP 地址
	if event.IpAddress != "" {
		ip := net.ParseIP(event.IpAddress)
		if ip != nil {
			// 仍然使用 City() 方法，因为它包含 Country 信息
			record, err := l.geoIPDB.City(ip)
			if err == nil && record != nil {
				if record.Country.Names["en"] != "" {
					country = record.Country.Names["en"]
				}
				// (已移除 city 的提取逻辑)
			}
		}
	}

	// b. 解析 User Agent
	if event.UserAgent != "" {
		client := l.uaParser.Parse(event.UserAgent)
		if client.UserAgent.Family != "Other" && client.UserAgent.Family != "" {
			browser = client.UserAgent.Family
		}
		if client.Os.Family != "Other" && client.Os.Family != "" {
			os = client.Os.Family
		}
		if client.Device.Family != "Other" && client.Device.Family != "" {
			device = client.Device.Family
		}
	}

	// c. 处理 Referer
	if event.Referer != "" {
		referer = event.Referer
		// 你未来可以在这里添加域名提取逻辑
	}
	// 如果 Referer 为空，它将保持 "unknown"

	// --- 2. 构造维度实体 ---
	dims := &entity.AnalyticsDimensions{
		Referer: referer,
		Country: country,
		Browser: browser,
		OS:      os,
		Device:  device,
	}

	// --- 3. 调用仓储执行聚合 ---
	// 仓储层 (analytics_repo_impl.go) 已经实现了事务处理
	// 我们使用 event.AccessedAt 作为聚合的 'date'
	if err := l.svcCtx.AnalyticsRepo.IncrementDimensions(ctx, event.LinkID, event.AccessedAt, dims); err != nil {
		logx.WithContext(ctx).Errorf("Analytics Consumer: Failed to increment dimensions for linkID %d: %v", event.LinkID, err)
		return err // 返回错误，消息将被重试
	}

	return nil
}
