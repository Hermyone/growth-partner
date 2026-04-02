// growth-partner/backend/internal/service/partner_service.go
// 伙伴核心业务服务：模板管理、成长值计算、进化逻辑
// 核心设计：纯正向激励，成长值只增不减，进化不可逆

package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"growth-partner/internal/model"
	"growth-partner/internal/repository"
)

// ─── 错误定义 ──────────────────────────────────────────────────

var (
	ErrPartnerNotFound     = errors.New("伙伴不存在")
	ErrTemplateNotFound    = errors.New("伙伴模板不存在")
	ErrPartnerAlreadyExist = errors.New("该学生已拥有伙伴，不可重复创建")
	ErrInvalidGrowthValue  = errors.New("成长值必须为正整数")
)

// ─── 接口定义 ──────────────────────────────────────────────────

// PartnerService 伙伴服务接口
type PartnerService interface {
	// 模板相关
	GetAllTemplates(ctx context.Context) ([]*model.PartnerTemplate, error)
	GetTemplateByID(ctx context.Context, id uint64) (*model.PartnerTemplate, error)

	// 伙伴实例相关
	CreatePartner(ctx context.Context, childID, templateID uint64, nickname string) (*model.Partner, error)
	GetPartnerByChildID(ctx context.Context, childID uint64) (*model.Partner, error)
	UpdateNickname(ctx context.Context, childID uint64, nickname string) error

	// 核心业务：增加成长值（触发进化检测）
	AddGrowthPoints(ctx context.Context, req AddGrowthPointsRequest) (*AddGrowthPointsResult, error)

	// 查询成长记录
	GetGrowthHistory(ctx context.Context, childID uint64, limit, offset int) ([]*model.GrowthRecord, int64, error)
}

// ─── DTO 定义 ──────────────────────────────────────────────────

// AddGrowthPointsRequest 增加成长值请求
type AddGrowthPointsRequest struct {
	ChildID    uint64
	BehaviorID uint64 // 来源行为记录ID
	SourceType model.GrowthSourceType
	Delta      int // 增加的成长值（必须>0）
	Remark     string
}

// AddGrowthPointsResult 增加成长值结果
type AddGrowthPointsResult struct {
	Partner        *model.Partner
	GrowthRecord   *model.GrowthRecord
	IsEvolved      bool                 // 是否触发了进化
	FromStage      model.EvolutionStage // 进化前阶段
	ToStage        model.EvolutionStage // 进化后阶段
	EvolutionMsg   string               // 进化祝贺语
	PartnerMessage string               // 伙伴说的鼓励话
	Milestones     []*model.Milestone   // 本次触发的里程碑
}

// ─── 实现 ──────────────────────────────────────────────────────

// partnerServiceImpl 伙伴服务实现
type partnerServiceImpl struct {
	partnerRepo   repository.PartnerRepository
	growthRepo    repository.GrowthRepository
	templateRepo  repository.TemplateRepository
	milestoneRepo repository.MilestoneRepository
	// 广播服务（进化时通知）
	broadcastService BroadcastService
}

// NewPartnerService 创建伙伴服务实例
func NewPartnerService(
	partnerRepo repository.PartnerRepository,
	growthRepo repository.GrowthRepository,
	templateRepo repository.TemplateRepository,
	milestoneRepo repository.MilestoneRepository,
	broadcastSvc BroadcastService,
) PartnerService {
	return &partnerServiceImpl{
		partnerRepo:      partnerRepo,
		growthRepo:       growthRepo,
		templateRepo:     templateRepo,
		milestoneRepo:    milestoneRepo,
		broadcastService: broadcastSvc,
	}
}

// GetAllTemplates 获取所有激活的伙伴模板
func (s *partnerServiceImpl) GetAllTemplates(ctx context.Context) ([]*model.PartnerTemplate, error) {
	return s.templateRepo.FindAllActive(ctx)
}

// GetTemplateByID 按 ID 获取模板
func (s *partnerServiceImpl) GetTemplateByID(ctx context.Context, id uint64) (*model.PartnerTemplate, error) {
	tmpl, err := s.templateRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if tmpl == nil {
		return nil, ErrTemplateNotFound
	}
	return tmpl, nil
}

// CreatePartner 为学生创建专属伙伴（每人只能有一只）
func (s *partnerServiceImpl) CreatePartner(ctx context.Context, childID, templateID uint64, nickname string) (*model.Partner, error) {
	// 检查是否已有伙伴
	existing, err := s.partnerRepo.FindByChildID(ctx, childID)
	if err != nil {
		return nil, fmt.Errorf("查询伙伴失败: %w", err)
	}
	if existing != nil {
		return nil, ErrPartnerAlreadyExist
	}

	// 验证模板是否存在且激活
	if _, err := s.GetTemplateByID(ctx, templateID); err != nil {
		return nil, err
	}

	partner := &model.Partner{
		ChildID:          childID,
		TemplateID:       templateID,
		SequenceNo:       1, // 第一个伙伴
		GrowthPoints:     0,
		CurrentStage:     model.StageLow,
		EvolutionCount:   0,
		Nickname:         nickname,
		InteractionLevel: 0,
		Status:           model.PartnerStatusActive,
		SchoolYear:       "2024-2025", // 暂时设置为当前学年，实际应该从请求中获取
	}

	if err := s.partnerRepo.Create(ctx, partner); err != nil {
		return nil, fmt.Errorf("创建伙伴失败: %w", err)
	}

	log.Printf("[PartnerService] 学生 %d 创建伙伴成功，模板: %d，昵称: %s", childID, templateID, nickname)
	return partner, nil
}

// GetPartnerByChildID 获取学生的伙伴
func (s *partnerServiceImpl) GetPartnerByChildID(ctx context.Context, childID uint64) (*model.Partner, error) {
	partner, err := s.partnerRepo.FindByChildID(ctx, childID)
	if err != nil {
		return nil, fmt.Errorf("查询伙伴失败: %w", err)
	}
	if partner == nil {
		return nil, ErrPartnerNotFound
	}
	return partner, nil
}

// UpdateNickname 更新伙伴昵称
func (s *partnerServiceImpl) UpdateNickname(ctx context.Context, childID uint64, nickname string) error {
	return s.partnerRepo.UpdateNickname(ctx, childID, nickname)
}

// AddGrowthPoints 核心业务：为伙伴增加成长值，自动检测并触发进化
// 设计原则：
//  1. 成长值只增不减（纯正向）
//  2. 进化检测：超过阈值自动进化
//  3. 进化不可逆：即使不再继续增加成长值，阶段不降级
//  4. 进化触发广播通知，让孩子有仪式感
func (s *partnerServiceImpl) AddGrowthPoints(ctx context.Context, req AddGrowthPointsRequest) (*AddGrowthPointsResult, error) {
	// 参数校验
	if req.Delta <= 0 {
		return nil, ErrInvalidGrowthValue
	}

	// 获取当前伙伴状态
	partner, err := s.GetPartnerByChildID(ctx, req.ChildID)
	if err != nil {
		return nil, err
	}

	// 记录进化前阶段
	fromStage := partner.CurrentStage
	oldPoints := partner.GrowthPoints

	// 计算新成长值
	newPoints := oldPoints + req.Delta

	// ─── 进化检测 ───────────────────────────────────────────
	newStage := calculateStage(newPoints)
	isEvolved := newStage > fromStage

	// 更新伙伴状态
	partner.GrowthPoints = newPoints
	partner.CurrentStage = newStage
	if isEvolved {
		now := time.Now()
		partner.EvolutionCount++
		partner.LastEvolvedAt = &now
		if partner.FirstEvolutionAt == nil {
			partner.FirstEvolutionAt = &now
		}
		// 高阶进化解锁中阶交互
		if newStage >= model.StageMid {
			partner.InteractionLevel = 1
		}
		// 高阶进化解锁大模型交互
		if newStage >= model.StageHigh {
			partner.InteractionLevel = 2
		}
	}

	// 持久化伙伴状态更新
	if err := s.partnerRepo.UpdateGrowthState(ctx, partner); err != nil {
		return nil, fmt.Errorf("更新伙伴成长状态失败: %w", err)
	}

	// 生成成长流水记录
	growthRecord := &model.GrowthRecord{
		ChildID:            req.ChildID,
		PartnerID:          partner.ID,
		SourceType:         req.SourceType,
		SourceID:           req.BehaviorID,
		Delta:              req.Delta,
		AfterPoints:        newPoints,
		IsEvolutionTrigger: isEvolved,
		Remark:             req.Remark,
	}
	if isEvolved {
		growthRecord.EvolutionFromStage = &fromStage
		growthRecord.EvolutionToStage = &newStage
	}
	if err := s.growthRepo.Create(ctx, growthRecord); err != nil {
		log.Printf("[PartnerService] 警告：创建成长记录失败（不影响主流程）: %v", err)
	}

	// ─── 构建结果 ────────────────────────────────────────────
	result := &AddGrowthPointsResult{
		Partner:      partner,
		GrowthRecord: growthRecord,
		IsEvolved:    isEvolved,
		FromStage:    fromStage,
		ToStage:      newStage,
	}

	if isEvolved {
		result.EvolutionMsg = buildEvolutionMessage(newStage)
		// 异步触发进化广播（不阻塞主流程）
		go func() {
			if err := s.broadcastService.BroadcastEvolution(context.Background(), req.ChildID, partner.ID, fromStage, newStage); err != nil {
				log.Printf("[PartnerService] 进化广播失败: %v", err)
			}
		}()
	}

	// 检查里程碑（异步）
	go s.checkMilestones(context.Background(), req.ChildID, partner, result)

	log.Printf("[PartnerService] 学生 %d 获得 +%d 成长值，当前: %d，阶段: %d→%d，进化: %v",
		req.ChildID, req.Delta, newPoints, fromStage, newStage, isEvolved)

	return result, nil
}

// GetGrowthHistory 获取成长值历史记录（分页）
func (s *partnerServiceImpl) GetGrowthHistory(ctx context.Context, childID uint64, limit, offset int) ([]*model.GrowthRecord, int64, error) {
	return s.growthRepo.FindByChildID(ctx, childID, limit, offset)
}

// ─── 内部辅助函数 ──────────────────────────────────────────────

// calculateStage 根据成长值计算所处阶段
func calculateStage(points int) model.EvolutionStage {
	switch {
	case points >= model.StageThresholds[model.StageHigh]:
		return model.StageHigh
	case points >= model.StageThresholds[model.StageMid]:
		return model.StageMid
	default:
		return model.StageLow
	}
}

// buildEvolutionMessage 构建进化祝贺语
func buildEvolutionMessage(stage model.EvolutionStage) string {
	switch stage {
	case model.StageMid:
		return "🎉 恭喜！你的伙伴成功进化到中阶形态！你的努力让它变得更强大了！"
	case model.StageHigh:
		return "✨ 哇！太厉害了！你的伙伴进化到最强高阶形态！这是属于你们的荣耀时刻！"
	default:
		return "你的伙伴在成长！"
	}
}

// checkMilestones 异步检查里程碑（不阻塞主流程）
func (s *partnerServiceImpl) checkMilestones(ctx context.Context, childID uint64, partner *model.Partner, result *AddGrowthPointsResult) {
	// 检查初次进化里程碑
	if result.IsEvolved && partner.EvolutionCount == 1 {
		cfg := model.MilestoneConfig[model.MilestoneFirstEvolution]
		milestone := &model.Milestone{
			ChildID:    childID,
			Type:       model.MilestoneFirstEvolution,
			Title:      cfg.Title,
			Content:    cfg.Content,
			SourceType: "partner",
			SourceID:   partner.ID,
		}
		if err := s.milestoneRepo.Create(ctx, milestone); err != nil {
			log.Printf("[PartnerService] 创建里程碑失败: %v", err)
		}
	}
	// TODO: 其他里程碑检查（连续7天打卡、累计10次德馨等）
}
