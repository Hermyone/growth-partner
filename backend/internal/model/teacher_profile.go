// growth-partner/backend/internal/model/teacher_profile.go
// 教师档案模型

package model

// TeacherProfile 教师档案表
type TeacherProfile struct {
	Base

	UserID        uint64 `gorm:"index;not null" json:"user_id"`       // 关联用户ID
	SchoolID      uint64 `gorm:"not null" json:"school_id"`           // 学校ID
	EmployeeNo    string `gorm:"size:50;not null" json:"employee_no"` // 工号
	Subject       string `gorm:"size:50;not null" json:"subject"`     // 教授科目
	SunshineColor string `gorm:"size:20" json:"sunshine_color"`       // 阳光颜色
	Title         string `gorm:"size:50" json:"title"`                // 职称
	IsHomeroom    bool   `gorm:"default:false" json:"is_homeroom"`    // 是否班主任
}

func (TeacherProfile) TableName() string { return "teacher_profiles" }
