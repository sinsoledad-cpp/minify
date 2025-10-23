package entity

import (
	"database/sql"
	"errors"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

const (
	StatusActive   = "active"
	StatusExpired  = "expired"
	StatusInactive = "inactive"
	StatusAll      = "all"
)

var (
	ErrLinkNotFound            = errors.New("link not found") // ⭐ 直接定义
	ErrInvalidOriginalURL      = errors.New("invalid original URL")
	ErrInvalidExpiresIn        = errors.New("invalid expires in duration format")
	ErrLinkExpired             = errors.New("link has expired")
	ErrLinkInactive            = errors.New("link is inactive")
	ErrLinkNotFoundOrForbidden = errors.New("link not found or forbidden") // 用于 ABAC 检查
)

// Link 是短链接领域的实体（Entity）
type Link struct {
	ID             int64
	UserID         uint64 // 匹配 PO 类型
	ShortCode      string
	OriginalUrl    string
	VisitCount     uint64 // 匹配 PO 类型
	IsActive       bool   // 使用 bool 类型更符合领域语义
	ExpirationTime sql.NullTime
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      sql.NullTime
}

// NewLink 工厂函数，用于创建新的 Link 实体
// 它接收原始输入，并应用业务规则（如计算过期时间）
func NewLink(userID uint64, originalURL, shortCode, expiresIn string) (*Link, error) {
	if originalURL == "" { // 简单校验，可以用 UrlValidator 服务做得更好
		return nil, ErrInvalidOriginalURL
	}

	var expirationTime sql.NullTime
	if expiresIn != "" && expiresIn != "null" {
		// 解析相对时间，例如 "1h", "7d", "30m"
		duration, err := time.ParseDuration(expiresIn)
		if err != nil {
			// 尝试解析 ISO 8601 格式 (兼容之前的设计，虽然创建时推荐用相对时间)
			t, isoErr := time.Parse(time.RFC3339, expiresIn)
			if isoErr != nil {
				logx.Errorf("Failed to parse duration '%s': %v, and ISO8601 '%s': %v", expiresIn, err, expiresIn, isoErr)
				return nil, ErrInvalidExpiresIn
			}
			expirationTime = sql.NullTime{Time: t, Valid: true}
		} else {
			expirationTime = sql.NullTime{Time: time.Now().Add(duration), Valid: true}
		}
	}

	now := time.Now()
	return &Link{
		UserID:         userID,
		OriginalUrl:    originalURL,
		ShortCode:      shortCode,
		IsActive:       true, // 默认激活
		ExpirationTime: expirationTime,
		CreatedAt:      now, // 由 Repo 实现写入数据库时更新
		UpdatedAt:      now, // 由 Repo 实现写入数据库时更新
	}, nil
}

// IsExpired 检查链接是否已过期
func (l *Link) IsExpired() bool {
	return l.ExpirationTime.Valid && l.ExpirationTime.Time.Before(time.Now())
}

// CanRedirect 检查链接是否可用于重定向 (激活且未过期)
func (l *Link) CanRedirect() error {
	if !l.IsActive {
		return ErrLinkInactive
	}
	if l.IsExpired() {
		return ErrLinkExpired
	}
	return nil
}

// Activate 激活链接
func (l *Link) Activate() {
	l.IsActive = true
	l.UpdatedAt = time.Now()
}

// Deactivate 禁用链接
func (l *Link) Deactivate() {
	l.IsActive = false
	l.UpdatedAt = time.Now()
}

// UpdateDetails 更新链接信息
func (l *Link) UpdateDetails(originalURL *string, isActive *bool, expirationTime *string) error {
	updated := false
	if originalURL != nil {
		if *originalURL == "" {
			return ErrInvalidOriginalURL
		}
		l.OriginalUrl = *originalURL
		updated = true
	}
	if isActive != nil {
		l.IsActive = *isActive
		updated = true
	}
	if expirationTime != nil {
		if *expirationTime == "" || *expirationTime == "null" {
			l.ExpirationTime = sql.NullTime{Valid: false}
		} else {
			/*
				loc, err := time.LoadLocation("Asia/Shanghai")
				t, err := time.ParseInLocation(time.RFC3339, *expirationTime, loc)
			*/
			// 更新时只接受 ISO 8601 格式
			t, err := time.Parse(time.RFC3339, *expirationTime)
			if err != nil {
				return errors.New("invalid expiration time format, use ISO 8601")
			}
			l.ExpirationTime = sql.NullTime{Time: t, Valid: true}
		}
		updated = true
	}

	if updated {
		l.UpdatedAt = time.Now()
	}
	return nil
}

// MarkDeleted 标记为软删除
func (l *Link) MarkDeleted() {
	if !l.DeletedAt.Valid {
		l.DeletedAt = sql.NullTime{Time: time.Now(), Valid: true}
		l.UpdatedAt = l.DeletedAt.Time // 同时更新 updated_at
	}
}
