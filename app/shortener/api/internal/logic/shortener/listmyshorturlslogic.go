package shortener

import (
	"context"
	"encoding/json"
	"fmt"
	"lucid/app/shortener/api/internal/svc"
	"lucid/app/shortener/api/internal/types"
	"lucid/data/model/shortener"
	"time"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type ListMyShortUrlsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListMyShortUrlsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListMyShortUrlsLogic {
	return &ListMyShortUrlsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// ListMyShortUrls 获取当前用户的所有短链接
func (l *ListMyShortUrlsLogic) ListMyShortUrls() (resp *types.ListShortUrlsResp, err error) {
	userId, err := l.getUserIdFromCtx()
	if err != nil {
		return nil, errors.Wrap(err, "获取用户ID失败")
	}

	// 1. 第一次DB查询：获取用户的所有短链接
	shortUrls, err := l.svcCtx.ShortUrlsModel.FindAllByUserId(l.ctx, userId)
	if err != nil {
		return nil, errors.Wrap(err, "查询用户短链接失败")
	}

	if len(shortUrls) == 0 {
		return &types.ListShortUrlsResp{Urls: []types.ShortUrlInfo{}}, nil
	}

	// 2. 从短链接列表中提取所有ID
	shortUrlIds := make([]uint64, 0, len(shortUrls))
	for _, su := range shortUrls {
		shortUrlIds = append(shortUrlIds, su.Id)
	}

	// 3. 第二次DB查询：一次性获取所有短链接的统计数据
	statsMap, err := l.svcCtx.UrlAnalyticsModel.GetAnalyticsStatsByShortUrlIds(l.ctx, shortUrlIds)
	if err != nil {
		// 即使统计查询失败，我们也可以选择优雅降级，只返回基本信息
		logx.Errorf("批量获取短链接统计数据失败, err: %v", err)
		statsMap = make(map[uint64]shortener.AnalyticsStats) // 创建一个空map以避免下面代码 panic
	}

	// 4. 在内存中组装最终结果，不再有DB查询
	var shortUrlInfos []types.ShortUrlInfo
	for _, su := range shortUrls {
		// 从map中查找统计数据
		stats, ok := statsMap[su.Id]
		if !ok {
			// 如果map中没有这个id的数据，说明它还没有任何访问记录
			stats = shortener.AnalyticsStats{TotalClicks: 0, UniqueVisitors: 0}
		}

		shortUrlInfo := types.ShortUrlInfo{
			ShortKey:       su.ShortKey,
			OriginalUrl:    su.OriginalUrl,
			ShortUrl:       fmt.Sprintf("%s/%s", l.svcCtx.Config.ShortDomain, su.ShortKey),
			CreatedAt:      su.CreatedAt.Format(time.RFC3339),
			TotalClicks:    stats.TotalClicks,
			UniqueVisitors: stats.UniqueVisitors,
		}

		if su.ExpiresAt.Valid {
			shortUrlInfo.ExpiresAt = su.ExpiresAt.Time.Format(time.RFC3339)
		} else {
			shortUrlInfo.ExpiresAt = "" // 表示永不过期
		}
		shortUrlInfos = append(shortUrlInfos, shortUrlInfo)
	}

	return &types.ListShortUrlsResp{
		Urls: shortUrlInfos,
	}, nil
}

// getUserIdFromCtx 从context中获取用户ID
func (l *ListMyShortUrlsLogic) getUserIdFromCtx() (uint64, error) {
	userId, ok := l.ctx.Value("userId").(json.Number)
	if !ok {
		return 0, errors.New("无法从context中获取userId")
	}
	uid, err := userId.Int64()
	if err != nil {
		return 0, errors.Wrap(err, "userId类型转换失败")
	}
	return uint64(uid), nil
}
