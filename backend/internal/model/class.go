// growth-partner/backend/internal/model/class.go
// 班级模型

package model

// Class 班级表
type Class struct {
	Base

	Name       string `gorm:"size:64;not null"   json:"name"`        // 如 "三年级四班"
	Grade      int    `gorm:"not null"           json:"grade"`       // 年级（1-6）
	SchoolYear string `gorm:"size:16;not null"   json:"school_year"` // 如 "2024-2025"
	TeacherID  uint64 `gorm:"index;not null"     json:"teacher_id"`  // 班主任用户ID

	// 班级设置
	MaxStudents int    `gorm:"default:50"         json:"max_students"`
	Description string `gorm:"size:256"           json:"description,omitempty"`
	IsActive    bool   `gorm:"default:true;index" json:"is_active"`
}

func (Class) TableName() string { return "classes" }
