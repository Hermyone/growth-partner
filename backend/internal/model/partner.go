// growth-partner/backend/internal/model/partner.go
// 学生专属伙伴实例：每个学生拥有一个伙伴

package model

import "time"

// PartnerStatus 伙伴状态
type PartnerStatus string

const (
	PartnerStatusActive   PartnerStatus = "active"   // 活跃状态
	PartnerStatusPending  PartnerStatus = "pending"  // 待激活
	PartnerStatusGraduated PartnerStatus = "graduated" // 已毕业
)

// Partner 学生专属伙伴实例
type Partner struct {
	Base

	ChildID    uint64        `gorm:"index;not null"       json:"child_id"`    // 学生ID
	TemplateID uint64        `gorm:"index;not null"       json:"template_id"` // 选择的模板
	SequenceNo int           `gorm:"not null"             json:"sequence_no"` // 伙伴序列号（学生可更换伙伴）

	// ─── 伙伴自定义 ───
	Nickname string `gorm:"size:50" json:"nickname"` // 小朋友给伙伴起的名字

	// ─── 成长状态 ───
	Status        PartnerStatus  `gorm:"size:20;default:'active'" json:"status"`
	GrowthPoints  int            `gorm:"default:0;not null"       json:"growth_points"`    // 当前成长值
	CurrentStage  EvolutionStage `gorm:"default:1;not null"       json:"current_stage"`    // 当前进化阶段
	EvolutionCount int           `gorm:"default:0"                json:"evolution_count"` // 进化次数

	// ─── 交互状态（中/高阶解锁功能） ───
	InteractionLevel int  `gorm:"default:0" json:"interaction_level"` // 0=无,1=预置,2=大模型

	// ─── 学年信息 ───
	SchoolYear string `gorm:"size:20;not null" json:"school_year"` // 所属学年

	// ─── 里程碑追踪 ───
	FirstEvolutionAt *time.Time `gorm:"index" json:"first_evolution_at,omitempty"` // 初次进化时间
	LastEvolvedAt    *time.Time `gorm:"index" json:"last_evolved_at,omitempty"`
	GraduatedAt      *time.Time `gorm:"index" json:"graduated_at,omitempty"`      // 毕业时间
	SelectedAt       *time.Time `gorm:"index" json:"selected_at,omitempty"`       // 选中时间
}

func (Partner) TableName() string { return "partners" }
