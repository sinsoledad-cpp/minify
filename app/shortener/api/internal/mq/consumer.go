package mq

import (
	"context"
	"minify/app/shortener/api/internal/config"
	"minify/app/shortener/api/internal/svc"
	"time"

	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/service"
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

func Consumers(c config.Config, ctx context.Context, svcContext *svc.ServiceContext) []service.Service {

	return []service.Service{
		//Listening for changes in consumption flow status
		kq.MustNewQueue(c.LinkEventConsumer, NewLinkAccessEventConsumer(ctx, svcContext)),
		kq.MustNewQueue(c.AnalyticsEventConsumer, NewAnalyticsEventConsumer(ctx, svcContext)),
		//.....
	}

}
