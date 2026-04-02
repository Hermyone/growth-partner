// growth-partner/backend/internal/model/sunshine.go
// 阳光章相关模型

package model

import (
	"time"
)

// SunshineColor 阳光章颜色配置
type SunshineColor struct {
	ID          uint64    `json:"id" gorm:"primaryKey"`
	SchoolID    uint64    `json:"school_id" gorm:"index;not null"`
	ColorName   string    `json:"color_name" gorm:"size:20;not null"`
	Subject     string    `json:"subject" gorm:"size:50;not null"`
	Description string    `json:"description" gorm:"size:200"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// SunshineStamp 阳光章盖章记录
type SunshineStamp struct {
	ID        uint64    `json:"id" gorm:"primaryKey"`
	StudentID uint64    `json:"student_id" gorm:"index;not null"`
	TeacherID uint64    `json:"teacher_id" gorm:"index;not null"`
	ColorID   uint64    `json:"color_id" gorm:"index;not null"`
	Subject   string    `json:"subject" gorm:"size:50;not null"`
	StampDate time.Time `json:"stamp_date" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}

// SunshineAward 阳光章之星评选
type SunshineAward struct {
	ID         uint64    `json:"id" gorm:"primaryKey"`
	StudentID  uint64    `json:"student_id" gorm:"index;not null"`
	ClassID    uint64    `json:"class_id" gorm:"index;not null"`
	ColorID    uint64    `json:"color_id" gorm:"index;not null"`
	AwardName  string    `json:"award_name" gorm:"size:50;not null"`
	Period     string    `json:"period" gorm:"size:20;not null"` // 月度/季度/年度
	PeriodYear int       `json:"period_year" gorm:"not null"`
	PeriodMonth int      `json:"period_month"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
}
