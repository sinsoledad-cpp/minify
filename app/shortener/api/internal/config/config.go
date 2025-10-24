// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.1

package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	Auth struct {
		AccessSecret string
		AccessExpire int64
	}
	Database struct { // ⭐ 数据库配置
		DataSource string
	}
	CacheRedis cache.CacheConf // ⭐ Redis 缓存配置 (用于 linksModel)
	// CasbinModelPath string       // ⭐ Casbin 模型文件路径 (如果需要在此初始化 Casbin)
	Snowflake struct {
		WorkerId int64
	}
	Casbin struct {
		ModelPath  string
		PolicyPath string
	}
}
