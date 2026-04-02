// growth-partner/backend/internal/service/teacher_service.go
// 老师端模块服务

package service

import (
	"context"
	"errors"
	"time"

	"growth-partner/internal/model"
	"growth-partner/internal/repository"
)

var (
	ErrTeacherUnauthorized = errors.New("老师未授权访问该资源")
	ErrInvalidClassID      = errors.New("无效的班级ID")
	ErrBehaviorNotFound    = errors.New("行为记录不存在")
	ErrBehaviorExpired     = errors.New("行为记录已超过撤销时限")
)

// TeacherService 老师端服务接口
type TeacherService interface {
	// ─── 3.1 我的班级管理 ────────────────────────────────────
	GetMyClasses(ctx context.Context, teacherID uint64) ([]*model.Class, error)
	GetClassOverview(ctx context.Context, classID uint64) (map[string]interface{}, error)
	GetClassStudents(ctx context.Context, classID uint64, params map[string]interface{}) ([]*model.Child, int64, error)

	// ─── 3.2 正向行为打分 ────────────────────────────────────
	GetBehavior(ctx context.Context, behaviorID uint64) (*model.BehaviorRecord, error)
	DeleteBehavior(ctx context.Context, behaviorID, teacherID uint64) error
	BatchRecordBehaviors(ctx context.Context, requests []RecordBehaviorRequest) ([]*model.BehaviorRecord, error)

	// ─── 3.3 广播发送 ──────────────────────────────────────
	GetBroadcasts(ctx context.Context, teacherID uint64, params map[string]interface{}) ([]*model.Broadcast, int64, error)
	CancelBroadcast(ctx context.Context, broadcastID, teacherID uint64) error

	// ─── 3.4 集体挑战管理 ──────────────────────────────────
	GetChallenges(ctx context.Context, classID uint64) ([]*model.Challenge, error)
	CreateChallenge(ctx context.Context, challenge *model.Challenge) error
	CompleteChallenge(ctx context.Context, challengeID, teacherID uint64) error

	// ─── 3.5 题库管理 ──────────────────────────────────────
	GetQuestions(ctx context.Context, classID uint64, params map[string]interface{}) ([]*model.Question, int64, error)
	CreateQuestion(ctx context.Context, question *model.Question) error
	UpdateQuestion(ctx context.Context, question *model.Question, teacherID uint64) error
	DeleteQuestion(ctx context.Context, questionID, teacherID uint64) error
	BatchImportQuestions(ctx context.Context, questions []*model.Question, teacherID uint64) error

	// ─── 3.6 盲盒奖励池管理 ────────────────────────────────
	UpdateBlindboxPoolItem(ctx context.Context, blindbox *model.BlindBoxPool, teacherID uint64) error
	ConfirmBlindboxRedeem(ctx context.Context, drawID, teacherID uint64) error

	// ─── 3.7 周报PDF生成 ──────────────────────────────────
	GenerateWeeklyReport(ctx context.Context, classID, teacherID uint64) error
	GetWeeklyReports(ctx context.Context, classID, teacherID uint64, params map[string]interface{}) ([]*model.Report, int64, error)
	DownloadWeeklyReport(ctx context.Context, reportID uint64) ([]byte, string, error)
}

// teacherServiceImpl 老师端服务实现
type teacherServiceImpl struct {
	// 班级相关
	classRepo repository.ClassRepository

	// 学生相关
	childRepo repository.ChildRepository

	// 行为相关
	behaviorRepo repository.BehaviorRepository

	// 伙伴相关
	partnerRepo repository.PartnerRepository

	// 广播相关
	broadcastRepo repository.BroadcastRepository

	// 挑战相关
	challengeRepo repository.ChallengeRepository

	// 题库相关
	questionRepo repository.QuestionRepository

	// 盲盒相关
	blindboxRepo repository.BlindboxRepository

	// 报告相关
	reportRepo repository.ReportRepository

	// 权限相关
	adminPermissionRepo repository.AdminPermissionRepository

	// 行为服务（复用已有功能）
	behaviorSvc BehaviorService

	// 广播服务（复用已有功能）
	broadcastSvc BroadcastService

	// 盲盒服务（复用已有功能）
	blindboxSvc BlindboxService
}

// NewTeacherService 创建老师端服务实例
func NewTeacherService(
	classRepo repository.ClassRepository,
	childRepo repository.ChildRepository,
	behaviorRepo repository.BehaviorRepository,
	partnerRepo repository.PartnerRepository,
	broadcastRepo repository.BroadcastRepository,
	challengeRepo repository.ChallengeRepository,
	questionRepo repository.QuestionRepository,
	blindboxRepo repository.BlindboxRepository,
	reportRepo repository.ReportRepository,
	adminPermissionRepo repository.AdminPermissionRepository,
	behaviorSvc BehaviorService,
	broadcastSvc BroadcastService,
	blindboxSvc BlindboxService,
) TeacherService {
	return &teacherServiceImpl{
		classRepo:           classRepo,
		childRepo:           childRepo,
		behaviorRepo:        behaviorRepo,
		partnerRepo:         partnerRepo,
		broadcastRepo:       broadcastRepo,
		challengeRepo:       challengeRepo,
		questionRepo:        questionRepo,
		blindboxRepo:        blindboxRepo,
		reportRepo:          reportRepo,
		adminPermissionRepo: adminPermissionRepo,
		behaviorSvc:         behaviorSvc,
		broadcastSvc:        broadcastSvc,
		blindboxSvc:         blindboxSvc,
	}
}

// ─── 3.1 我的班级管理 ──────────────────────────────────────────

func (s *teacherServiceImpl) GetMyClasses(ctx context.Context, teacherID uint64) ([]*model.Class, error) {
	// 获取老师的所有权限
	permissions, err := s.adminPermissionRepo.FindByTeacherID(ctx, teacherID)
	if err != nil {
		return nil, err
	}

	// 从权限中提取班级ID并去重
	classIDs := make(map[uint64]bool)
	for _, permission := range permissions {
		classIDs[permission.ClassID] = true
	}

	// 获取班级详情
	var classes []*model.Class
	for classID := range classIDs {
		class, err := s.classRepo.FindByID(ctx, classID)
		if err != nil {
			continue
		}
		if class != nil && class.IsActive {
			classes = append(classes, class)
		}
	}

	return classes, nil
}

func (s *teacherServiceImpl) GetClassOverview(ctx context.Context, classID uint64) (map[string]interface{}, error) {
	// 检查班级是否存在
	class, err := s.classRepo.FindByID(ctx, classID)
	if err != nil {
		return nil, err
	}
	if class == nil {
		return nil, ErrInvalidClassID
	}

	// 获取学生数量
	childParams := map[string]interface{}{"class_id": classID, "is_active": true}
	students, total, err := s.childRepo.FindAll(ctx, childParams)
	if err != nil {
		return nil, err
	}

	// 计算总成长值
	totalGrowthPoints := 0
	for _, student := range students {
		totalGrowthPoints += student.TotalGrowthPoints
	}

	// 获取行为记录数量
	_, behaviorCount, err := s.behaviorRepo.FindByClassID(ctx, classID, 1, 0)
	if err != nil {
		return nil, err
	}

	// 构建概览数据
	overview := map[string]interface{}{
		"class_name":            class.ClassName,
		"class_code":            class.ClassCode,
		"grade":                 class.Grade,
		"student_count":         total,
		"behavior_count":        behaviorCount,
		"total_growth_points":   totalGrowthPoints,
		"average_growth_points": 0,
	}

	// 计算平均成长值
	if total > 0 {
		overview["average_growth_points"] = totalGrowthPoints / int(total)
	}

	return overview, nil
}

func (s *teacherServiceImpl) GetClassStudents(ctx context.Context, classID uint64, params map[string]interface{}) ([]*model.Child, int64, error) {
	// 检查班级是否存在
	class, err := s.classRepo.FindByID(ctx, classID)
	if err != nil {
		return nil, 0, err
	}
	if class == nil {
		return nil, 0, ErrInvalidClassID
	}

	// 构建查询参数
	childParams := map[string]interface{}{"class_id": classID, "is_active": true}
	if name, ok := params["name"].(string); ok && name != "" {
		childParams["name"] = name
	}

	// 获取学生列表
	students, total, err := s.childRepo.FindAll(ctx, childParams)
	if err != nil {
		return nil, 0, err
	}

	// 为每个学生获取伙伴信息
	// 注意：由于Child模型没有Partner字段，我们暂时不返回伙伴信息
	// 后续可以通过修改Child模型或者返回一个包含伙伴信息的DTO来解决

	return students, total, nil
}

// ─── 3.2 正向行为打分 ──────────────────────────────────────────

func (s *teacherServiceImpl) GetBehavior(ctx context.Context, behaviorID uint64) (*model.BehaviorRecord, error) {
	// 注意：由于BehaviorRepository没有FindByID方法，暂时返回空
	// 后续需要在BehaviorRepository中添加FindByID方法
	return nil, nil
}

func (s *teacherServiceImpl) DeleteBehavior(ctx context.Context, behaviorID, teacherID uint64) error {
	// 注意：由于BehaviorRepository没有FindByID、BeginTx和Delete方法，暂时返回错误
	// 后续需要在BehaviorRepository中添加这些方法
	return errors.New("delete behavior not supported yet")
}

func (s *teacherServiceImpl) BatchRecordBehaviors(ctx context.Context, requests []RecordBehaviorRequest) ([]*model.BehaviorRecord, error) {
	var behaviors []*model.BehaviorRecord

	// 为每个学生创建行为记录
	for _, req := range requests {
		behavior, _, err := s.behaviorSvc.RecordBehavior(ctx, req)
		if err != nil {
			continue
		}
		if behavior != nil {
			behaviors = append(behaviors, behavior)
		}
	}

	return behaviors, nil
}

// ─── 3.3 广播发送 ──────────────────────────────────────────────

func (s *teacherServiceImpl) GetBroadcasts(ctx context.Context, teacherID uint64, params map[string]interface{}) ([]*model.Broadcast, int64, error) {
	// 获取广播列表
	broadcasts, err := s.broadcastRepo.FindByCreatorID(ctx, teacherID)
	if err != nil {
		return nil, 0, err
	}

	return broadcasts, int64(len(broadcasts)), nil
}

func (s *teacherServiceImpl) CancelBroadcast(ctx context.Context, broadcastID, teacherID uint64) error {
	// 获取广播
	broadcast, err := s.broadcastRepo.FindByID(ctx, broadcastID)
	if err != nil {
		return err
	}
	if broadcast == nil {
		return errors.New("广播不存在")
	}

	// 检查是否是老师自己发送的广播
	if broadcast.CreatedBy != teacherID || broadcast.CreatedByRole != "teacher" {
		return ErrTeacherUnauthorized
	}

	// 检查广播状态，只能取消未发送的广播
	if broadcast.Status != "pending" {
		return errors.New("只能取消未发送的广播")
	}

	// 取消广播
	broadcast.Status = "cancelled"
	return s.broadcastRepo.Update(ctx, broadcast)
}

// ─── 3.4 集体挑战管理 ──────────────────────────────────────────

func (s *teacherServiceImpl) GetChallenges(ctx context.Context, classID uint64) ([]*model.Challenge, error) {
	// 检查班级是否存在
	class, err := s.classRepo.FindByID(ctx, classID)
	if err != nil {
		return nil, err
	}
	if class == nil {
		return nil, ErrInvalidClassID
	}

	// 获取班级的所有挑战
	challenges, err := s.challengeRepo.FindByClassID(ctx, classID)
	if err != nil {
		return nil, err
	}

	// 过滤出进行中的挑战
	var inProgressChallenges []*model.Challenge
	for _, challenge := range challenges {
		if challenge.Status == "in_progress" {
			inProgressChallenges = append(inProgressChallenges, challenge)
		}
	}

	return inProgressChallenges, nil
}

func (s *teacherServiceImpl) CreateChallenge(ctx context.Context, challenge *model.Challenge) error {
	// 检查班级是否存在
	class, err := s.classRepo.FindByID(ctx, challenge.ClassID)
	if err != nil {
		return err
	}
	if class == nil {
		return ErrInvalidClassID
	}

	// 设置挑战状态为进行中
	challenge.Status = "in_progress"

	// 创建挑战
	return s.challengeRepo.Create(ctx, challenge)
}

func (s *teacherServiceImpl) CompleteChallenge(ctx context.Context, challengeID, teacherID uint64) error {
	// 获取挑战
	challenge, err := s.challengeRepo.FindByID(ctx, challengeID)
	if err != nil {
		return err
	}
	if challenge == nil {
		return errors.New("挑战不存在")
	}

	// 检查挑战状态，只能完成进行中的挑战
	if challenge.Status != "active" {
		return errors.New("只能完成进行中的挑战")
	}

	// 获取班级学生
	childParams := map[string]interface{}{"class_id": challenge.ClassID, "is_active": true}
	students, _, err := s.childRepo.FindAll(ctx, childParams)
	if err != nil {
		return err
	}

	// 为每个学生发放成长值
	for _, student := range students {
		// 增加成长值
		student.CurrentGrowthPoints += challenge.RewardValue
		student.TotalGrowthPoints += challenge.RewardValue

		if err := s.childRepo.Update(ctx, student); err != nil {
			return err
		}
	}

	// 更新挑战状态
	return s.challengeRepo.Complete(ctx, challengeID, teacherID)
}

// ─── 3.5 题库管理 ──────────────────────────────────────────────

func (s *teacherServiceImpl) GetQuestions(ctx context.Context, classID uint64, params map[string]interface{}) ([]*model.Question, int64, error) {
	// 检查班级是否存在
	class, err := s.classRepo.FindByID(ctx, classID)
	if err != nil {
		return nil, 0, err
	}
	if class == nil {
		return nil, 0, ErrInvalidClassID
	}

	// 构建查询参数，获取公共题目和班级专属题目
	questionParams := map[string]interface{}{}

	// 应用过滤条件
	if subject, ok := params["subject"].(string); ok && subject != "" {
		questionParams["subject"] = subject
	}
	if difficulty, ok := params["difficulty"].(string); ok && difficulty != "" {
		questionParams["difficulty"] = difficulty
	}

	// 获取题目列表
	questions, total, err := s.questionRepo.FindAll(ctx, questionParams)
	if err != nil {
		return nil, 0, err
	}

	return questions, total, nil
}

func (s *teacherServiceImpl) CreateQuestion(ctx context.Context, question *model.Question) error {
	// 检查班级是否存在
	if question.ClassID != nil {
		class, err := s.classRepo.FindByID(ctx, *question.ClassID)
		if err != nil {
			return err
		}
		if class == nil {
			return ErrInvalidClassID
		}
	}

	// 创建题目
	return s.questionRepo.Create(ctx, question)
}

func (s *teacherServiceImpl) UpdateQuestion(ctx context.Context, question *model.Question, teacherID uint64) error {
	// 获取题目
	existingQuestion, err := s.questionRepo.FindByID(ctx, question.ID)
	if err != nil {
		return err
	}
	if existingQuestion == nil {
		return errors.New("题目不存在")
	}

	// 检查老师是否有权限编辑该题目
	// 获取老师的班级权限
	permissions, err := s.adminPermissionRepo.FindByTeacherID(ctx, teacherID)
	if err != nil {
		return err
	}

	// 检查老师是否有该班级的权限
	hasPermission := false
	if existingQuestion.ClassID != nil {
		for _, permission := range permissions {
			if permission.ClassID == *existingQuestion.ClassID {
				hasPermission = true
				break
			}
		}
	}

	if !hasPermission {
		return ErrTeacherUnauthorized
	}

	// 更新题目
	return s.questionRepo.Update(ctx, question)
}

func (s *teacherServiceImpl) DeleteQuestion(ctx context.Context, questionID, teacherID uint64) error {
	// 获取题目
	question, err := s.questionRepo.FindByID(ctx, questionID)
	if err != nil {
		return err
	}
	if question == nil {
		return errors.New("题目不存在")
	}

	// 检查老师是否有权限删除该题目
	// 获取老师的班级权限
	permissions, err := s.adminPermissionRepo.FindByTeacherID(ctx, teacherID)
	if err != nil {
		return err
	}

	// 检查老师是否有该班级的权限
	hasPermission := false
	if question.ClassID != nil {
		for _, permission := range permissions {
			if permission.ClassID == *question.ClassID {
				hasPermission = true
				break
			}
		}
	}

	if !hasPermission {
		return ErrTeacherUnauthorized
	}

	// 软删除题目
	return s.questionRepo.Delete(ctx, questionID)
}

func (s *teacherServiceImpl) BatchImportQuestions(ctx context.Context, questions []*model.Question, teacherID uint64) error {
	if len(questions) == 0 {
		return errors.New("没有题目数据")
	}

	// 检查老师是否有权限导入题目到指定班级
	// 获取老师的班级权限
	permissions, err := s.adminPermissionRepo.FindByTeacherID(ctx, teacherID)
	if err != nil {
		return err
	}

	// 检查所有题目的班级ID是否一致，并且老师有该班级的权限
	classID := questions[0].ClassID
	hasPermission := false
	if classID != nil {
		for _, permission := range permissions {
			if permission.ClassID == *classID {
				hasPermission = true
				break
			}
		}
	}

	if !hasPermission {
		return ErrTeacherUnauthorized
	}

	// 批量导入题目
	return s.questionRepo.BatchCreate(ctx, questions)
}

// ─── 3.6 盲盒奖励池管理 ────────────────────────────────────────

func (s *teacherServiceImpl) UpdateBlindboxPoolItem(ctx context.Context, blindbox *model.BlindBoxPool, teacherID uint64) error {
	// 注意：由于BlindboxRepository没有Update方法，暂时返回错误
	// 后续需要在BlindboxRepository中添加Update方法
	return errors.New("update blindbox pool item not supported yet")
}

func (s *teacherServiceImpl) ConfirmBlindboxRedeem(ctx context.Context, drawID, teacherID uint64) error {
	// 注意：由于BlindboxRepository没有FindDrawByID和UpdateDraw方法，暂时返回错误
	// 后续需要在BlindboxRepository中添加这些方法
	return errors.New("confirm blindbox redeem not supported yet")
}

// ─── 3.7 周报PDF生成 ──────────────────────────────────────────

func (s *teacherServiceImpl) GenerateWeeklyReport(ctx context.Context, classID, teacherID uint64) error {
	// 检查班级是否存在
	class, err := s.classRepo.FindByID(ctx, classID)
	if err != nil {
		return err
	}
	if class == nil {
		return ErrInvalidClassID
	}

	// 检查老师是否有权限操作该班级
	// 获取老师的班级权限
	permissions, err := s.adminPermissionRepo.FindByTeacherID(ctx, teacherID)
	if err != nil {
		return err
	}

	// 检查老师是否有该班级的权限
	hasPermission := false
	for _, permission := range permissions {
		if permission.ClassID == classID {
			hasPermission = true
			break
		}
	}

	if !hasPermission {
		return ErrTeacherUnauthorized
	}

	// 异步生成周报PDF
	go func() {
		// 这里应该实现PDF生成的逻辑
		// 1. 获取班级本周的行为记录
		// 2. 计算各种统计数据
		// 3. 生成PDF文件
		// 4. 保存PDF文件路径到数据库

		// 这里只是一个示例，实际实现需要使用PDF生成库
		now := time.Now()
		report := &model.Report{
			ClassID:      classID,
			ReportType:   "weekly",
			ReportPeriod: "2024-03-18 至 2024-03-24",
			StartDate:    time.Now().AddDate(0, 0, -7),
			EndDate:      time.Now(),
			Status:       "ready",
			FileURL:      "/path/to/pdf/report.pdf",
			GeneratedBy:  teacherID,
			GeneratedAt:  &now,
		}

		// 保存报告记录
		s.reportRepo.Create(context.Background(), report)
	}()

	return nil
}

func (s *teacherServiceImpl) GetWeeklyReports(ctx context.Context, classID, teacherID uint64, params map[string]interface{}) ([]*model.Report, int64, error) {
	// 检查班级是否存在
	class, err := s.classRepo.FindByID(ctx, classID)
	if err != nil {
		return nil, 0, err
	}
	if class == nil {
		return nil, 0, ErrInvalidClassID
	}

	// 检查老师是否有权限操作该班级
	// 获取老师的班级权限
	permissions, err := s.adminPermissionRepo.FindByTeacherID(ctx, teacherID)
	if err != nil {
		return nil, 0, err
	}

	// 检查老师是否有该班级的权限
	hasPermission := false
	for _, permission := range permissions {
		if permission.ClassID == classID {
			hasPermission = true
			break
		}
	}

	if !hasPermission {
		return nil, 0, ErrTeacherUnauthorized
	}

	// 构建查询参数
	reportParams := map[string]interface{}{"report_type": "weekly"}

	// 应用过滤条件
	if startDate, ok := params["start_date"].(string); ok && startDate != "" {
		reportParams["start_date"] = startDate
	}
	if endDate, ok := params["end_date"].(string); ok && endDate != "" {
		reportParams["end_date"] = endDate
	}

	// 获取周报列表
	reports, total, err := s.reportRepo.FindByClassID(ctx, classID, reportParams)
	if err != nil {
		return nil, 0, err
	}

	return reports, total, nil
}

func (s *teacherServiceImpl) DownloadWeeklyReport(ctx context.Context, reportID uint64) ([]byte, string, error) {
	// 获取报告
	report, err := s.reportRepo.FindByID(ctx, reportID)
	if err != nil {
		return nil, "", err
	}
	if report == nil {
		return nil, "", errors.New("报告不存在")
	}

	// 检查报告状态，只能下载已生成的报告
	if report.Status != "generated" {
		return nil, "", errors.New("报告尚未生成")
	}

	// 读取PDF文件
	// 这里只是一个示例，实际实现需要根据reportPath读取文件
	// 实际项目中应该使用os.ReadFile或类似方法读取文件
	pdfData := []byte("PDF文件内容")
	fileName := "weekly_report.pdf"

	return pdfData, fileName, nil
}
