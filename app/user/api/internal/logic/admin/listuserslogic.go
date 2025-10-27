// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.1

package admin

import (
	"context"
	"minify/app/user/api/internal/logic/errcode"
	"minify/app/user/api/internal/svc"
	"minify/app/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListUsersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取所有用户列表 (管理员)
func NewListUsersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListUsersLogic {
	return &ListUsersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListUsersLogic) ListUsers(req *types.ListUsersRequest) (resp *types.ListUsersResponse, err error) {
	// 1. (鉴权) Casbin 中间件已经确保了只有 admin 才能访问
	// 我们不需要在这里检查 claims.Role

	// 2. 处理分页参数
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	switch {
	case pageSize <= 0:
		pageSize = 20 // 默认值
	case pageSize > 100:
		pageSize = 100 // 最大值
	}

	// 3. 调用仓储层(Repository)
	users, total, err := l.svcCtx.UserRepo.ListAll(l.ctx, page, pageSize)
	if err != nil {
		l.Logger.Errorf("ListAll error: %v", err)
		return nil, errcode.ErrInternalError
	}

	// 4. 将领域实体(Entity)列表转换为 DTO 列表
	dtoUsers := l.svcCtx.Converter.ToUserInfoResponseList(users)

	// 5. 返回响应
	return &types.ListUsersResponse{
		Users: dtoUsers,
		Total: total,
	}, nil
}
