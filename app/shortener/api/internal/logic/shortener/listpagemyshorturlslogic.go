// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.1

package shortener

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"lucid/app/shortener/api/internal/svc"
	"lucid/app/shortener/api/internal/types"
	"lucid/data/model/shortener"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type ListPageMyShortUrlsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取当前用户的所有短链接
func NewListPageMyShortUrlsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListPageMyShortUrlsLogic {
	return &ListPageMyShortUrlsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListPageMyShortUrlsLogic) ListPageMyShortUrls(req *types.ListPageMyShortUrlsReq) (resp *types.ListPageMyShortUrlsResp, err error) {
	userId, err := l.getUserIdFromCtx()
	if err != nil {
		return nil, errors.Wrap(err, "获取用户ID失败")
	}

	// 1. 调用新的分页查询方法
	shortUrls, total, err := l.svcCtx.ShortUrlsModel.FindPagedByUserId(l.ctx, userId, req.Page, req.PageSize)
	if err != nil {
		return nil, errors.Wrap(err, "查询用户短链接失败")
	}

	if len(shortUrls) == 0 {
		return &types.ListPageMyShortUrlsResp{
			Urls:  []types.ShortUrlInfo{},
			Total: 0,
		}, nil
	}

	// 2. 批量获取统计数据的逻辑保持不变
	shortUrlIds := make([]uint64, 0, len(shortUrls))
	for _, su := range shortUrls {
		shortUrlIds = append(shortUrlIds, su.Id)
	}

	statsMap, err := l.svcCtx.UrlAnalyticsModel.GetAnalyticsStatsByShortUrlIds(l.ctx, shortUrlIds)
	if err != nil {
		logx.Errorf("批量获取短链接统计数据失败, err: %v", err)
		statsMap = make(map[uint64]shortener.AnalyticsStats)
	}

	// 3. 组装当前页的数据
	var shortUrlInfos []types.ShortUrlInfo
	for _, su := range shortUrls {
		stats, ok := statsMap[su.Id]
		if !ok {
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
			shortUrlInfo.ExpiresAt = ""
		}
		shortUrlInfos = append(shortUrlInfos, shortUrlInfo)
	}

	// 4. 返回最终的分页响应
	return &types.ListPageMyShortUrlsResp{
		Urls:  shortUrlInfos,
		Total: total,
	}, nil
}

// getUserIdFromCtx 从context中获取用户ID
func (l *ListPageMyShortUrlsLogic) getUserIdFromCtx() (uint64, error) {
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
