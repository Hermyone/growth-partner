// growth-partner/backend/internal/model/partner.go
// 学生专属伙伴实例：每个学生拥有一个伙伴

package model

import "time"

// Partner 学生专属伙伴实例
// 每个学生只有一只伙伴，伙伴随成长值自动进化
type Partner struct {
	Base

	ChildID    uint64 `gorm:"uniqueIndex;not null" json:"child_id"`    // 一个学生只有一只伙伴
	TemplateID uint64 `gorm:"index;not null"       json:"template_id"` // 选择的模板

	// ─── 成长状态 ───
	GrowthPoints   int            `gorm:"default:0;not null"   json:"growth_points"`    // 当前成长值
	CurrentStage   EvolutionStage `gorm:"default:1;not null"   json:"current_stage"`    // 当前进化阶段
	EvolutionCount int            `gorm:"default:0"             json:"evolution_count"` // 进化次数

	// ─── 伙伴自定义 ───
	Nickname string `gorm:"size:32" json:"nickname"` // 小朋友给伙伴起的名字

	// ─── 里程碑追踪 ───
	FirstEvolutionAt *time.Time `gorm:"index" json:"first_evolution_at,omitempty"` // 初次进化时间
	LastEvolvedAt    *time.Time `gorm:"index" json:"last_evolved_at,omitempty"`

	// ─── 交互状态（中/高阶解锁功能） ───
	InteractionLevel int  `gorm:"default:0" json:"interaction_level"` // 0=无,1=预置,2=大模型
	IsActive         bool `gorm:"default:true" json:"is_active"`
}

func (Partner) TableName() string { return "partners" }
