// growth-partner/backend/internal/model/user.go
// 用户模型：支持三种角色——学生、班主任/园长、家长
// 严格遵循未成年人隐私保护：真实姓名 AES 加密存储

package model

import "time"

// UserRole 用户角色枚举
type UserRole string

const (
	RoleStudent UserRole = "student" // 学生（小朋友端）
	RoleTeacher UserRole = "teacher" // 班主任/科任老师/园长端
	RoleParent  UserRole = "parent"  // 家长端
	RoleAdmin   UserRole = "admin"   // 系统管理员
)

// User 用户基础账号表
// 真实姓名使用 AES-256 加密存储，数据库内不存明文
type User struct {
	Base

	// ─── 账号信息 ───
	Username     string   `gorm:"uniqueIndex;size:64;not null"  json:"username"`
	PasswordHash string   `gorm:"size:128;not null"             json:"-"` // bcrypt 哈希，永不返回前端
	Role         UserRole `gorm:"size:16;index;not null"        json:"role"`

	// ─── 个人信息（加密存储，PIPL 合规） ───
	RealNameEnc string `gorm:"size:256"                      json:"-"` // AES 加密的真实姓名
	PhoneEnc    string `gorm:"size:256"                      json:"-"` // AES 加密的手机号
	AvatarURL   string `gorm:"size:512"                      json:"avatar_url,omitempty"`

	// ─── 账号状态 ───
	IsActive    bool       `gorm:"default:true;index"            json:"is_active"`
	LastLoginAt *time.Time `gorm:"index"                       json:"last_login_at,omitempty"`

	// ─── 关联（不使用外键约束，保持灵活性） ───
	// 一个用户只有一种角色，通过角色决定关联哪个实体
	// 学生用户关联 Child，老师关联 Class，家长通过 ParentChild 关联 Child
}

// TableName 指定表名
func (User) TableName() string { return "users" }

// ParentChildRelation 家长-学生绑定关系表
// 一个家长可以绑定多个孩子（如兄弟姐妹），一个孩子可有多个家长（父/母）
type ParentChildRelation struct {
	Base
	ParentUserID uint64 `gorm:"index;not null"  json:"parent_user_id"`
	ChildID      uint64 `gorm:"index;not null"  json:"child_id"`
	Relation     string `gorm:"size:16"         json:"relation"` // "father"/"mother"/"guardian"
}

func (ParentChildRelation) TableName() string { return "parent_child_relations" }
