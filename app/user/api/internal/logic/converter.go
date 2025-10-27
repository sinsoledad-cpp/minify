package logic

import (
	"minify/app/user/api/internal/types"
	"minify/app/user/domain/entity"
	"time"
)

// Converter 负责 user 服务的 DTO 与 Entity 转换
type Converter struct {
	// 目前 user 服务的转换不需要依赖（例如 shortDomain）
	// 如果未来需要，可以在这里添加
}

// NewConverter 创建一个新的转换器实例
func NewConverter() *Converter {
	return &Converter{}
}

// ToUserInfoResponse 将 User 实体转换为 UserInfoResponse DTO
func (c *Converter) ToUserInfoResponse(e *entity.User) *types.UserInfoResponse {
	if e == nil {
		return nil
	}

	return &types.UserInfoResponse{
		Id:        e.ID,
		Username:  e.Username,
		Email:     e.Email,
		Role:      e.Role,
		CreatedAt: e.CreatedAt.Format(time.RFC3339),
	}
}

// ToUserInfoResponseList (辅助函数) 转换实体列表
func (c *Converter) ToUserInfoResponseList(entities []*entity.User) []types.UserInfoResponse {
	dtos := make([]types.UserInfoResponse, len(entities))
	for i, e := range entities {
		// 复用单个转换逻辑
		dto := c.ToUserInfoResponse(e)
		if dto != nil {
			dtos[i] = *dto
		}
	}
	return dtos
}
