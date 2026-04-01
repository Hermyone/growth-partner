// growth-partner/backend/internal/model/battle_record.go
// 对战记录：知识竞技场，纯正向，双方均获成长值

package model

import "time"

// BattleSubject 对战科目
type BattleSubject string

const (
	BattlePoetry  BattleSubject = "poetry"  // 诗词
	BattleMath    BattleSubject = "math"    // 速算
	BattleGeneral BattleSubject = "general" // 百科
	BattleEnglish BattleSubject = "english" // 英语
	BattleCustom  BattleSubject = "custom"  // 自定义题库
)

// BattleMode 对战模式
type BattleMode string

const (
	BattleModeNormal     BattleMode = "normal"     // 正式对战（获得参与成长值）
	BattleModeFriendship BattleMode = "friendship" // 友谊战（纯娱乐，不计记录）
)

// BattleRecord 对战记录表
// 注意：不公开排名，不记录胜负到学生可见界面
type BattleRecord struct {
	Base

	ClassID uint64        `gorm:"index;not null" json:"class_id"`
	Subject BattleSubject `gorm:"size:16;not null" json:"subject"`
	Mode    BattleMode    `gorm:"size:16;default:'normal'" json:"mode"`

	// ─── 参与者（两方） ───
	PlayerAChildID uint64 `gorm:"index;not null" json:"player_a_child_id"`
	PlayerBChildID uint64 `gorm:"index;not null" json:"player_b_child_id"`

	// ─── 成绩（仅后台可见，前端不显示胜负） ───
	PlayerAScore int `gorm:"default:0" json:"-"` // 隐藏分数，不暴露给学生
	PlayerBScore int `gorm:"default:0" json:"-"`

	// ─── 成长值奖励（双方均获得） ───
	PlayerAGrowthGained int `gorm:"default:3" json:"player_a_growth_gained"` // 参与奖励
	PlayerBGrowthGained int `gorm:"default:3" json:"player_b_growth_gained"`

	// ─── 荣誉徽章（胜方，纯展示性） ───
	WinnerChildID *uint64 `gorm:"index" json:"-"`                       // 后台记录但不前端展示排名
	HonorBadge    string  `gorm:"size:32" json:"honor_badge,omitempty"` // 徽章展示给获奖者本人

	// ─── 时间 ───
	StartedAt  time.Time  `gorm:"not null" json:"started_at"`
	FinishedAt *time.Time `gorm:"index"    json:"finished_at,omitempty"`
	Duration   int        `gorm:"default:0" json:"duration"` // 秒

	RoomID string `gorm:"size:64;index" json:"room_id"` // WebSocket 房间ID
}

func (BattleRecord) TableName() string { return "battle_records" }
