// growth-partner/backend/internal/repository/sunshine_repo.go
// 阳光章仓库

package repository

import (
	"context"

	"growth-partner/internal/model"

	"gorm.io/gorm"
)

// SunshineRepository 阳光章仓库接口
type SunshineRepository interface {
	// 七色配置管理
	FindColorsBySchoolID(ctx context.Context, schoolID uint64) ([]*model.SunshineColor, error)
	CreateColor(ctx context.Context, color *model.SunshineColor) error
	UpdateColor(ctx context.Context, color *model.SunshineColor) error

	// 盖章操作
	CreateStamp(ctx context.Context, stamp *model.SunshineStamp) error
	FindStampsByClassID(ctx context.Context, classID uint64, params map[string]interface{}, limit, offset int) ([]*model.SunshineStamp, int64, error)
	FindStampsByStudentID(ctx context.Context, studentID uint64) ([]*model.SunshineStamp, error)

	// 之星评选
	CreateAward(ctx context.Context, award *model.SunshineAward) error
	FindAwardsByClassID(ctx context.Context, classID uint64) ([]*model.SunshineAward, error)
	FindAwardsByStudentID(ctx context.Context, studentID uint64) ([]*model.SunshineAward, error)
}

// sunshineRepositoryImpl 阳光章仓库实现
type sunshineRepositoryImpl struct {
	db *gorm.DB
}

// NewSunshineRepository 创建阳光章仓库实例
func NewSunshineRepository(db *gorm.DB) SunshineRepository {
	return &sunshineRepositoryImpl{
		db: db,
	}
}

// ─── 七色配置管理 ──────────────────────────────────────────

func (r *sunshineRepositoryImpl) FindColorsBySchoolID(ctx context.Context, schoolID uint64) ([]*model.SunshineColor, error) {
	var colors []*model.SunshineColor
	err := r.db.WithContext(ctx).Where("school_id = ?", schoolID).Find(&colors).Error
	return colors, err
}

func (r *sunshineRepositoryImpl) CreateColor(ctx context.Context, color *model.SunshineColor) error {
	return r.db.WithContext(ctx).Create(color).Error
}

func (r *sunshineRepositoryImpl) UpdateColor(ctx context.Context, color *model.SunshineColor) error {
	return r.db.WithContext(ctx).Save(color).Error
}

// ─── 盖章操作 ──────────────────────────────────────────────

func (r *sunshineRepositoryImpl) CreateStamp(ctx context.Context, stamp *model.SunshineStamp) error {
	return r.db.WithContext(ctx).Create(stamp).Error
}

func (r *sunshineRepositoryImpl) FindStampsByClassID(ctx context.Context, classID uint64, params map[string]interface{}, limit, offset int) ([]*model.SunshineStamp, int64, error) {
	var stamps []*model.SunshineStamp
	var total int64

	db := r.db.WithContext(ctx).Model(&model.SunshineStamp{})

	// 构建查询条件
	if month, ok := params["month"].(string); ok && month != "" {
		db = db.Where("DATE_FORMAT(stamp_date, '%Y-%m') = ?", month)
	}
	if colorID, ok := params["color_id"].(uint64); ok && colorID > 0 {
		db = db.Where("color_id = ?", colorID)
	}

	// 统计总数
	db.Count(&total)

	// 查询数据
	err := db.Order("stamp_date DESC").Limit(limit).Offset(offset).Find(&stamps).Error

	return stamps, total, err
}

func (r *sunshineRepositoryImpl) FindStampsByStudentID(ctx context.Context, studentID uint64) ([]*model.SunshineStamp, error) {
	var stamps []*model.SunshineStamp
	err := r.db.WithContext(ctx).Where("student_id = ?", studentID).Order("stamp_date DESC").Find(&stamps).Error
	return stamps, err
}

// ─── 之星评选 ──────────────────────────────────────────────

func (r *sunshineRepositoryImpl) CreateAward(ctx context.Context, award *model.SunshineAward) error {
	return r.db.WithContext(ctx).Create(award).Error
}

func (r *sunshineRepositoryImpl) FindAwardsByClassID(ctx context.Context, classID uint64) ([]*model.SunshineAward, error) {
	var awards []*model.SunshineAward
	err := r.db.WithContext(ctx).Where("class_id = ?", classID).Order("created_at DESC").Find(&awards).Error
	return awards, err
}

func (r *sunshineRepositoryImpl) FindAwardsByStudentID(ctx context.Context, studentID uint64) ([]*model.SunshineAward, error) {
	var awards []*model.SunshineAward
	err := r.db.WithContext(ctx).Where("student_id = ?", studentID).Order("created_at DESC").Find(&awards).Error
	return awards, err
}
