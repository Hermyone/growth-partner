// growth-partner/backend/internal/model/admin_permission.go
// 管理员权限模型

package model

import "time"

// PermissionType 权限类型
type PermissionType string

const (
	PermissionManage PermissionType = "manage" // 管理权限
	PermissionScore  PermissionType = "score"  // 打分权限
	PermissionView   PermissionType = "view"   // 查看权限
)

// AdminPermission 管理员权限分配表
type AdminPermission struct {
	Base

	GrantedBy       uint64         `gorm:"index;not null" json:"granted_by"`       // 授权人用户ID
	TeacherUserID   uint64         `gorm:"index;not null" json:"teacher_user_id"`   // 被授权教师用户ID
	ClassID         uint64         `gorm:"not null" json:"class_id"`              // 班级ID
	PermissionType  PermissionType `gorm:"size:20;not null" json:"permission_type"`  // 权限类型
	SchoolYear      string         `gorm:"size:20;not null" json:"school_year"`     // 学年
	IsActive        bool           `gorm:"default:true" json:"is_active"`          // 是否激活
	ExpiresAt       *time.Time     `gorm:"index" json:"expires_at,omitempty"`     // 过期时间
	RevokedAt       *time.Time     `gorm:"index" json:"revoked_at,omitempty"`     // 撤销时间
}

func (AdminPermission) TableName() string { return "admin_permissions" }
