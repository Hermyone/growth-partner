// growth-partner/backend/internal/model/audit_log.go
// 管理员操作审计日志模型

package model

// AuditAction 审计操作类型
type AuditAction string

const (
	AuditActionCreate AuditAction = "create"
	AuditActionUpdate AuditAction = "update"
	AuditActionDelete AuditAction = "delete"
	AuditActionLogin  AuditAction = "login"
	AuditActionLogout AuditAction = "logout"
	AuditActionOther  AuditAction = "other"
)

// AuditLog 管理员操作审计日志表
type AuditLog struct {
	Base

	UserID       uint64      `gorm:"index;not null" json:"user_id"`               // 操作用户ID
	Username     string      `gorm:"size:50;index;not null" json:"username"`      // 操作用户名
	Role         UserRole    `gorm:"size:20;not null" json:"role"`                // 操作用户角色
	Action       AuditAction `gorm:"size:20;not null" json:"action"`              // 操作类型
	ResourceType string      `gorm:"size:50;index;not null" json:"resource_type"` // 资源类型
	ResourceID   uint64      `gorm:"index" json:"resource_id"`                    // 资源ID
	Details      string      `gorm:"type:text" json:"details"`                    // 操作详情（JSON格式）
	IPAddress    string      `gorm:"size:50" json:"ip_address"`                   // 操作IP地址
	UserAgent    string      `gorm:"size:255" json:"user_agent"`                  // 用户代理
}

func (AuditLog) TableName() string { return "audit_logs" }
