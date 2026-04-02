// growth-partner/backend/internal/model/challenge.go
// 集体挑战模型

package model

import "time"

// ChallengeStatus 挑战状态
type ChallengeStatus string

const (
	ChallengeStatusActive   ChallengeStatus = "active"   // 进行中
	ChallengeStatusCompleted ChallengeStatus = "completed" // 已完成
	ChallengeStatusCancelled ChallengeStatus = "cancelled" // 已取消
)

// Challenge 班级集体挑战表
type Challenge struct {
	Base

	ClassID        uint64           `gorm:"index;not null" json:"class_id"`          // 班级ID
	Title          string           `gorm:"size:100;not null" json:"title"`          // 挑战标题
	Description    string           `gorm:"type:text;not null" json:"description"`    // 挑战描述
	StartAt        time.Time        `gorm:"not null" json:"start_at"`                // 开始时间
	EndAt          *time.Time       `json:"end_at,omitempty"`                        // 结束时间
	Status         ChallengeStatus  `gorm:"size:20;default:active" json:"status"`    // 挑战状态
	TargetType     string           `gorm:"size:20;not null" json:"target_type"`     // 目标类型（如：behavior_count, total_growth_points）
	TargetValue    int              `gorm:"not null" json:"target_value"`           // 目标值
	CurrentValue   int              `gorm:"default:0" json:"current_value"`         // 当前进度
	RewardType     string           `gorm:"size:20;not null" json:"reward_type"`     // 奖励类型（如：growth_points, glory_coins）
	RewardValue    int              `gorm:"not null" json:"reward_value"`           // 奖励值
	CreatedBy      uint64           `gorm:"index;not null" json:"created_by"`        // 创建人ID（老师）
	CompletedAt    *time.Time       `json:"completed_at,omitempty"`                  // 完成时间
	CompletedBy    *uint64          `json:"completed_by,omitempty"`                  // 完成操作人ID
}

func (Challenge) TableName() string { return "challenges" }
