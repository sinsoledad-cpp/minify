package entity

import (
	"errors"
	"minify/app/user/data/model" // 引用 goctl 生成的 model
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrPasswordMismatch = errors.New("invalid username or password")
	ErrUserNotFound     = model.ErrNotFound // 复用 model 的错误
)

// User 是用户领域的实体（Entity）
type User struct {
	ID           int64
	Username     string
	Email        string
	PasswordHash string
	Role         string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// NewUser 是创建用户的工厂函数，确保业务不变量
func NewUser(username, email, plainPassword string) (*User, error) {
	if username == "" || email == "" || plainPassword == "" {
		return nil, errors.New("username, email, and password are required")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &User{
		Username:     username,
		Email:        email,
		PasswordHash: string(hash),
		Role:         "shortener", // 匹配你的 SQL 默认值
	}, nil
}

// CheckPassword 是属于 User 实体的业务方法
func (u *User) CheckPassword(plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(plainPassword))
	return err == nil
}

// --- DTO 转换 ---
// (DDD 中，实体与持久化对象(PO)的转换是必要的)

// ToModel 将领域实体(Entity)转换为数据模型(PO)
func (u *User) ToModel() *model.Users {
	return &model.Users{
		Id:           uint64(u.ID), // 注意类型转换
		Username:     u.Username,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		Role:         u.Role,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}

// FromModel 将数据模型(PO)转换为领域实体(Entity)
func FromModel(m *model.Users) *User {
	return &User{
		ID:           int64(m.Id), // 注意类型转换
		Username:     m.Username,
		Email:        m.Email,
		PasswordHash: m.PasswordHash,
		Role:         m.Role,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}
