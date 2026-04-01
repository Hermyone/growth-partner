// growth-partner/backend/internal/model/blind_box.go
// 班级盲盒系统

package model

import "time"

// BlindBoxRewardType 盲盒奖励类型
type BlindBoxRewardType string

const (
	RewardPrivilege  BlindBoxRewardType = "privilege"  // 特权类
	RewardExperience BlindBoxRewardType = "experience" // 体验类
	RewardHonor      BlindBoxRewardType = "honor"      // 荣誉类
	RewardService    BlindBoxRewardType = "service"    // 服务类
	RewardCollective BlindBoxRewardType = "collective" // 集体类
)

// BlindBoxPool 班级盲盒奖励池（老师配置）
type BlindBoxPool struct {
	Base

	ClassID     uint64             `gorm:"index;not null"    json:"class_id"`
	Type        BlindBoxRewardType `gorm:"size:16;not null"  json:"type"`
	Title       string             `gorm:"size:64;not null"  json:"title"` // 如"免作业一次"
	Description string             `gorm:"size:256"          json:"description"`
	Stock       int                `gorm:"default:1"         json:"stock"` // 剩余库存（-1无限）
	IsActive    bool               `gorm:"default:true"      json:"is_active"`
}

// BlindBoxDraw 抽盲盒记录
type BlindBoxDraw struct {
	Base

	ChildID uint64 `gorm:"index;not null" json:"child_id"`
	ClassID uint64 `gorm:"index;not null" json:"class_id"`
	PoolID  uint64 `gorm:"index;not null" json:"pool_id"` // 抽到的奖励

	// 兑换状态
	DrawnAt    time.Time  `gorm:"not null;index" json:"drawn_at"`
	RedeemedAt *time.Time `gorm:"index"          json:"redeemed_at,omitempty"` // 兑换时间
	ExpiresAt  *time.Time `gorm:"index"          json:"expires_at,omitempty"`  // 过期时间

	IsRedeemed bool `gorm:"default:false;index" json:"is_redeemed"`
}

func (BlindBoxPool) TableName() string { return "blind_box_pools" }
func (BlindBoxDraw) TableName() string { return "blind_box_draws" }
