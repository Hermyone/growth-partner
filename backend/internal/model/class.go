// growth-partner/backend/internal/model/class.go
// 班级模型

package model

// Class 班级表
type Class struct {
	Base

	SchoolID     uint64 `gorm:"index;not null" json:"school_id"`     // 学校ID
	ClassName    string `gorm:"size:100;not null" json:"class_name"`  // 班级名称，如 "三年级四班"
	ClassCode    string `gorm:"size:50;uniqueIndex;not null" json:"class_code"` // 班级代码，唯一
	Grade        int    `gorm:"not null" json:"grade"`           // 年级（1-6）
	ClassNo      int    `gorm:"not null" json:"class_no"`         // 班级号
	SchoolYear   string `gorm:"size:20;not null" json:"school_year"` // 学年，如 "2024-2025"
	HomeroomTeacherID uint64 `gorm:"index" json:"homeroom_teacher_id"` // 班主任用户ID
	StudentCount  int   `gorm:"default:0" json:"student_count"`     // 学生数量

	IsActive     bool   `gorm:"default:true;index" json:"is_active"`
}

func (Class) TableName() string { return "classes" }
