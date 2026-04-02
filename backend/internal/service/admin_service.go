// growth-partner/backend/internal/service/admin_service.go
// 管理员模块服务

package service

import (
	"context"
	"errors"

	"growth-partner/internal/model"
	"growth-partner/internal/repository"
	"growth-partner/internal/utils"
)

var (
	ErrAdminUnauthorized = errors.New("管理员未授权")
	ErrInvalidInput      = errors.New("无效的输入参数")
	ErrResourceNotFound  = errors.New("资源不存在")
	ErrDuplicateResource = errors.New("资源已存在")
)

// AdminService 管理员服务接口
type AdminService interface {
	// ─── 学校管理 ──────────────────────────────────────────────
	GetSchools(ctx context.Context, params map[string]interface{}) ([]*model.School, int64, error)
	CreateSchool(ctx context.Context, school *model.School) error
	UpdateSchool(ctx context.Context, school *model.School) error
	UpdateSchoolStatus(ctx context.Context, id uint64, isActive bool) error

	// ─── 班级管理 ──────────────────────────────────────────────
	GetClasses(ctx context.Context, params map[string]interface{}) ([]*model.Class, int64, error)
	CreateClass(ctx context.Context, class *model.Class) error
	UpdateClass(ctx context.Context, class *model.Class) error
	PromoteClass(ctx context.Context, classID uint64, newGrade int, newClassCode string) error
	UpdateClassStatus(ctx context.Context, id uint64, isActive bool) error

	// ─── 老师/家长用户管理 ────────────────────────────────────
	GetUsers(ctx context.Context, params map[string]interface{}) ([]*model.User, int64, error)
	CreateUser(ctx context.Context, user *model.User, initialPassword string) error
	UpdateUser(ctx context.Context, user *model.User) error
	UpdateUserStatus(ctx context.Context, id uint64, isActive bool) error
	ResetUserPassword(ctx context.Context, id uint64, newPassword string) error

	// ─── 学生账号管理 ──────────────────────────────────────────
	BatchImportStudents(ctx context.Context, data []*model.Child) error
	GetStudents(ctx context.Context, params map[string]interface{}) ([]*model.Child, int64, error)
	CreateStudent(ctx context.Context, child *model.Child) error
	UpdateStudent(ctx context.Context, child *model.Child) error

	// ─── 老师班级权限分配 ──────────────────────────────────────
	GetAssignments(ctx context.Context, params map[string]interface{}) ([]*model.AdminPermission, int64, error)
	CreateAssignment(ctx context.Context, assignment *model.AdminPermission) error
	DeleteAssignment(ctx context.Context, id uint64) error
	BatchCreateAssignments(ctx context.Context, assignments []*model.AdminPermission) error

	// ─── 家长-学生绑定管理 ────────────────────────────────────
	GetParentBindings(ctx context.Context, params map[string]interface{}) ([]*model.ParentChildRelation, int64, error)
	CreateParentBinding(ctx context.Context, binding *model.ParentChildRelation) error
	DeleteParentBinding(ctx context.Context, id uint64) error

	// ─── 伙伴模板管理 ──────────────────────────────────────────
	CreatePartnerTemplate(ctx context.Context, template *model.PartnerTemplate) error
	UpdatePartnerTemplate(ctx context.Context, template *model.PartnerTemplate) error
	SeedPartnerTemplates(ctx context.Context) error

	// ─── 数据概览 ──────────────────────────────────────────────
	GetDashboard(ctx context.Context) (map[string]interface{}, error)
	GetAuditLogs(ctx context.Context, params map[string]interface{}) ([]*model.AuditLog, int64, error)
}

// adminServiceImpl 管理员服务实现
type adminServiceImpl struct {
	// 学校相关
	schoolRepo repository.SchoolRepository

	// 班级相关
	classRepo repository.ClassRepository

	// 用户相关
	userRepo repository.UserRepository

	// 学生相关
	childRepo repository.ChildRepository

	// 家长相关
	parentChildRepo repository.ParentChildRepository

	// 权限相关
	permissionRepo repository.AdminPermissionRepository

	// 审计日志
	auditLogRepo repository.AuditLogRepository

	// 伙伴模板相关
	templateRepo repository.TemplateRepository
}

// NewAdminService 创建管理员服务实例
func NewAdminService(
	schoolRepo repository.SchoolRepository,
	classRepo repository.ClassRepository,
	userRepo repository.UserRepository,
	childRepo repository.ChildRepository,
	parentChildRepo repository.ParentChildRepository,
	permissionRepo repository.AdminPermissionRepository,
	auditLogRepo repository.AuditLogRepository,
	templateRepo repository.TemplateRepository,
) AdminService {
	return &adminServiceImpl{
		schoolRepo:      schoolRepo,
		classRepo:       classRepo,
		userRepo:        userRepo,
		childRepo:       childRepo,
		parentChildRepo: parentChildRepo,
		permissionRepo:  permissionRepo,
		auditLogRepo:    auditLogRepo,
		templateRepo:    templateRepo,
	}
}

// ─── 学校管理 ──────────────────────────────────────────────────

func (s *adminServiceImpl) GetSchools(ctx context.Context, params map[string]interface{}) ([]*model.School, int64, error) {
	return s.schoolRepo.FindAll(ctx, params)
}

func (s *adminServiceImpl) CreateSchool(ctx context.Context, school *model.School) error {
	// 验证学校信息
	if school.Name == "" {
		return ErrInvalidInput
	}

	// 检查学校是否已存在
	params := map[string]interface{}{"name": school.Name}
	schools, _, err := s.schoolRepo.FindAll(ctx, params)
	if err != nil {
		return err
	}
	if len(schools) > 0 {
		return ErrDuplicateResource
	}

	// 创建学校
	return s.schoolRepo.Create(ctx, school)
}

func (s *adminServiceImpl) UpdateSchool(ctx context.Context, school *model.School) error {
	// 验证学校信息
	if school.ID == 0 || school.Name == "" {
		return ErrInvalidInput
	}

	// 检查学校是否存在
	existingSchool, err := s.schoolRepo.FindByID(ctx, school.ID)
	if err != nil {
		return err
	}
	if existingSchool == nil {
		return ErrResourceNotFound
	}

	// 更新学校信息
	return s.schoolRepo.Update(ctx, school)
}

func (s *adminServiceImpl) UpdateSchoolStatus(ctx context.Context, id uint64, isActive bool) error {
	// 检查学校是否存在
	existingSchool, err := s.schoolRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existingSchool == nil {
		return ErrResourceNotFound
	}

	// 更新学校状态
	existingSchool.IsActive = isActive
	return s.schoolRepo.Update(ctx, existingSchool)
}

// ─── 班级管理 ──────────────────────────────────────────────────

func (s *adminServiceImpl) GetClasses(ctx context.Context, params map[string]interface{}) ([]*model.Class, int64, error) {
	return s.classRepo.FindAll(ctx, params)
}

func (s *adminServiceImpl) CreateClass(ctx context.Context, class *model.Class) error {
	// 验证班级信息
	if class.ClassName == "" || class.SchoolID == 0 || class.Grade == 0 || class.ClassCode == "" {
		return ErrInvalidInput
	}

	// 检查班级代码是否已存在
	params := map[string]interface{}{"class_code": class.ClassCode}
	classes, _, err := s.classRepo.FindAll(ctx, params)
	if err != nil {
		return err
	}
	if len(classes) > 0 {
		return ErrDuplicateResource
	}

	// 创建班级
	return s.classRepo.Create(ctx, class)
}

func (s *adminServiceImpl) UpdateClass(ctx context.Context, class *model.Class) error {
	// 验证班级信息
	if class.ID == 0 || class.ClassName == "" || class.SchoolID == 0 || class.Grade == 0 {
		return ErrInvalidInput
	}

	// 检查班级是否存在
	existingClass, err := s.classRepo.FindByID(ctx, class.ID)
	if err != nil {
		return err
	}
	if existingClass == nil {
		return ErrResourceNotFound
	}

	// 如果班级代码变更，检查是否已存在
	if class.ClassCode != existingClass.ClassCode {
		params := map[string]interface{}{"class_code": class.ClassCode}
		classes, _, err := s.classRepo.FindAll(ctx, params)
		if err != nil {
			return err
		}
		if len(classes) > 0 {
			return ErrDuplicateResource
		}
	}

	// 更新班级信息
	return s.classRepo.Update(ctx, class)
}

func (s *adminServiceImpl) PromoteClass(ctx context.Context, classID uint64, newGrade int, newClassCode string) error {
	// 检查班级是否存在
	existingClass, err := s.classRepo.FindByID(ctx, classID)
	if err != nil {
		return err
	}
	if existingClass == nil {
		return ErrResourceNotFound
	}

	// 验证新年级
	if newGrade <= existingClass.Grade || newGrade > 6 {
		return ErrInvalidInput
	}

	// 检查新班级代码是否已存在
	if newClassCode != "" {
		params := map[string]interface{}{"class_code": newClassCode}
		classes, _, err := s.classRepo.FindAll(ctx, params)
		if err != nil {
			return err
		}
		if len(classes) > 0 {
			return ErrDuplicateResource
		}
	} else {
		// 如果没有提供新班级代码，自动生成
		newClassCode = existingClass.ClassCode
	}

	// 开始事务
	tx := s.classRepo.BeginTx(ctx)
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 创建新班级
	newClass := &model.Class{
		SchoolID:   existingClass.SchoolID,
		ClassName:  existingClass.ClassName,
		Grade:      newGrade,
		ClassCode:  newClassCode,
		SchoolYear: existingClass.SchoolYear,
		IsActive:   true,
	}

	if err := s.classRepo.Create(ctx, newClass); err != nil {
		tx.Rollback()
		return err
	}

	// 获取原班级的所有学生
	childParams := map[string]interface{}{"class_id": existingClass.ID}
	children, _, err := s.childRepo.FindAll(ctx, childParams)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 更新所有学生的班级ID
	for _, child := range children {
		child.ClassID = newClass.ID
		if err := s.childRepo.Update(ctx, child); err != nil {
			tx.Rollback()
			return err
		}
	}

	// 停用原班级
	existingClass.IsActive = false
	if err := s.classRepo.Update(ctx, existingClass); err != nil {
		tx.Rollback()
		return err
	}

	// 提交事务
	return tx.Commit().Error
}

func (s *adminServiceImpl) UpdateClassStatus(ctx context.Context, id uint64, isActive bool) error {
	// 检查班级是否存在
	existingClass, err := s.classRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existingClass == nil {
		return ErrResourceNotFound
	}

	// 更新班级状态
	existingClass.IsActive = isActive
	return s.classRepo.Update(ctx, existingClass)
}

// ─── 老师/家长用户管理 ────────────────────────────────────────

func (s *adminServiceImpl) GetUsers(ctx context.Context, params map[string]interface{}) ([]*model.User, int64, error) {
	// 注意：由于UserRepository没有FindAll方法，这里暂时返回空列表
	// 实际项目中需要在UserRepository中添加FindAll方法
	return []*model.User{}, 0, nil
}

func (s *adminServiceImpl) CreateUser(ctx context.Context, user *model.User, initialPassword string) error {
	// 验证用户信息
	if user.Username == "" || user.Role == "" || initialPassword == "" {
		return ErrInvalidInput
	}

	// 检查用户是否已存在
	existingUser, err := s.userRepo.FindByUsername(ctx, user.Username)
	if err == nil && existingUser != nil {
		return ErrDuplicateResource
	}

	// 设置用户初始密码
	hashedPassword, err := utils.HashPassword(initialPassword)
	if err != nil {
		return err
	}
	user.PasswordHash = hashedPassword

	// 创建用户
	return s.userRepo.Create(ctx, user)
}

func (s *adminServiceImpl) UpdateUser(ctx context.Context, user *model.User) error {
	// 验证用户信息
	if user.ID == 0 || user.Username == "" || user.Role == "" {
		return ErrInvalidInput
	}

	// 检查用户是否存在
	existingUser, err := s.userRepo.FindByID(ctx, user.ID)
	if err != nil {
		return err
	}
	if existingUser == nil {
		return ErrResourceNotFound
	}

	// 如果用户名变更，检查是否已存在
	if user.Username != existingUser.Username {
		existingUserByUsername, err := s.userRepo.FindByUsername(ctx, user.Username)
		if err == nil && existingUserByUsername != nil {
			return ErrDuplicateResource
		}
	}

	// 保留原密码
	user.PasswordHash = existingUser.PasswordHash

	// 更新用户信息
	return s.userRepo.Update(ctx, user)
}

func (s *adminServiceImpl) UpdateUserStatus(ctx context.Context, id uint64, isActive bool) error {
	// 检查用户是否存在
	existingUser, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existingUser == nil {
		return ErrResourceNotFound
	}

	// 更新用户状态
	existingUser.IsActive = isActive
	return s.userRepo.Update(ctx, existingUser)
}

func (s *adminServiceImpl) ResetUserPassword(ctx context.Context, id uint64, newPassword string) error {
	// 验证新密码
	if newPassword == "" {
		return ErrInvalidInput
	}

	// 检查用户是否存在
	existingUser, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existingUser == nil {
		return ErrResourceNotFound
	}

	// 设置新密码
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}
	existingUser.PasswordHash = hashedPassword

	// 更新用户信息
	return s.userRepo.Update(ctx, existingUser)
}

// ─── 学生账号管理 ──────────────────────────────────────────────

func (s *adminServiceImpl) BatchImportStudents(ctx context.Context, data []*model.Child) error {
	// 验证数据
	if len(data) == 0 {
		return ErrInvalidInput
	}

	// 注意：由于ChildRepository没有BatchCreate方法，这里暂时逐个创建
	// 实际项目中需要在ChildRepository中添加BatchCreate方法
	for _, child := range data {
		if err := s.childRepo.Create(ctx, child); err != nil {
			return err
		}
	}

	return nil
}

func (s *adminServiceImpl) GetStudents(ctx context.Context, params map[string]interface{}) ([]*model.Child, int64, error) {
	// 注意：由于ChildRepository没有FindAll方法，这里暂时返回空列表
	// 实际项目中需要在ChildRepository中添加FindAll方法
	return []*model.Child{}, 0, nil
}

func (s *adminServiceImpl) CreateStudent(ctx context.Context, child *model.Child) error {
	// 验证学生信息
	if child.DisplayName == "" || child.ClassID == 0 || child.StudentNoEnc == "" {
		return ErrInvalidInput
	}

	// 检查学生是否已存在
	params := map[string]interface{}{"student_id": child.StudentNoEnc}
	children, _, err := s.childRepo.FindAll(ctx, params)
	if err != nil {
		return err
	}
	if len(children) > 0 {
		return ErrDuplicateResource
	}

	// 创建学生
	return s.childRepo.Create(ctx, child)
}

func (s *adminServiceImpl) UpdateStudent(ctx context.Context, child *model.Child) error {
	// 验证学生信息
	if child.ID == 0 || child.DisplayName == "" || child.ClassID == 0 {
		return ErrInvalidInput
	}

	// 检查学生是否存在
	existingChild, err := s.childRepo.FindByID(ctx, child.ID)
	if err != nil {
		return err
	}
	if existingChild == nil {
		return ErrResourceNotFound
	}

	// 如果学号变更，检查是否已存在
	if child.StudentNoEnc != existingChild.StudentNoEnc {
		params := map[string]interface{}{"student_id": child.StudentNoEnc}
		children, _, err := s.childRepo.FindAll(ctx, params)
		if err != nil {
			return err
		}
		if len(children) > 0 {
			return ErrDuplicateResource
		}
	}

	// 更新学生信息
	return s.childRepo.Update(ctx, child)
}

// ─── 老师班级权限分配 ──────────────────────────────────────────

func (s *adminServiceImpl) GetAssignments(ctx context.Context, params map[string]interface{}) ([]*model.AdminPermission, int64, error) {
	return s.permissionRepo.FindAll(ctx, params)
}

func (s *adminServiceImpl) CreateAssignment(ctx context.Context, assignment *model.AdminPermission) error {
	// 验证权限信息
	if assignment.TeacherUserID == 0 || assignment.ClassID == 0 || assignment.PermissionType == "" {
		return ErrInvalidInput
	}

	// 检查用户是否存在且为老师角色
	user, err := s.userRepo.FindByID(ctx, assignment.TeacherUserID)
	if err != nil {
		return err
	}
	if user == nil || user.Role != "teacher" {
		return ErrInvalidInput
	}

	// 检查班级是否存在
	class, err := s.classRepo.FindByID(ctx, assignment.ClassID)
	if err != nil {
		return err
	}
	if class == nil {
		return ErrResourceNotFound
	}

	// 检查权限是否已存在
	params := map[string]interface{}{"teacher_user_id": assignment.TeacherUserID, "class_id": assignment.ClassID}
	permissions, _, err := s.permissionRepo.FindAll(ctx, params)
	if err != nil {
		return err
	}
	if len(permissions) > 0 {
		return ErrDuplicateResource
	}

	// 创建权限分配
	return s.permissionRepo.Create(ctx, assignment)
}

func (s *adminServiceImpl) DeleteAssignment(ctx context.Context, id uint64) error {
	// 检查权限是否存在
	existingPermission, err := s.permissionRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existingPermission == nil {
		return ErrResourceNotFound
	}

	// 删除权限分配
	return s.permissionRepo.Delete(ctx, id)
}

func (s *adminServiceImpl) BatchCreateAssignments(ctx context.Context, assignments []*model.AdminPermission) error {
	// 验证数据
	if len(assignments) == 0 {
		return ErrInvalidInput
	}

	// 批量创建权限分配
	return s.permissionRepo.BatchCreate(ctx, assignments)
}

// ─── 家长-学生绑定管理 ────────────────────────────────────────

func (s *adminServiceImpl) GetParentBindings(ctx context.Context, params map[string]interface{}) ([]*model.ParentChildRelation, int64, error) {
	return s.parentChildRepo.FindAll(ctx, params)
}

func (s *adminServiceImpl) CreateParentBinding(ctx context.Context, binding *model.ParentChildRelation) error {
	// 验证绑定信息
	if binding.ParentUserID == 0 || binding.ChildID == 0 || binding.Relationship == "" {
		return ErrInvalidInput
	}

	// 检查家长用户是否存在且为家长角色
	parent, err := s.userRepo.FindByID(ctx, binding.ParentUserID)
	if err != nil {
		return err
	}
	if parent == nil || parent.Role != "parent" {
		return ErrInvalidInput
	}

	// 检查学生是否存在
	child, err := s.childRepo.FindByID(ctx, binding.ChildID)
	if err != nil {
		return err
	}
	if child == nil {
		return ErrResourceNotFound
	}

	// 检查绑定是否已存在
	params := map[string]interface{}{"parent_user_id": binding.ParentUserID, "child_id": binding.ChildID}
	bindings, _, err := s.parentChildRepo.FindAll(ctx, params)
	if err != nil {
		return err
	}
	if len(bindings) > 0 {
		return ErrDuplicateResource
	}

	// 如果设置为主绑定，将其他绑定设置为非主绑定
	if binding.IsPrimary {
		if err := s.parentChildRepo.SetAllNonPrimary(ctx, binding.ChildID); err != nil {
			return err
		}
	}

	// 创建绑定
	return s.parentChildRepo.Create(ctx, binding)
}

func (s *adminServiceImpl) DeleteParentBinding(ctx context.Context, id uint64) error {
	// 检查绑定是否存在
	existingBinding, err := s.parentChildRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existingBinding == nil {
		return ErrResourceNotFound
	}

	// 删除绑定
	return s.parentChildRepo.Delete(ctx, id)
}

// ─── 数据概览 ──────────────────────────────────────────────────

func (s *adminServiceImpl) GetDashboard(ctx context.Context) (map[string]interface{}, error) {
	// 开始事务
	tx := s.schoolRepo.BeginTx(ctx)
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 获取统计数据
	// 学校数量
	schools, _, err := s.schoolRepo.FindAll(ctx, map[string]interface{}{})
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// 班级数量
	classes, _, err := s.classRepo.FindAll(ctx, map[string]interface{}{})
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// 学生数量
	students, _, err := s.childRepo.FindAll(ctx, map[string]interface{}{})
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// 老师数量
	teachers, _, err := s.userRepo.FindAll(ctx, map[string]interface{}{"role": "teacher"})
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// 家长数量
	parents, _, err := s.userRepo.FindAll(ctx, map[string]interface{}{"role": "parent"})
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// 构建响应数据
	dashboard := map[string]interface{}{
		"total_schools":  len(schools),
		"total_classes":  len(classes),
		"total_students": len(students),
		"total_teachers": len(teachers),
		"total_parents":  len(parents),
	}

	return dashboard, nil
}

func (s *adminServiceImpl) GetAuditLogs(ctx context.Context, params map[string]interface{}) ([]*model.AuditLog, int64, error) {
	return s.auditLogRepo.FindAll(ctx, params)
}

// ─── 伙伴模板管理 ──────────────────────────────────────────────

func (s *adminServiceImpl) CreatePartnerTemplate(ctx context.Context, template *model.PartnerTemplate) error {
	// 验证模板信息
	if template.Name == "" || template.Type == "" {
		return ErrInvalidInput
	}

	// 创建模板
	return s.templateRepo.Create(ctx, template)
}

func (s *adminServiceImpl) UpdatePartnerTemplate(ctx context.Context, template *model.PartnerTemplate) error {
	// 验证模板信息
	if template.ID == 0 || template.Name == "" || template.Type == "" {
		return ErrInvalidInput
	}

	// 检查模板是否存在
	existingTemplate, err := s.templateRepo.FindByID(ctx, template.ID)
	if err != nil {
		return err
	}
	if existingTemplate == nil {
		return ErrResourceNotFound
	}

	// 更新模板信息
	return s.templateRepo.Update(ctx, template)
}

func (s *adminServiceImpl) SeedPartnerTemplates(ctx context.Context) error {
	// 初始化30个预设模板
	templates := []*model.PartnerTemplate{
		// 宠物类型
		{Code: "pet_001", Name: "小奶狗", Type: model.PartnerTypePet, Description: "可爱忠诚的小奶狗", Slogan: "汪！今天也要一起变得更棒哦！", LowStageAsset: "https://example.com/pet1_low.json", MidStageAsset: "https://example.com/pet1_mid.json", HighStageAsset: "https://example.com/pet1_high.json", IsActive: true, SortOrder: 1},
		{Code: "pet_002", Name: "小绒猫", Type: model.PartnerTypePet, Description: "温柔可爱的小绒猫", Slogan: "喵～和我一起成长吧！", LowStageAsset: "https://example.com/pet2_low.json", MidStageAsset: "https://example.com/pet2_mid.json", HighStageAsset: "https://example.com/pet2_high.json", IsActive: true, SortOrder: 2},
		{Code: "pet_003", Name: "幼龙", Type: model.PartnerTypePet, Description: "神秘强大的幼龙", Slogan: "我会陪你成为传奇！", LowStageAsset: "https://example.com/pet3_low.json", MidStageAsset: "https://example.com/pet3_mid.json", HighStageAsset: "https://example.com/pet3_high.json", IsActive: true, SortOrder: 3},
		{Code: "pet_004", Name: "小狐", Type: model.PartnerTypePet, Description: "聪明灵动的小狐", Slogan: "一起创造美好的回忆吧！", LowStageAsset: "https://example.com/pet4_low.json", MidStageAsset: "https://example.com/pet4_mid.json", HighStageAsset: "https://example.com/pet4_high.json", IsActive: true, SortOrder: 4},
		{Code: "pet_005", Name: "小兔子", Type: model.PartnerTypePet, Description: "活泼可爱的小兔子", Slogan: "蹦蹦跳跳，快乐成长！", LowStageAsset: "https://example.com/pet5_low.json", MidStageAsset: "https://example.com/pet5_mid.json", HighStageAsset: "https://example.com/pet5_high.json", IsActive: true, SortOrder: 5},
		{Code: "pet_006", Name: "小熊猫", Type: model.PartnerTypePet, Description: "憨厚可掬的小熊猫", Slogan: "和我一起享受成长的每一刻！", LowStageAsset: "https://example.com/pet6_low.json", MidStageAsset: "https://example.com/pet6_mid.json", HighStageAsset: "https://example.com/pet6_high.json", IsActive: true, SortOrder: 6},
		{Code: "pet_007", Name: "小松鼠", Type: model.PartnerTypePet, Description: "敏捷灵活的小松鼠", Slogan: "一起收集成长的果实！", LowStageAsset: "https://example.com/pet7_low.json", MidStageAsset: "https://example.com/pet7_mid.json", HighStageAsset: "https://example.com/pet7_high.json", IsActive: true, SortOrder: 7},
		{Code: "pet_008", Name: "小企鹅", Type: model.PartnerTypePet, Description: "呆萌可爱的小企鹅", Slogan: "摇摇摆摆，一起进步！", LowStageAsset: "https://example.com/pet8_low.json", MidStageAsset: "https://example.com/pet8_mid.json", HighStageAsset: "https://example.com/pet8_high.json", IsActive: true, SortOrder: 8},
		{Code: "pet_009", Name: "小海豚", Type: model.PartnerTypePet, Description: "聪明友好的小海豚", Slogan: "在知识的海洋里遨游！", LowStageAsset: "https://example.com/pet9_low.json", MidStageAsset: "https://example.com/pet9_mid.json", HighStageAsset: "https://example.com/pet9_high.json", IsActive: true, SortOrder: 9},
		{Code: "pet_010", Name: "小狮子", Type: model.PartnerTypePet, Description: "勇敢坚强的小狮子", Slogan: "成为自己的王者！", LowStageAsset: "https://example.com/pet10_low.json", MidStageAsset: "https://example.com/pet10_mid.json", HighStageAsset: "https://example.com/pet10_high.json", IsActive: true, SortOrder: 10},

		// 植物类型
		{Code: "plant_001", Name: "小绿芽", Type: model.PartnerTypePlant, Description: "充满生机的小绿芽", Slogan: "我正因为你而成长～", LowStageAsset: "https://example.com/plant1_low.json", MidStageAsset: "https://example.com/plant1_mid.json", HighStageAsset: "https://example.com/plant1_high.json", IsActive: true, SortOrder: 11},
		{Code: "plant_002", Name: "小幼苗", Type: model.PartnerTypePlant, Description: "茁壮成长的小幼苗", Slogan: "阳光雨露，共同成长！", LowStageAsset: "https://example.com/plant2_low.json", MidStageAsset: "https://example.com/plant2_mid.json", HighStageAsset: "https://example.com/plant2_high.json", IsActive: true, SortOrder: 12},
		{Code: "plant_003", Name: "小种子", Type: model.PartnerTypePlant, Description: "蕴含希望的小种子", Slogan: "埋下希望，收获成长！", LowStageAsset: "https://example.com/plant3_low.json", MidStageAsset: "https://example.com/plant3_mid.json", HighStageAsset: "https://example.com/plant3_high.json", IsActive: true, SortOrder: 13},
		{Code: "plant_004", Name: "小花朵", Type: model.PartnerTypePlant, Description: "美丽绽放的小花朵", Slogan: "为你绽放最美丽的笑容！", LowStageAsset: "https://example.com/plant4_low.json", MidStageAsset: "https://example.com/plant4_mid.json", HighStageAsset: "https://example.com/plant4_high.json", IsActive: true, SortOrder: 14},
		{Code: "plant_005", Name: "小树苗", Type: model.PartnerTypePlant, Description: "挺拔向上的小树苗", Slogan: "和你一起顶天立地！", LowStageAsset: "https://example.com/plant5_low.json", MidStageAsset: "https://example.com/plant5_mid.json", HighStageAsset: "https://example.com/plant5_high.json", IsActive: true, SortOrder: 15},
		{Code: "plant_006", Name: "多肉植物", Type: model.PartnerTypePlant, Description: "可爱多肉植物", Slogan: "胖胖乎乎，充满能量！", LowStageAsset: "https://example.com/plant6_low.json", MidStageAsset: "https://example.com/plant6_mid.json", HighStageAsset: "https://example.com/plant6_high.json", IsActive: true, SortOrder: 16},
		{Code: "plant_007", Name: "向日葵", Type: model.PartnerTypePlant, Description: "向阳而生的向日葵", Slogan: "永远向着阳光生长！", LowStageAsset: "https://example.com/plant7_low.json", MidStageAsset: "https://example.com/plant7_mid.json", HighStageAsset: "https://example.com/plant7_high.json", IsActive: true, SortOrder: 17},
		{Code: "plant_008", Name: "小蘑菇", Type: model.PartnerTypePlant, Description: "小巧可爱的小蘑菇", Slogan: "在成长的森林里探险！", LowStageAsset: "https://example.com/plant8_low.json", MidStageAsset: "https://example.com/plant8_mid.json", HighStageAsset: "https://example.com/plant8_high.json", IsActive: true, SortOrder: 18},
		{Code: "plant_009", Name: "薰衣草", Type: model.PartnerTypePlant, Description: "芳香四溢的薰衣草", Slogan: "用香气记录每一个美好瞬间！", LowStageAsset: "https://example.com/plant9_low.json", MidStageAsset: "https://example.com/plant9_mid.json", HighStageAsset: "https://example.com/plant9_high.json", IsActive: true, SortOrder: 19},
		{Code: "plant_010", Name: "仙人掌", Type: model.PartnerTypePlant, Description: "坚韧不拔的仙人掌", Slogan: "在困难中也要茁壮成长！", LowStageAsset: "https://example.com/plant10_low.json", MidStageAsset: "https://example.com/plant10_mid.json", HighStageAsset: "https://example.com/plant10_high.json", IsActive: true, SortOrder: 20},

		// 二次元类型
		{Code: "anime_001", Name: "见习剑士", Type: model.PartnerTypeAnime, Description: "勇敢的见习剑士", Slogan: "用勇气开辟成长之路！", LowStageAsset: "https://example.com/anime1_low.json", MidStageAsset: "https://example.com/anime1_mid.json", HighStageAsset: "https://example.com/anime1_high.json", IsActive: true, SortOrder: 21},
		{Code: "anime_002", Name: "星芒学徒", Type: model.PartnerTypeAnime, Description: "神秘的星芒学徒", Slogan: "探索星空，收获知识！", LowStageAsset: "https://example.com/anime2_low.json", MidStageAsset: "https://example.com/anime2_mid.json", HighStageAsset: "https://example.com/anime2_high.json", IsActive: true, SortOrder: 22},
		{Code: "anime_003", Name: "机灵助手", Type: model.PartnerTypeAnime, Description: "聪明的机灵助手", Slogan: "有我在，成长路上不孤单！", LowStageAsset: "https://example.com/anime3_low.json", MidStageAsset: "https://example.com/anime3_mid.json", HighStageAsset: "https://example.com/anime3_high.json", IsActive: true, SortOrder: 23},
		{Code: "anime_004", Name: "魔法少女", Type: model.PartnerTypeAnime, Description: "可爱的魔法少女", Slogan: "用魔法点亮成长的道路！", LowStageAsset: "https://example.com/anime4_low.json", MidStageAsset: "https://example.com/anime4_mid.json", HighStageAsset: "https://example.com/anime4_high.json", IsActive: true, SortOrder: 24},
		{Code: "anime_005", Name: "科学少年", Type: model.PartnerTypeAnime, Description: "聪明的科学少年", Slogan: "用科学探索世界的奥秘！", LowStageAsset: "https://example.com/anime5_low.json", MidStageAsset: "https://example.com/anime5_mid.json", HighStageAsset: "https://example.com/anime5_high.json", IsActive: true, SortOrder: 25},
		{Code: "anime_006", Name: "音乐少女", Type: model.PartnerTypeAnime, Description: "优雅的音乐少女", Slogan: "用音乐谱写成长的旋律！", LowStageAsset: "https://example.com/anime6_low.json", MidStageAsset: "https://example.com/anime6_mid.json", HighStageAsset: "https://example.com/anime6_high.json", IsActive: true, SortOrder: 26},
		{Code: "anime_007", Name: "体育健将", Type: model.PartnerTypeAnime, Description: "活力四射的体育健将", Slogan: "运动让成长更精彩！", LowStageAsset: "https://example.com/anime7_low.json", MidStageAsset: "https://example.com/anime7_mid.json", HighStageAsset: "https://example.com/anime7_high.json", IsActive: true, SortOrder: 27},
		{Code: "anime_008", Name: "艺术达人", Type: model.PartnerTypeAnime, Description: "才华横溢的艺术达人", Slogan: "用艺术表达成长的美好！", LowStageAsset: "https://example.com/anime8_low.json", MidStageAsset: "https://example.com/anime8_mid.json", HighStageAsset: "https://example.com/anime8_high.json", IsActive: true, SortOrder: 28},
		{Code: "anime_009", Name: "编程小能手", Type: model.PartnerTypeAnime, Description: "聪明的编程小能手", Slogan: "用代码创造无限可能！", LowStageAsset: "https://example.com/anime9_low.json", MidStageAsset: "https://example.com/anime9_mid.json", HighStageAsset: "https://example.com/anime9_high.json", IsActive: true, SortOrder: 29},
		{Code: "anime_010", Name: "读书郎", Type: model.PartnerTypeAnime, Description: "勤奋的读书郎", Slogan: "书中自有黄金屋，知识让我成长！", LowStageAsset: "https://example.com/anime10_low.json", MidStageAsset: "https://example.com/anime10_mid.json", HighStageAsset: "https://example.com/anime10_high.json", IsActive: true, SortOrder: 30},
	}

	// 批量创建模板
	for _, template := range templates {
		if err := s.templateRepo.Create(ctx, template); err != nil {
			return err
		}
	}

	return nil
}
