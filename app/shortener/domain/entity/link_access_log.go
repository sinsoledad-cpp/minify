package entity

import "time"

// LinkAccessLog 是访问日志的领域实体
type LinkAccessLog struct {
	ID          int64
	LinkID      int64
	ShortCode   string
	AccessedAt  time.Time
	IpAddress   string
	UserAgent   string
	Referer     string
	GeoCountry  string // 这些字段由 ETL 填充，但在实体中定义
	GeoCity     string
	DeviceType  string
	BrowserName string
	OsName      string
}

// NewLinkAccessLog 是创建新日志实体的工厂函数
// 消费者将调用这个函数
func NewLinkAccessLog(linkID int64, shortCode string, accessedAt time.Time, ip, ua, referer string) *LinkAccessLog {
	return &LinkAccessLog{
		LinkID:     linkID,
		ShortCode:  shortCode,
		AccessedAt: accessedAt,
		IpAddress:  ip,
		UserAgent:  ua,
		Referer:    referer,
		// 其他 Geo/UA 字段默认为空
	}
}
