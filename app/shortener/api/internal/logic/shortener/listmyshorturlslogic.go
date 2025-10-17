package shortener

import (
	"context"
	"encoding/json"
	"fmt"
	"lucid/app/shortener/api/internal/svc"
	"lucid/app/shortener/api/internal/types"
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

	shortUrls, err := l.svcCtx.ShortUrlsModel.FindAllByUserId(l.ctx, userId)
	if err != nil {
		return nil, errors.Wrap(err, "查询用户短链接失败")
	}

	if len(shortUrls) == 0 {
		return &types.ListShortUrlsResp{Urls: []types.ShortUrlInfo{}}, nil
	}

	var shortUrlInfos []types.ShortUrlInfo
	for _, su := range shortUrls {
		totalClicks, uniqueVisitors, err := l.svcCtx.UrlAnalyticsModel.GetTotalClicksAndUniqueVisitorsByShortUrlId(l.ctx, su.Id)
		if err != nil {
			logx.Errorf("获取短链接统计数据失败, shortUrlId: %d, err: %v", su.Id, err)
			// 即使统计数据获取失败，也继续处理其他短链接，并将点击量设为0
			totalClicks = 0
			uniqueVisitors = 0
		}

		shortUrlInfo := types.ShortUrlInfo{
			ShortKey:       su.ShortKey,
			OriginalUrl:    su.OriginalUrl,
			ShortUrl:       fmt.Sprintf("%s/%s", l.svcCtx.Config.ShortDomain, su.ShortKey),
			CreatedAt:      su.CreatedAt.Format(time.RFC3339),
			TotalClicks:    totalClicks,
			UniqueVisitors: uniqueVisitors,
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
