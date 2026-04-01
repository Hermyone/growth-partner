// growth-partner/backend/internal/model/base.go
// 基础模型：所有数据表的公共字段，包含软删除支持

package model

import (
	"time"

	"gorm.io/gorm"
)

// Base 所有模型的基础字段
// 使用 GORM 约定，自动管理 created_at / updated_at / deleted_at
type Base struct {
	ID        uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt time.Time      `gorm:"index;not null"           json:"created_at"`
	UpdatedAt time.Time      `gorm:"not null"                  json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"                     json:"-"` // 软删除，不暴露给前端
}
