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

type GetShortUrlDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取单个短链接的详细信息和统计数据
func NewGetShortUrlDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetShortUrlDetailLogic {
	return &GetShortUrlDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetShortUrlDetailLogic) GetShortUrlDetail(req *types.GetShortUrlDetailReq) (resp *types.GetShortUrlDetailResp, err error) {
	// 1. 获取用户ID，用于权限校验
	// 此接口受JWT保护，因此必须能获取到userId
	userId, err := l.getUserIdFromCtx()
	if err != nil {
		return nil, errors.Wrap(err, "获取用户信息失败")
	}

	// 2. 查询短链接基本信息
	shortUrl, err := l.svcCtx.ShortUrlsModel.FindOneByShortKey(l.ctx, req.ShortKey)
	if err != nil {
		if err == shortener.ErrNotFound {
			return nil, errors.New("短链接不存在")
		}
		l.Logger.Errorf("查询短链接失败: %v", err)
		return nil, errors.New("查询失败")
	}

	// 3. 权限校验：确保该短链接属于当前用户
	if shortUrl.UserId != userId {
		// 出于安全考虑，返回“不存在”而不是“无权限”
		return nil, errors.New("短链接不存在")
	}

	// 4. 查询总体统计数据 (TotalClicks, UniqueVisitors)
	// 我们使用 data/model/shortener/urlanalyticsmodel.go 中已有的方法
	totalClicks, uniqueVisitors, err := l.svcCtx.UrlAnalyticsModel.GetTotalClicksAndUniqueVisitorsByShortUrlId(l.ctx, shortUrl.Id)
	if err != nil {
		// 统计查询失败可以降级处理，不阻塞主流程
		l.Logger.Errorf("查询总体统计数据失败: %v", err)
		// 即使失败，也继续执行，只是统计数据为0
		totalClicks = 0
		uniqueVisitors = 0
	}

	// 5. 组装基本信息 (Info)
	info := types.ShortUrlInfo{
		ShortKey:       shortUrl.ShortKey,
		OriginalUrl:    shortUrl.OriginalUrl,
		ShortUrl:       fmt.Sprintf("%s/%s", l.svcCtx.Config.ShortDomain, shortUrl.ShortKey),
		CreatedAt:      shortUrl.CreatedAt.Format(time.RFC3339),
		TotalClicks:    totalClicks,
		UniqueVisitors: uniqueVisitors,
	}
	if shortUrl.ExpiresAt.Valid {
		info.ExpiresAt = shortUrl.ExpiresAt.Time.Format(time.RFC3339)
	} else {
		info.ExpiresAt = "" // 表示永不过期
	}

	var dailySummaryList []types.DailySummary

	// 6. 查询每日统计 (DailySummary)
	// !!! 注意：这需要您在 AggDailySummaryModel 中实现 FindAllByShortUrlId 方法
	dailyStats, err := l.svcCtx.AggDailySummaryModel.FindAllByShortUrlId(l.ctx, shortUrl.Id)
	if err != nil {
		l.Logger.Errorf("查询每日统计失败: %v", err)
		// 降级处理，返回空列表
		dailySummaryList = []types.DailySummary{}
	} else {
		dailySummaryList = make([]types.DailySummary, 0, len(dailyStats))
		for _, stat := range dailyStats {
			dailySummaryList = append(dailySummaryList, types.DailySummary{
				Date:           stat.SummaryDate.Format("2006-01-02"), // YYYY-MM-DD
				TotalClicks:    int64(stat.TotalClicks),               // AggDailySummary.TotalClicks 是 uint64
				UniqueVisitors: int64(stat.UniqueVisitors),            // AggDailySummary.UniqueVisitors 是 uint64
			})
		}
	}

	var recentRecordsList []types.AnalyticsRecord

	// 7. 查询最近访问记录 (RecentRecords)
	// !!! 注意：这需要您在 UrlAnalyticsModel 中实现 FindRecentRecords 方法
	const recentRecordLimit = 20 // 如 types.go 中注释的示例
	recentLogs, err := l.svcCtx.UrlAnalyticsModel.FindRecentRecords(l.ctx, shortUrl.Id, recentRecordLimit)
	if err != nil {
		l.Logger.Errorf("查询最近访问记录失败: %v", err)
		// 降级处理，返回空列表
		recentRecordsList = []types.AnalyticsRecord{}
	} else {
		recentRecordsList = make([]types.AnalyticsRecord, 0, len(recentLogs))
		for _, log := range recentLogs {
			recentRecordsList = append(recentRecordsList, types.AnalyticsRecord{
				IpAddress: log.IpAddress,
				UserAgent: log.UserAgent.String, // 从 sql.NullString 获取
				Referer:   log.Referer.String,   // 从 sql.NullString 获取
				CreatedAt: log.CreatedAt.Format(time.RFC3339),
			})
		}
	}

	// 8. 组装最终响应
	return &types.GetShortUrlDetailResp{
		Info:          info,
		DailySummary:  dailySummaryList,
		RecentRecords: recentRecordsList,
	}, nil
}

// getUserIdFromCtx 从context中获取用户ID
// (这个辅助函数在 ListMyShortUrlsLogic 中已存在，复制过来)
func (l *GetShortUrlDetailLogic) getUserIdFromCtx() (uint64, error) {
	userIdVal := l.ctx.Value("userId")
	if userIdVal == nil {
		return 0, errors.New("无法从context中获取userId, value is nil")
	}

	userId, ok := userIdVal.(json.Number)
	if !ok {
		return 0, errors.New(fmt.Sprintf("userId类型不是json.Number, 实际类型: %T", userIdVal))
	}

	uid, err := userId.Int64()
	if err != nil {
		return 0, errors.Wrap(err, "userId类型转换失败")
	}
	return uint64(uid), nil
}
