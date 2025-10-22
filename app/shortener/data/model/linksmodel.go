package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ LinksModel = (*customLinksModel)(nil)

type (
	// LinksModel is an interface to be customized, add more methods here,
	// and implement the added methods in customLinksModel.
	LinksModel interface {
		linksModel
	}

	customLinksModel struct {
		*defaultLinksModel
	}
)

// NewLinksModel returns a model for the database table.
func NewLinksModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) LinksModel {
	return &customLinksModel{
		defaultLinksModel: newLinksModel(conn, c, opts...),
	}
}
