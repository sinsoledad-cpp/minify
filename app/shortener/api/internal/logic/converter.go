package logic

import (
	"minify/app/shortener/api/internal/types"
	"minify/app/shortener/domain/entity"
	"time"
)

// ToTypesLink 将 Link 实体转换为 Link DTO
func ToTypesLink(e *entity.Link) *types.Link {
	if e == nil {
		return nil
	}

	expTime := ""
	if e.ExpirationTime.Valid {
		// 格式化为 ISO 8601，与 API 定义一致
		expTime = e.ExpirationTime.Time.Format(time.RFC3339)
	}

	return &types.Link{
		Id:             e.ID,
		ShortCode:      e.ShortCode,
		OriginalUrl:    e.OriginalUrl,
		VisitCount:     int64(e.VisitCount), // 类型转换
		IsActive:       e.IsActive,
		ExpirationTime: expTime,
		CreatedAt:      e.CreatedAt.Format(time.RFC3339),
	}
}
