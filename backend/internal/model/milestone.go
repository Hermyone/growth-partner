// growth-partner/backend/internal/model/milestone.go
// 里程碑记录：成长画卷的亮点贴纸

package model

// MilestoneType 里程碑类型
type MilestoneType string

const (
	MilestoneFirstEvolution MilestoneType = "first_evolution" // 初次进化
	MilestoneConsecutive7   MilestoneType = "consecutive_7"   // 连续7天打卡
	MilestoneVirtue10       MilestoneType = "virtue_10"       // 累计10次德馨记录
	MilestoneBattle20       MilestoneType = "battle_20"       // 参与20次对战
	MilestoneCollective     MilestoneType = "collective_hero" // 集体英雄
	MilestoneAnnual500      MilestoneType = "annual_500"      // 全年500成长值
)

// MilestoneConfig 里程碑配置
var MilestoneConfig = map[MilestoneType]struct {
	Title   string
	Content string
}{
	MilestoneFirstEvolution: {"初次进化", "第一次蜕变，你做到了！"},
	MilestoneConsecutive7:   {"连续打卡", "七天不间断，坚持的力量真强大！"},
	MilestoneVirtue10:       {"助人达人", "你的善良，温暖了好多人。"},
	MilestoneBattle20:       {"对战勇士", "挑战从不停歇，勇士就是你！"},
	MilestoneCollective:     {"集体英雄", "因为有你，整个班级都更幸福了。"},
	MilestoneAnnual500:      {"全年成长王", "这一年，你真的很了不起。"},
}

// Milestone 里程碑记录表
type Milestone struct {
	Base

	ChildID uint64        `gorm:"index;not null" json:"child_id"`
	Type    MilestoneType `gorm:"size:32;not null" json:"type"`
	Title   string        `gorm:"size:64;not null"  json:"title"`
	Content string        `gorm:"size:256;not null" json:"content"`

	// 触发来源
	SourceType string `gorm:"size:32" json:"source_type"` // "behavior"/"battle"/"system"
	SourceID   uint64 `gorm:"index"   json:"source_id"`

	// 已通知标记
	IsNotified bool `gorm:"default:false;index" json:"is_notified"`
}

func (Milestone) TableName() string { return "milestones" }
