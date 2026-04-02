// growth-partner/backend/internal/model/parent_child_relation.go
// 家长-学生绑定关系模型

package model

import "time"

// ParentChildRelation 家长-学生绑定关系表
type ParentChildRelation struct {
	Base

	ParentUserID uint64     `gorm:"index;not null" json:"parent_user_id"` // 家长用户ID
	ChildID      uint64     `gorm:"index;not null" json:"child_id"`       // 学生档案ID
	Relationship string     `gorm:"size:20" json:"relationship"`          // 关系（父亲/母亲/其他）
	IsPrimary    bool       `gorm:"default:false" json:"is_primary"`      // 是否为主要监护人
	ValidFrom    time.Time  `gorm:"not null" json:"valid_from"`           // 绑定生效时间
	ValidUntil   *time.Time `json:"valid_until,omitempty"`                // 绑定失效时间
	IsActive     bool       `gorm:"default:true" json:"is_active"`        // 是否激活
}

func (ParentChildRelation) TableName() string { return "parent_child_relations" }
