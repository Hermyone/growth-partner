// growth-partner/backend/internal/model/battle.go
// 战斗相关模型

package model

import (
	"time"
)

// BattleRoom 对战房间
type BattleRoom struct {
	ID            uint64        `json:"id" gorm:"primaryKey"`
	RoomCode      string        `json:"room_code" gorm:"size:6;uniqueIndex;not null"`
	Subject       BattleSubject `json:"subject" gorm:"size:16;not null"`
	Status        string        `json:"status" gorm:"size:20;not null"` // waiting, battle, ended
	CreatorID     uint64        `json:"creator_id" gorm:"not null"`
	ParticipantID uint64        `json:"participant_id"`
	CreatedAt     time.Time     `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time     `json:"updated_at" gorm:"autoUpdateTime"`
}

// BattleQuestion 对战题目
type BattleQuestion struct {
	ID         uint64    `json:"id" gorm:"primaryKey"`
	RoomID     uint64    `json:"room_id" gorm:"not null"`
	QuestionID uint64    `json:"question_id" gorm:"not null"`
	Order      int       `json:"order" gorm:"not null"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
}

// BattleReview 对战复盘
type BattleReview struct {
	RoomID           uint64            `json:"room_id"`
	Subject          BattleSubject     `json:"subject"`
	StartTime        time.Time         `json:"start_time"`
	EndTime          time.Time         `json:"end_time"`
	CreatorID        uint64            `json:"creator_id"`
	ParticipantID    uint64            `json:"participant_id"`
	CreatorScore     int               `json:"creator_score"`
	ParticipantScore int               `json:"participant_score"`
	WinnerID         uint64            `json:"winner_id"`
	Questions        []*BattleQuestion `json:"questions"`
}

// 战斗房间状态常量
const (
	BattleRoomStatusWaiting = "waiting" // 等待中
	BattleRoomStatusBattle  = "battle"  // 对战中
	BattleRoomStatusEnded   = "ended"   // 已结束
)
