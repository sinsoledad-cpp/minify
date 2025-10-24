// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.1

package shortener

import (
	"context"
	"errors"
	"minify/app/shortener/api/internal/logic"
	"minify/app/shortener/api/internal/logic/errcode"
	"minify/app/shortener/api/internal/svc"
	"minify/app/shortener/api/internal/types"
	"minify/app/shortener/domain/entity"
	"minify/common/utils/codec"
	"minify/common/utils/jwtx"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLinkLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建短链接
func NewCreateLinkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLinkLogic {
	return &CreateLinkLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateLinkLogic) CreateLink(req *types.CreateLinkRequest) (resp *types.CreateLinkResponse, err error) {
	// 1. 从 JWT Context 获取用户 ID
	claims, err := jwtx.GetClaimsFromCtx(l.ctx) //
	if err != nil {
		return nil, errcode.ErrInvalidToken
	}
	userId := uint64(claims.UserID)

	var shortCode string

	// 2. 检查是自定义短码还是自动生成
	if req.CustomCode != "" {
		// --- 策略二：用户自定义 ---
		shortCode = req.CustomCode
		_, err := l.svcCtx.LinkRepo.FindByCode(l.ctx, shortCode)
		if err == nil {
			return nil, errcode.ErrCustomCodeExists
		}
		if !errors.Is(err, entity.ErrLinkNotFound) {
			l.Logger.Errorf("FindByCode for custom code error: %v", err)
			return nil, errcode.ErrInternalError
		}
	} else {
		// --- 策略一：自动生成 (Snowflake) ---
		// a. 调用 IdGenerator 获取唯一 ID
		id, err := l.svcCtx.IdGenerator.NextID(l.ctx) //
		if err != nil {
			l.Logger.Errorf("IdGenerator.NextID (Snowflake) error: %v", err)
			// ⭐ 修正：使用 httpx.NewCodeError
			return nil, errcode.ErrIdGenerateError
		}

		// b. 使用 Base62 编码 (工具来自 common/utils/codec)
		shortCode = codec.Base62Encode(id) //
	}

	// 3. 创建 Link 实体
	link, err := entity.NewLink(userId, req.OriginalUrl, shortCode, req.ExpiresIn) //
	if err != nil {
		l.Logger.Infof("NewLink validation error: %v", err)
		// ⭐ 修正：使用 httpx.NewCodeError
		return nil, errcode.ErrInternalError
	}

	// 4. 保存到仓储 (Repository)
	if err := l.svcCtx.LinkRepo.Create(l.ctx, link); err != nil { //
		l.Logger.Errorf("LinkRepo.Create error: %v", err)
		// ⭐ 修正：使用 httpx.NewCodeError
		return nil, errcode.ErrInternalError
	}

	// 5. 转换 DTO 并返回
	// (确保 converter.go 也在 logic/shortener 目录下)
	dtoLink := logic.ToTypesLink(link)

	return &types.CreateLinkResponse{
		Link: *dtoLink,
	}, nil
}
