// growth-partner/backend/internal/model/behavior_record.go
// 正向行为记录：七个维度，只记录进步，不记录失败

package model

import (
	"time"
)

// BehaviorDimension 行为维度枚举
type BehaviorDimension string

const (
	DimVirtue BehaviorDimension = "virtue" // 德馨（品德）
	DimStudy  BehaviorDimension = "study"  // 智睿（学习）
	DimSport  BehaviorDimension = "sport"  // 体健（运动）
	DimArt    BehaviorDimension = "art"    // 美雅（艺术）
	DimLabor  BehaviorDimension = "labor"  // 劳朴（劳动）
)

// RecorderRole 记录者角色（谁给的鼓励）
type RecorderRole string

const (
	RecorderTeacher RecorderRole = "teacher" // 班主任/科任老师
	RecorderParent  RecorderRole = "parent"  // 家长
	RecorderSystem  RecorderRole = "system"  // 系统自动（如连续打卡）
)

// DimensionGrowthRange 各维度成长值范围（最小值-最大值）
var DimensionGrowthRange = map[BehaviorDimension][2]int{
	DimVirtue: {2, 5},
	DimStudy:  {1, 3},
	DimSport:  {1, 3},
	DimArt:    {1, 3},
	DimLabor:  {1, 2},
}

// BehaviorRecord 正向行为记录表
type BehaviorRecord struct {
	Base

	ChildID        uint64       `gorm:"index;not null" json:"child_id"`
	ClassID        uint64       `gorm:"index;not null" json:"class_id"`
	SchoolYear     string       `gorm:"size:20;not null" json:"school_year"`    // 所属学年
	RecorderUserID uint64       `gorm:"index;not null" json:"recorder_user_id"` // 记录者用户ID
	RecorderRole   RecorderRole `gorm:"size:20;not null" json:"recorder_role"`

	// ─── 行为详情 ───
	Dimension   BehaviorDimension `gorm:"size:20;index;not null" json:"dimension"`
	Description string            `gorm:"type:text;not null"      json:"description"` // 具体行为描述
	GrowthValue int               `gorm:"not null"               json:"growth_value"` // 本次获得的成长值

	// ─── 伙伴反馈 ───
	PartnerMessage string `gorm:"type:text" json:"partner_message"` // 伙伴说的鼓励话

	// ─── 推送状态 ───
	IsPushed bool       `gorm:"default:false" json:"is_pushed"`
	PushedAt *time.Time `gorm:"index" json:"pushed_at,omitempty"`

	// ─── 审计字段 ───
	IsAudited bool `gorm:"default:true;index" json:"is_audited"`
}

func (BehaviorRecord) TableName() string { return "behavior_records" }
