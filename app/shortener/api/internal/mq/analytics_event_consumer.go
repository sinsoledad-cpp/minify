package mq

import (
	"context"
	"encoding/json"
	"minify/app/shortener/api/internal/svc"
	"minify/app/shortener/domain/entity"
	"net"     // (不变)
	"net/url" // (不变)

	"github.com/oschwald/geoip2-golang"    // (不变)
	"github.com/ua-parser/uap-go/uaparser" // (不变)
	"github.com/zeromicro/go-zero/core/logx"
)

// (AnalyticsEventConsumer 结构体... 保持不变)
type AnalyticsEventConsumer struct {
	ctx      context.Context
	svcCtx   *svc.ServiceContext
	uaParser *uaparser.Parser
	geoIPDB  *geoip2.Reader
}

// (NewAnalyticsEventConsumer 函数... 保持不变)
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
		uaParser: ua,
		geoIPDB:  db,
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

	// ⭐ 改进：所有维度首先默认为 "unknown"
	country := "unknown"
	referer := "unknown"
	browser := "unknown"
	os := "unknown"
	device := "unknown"

	// --- 1a. 处理 Country ---
	if event.IpAddress != "" {
		ip := net.ParseIP(event.IpAddress)
		if ip != nil {
			if ip.IsLoopback() || ip.IsPrivate() {
				country = "Local Network" // 明确标记本地/私有网络
			} else {
				record, err := l.geoIPDB.City(ip)
				if err == nil && record != nil && record.Country.Names["en"] != "" {
					country = record.Country.Names["en"]
				}
				// (如果 GeoIP 查询失败, 保持 "unknown")
			}
		}
		// (如果 IP 解析失败, 保持 "unknown")
	}

	// --- 1b. 处理 Referer ---
	if event.Referer == "" {
		referer = "Direct Entry" // 明确标记直接访问
	} else {
		parsedUrl, err := url.Parse(event.Referer)
		if err == nil && parsedUrl.Hostname() != "" {
			referer = parsedUrl.Hostname() // 提取域名
		} else {
			referer = "Invalid Referer" // 标记无效的 Referer
		}
	}

	// --- 1c. 处理 User Agent (Browser, OS, Device) ---
	if event.UserAgent != "" {
		client := l.uaParser.Parse(event.UserAgent)

		// 浏览器
		if client.UserAgent.Family != "Other" && client.UserAgent.Family != "" {
			browser = client.UserAgent.Family
		}

		// 操作系统
		if client.Os.Family != "Other" && client.Os.Family != "" {
			os = client.Os.Family
		}

		// 设备
		if client.Device.Family != "Other" && client.Device.Family != "" {
			device = client.Device.Family // 明确的设备 (e.g., "iPhone")
		} else if client.Device.Family == "Other" {
			// "Other" 是 uap-go 对 PC/Mac/Linux 桌面端的通用标识
			device = "Desktop"
		}
		// (如果 Family 是 "", 保持 "unknown")
	}

	// --- 2. 构造维度实体 ---
	dims := &entity.AnalyticsDimensions{
		Referer: referer,
		Country: country,
		Browser: browser,
		OS:      os,
		Device:  device,
	}

	// --- 3. 调用仓储执行聚合 ---
	if err := l.svcCtx.AnalyticsRepo.IncrementDimensions(ctx, event.LinkID, event.AccessedAt, dims); err != nil {
		logx.WithContext(ctx).Errorf("Analytics Consumer: Failed to increment dimensions for linkID %d: %v", event.LinkID, err)
		return err // 返回错误，消息将被重试
	}

	return nil
}
