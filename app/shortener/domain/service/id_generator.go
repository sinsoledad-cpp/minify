package service

import "context"

// IdGenerator 定义分布式唯一 ID 生成器接口
type IdGenerator interface {
	// NextID 生成下一个全局唯一的 ID
	NextID(ctx context.Context) (uint64, error)
}
