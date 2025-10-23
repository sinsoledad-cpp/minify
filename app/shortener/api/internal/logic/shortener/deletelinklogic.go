// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.1

package shortener

import (
	"context"
	"errors"
	"minify/app/shortener/domain/entity"
	"minify/common/utils/jwtx"

	"minify/app/shortener/api/internal/svc"
	"minify/app/shortener/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteLinkLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除短链接 (软删除)
func NewDeleteLinkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLinkLogic {
	return &DeleteLinkLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteLinkLogic) DeleteLink(req *types.DeleteLinkRequest) error {
	// 1. 从 JWT Context 获取用户 ID (身份认证)
	claims, err := jwtx.GetClaimsFromCtx(l.ctx)
	if err != nil {
		return errors.New("invalid token")
	}
	userId := uint64(claims.UserID)

	// 2. 查找链接实体 (获取聚合根)
	// LinkRepo.FindByCode 已经处理了 "已软删除" 的情况，会返回 ErrLinkNotFound
	//
	link, err := l.svcCtx.LinkRepo.FindByCode(l.ctx, req.Code)
	if err != nil {
		if errors.Is(err, entity.ErrLinkNotFound) {
			// 如果链接已经不存在 (或已被删除)，
			// 那么 "删除" 这个操作已经完成了 (保持幂等性)，直接返回成功。
			logx.Infof("Link %s not found or has deleted", req.Code)
			return nil
		}
		// 其他数据库错误
		l.Logger.Errorf("FindByCode error: %v", err)
		return err
	}

	// 3. 检查所有权 (DDD 核心：应用层执行授权策略)
	// 只有链接的所有者才能删除它
	if link.UserID != userId {
		// 出于安全考虑，我们不告诉用户“链接存在但你没权限”
		// 而是假装它不存在（同样为了保持幂等性）。
		return nil
	}

	// 4. 调用仓储执行软删除
	// 仓储(Repository)的 Delete 方法会封装“调用实体的 MarkDeleted()”
	// 和“调用 model 的 Update()”这两个基础设施操作
	//
	if err := l.svcCtx.LinkRepo.Delete(l.ctx, link); err != nil {
		l.Logger.Errorf("LinkRepo.Delete error: %v", err)
		return errors.New("failed to delete link")
	}

	return nil
}
