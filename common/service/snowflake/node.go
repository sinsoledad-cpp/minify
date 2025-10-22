package snowflake

import (
	"context" // 保留 context 占位符，尽管 snowflake.Node.Generate() 不使用它

	"github.com/bwmarrin/snowflake"
)

// Generator 是一个具体的 Snowflake ID 生成器实现
type Generator struct {
	node *snowflake.Node
}

// NewGenerator 创建一个新的 Snowflake Generator
// 它需要调用方（如 shortener-api）传入其唯一的 workerId
func NewGenerator(workerId int64) (*Generator, error) {
	node, err := snowflake.NewNode(workerId)
	if err != nil {
		return nil, err
	}

	return &Generator{
		node: node,
	}, nil
}

// NextID 生成下一个 ID
// (注意：这个方法签名必须与 app/shortener/domain/service/id_generator.go 中的接口匹配)
func (g *Generator) NextID(ctx context.Context) (uint64, error) {
	id := g.node.Generate()
	return uint64(id), nil
}
