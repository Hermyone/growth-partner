// growth-partner/backend/internal/model/broadcast_record.go
// 广播记录：园长广播持久化；同伴互评消息纯内存广播，不落库（保护隐私）

package model

import "time"

// BroadcastType 广播类型
type BroadcastType string

const (
	BroadcastPartnerToChild BroadcastType = "partner_to_child" // 伙伴对主人的鼓励消息
	BroadcastTeacherToClass BroadcastType = "teacher_to_class" // 园长对全班广播
	BroadcastSystemMorning  BroadcastType = "system_morning"   // 早安广播
	BroadcastSystemEvening  BroadcastType = "system_evening"   // 晚安广播
	// 注意：同伴互评（peer_gift）消息不落库，Redis 内存广播后自然消失
)

// BroadcastRecord 广播记录表（仅持久化需要保留的广播）
type BroadcastRecord struct {
	Base

	Type     BroadcastType `gorm:"size:32;index;not null"  json:"type"`
	SenderID uint64        `gorm:"index;not null"          json:"sender_id"` // 发送者用户ID（系统广播时为0）
	ClassID  uint64        `gorm:"index"                   json:"class_id"`  // 目标班级（0表示全校）
	ChildID  uint64        `gorm:"index"                   json:"child_id"`  // 目标学生（0表示全班）

	Content    string `gorm:"size:512;not null"     json:"content"`     // 广播内容
	TemplateID string `gorm:"size:64"               json:"template_id"` // 使用的模板ID

	// ─── 定时发送 ───
	ScheduledAt *time.Time `gorm:"index" json:"scheduled_at,omitempty"` // 定时发送时间
	SentAt      *time.Time `gorm:"index" json:"sent_at,omitempty"`      // 实际发送时间

	IsRead bool `gorm:"default:false;index" json:"is_read"`
}

func (BroadcastRecord) TableName() string { return "broadcast_records" }
