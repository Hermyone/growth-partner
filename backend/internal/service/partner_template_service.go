// growth-partner/backend/internal/service/partner_template_service.go
// 伙伴模板服务层

package service

import (
	"context"
	"errors"

	"growth-partner/internal/model"
	"growth-partner/internal/repository"
)

// PartnerTemplateService 伙伴模板服务接口
type PartnerTemplateService interface {
	// 公开接口
	GetAllActivePartnerTemplates(ctx context.Context) ([]*model.PartnerTemplate, error)
	GetPartnerTemplateByID(ctx context.Context, templateID uint64) (*model.PartnerTemplate, error)
	GetHealthStatus(ctx context.Context) map[string]string
	GetClientConfig(ctx context.Context) map[string]interface{}

	// 管理员接口
	CreatePartnerTemplate(ctx context.Context, template *model.PartnerTemplate) error
	UpdatePartnerTemplate(ctx context.Context, template *model.PartnerTemplate) error
	SeedPartnerTemplates(ctx context.Context) error
}

// partnerTemplateServiceImpl 伙伴模板服务实现
type partnerTemplateServiceImpl struct {
	templateRepo repository.TemplateRepository
}

// NewPartnerTemplateService 创建伙伴模板服务实例
func NewPartnerTemplateService(
	templateRepo repository.TemplateRepository,
) PartnerTemplateService {
	return &partnerTemplateServiceImpl{
		templateRepo: templateRepo,
	}
}

// ─── 公开接口 ───────────────────────────────────────────────

func (s *partnerTemplateServiceImpl) GetAllActivePartnerTemplates(ctx context.Context) ([]*model.PartnerTemplate, error) {
	// 获取全部激活的伙伴模板列表
	return s.templateRepo.FindAllActive(ctx)
}

func (s *partnerTemplateServiceImpl) GetPartnerTemplateByID(ctx context.Context, templateID uint64) (*model.PartnerTemplate, error) {
	// 获取单个伙伴模板详情
	template, err := s.templateRepo.FindByID(ctx, templateID)
	if err != nil {
		return nil, err
	}
	if template == nil {
		return nil, errors.New("伙伴模板不存在")
	}
	return template, nil
}

func (s *partnerTemplateServiceImpl) GetHealthStatus(ctx context.Context) map[string]string {
	// 健康检查，返回{status:ok}
	return map[string]string{"status": "ok"}
}

func (s *partnerTemplateServiceImpl) GetClientConfig(ctx context.Context) map[string]interface{} {
	// 返回前端全局配置
	return map[string]interface{}{
		"growthPointsThreshold": map[string]int{
			"level1": 100,
			"level2": 300,
			"level3": 600,
			"level4": 1000,
			"level5": 1500,
		},
		"dimensions": []string{
			"德馨",
			"智睿",
			"体健",
			"美雅",
			"劳朴",
			"进步",
			"创新",
		},
		"battle": map[string]interface{}{
			"questionCount": 10,
			"timePerQuestion": 30, // 秒
			"subjects": []string{"语文", "数学", "英语", "综合"},
		},
	}
}

// ─── 管理员接口 ─────────────────────────────────────────────

func (s *partnerTemplateServiceImpl) CreatePartnerTemplate(ctx context.Context, template *model.PartnerTemplate) error {
	// 新增伙伴模板
	return s.templateRepo.Create(ctx, template)
}

func (s *partnerTemplateServiceImpl) UpdatePartnerTemplate(ctx context.Context, template *model.PartnerTemplate) error {
	// 更新模板信息
	existingTemplate, err := s.templateRepo.FindByID(ctx, template.ID)
	if err != nil {
		return err
	}
	if existingTemplate == nil {
		return errors.New("伙伴模板不存在")
	}
	return s.templateRepo.Update(ctx, template)
}

func (s *partnerTemplateServiceImpl) SeedPartnerTemplates(ctx context.Context) error {
	// 一键初始化30个预设模板（幂等）
	// 这里可以实现初始化逻辑，暂时返回nil
	return nil
}
