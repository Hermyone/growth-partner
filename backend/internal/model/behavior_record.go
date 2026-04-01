// growth-partner/backend/internal/model/behavior_record.go
// 正向行为记录：七个维度，只记录进步，不记录失败

package model

// BehaviorDimension 行为维度枚举（七维体系）
type BehaviorDimension string

const (
	DimVirtue     BehaviorDimension = "virtue"     // 德馨（品德）
	DimStudy      BehaviorDimension = "study"      // 智睿（学习）
	DimSport      BehaviorDimension = "sport"      // 体健（运动）
	DimArt        BehaviorDimension = "art"        // 美雅（艺术）
	DimLabor      BehaviorDimension = "labor"      // 劳朴（劳动）
	DimProgress   BehaviorDimension = "progress"   // 进步（努力）
	DimInnovation BehaviorDimension = "innovation" // 创新（创造力）
)

// RecorderRole 记录者角色（谁给的鼓励）
type RecorderRole string

const (
	RecorderTeacher RecorderRole = "teacher" // 班主任/科任老师
	RecorderCoach   RecorderRole = "coach"   // 专业教练
	RecorderParent  RecorderRole = "parent"  // 家长
	RecorderSystem  RecorderRole = "system"  // 系统自动（如连续打卡）
)

// DimensionGrowthRange 各维度成长值范围（最小值-最大值）
var DimensionGrowthRange = map[BehaviorDimension][2]int{
	DimVirtue:     {2, 5},
	DimStudy:      {1, 3},
	DimSport:      {1, 3},
	DimArt:        {1, 3},
	DimLabor:      {1, 2},
	DimProgress:   {2, 5},
	DimInnovation: {2, 5},
}

// BehaviorRecord 正向行为记录表
// 每条记录代表一次被看见的进步，永久保留，不可删除（PIPL 合规）
type BehaviorRecord struct {
	Base

	ChildID      uint64       `gorm:"index;not null" json:"child_id"`
	ClassID      uint64       `gorm:"index;not null" json:"class_id"`
	RecorderID   uint64       `gorm:"index;not null" json:"recorder_id"` // 记录者用户ID
	RecorderRole RecorderRole `gorm:"size:16;not null" json:"recorder_role"`

	// ─── 行为详情 ───
	Dimension   BehaviorDimension `gorm:"size:16;index;not null" json:"dimension"`
	Description string            `gorm:"size:256;not null"      json:"description"`  // 具体行为描述
	GrowthValue int               `gorm:"not null"               json:"growth_value"` // 本次获得的成长值

	// ─── 伙伴反馈 ───
	PartnerMessage string `gorm:"size:256" json:"partner_message"` // 伙伴说的鼓励话
	MessageSentAt  *int64 `gorm:"index"    json:"-"`               // WebSocket 推送时间戳

	// 审计字段（不暴露给学生端）
	IsAudited bool `gorm:"default:false;index" json:"-"`
}

func (BehaviorRecord) TableName() string { return "behavior_records" }
