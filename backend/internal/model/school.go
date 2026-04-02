// growth-partner/backend/internal/model/school.go
// 学校模型

package model

// School 学校信息表
type School struct {
	Base

	Name          string `gorm:"size:100;not null" json:"name"`          // 学校名称
	District      string `gorm:"size:50;not null" json:"district"`      // 区域
	Address       string `gorm:"size:255;not null" json:"address"`       // 地址
	ContactPhone  string `gorm:"size:50" json:"contact_phone"`          // 联系电话
	IsActive      bool   `gorm:"default:true;index" json:"is_active"`   // 是否激活
}

func (School) TableName() string { return "schools" }
