// growth-partner/backend/internal/service/student_service.go
// 学生端服务层

package service

import (
	"context"
	"errors"
	"time"

	"growth-partner/config"
	"growth-partner/internal/model"
	"growth-partner/internal/repository"
)

// StudentService 学生端服务接口
type StudentService interface {
	// 伙伴系统
	GetCurrentPartner(ctx context.Context, studentID uint64) (*model.Partner, error)
	GetAllPartners(ctx context.Context, studentID uint64) ([]*model.Partner, error)
	SelectNewPartner(ctx context.Context, studentID, templateID uint64) error
	UpdatePartnerNickname(ctx context.Context, studentID uint64, nickname string) error
	GetPartnerGrowthHistory(ctx context.Context, studentID uint64, params map[string]interface{}) ([]*model.GrowthRecord, int64, error)
	GetPartnerTemplates(ctx context.Context) ([]*model.PartnerTemplate, error)

	// 行为记录
	GetMyBehaviors(ctx context.Context, studentID uint64, params map[string]interface{}) ([]*model.BehaviorRecord, int64, error)
	GetBehaviorStats(ctx context.Context, studentID uint64) (map[string]interface{}, error)

	// 广播中心
	GetMyBroadcasts(ctx context.Context, studentID uint64) ([]*model.Broadcast, error)
	MarkBroadcastAsRead(ctx context.Context, broadcastID, studentID uint64) error
	MarkAllBroadcastsAsRead(ctx context.Context, studentID uint64) error

	// 成长年历
	GetGrowthCalendarMonths(ctx context.Context, studentID uint64) (map[string]interface{}, error)
	GetGrowthCalendarMonth(ctx context.Context, studentID uint64, month string) (map[string]interface{}, error)
	GetGrowthCalendarAnnual(ctx context.Context, studentID uint64, year string) (map[string]interface{}, error)
	GetMilestones(ctx context.Context, studentID uint64) ([]*model.Milestone, error)

	// 盲盒查看
	GetMyBlindboxDraws(ctx context.Context, studentID uint64) ([]*model.BlindBoxDraw, error)
}

// studentServiceImpl 学生端服务实现
type studentServiceImpl struct {
	partnerRepo   repository.PartnerRepository
	growthRepo    repository.GrowthRepository
	templateRepo  repository.TemplateRepository
	behaviorRepo  repository.BehaviorRepository
	broadcastRepo repository.BroadcastRepository
	milestoneRepo repository.MilestoneRepository
	blindboxRepo  repository.BlindboxRepository
	childRepo     repository.ChildRepository
}

// NewStudentService 创建学生端服务实例
func NewStudentService(
	partnerRepo repository.PartnerRepository,
	growthRepo repository.GrowthRepository,
	templateRepo repository.TemplateRepository,
	behaviorRepo repository.BehaviorRepository,
	broadcastRepo repository.BroadcastRepository,
	milestoneRepo repository.MilestoneRepository,
	blindboxRepo repository.BlindboxRepository,
	childRepo repository.ChildRepository,
) StudentService {
	return &studentServiceImpl{
		partnerRepo:   partnerRepo,
		growthRepo:    growthRepo,
		templateRepo:  templateRepo,
		behaviorRepo:  behaviorRepo,
		broadcastRepo: broadcastRepo,
		milestoneRepo: milestoneRepo,
		blindboxRepo:  blindboxRepo,
		childRepo:     childRepo,
	}
}

// ─── 4.1 伙伴系统 ──────────────────────────────────────────

func (s *studentServiceImpl) GetCurrentPartner(ctx context.Context, studentID uint64) (*model.Partner, error) {
	// 获取学生当前的伙伴
	partner, err := s.partnerRepo.FindByChildID(ctx, studentID)
	if err != nil {
		return nil, err
	}
	if partner == nil {
		return nil, errors.New("当前没有伙伴")
	}
	return partner, nil
}

func (s *studentServiceImpl) GetAllPartners(ctx context.Context, studentID uint64) ([]*model.Partner, error) {
	// 获取学生的伙伴
	partner, err := s.partnerRepo.FindByChildID(ctx, studentID)
	if err != nil {
		return nil, err
	}
	if partner == nil {
		return []*model.Partner{}, nil
	}
	return []*model.Partner{partner}, nil
}

func (s *studentServiceImpl) SelectNewPartner(ctx context.Context, studentID, templateID uint64) error {
	// 检查学生是否存在
	child, err := s.childRepo.FindByID(ctx, studentID)
	if err != nil {
		return err
	}
	if child == nil {
		return errors.New("学生不存在")
	}

	// 检查模板是否存在
	template, err := s.templateRepo.FindByID(ctx, templateID)
	if err != nil {
		return err
	}
	if template == nil {
		return errors.New("伙伴模板不存在")
	}

	// 检查学生是否有权限选择该模板
	// 这里可以根据实际需求添加权限检查逻辑

	// 创建新伙伴
	now := time.Now()
	newPartner := &model.Partner{
		ChildID:      studentID,
		TemplateID:   templateID,
		SequenceNo:   config.Get().Partner.InitialSequenceNo, // 初始序列号
		Nickname:     template.Name,
		Status:       model.PartnerStatusActive,
		GrowthPoints: config.Get().Partner.InitialGrowthPoints,
		CurrentStage: model.EvolutionStage(config.Get().Partner.InitialStage),           // 初始阶段
		SchoolYear:   config.Get().Partner.DefaultSchoolYear, // 从配置中获取
		SelectedAt:   &now,
	}

	return s.partnerRepo.Create(ctx, newPartner)
}

func (s *studentServiceImpl) UpdatePartnerNickname(ctx context.Context, studentID uint64, nickname string) error {
	// 更新伙伴昵称
	return s.partnerRepo.UpdateNickname(ctx, studentID, nickname)
}

func (s *studentServiceImpl) GetPartnerGrowthHistory(ctx context.Context, studentID uint64, params map[string]interface{}) ([]*model.GrowthRecord, int64, error) {
	// 解析分页参数
	limit := config.Get().Behavior.DefaultLimit
	offset := config.Get().Behavior.DefaultOffset
	if l, ok := params["limit"].(int); ok && l > 0 {
		limit = l
	}
	if page, ok := params["page"].(int); ok && page > 0 {
		offset = (page - 1) * limit
	}
	// 获取学生的成长值流水
	return s.growthRepo.FindByChildID(ctx, studentID, limit, offset)
}

func (s *studentServiceImpl) GetPartnerTemplates(ctx context.Context) ([]*model.PartnerTemplate, error) {
	// 获取所有可用的伙伴模板
	return s.templateRepo.FindAllActive(ctx)
}

// ─── 4.2 行为记录查看 ───────────────────────────────────────

func (s *studentServiceImpl) GetMyBehaviors(ctx context.Context, studentID uint64, params map[string]interface{}) ([]*model.BehaviorRecord, int64, error) {
	// 解析分页参数
	limit := config.Get().Behavior.DefaultLimit
	offset := config.Get().Behavior.DefaultOffset
	if l, ok := params["limit"].(int); ok && l > 0 {
		limit = l
	}
	if page, ok := params["page"].(int); ok && page > 0 {
		offset = (page - 1) * limit
	}
	// 获取学生的行为记录
	return s.behaviorRepo.FindByChildID(ctx, studentID, limit, offset)
}

func (s *studentServiceImpl) GetBehaviorStats(ctx context.Context, studentID uint64) (map[string]interface{}, error) {
	// 获取学生的行为统计数据
	// 这里可以根据实际需求实现统计逻辑
	// 例如：按维度统计行为次数、成长值等
	return map[string]interface{}{
		"dimensions":        config.Get().Behavior.Dimensions,
		"totalBehaviors":    config.Get().Student.DefaultTotalBehaviors,
		"totalGrowthPoints": config.Get().Student.DefaultTotalGrowthPoints,
	}, nil
}

// ─── 4.3 广播中心 ──────────────────────────────────────────

func (s *studentServiceImpl) GetMyBroadcasts(ctx context.Context, studentID uint64) ([]*model.Broadcast, error) {
	// 获取学生收到的广播消息
	// 这里需要根据实际的广播记录模型来实现
	// 暂时返回空数组
	return []*model.Broadcast{}, nil
}

func (s *studentServiceImpl) MarkBroadcastAsRead(ctx context.Context, broadcastID, studentID uint64) error {
	// 标记广播为已读
	// 这里需要根据实际的广播记录模型来实现
	// 暂时返回nil
	return nil
}

func (s *studentServiceImpl) MarkAllBroadcastsAsRead(ctx context.Context, studentID uint64) error {
	// 一键标记所有广播为已读
	// 这里需要根据实际的广播记录模型来实现
	// 暂时返回nil
	return nil
}

// ─── 4.4 成长年历 ──────────────────────────────────────────

func (s *studentServiceImpl) GetGrowthCalendarMonths(ctx context.Context, studentID uint64) (map[string]interface{}, error) {
	// 获取学生的月度成长卡列表（按学年分组）
	// 这里需要根据实际的月度成长卡模型来实现
	// 暂时返回空map
	return map[string]interface{}{"years": map[string][]interface{}{}}, nil
}

func (s *studentServiceImpl) GetGrowthCalendarMonth(ctx context.Context, studentID uint64, month string) (map[string]interface{}, error) {
	// 获取指定月份的成长卡详情
	// 这里需要根据实际的月度成长卡模型来实现
	// 暂时返回空map
	return map[string]interface{}{"month": month, "data": nil}, nil
}

func (s *studentServiceImpl) GetGrowthCalendarAnnual(ctx context.Context, studentID uint64, year string) (map[string]interface{}, error) {
	// 获取年度成长画卷数据
	// 这里需要根据实际的年度报告模型来实现
	// 暂时返回空map
	return map[string]interface{}{"year": year, "data": nil}, nil
}

func (s *studentServiceImpl) GetMilestones(ctx context.Context, studentID uint64) ([]*model.Milestone, error) {
	// 获取学生的里程碑列表（勋章墙）
	// 这里需要根据实际的里程碑模型来实现
	// 暂时返回空数组
	return []*model.Milestone{}, nil
}

// ─── 4.5 盲盒查看 ──────────────────────────────────────────

func (s *studentServiceImpl) GetMyBlindboxDraws(ctx context.Context, studentID uint64) ([]*model.BlindBoxDraw, error) {
	// 获取学生已抽到的盲盒奖励
	// 这里需要根据实际的盲盒抽取记录模型来实现
	// 暂时返回空数组
	return []*model.BlindBoxDraw{}, nil
}
