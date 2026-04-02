// growth-partner/backend/internal/model/growth_record.go
// 成长值流水账：记录每一次成长值变动，支持成长轨迹回溯

package model

// GrowthSourceType 成长值来源类型
type GrowthSourceType string

const (
	GrowthSourceBehavior GrowthSourceType = "behavior" // 正向行为
	GrowthSourceBattle   GrowthSourceType = "battle"   // 对战参与
	GrowthSourceSystem   GrowthSourceType = "system"   // 系统奖励（如连续打卡）
	GrowthSourceStamp    GrowthSourceType = "stamp"    // 阳光章奖励
)

// GrowthRecord 成长值流水记录
type GrowthRecord struct {
	Base

	ChildID    uint64 `gorm:"index;not null" json:"child_id"`
	PartnerID  uint64 `gorm:"index;not null" json:"partner_id"`
	SchoolYear string `gorm:"size:20;not null" json:"school_year"` // 所属学年

	// ─── 变动详情 ───
	SourceType  GrowthSourceType `gorm:"size:20;index;not null" json:"source_type"`
	SourceID    uint64           `gorm:"index;not null"        json:"source_id"`     // 来源记录的ID
	Delta       int              `gorm:"not null"               json:"delta"`        // 变化量（正数）
	AfterPoints int              `gorm:"not null"               json:"after_points"` // 变化后的总成长值

	// ─── 进化事件标记 ───
	IsEvolutionTrigger bool            `gorm:"default:false;index"    json:"is_evolution_trigger"` // 是否触发进化
	EvolutionFromStage *EvolutionStage `gorm:"default:null"          json:"evolution_from_stage,omitempty"`
	EvolutionToStage   *EvolutionStage `gorm:"default:null"          json:"evolution_to_stage,omitempty"`

	Remark string `gorm:"size:255" json:"remark,omitempty"`
}

func (GrowthRecord) TableName() string { return "growth_records" }
