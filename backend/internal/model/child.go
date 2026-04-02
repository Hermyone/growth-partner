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
	DisplayName   string `gorm:"size:50;not null" json:"display_name"`   // 对外展示名称（不含真实姓名）
	RealNameEnc   string `gorm:"size:255;not null" json:"-"`            // AES 加密的真实姓名
	StudentNoEnc  string `gorm:"size:255;not null" json:"-"`            // AES 加密的学号
	Gender        string `gorm:"type:char(1);not null;check:gender IN ('M', 'F')" json:"gender"` // 性别
	BirthYear     int    `gorm:"not null" json:"birth_year"`            // 出生年份

	// ─── 学籍信息 ───
	EnrollYear    int    `gorm:"not null" json:"enroll_year"`           // 入学年份
	CurrentGrade  int    `gorm:"not null" json:"current_grade"`         // 当前年级

	// ─── 成长统计（冗余字段，避免频繁聚合查询） ───
	TotalGrowthPoints   int `gorm:"default:0" json:"total_growth_points"`   // 历史总成长值
	CurrentGrowthPoints int `gorm:"default:0" json:"current_growth_points"` // 当前可用成长值
	CurrentGloryCoins   int `gorm:"default:0" json:"current_glory_coins"`   // 当前荣耀币
	TotalBehaviorCount  int `gorm:"default:0" json:"total_behavior_count"`  // 总行为记录次数
	BattleCount         int `gorm:"default:0" json:"battle_count"`          // 参与对战次数
	ConsecutiveDays     int `gorm:"default:0" json:"consecutive_days"`      // 当前连续天数
	MaxConsecutiveDays  int `gorm:"default:0" json:"max_consecutive_days"`  // 最大连续天数

	IsActive bool `gorm:"default:true;index" json:"is_active"`
}

func (Child) TableName() string { return "children" }
