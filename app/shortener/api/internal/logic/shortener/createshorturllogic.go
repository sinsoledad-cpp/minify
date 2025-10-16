// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.1

package shortener

import (
	"context"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"lucid/app/shortener/api/internal/svc"
	"lucid/app/shortener/api/internal/types"
	"lucid/data/model/shortener"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateShortUrlLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建新的短链接
func NewCreateShortUrlLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateShortUrlLogic {
	return &CreateShortUrlLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateShortUrlLogic) CreateShortUrl(req *types.CreateShortUrlReq) (resp *types.CreateShortUrlResp, err error) {
	// 1. 获取当前用户ID
	userId, err := l.ctx.Value("userId").(json.Number).Int64()
	// if err != nil {
	// 	return nil, fmt.Errorf("获取用户信息失败")
	// }
	// userId, err := l.svcCtx.Auth.GetUserId(l.ctx)
	if err != nil {
		return nil, errors.Wrap(err, "获取用户ID失败")
	}

	// 2. 验证URL有效性
	valid, err := shortener.ValidateURL(req.OriginalUrl)
	if err != nil {
		return nil, errors.Wrap(err, "URL解析失败")
	}
	if !valid {
		return nil, errors.New("无效的URL格式，必须是有效的http或https URL")
	}

	// 3. 检查是否是短链接
	if shortener.IsShortURL(req.OriginalUrl, l.svcCtx.Config.ShortDomain) {
		return nil, errors.New("不能缩短已经是短链接的URL")
	}

	// 4. 计算原始URL的MD5哈希值用于去重
	urlMd5 := md5.Sum([]byte(req.OriginalUrl))
	urlMd5Hex := hex.EncodeToString(urlMd5[:])

	// 5. 检查当前用户是否已经为此URL创建过短链接
	existingUrl, err := l.svcCtx.ShortUrlsModel.FindOneByUserIdOriginalUrlMd5(l.ctx, uint64(userId), urlMd5Hex)
	if err == nil {
		// 已存在，直接返回
		return &types.CreateShortUrlResp{
			ShortKey:    existingUrl.ShortKey,
			OriginalUrl: existingUrl.OriginalUrl,
			ShortUrl:    fmt.Sprintf("%s/%s", l.svcCtx.Config.ShortDomain, existingUrl.ShortKey),
		}, nil
	} else if err != shortener.ErrNotFound {
		// 发生了其他错误
		return nil, errors.Wrap(err, "查询数据库失败")
	}

	// 6. 生成雪花算法ID
	snowflake, err := shortener.NewSnowflake(1) // 使用节点ID 1
	if err != nil {
		return nil, errors.Wrap(err, "创建雪花ID生成器失败")
	}
	snowflakeID := snowflake.NextID()

	// 7. 将雪花ID转换为Base62编码作为短链接key
	shortKey := shortener.ToBase62(snowflakeID)

	// 8. 处理过期时间
	var expiresAt sql.NullTime
	if req.ExpiresAt != "" {
		expTime, err := time.Parse(time.RFC3339, req.ExpiresAt)
		if err != nil {
			return nil, errors.Wrap(err, "过期时间格式无效，请使用RFC3339格式")
		}
		expiresAt = sql.NullTime{
			Time:  expTime,
			Valid: true,
		}
	}

	// 9. 创建短链接记录
	now := time.Now()
	shortUrl := &shortener.ShortUrls{
		UserId:         uint64(userId),
		ShortKey:       shortKey,
		OriginalUrl:    req.OriginalUrl,
		OriginalUrlMd5: urlMd5Hex,
		CreatedAt:      now,
		UpdatedAt:      now,
		ExpiresAt:      expiresAt,
	}

	// 10. 插入数据库
	_, err = l.svcCtx.ShortUrlsModel.Insert(l.ctx, shortUrl)
	if err != nil {
		return nil, errors.Wrap(err, "创建短链接失败")
	}

	// 11. 返回结果
	return &types.CreateShortUrlResp{
		ShortKey:    shortKey,
		OriginalUrl: req.OriginalUrl,
		ShortUrl:    fmt.Sprintf("%s/%s", l.svcCtx.Config.ShortDomain, shortKey),
	}, nil
}
