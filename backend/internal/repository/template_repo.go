// growth-partner/backend/internal/repository/template_repo.go
// 伙伴模板仓储：管理预设的 30 种伙伴形象及进化配置

package repository

import (
	"context"
	"gorm.io/gorm"
	"growth-partner/internal/model"
)

type TemplateRepository interface {
	FindAllActive(ctx context.Context) ([]*model.PartnerTemplate, error)
	FindByID(ctx context.Context, id uint64) (*model.PartnerTemplate, error)
	Create(ctx context.Context, template *model.PartnerTemplate) error
}

type templateRepositoryImpl struct {
	db *gorm.DB
}

func NewTemplateRepository(db *gorm.DB) TemplateRepository {
	return &templateRepositoryImpl{db: db}
}

// FindAllActive 获取所有可供学生领养的活跃模板（按排序权重排列）
func (r *templateRepositoryImpl) FindAllActive(ctx context.Context) ([]*model.PartnerTemplate, error) {
	var list []*model.PartnerTemplate
	err := r.db.WithContext(ctx).
		Where("is_active = ?", true).
		Order("sort_order ASC").
		Find(&list).Error
	return list, err
}

// FindByID 获取特定模板详情
func (r *templateRepositoryImpl) FindByID(ctx context.Context, id uint64) (*model.PartnerTemplate, error) {
	var template model.PartnerTemplate
	err := r.db.WithContext(ctx).First(&template, id).Error
	if err != nil {
		return nil, err
	}
	return &template, nil
}

// Create 初始化模板数据（用于系统启动时的 Seed 数据）
func (r *templateRepositoryImpl) Create(ctx context.Context, template *model.PartnerTemplate) error {
	return r.db.WithContext(ctx).Create(template).Error
}
