// growth-partner/backend/internal/model/child.go
// 学生档案模型
// 注意：学生的个人信息全部加密，严格保护未成年人隐私

package model

// Child 学生档案表（对应一个学生用户账号）
type Child struct {
	Base

	// 关联用户账号（1:1）
	UserID  uint64 `gorm:"uniqueIndex;not null" json:"user_id"`
	ClassID uint64 `gorm:"index;not null"       json:"class_id"`

	// ─── 个人信息（加密存储） ───
	NicknameEnc string `gorm:"size:256" json:"-"` // 昵称（加密）
	// 对外展示使用 DisplayName（可以是昵称脱敏版本）
	DisplayName string `gorm:"size:32;not null" json:"display_name"` // 如 "小明同学"（不含真实姓名）

	// ─── 学籍信息 ───
	StudentNumber string `gorm:"uniqueIndex;size:32" json:"-"` // 学号（加密，不暴露）
	Grade         int    `gorm:"not null"             json:"grade"`
	EnrollYear    int    `gorm:"not null"             json:"enroll_year"` // 入学年份

	// ─── 成长统计（冗余字段，避免频繁聚合查询） ───
	TotalGrowthPoints   int `gorm:"default:0" json:"total_growth_points"`   // 历史总成长值
	CurrentGrowthPoints int `gorm:"default:0" json:"current_growth_points"` // 当前可用成长值
	TotalBehaviorCount  int `gorm:"default:0" json:"total_behavior_count"`  // 总行为记录次数
	BattleCount         int `gorm:"default:0" json:"battle_count"`          // 参与对战次数
	ConsecutiveDays     int `gorm:"default:0" json:"consecutive_days"`      // 连续打卡天数

	IsActive bool `gorm:"default:true;index" json:"is_active"`
}

func (Child) TableName() string { return "children" }
