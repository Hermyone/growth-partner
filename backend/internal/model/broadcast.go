// growth-partner/backend/internal/model/broadcast.go
// 广播模型

package model

import "time"

// BroadcastStatus 广播状态
type BroadcastStatus string

const (
	BroadcastStatusPending    BroadcastStatus = "pending"    // 待发送
	BroadcastStatusSent       BroadcastStatus = "sent"       // 已发送
	BroadcastStatusCancelled  BroadcastStatus = "cancelled"  // 已取消
	BroadcastStatusFailed     BroadcastStatus = "failed"     // 发送失败
)

// BroadcastTargetType 广播目标类型
type BroadcastTargetType string

const (
	BroadcastTargetTypeClass   BroadcastTargetType = "class"   // 班级
	BroadcastTargetTypeSchool  BroadcastTargetType = "school"  // 学校
	BroadcastTargetTypeAll     BroadcastTargetType = "all"     // 全体
)

// Broadcast 广播表
type Broadcast struct {
	Base

	Title          string               `gorm:"size:100;not null" json:"title"`          // 广播标题
	Content        string               `gorm:"type:text;not null" json:"content"`        // 广播内容
	TargetType     BroadcastTargetType  `gorm:"size:20;not null" json:"target_type"`     // 目标类型
	TargetID       *uint64              `json:"target_id,omitempty"`                        // 目标ID（班级ID或学校ID）
	BroadcastType  string               `gorm:"size:20;not null" json:"broadcast_type"`   // 广播类型（如：text, image, video）
	Status         BroadcastStatus      `gorm:"size:20;default:pending" json:"status"`     // 广播状态
	SendAt         *time.Time           `json:"send_at,omitempty"`                          // 发送时间
	CreatedBy      uint64               `gorm:"index;not null" json:"created_by"`          // 创建人ID（老师或管理员）
	CreatedByRole  string               `gorm:"size:20;not null" json:"created_by_role"`   // 创建人角色
	RecipientCount int                  `gorm:"default:0" json:"recipient_count"`          // 接收人数
	ErrorMsg       string               `gorm:"size:512" json:"error_msg,omitempty"`       // 发送失败时的错误信息
}

func (Broadcast) TableName() string { return "broadcasts" }
