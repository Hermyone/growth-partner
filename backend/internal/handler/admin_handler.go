// growth-partner/backend/internal/handler/admin_handler.go
// 管理员模块控制器

package handler

import (
	"fmt"
	"strconv"

	"growth-partner/internal/middleware"
	"growth-partner/internal/model"
	"growth-partner/internal/service"

	"github.com/gin-gonic/gin"
)

// 辅助函数：将字符串转换为uint64
func parseUint64(s string) (uint64, error) {
	id, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid uint64: %w", err)
	}
	return id, nil
}

// 辅助函数：将字符串转换为int
func parseInt(s string) (int, error) {
	id, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("invalid int: %w", err)
	}
	return id, nil
}

// AdminHandler 管理员控制器
type AdminHandler struct {
	adminSvc service.AdminService
}

// NewAdminHandler 创建管理员控制器实例
func NewAdminHandler(adminSvc service.AdminService) *AdminHandler {
	return &AdminHandler{
		adminSvc: adminSvc,
	}
}

// ─── 学校管理 ──────────────────────────────────────────────────

// SchoolCreateReq 创建学校请求参数
type SchoolCreateReq struct {
	Name     string `json:"name" binding:"required"`
	Address  string `json:"address"`
	Contact  string `json:"contact"`
	IsActive bool   `json:"is_active"`
}

// SchoolUpdateReq 更新学校请求参数
type SchoolUpdateReq struct {
	Name     string `json:"name" binding:"required"`
	Address  string `json:"address"`
	Contact  string `json:"contact"`
	IsActive bool   `json:"is_active"`
}

// SchoolStatusReq 启用/停用学校请求参数
type SchoolStatusReq struct {
	IsActive bool `json:"is_active" binding:"required"`
}

// GetSchools 获取学校列表（分页+搜索）
// GET /api/v1/admin/schools
func (h *AdminHandler) GetSchools(c *gin.Context) {
	// 构建查询参数
	params := make(map[string]interface{})
	if name := c.Query("name"); name != "" {
		params["name"] = name
	}
	if isActive := c.Query("is_active"); isActive != "" {
		params["is_active"] = isActive == "true"
	}

	// 调用服务获取学校列表
	schools, total, err := h.adminSvc.GetSchools(c.Request.Context(), params)
	if err != nil {
		middleware.ResponseError(c, 500, "GET_SCHOOLS_FAILED", "获取学校列表失败")
		return
	}

	middleware.ResponseOK(c, gin.H{
		"list":  schools,
		"total": total,
	})
}

// CreateSchool 创建学校
// POST /api/v1/admin/schools
func (h *AdminHandler) CreateSchool(c *gin.Context) {
	var req SchoolCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResponseValidationError(c, err.Error())
		return
	}

	// 构建学校对象
	school := &model.School{
		Name:         req.Name,
		Address:      req.Address,
		ContactPhone: req.Contact,
		IsActive:     req.IsActive,
	}

	// 调用服务创建学校
	if err := h.adminSvc.CreateSchool(c.Request.Context(), school); err != nil {
		if err == service.ErrDuplicateResource {
			middleware.ResponseError(c, 400, "SCHOOL_EXISTS", "学校已存在")
			return
		}
		middleware.ResponseError(c, 500, "CREATE_SCHOOL_FAILED", "创建学校失败")
		return
	}

	middleware.ResponseOKWithMessage(c, "创建学校成功", school)
}

// UpdateSchool 更新学校信息
// PUT /api/v1/admin/schools/:id
func (h *AdminHandler) UpdateSchool(c *gin.Context) {
	var req SchoolUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResponseValidationError(c, err.Error())
		return
	}

	// 获取学校ID
	id, err := parseUint64(c.Param("id"))
	if err != nil {
		middleware.ResponseValidationError(c, "无效的学校ID")
		return
	}

	// 构建学校对象
	school := &model.School{
		Base: model.Base{
			ID: id,
		},
		Name:         req.Name,
		Address:      req.Address,
		ContactPhone: req.Contact,
		IsActive:     req.IsActive,
	}

	// 调用服务更新学校
	if err := h.adminSvc.UpdateSchool(c.Request.Context(), school); err != nil {
		if err == service.ErrResourceNotFound {
			middleware.ResponseError(c, 404, "SCHOOL_NOT_FOUND", "学校不存在")
			return
		}
		middleware.ResponseError(c, 500, "UPDATE_SCHOOL_FAILED", "更新学校失败")
		return
	}

	middleware.ResponseOKWithMessage(c, "更新学校成功", school)
}

// UpdateSchoolStatus 启用/停用学校
// PATCH /api/v1/admin/schools/:id/status
func (h *AdminHandler) UpdateSchoolStatus(c *gin.Context) {
	var req SchoolStatusReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResponseValidationError(c, err.Error())
		return
	}

	// 获取学校ID
	id, err := parseUint64(c.Param("id"))
	if err != nil {
		middleware.ResponseValidationError(c, "无效的学校ID")
		return
	}

	// 调用服务更新学校状态
	if err := h.adminSvc.UpdateSchoolStatus(c.Request.Context(), id, req.IsActive); err != nil {
		if err == service.ErrResourceNotFound {
			middleware.ResponseError(c, 404, "SCHOOL_NOT_FOUND", "学校不存在")
			return
		}
		middleware.ResponseError(c, 500, "UPDATE_SCHOOL_STATUS_FAILED", "更新学校状态失败")
		return
	}

	middleware.ResponseOKWithMessage(c, "更新学校状态成功", nil)
}

// ─── 班级管理 ──────────────────────────────────────────────────

// ClassCreateReq 创建班级请求参数
type ClassCreateReq struct {
	SchoolID          uint64 `json:"school_id" binding:"required"`
	ClassName         string `json:"class_name" binding:"required"`
	ClassCode         string `json:"class_code" binding:"required"`
	Grade             int    `json:"grade" binding:"required"`
	ClassNo           int    `json:"class_no" binding:"required"`
	SchoolYear        string `json:"school_year" binding:"required"`
	HomeroomTeacherID uint64 `json:"homeroom_teacher_id"`
	IsActive          bool   `json:"is_active"`
}

// ClassUpdateReq 更新班级请求参数
type ClassUpdateReq struct {
	SchoolID          uint64 `json:"school_id" binding:"required"`
	ClassName         string `json:"class_name" binding:"required"`
	ClassCode         string `json:"class_code" binding:"required"`
	Grade             int    `json:"grade" binding:"required"`
	ClassNo           int    `json:"class_no" binding:"required"`
	SchoolYear        string `json:"school_year" binding:"required"`
	HomeroomTeacherID uint64 `json:"homeroom_teacher_id"`
	IsActive          bool   `json:"is_active"`
}

// ClassPromoteReq 升年级请求参数
type ClassPromoteReq struct {
	NewGrade     int    `json:"new_grade" binding:"required"`
	NewClassCode string `json:"new_class_code"`
}

// ClassStatusReq 启用/停用班级请求参数
type ClassStatusReq struct {
	IsActive bool `json:"is_active" binding:"required"`
}

// GetClasses 获取班级列表（按学校/学年/年级筛选）
// GET /api/v1/admin/classes
func (h *AdminHandler) GetClasses(c *gin.Context) {
	// 构建查询参数
	params := make(map[string]interface{})
	if schoolID := c.Query("school_id"); schoolID != "" {
		if id, err := parseUint64(schoolID); err == nil {
			params["school_id"] = id
		}
	}
	if schoolYear := c.Query("school_year"); schoolYear != "" {
		params["school_year"] = schoolYear
	}
	if grade := c.Query("grade"); grade != "" {
		if g, err := parseInt(grade); err == nil {
			params["grade"] = g
		}
	}
	if className := c.Query("class_name"); className != "" {
		params["class_name"] = className
	}
	if isActive := c.Query("is_active"); isActive != "" {
		params["is_active"] = isActive == "true"
	}

	// 调用服务获取班级列表
	classes, total, err := h.adminSvc.GetClasses(c.Request.Context(), params)
	if err != nil {
		middleware.ResponseError(c, 500, "GET_CLASSES_FAILED", "获取班级列表失败")
		return
	}

	middleware.ResponseOK(c, gin.H{
		"list":  classes,
		"total": total,
	})
}

// CreateClass 创建新班级（class_code唯一校验）
// POST /api/v1/admin/classes
func (h *AdminHandler) CreateClass(c *gin.Context) {
	var req ClassCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResponseValidationError(c, err.Error())
		return
	}

	// 构建班级对象
	class := &model.Class{
		SchoolID:          req.SchoolID,
		ClassName:         req.ClassName,
		ClassCode:         req.ClassCode,
		Grade:             req.Grade,
		ClassNo:           req.ClassNo,
		SchoolYear:        req.SchoolYear,
		HomeroomTeacherID: req.HomeroomTeacherID,
		IsActive:          req.IsActive,
	}

	// 调用服务创建班级
	if err := h.adminSvc.CreateClass(c.Request.Context(), class); err != nil {
		if err == service.ErrDuplicateResource {
			middleware.ResponseError(c, 400, "CLASS_CODE_EXISTS", "班级代码已存在")
			return
		}
		middleware.ResponseError(c, 500, "CREATE_CLASS_FAILED", "创建班级失败")
		return
	}

	middleware.ResponseOKWithMessage(c, "创建班级成功", class)
}

// UpdateClass 更新班级信息（班级名/班主任）
// PUT /api/v1/admin/classes/:id
func (h *AdminHandler) UpdateClass(c *gin.Context) {
	var req ClassUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResponseValidationError(c, err.Error())
		return
	}

	// 获取班级ID
	id, err := parseUint64(c.Param("id"))
	if err != nil {
		middleware.ResponseValidationError(c, "无效的班级ID")
		return
	}

	// 构建班级对象
	class := &model.Class{
		Base: model.Base{
			ID: id,
		},
		SchoolID:          req.SchoolID,
		ClassName:         req.ClassName,
		ClassCode:         req.ClassCode,
		Grade:             req.Grade,
		ClassNo:           req.ClassNo,
		SchoolYear:        req.SchoolYear,
		HomeroomTeacherID: req.HomeroomTeacherID,
		IsActive:          req.IsActive,
	}

	// 调用服务更新班级
	if err := h.adminSvc.UpdateClass(c.Request.Context(), class); err != nil {
		if err == service.ErrResourceNotFound {
			middleware.ResponseError(c, 404, "CLASS_NOT_FOUND", "班级不存在")
			return
		}
		if err == service.ErrDuplicateResource {
			middleware.ResponseError(c, 400, "CLASS_CODE_EXISTS", "班级代码已存在")
			return
		}
		middleware.ResponseError(c, 500, "UPDATE_CLASS_FAILED", "更新班级失败")
		return
	}

	middleware.ResponseOKWithMessage(c, "更新班级成功", class)
}

// PromoteClass 升年级操作，批量新建学生班级关联
// POST /api/v1/admin/classes/:id/promote
func (h *AdminHandler) PromoteClass(c *gin.Context) {
	var req ClassPromoteReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResponseValidationError(c, err.Error())
		return
	}

	// 获取班级ID
	id, err := parseUint64(c.Param("id"))
	if err != nil {
		middleware.ResponseValidationError(c, "无效的班级ID")
		return
	}

	// 调用服务升年级
	if err := h.adminSvc.PromoteClass(c.Request.Context(), id, req.NewGrade, req.NewClassCode); err != nil {
		if err == service.ErrResourceNotFound {
			middleware.ResponseError(c, 404, "CLASS_NOT_FOUND", "班级不存在")
			return
		}
		if err == service.ErrInvalidInput {
			middleware.ResponseError(c, 400, "INVALID_INPUT", "无效的输入参数")
			return
		}
		if err == service.ErrDuplicateResource {
			middleware.ResponseError(c, 400, "CLASS_CODE_EXISTS", "班级代码已存在")
			return
		}
		middleware.ResponseError(c, 500, "PROMOTE_CLASS_FAILED", "升年级操作失败")
		return
	}

	middleware.ResponseOKWithMessage(c, "升年级操作成功", nil)
}

// UpdateClassStatus 启用/停用班级
// PATCH /api/v1/admin/classes/:id/status
func (h *AdminHandler) UpdateClassStatus(c *gin.Context) {
	var req ClassStatusReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResponseValidationError(c, err.Error())
		return
	}

	// 获取班级ID
	id, err := parseUint64(c.Param("id"))
	if err != nil {
		middleware.ResponseValidationError(c, "无效的班级ID")
		return
	}

	// 调用服务更新班级状态
	if err := h.adminSvc.UpdateClassStatus(c.Request.Context(), id, req.IsActive); err != nil {
		if err == service.ErrResourceNotFound {
			middleware.ResponseError(c, 404, "CLASS_NOT_FOUND", "班级不存在")
			return
		}
		middleware.ResponseError(c, 500, "UPDATE_CLASS_STATUS_FAILED", "更新班级状态失败")
		return
	}

	middleware.ResponseOKWithMessage(c, "更新班级状态成功", nil)
}

// ─── 老师/家长用户管理 ────────────────────────────────────────

// UserCreateReq 创建用户请求参数
type UserCreateReq struct {
	Username        string `json:"username" binding:"required"`
	Role            string `json:"role" binding:"required,oneof=teacher parent"`
	InitialPassword string `json:"initial_password" binding:"required,min=6,max=32"`
	AvatarURL       string `json:"avatar_url"`
}

// UserUpdateReq 更新用户请求参数
type UserUpdateReq struct {
	Username  string `json:"username" binding:"required"`
	Role      string `json:"role" binding:"required,oneof=teacher parent"`
	AvatarURL string `json:"avatar_url"`
	IsActive  bool   `json:"is_active"`
}

// UserStatusReq 启用/停用账号请求参数
type UserStatusReq struct {
	IsActive bool `json:"is_active" binding:"required"`
}

// UserResetPasswordReq 重置用户密码请求参数
type UserResetPasswordReq struct {
	NewPassword string `json:"new_password" binding:"required,min=6,max=32"`
}

// GetUsers 获取用户列表（按角色/学校筛选，分页）
// GET /api/v1/admin/users
func (h *AdminHandler) GetUsers(c *gin.Context) {
	// 构建查询参数
	params := make(map[string]interface{})
	if role := c.Query("role"); role != "" {
		params["role"] = role
	}
	if username := c.Query("username"); username != "" {
		params["username"] = username
	}
	if isActive := c.Query("is_active"); isActive != "" {
		params["is_active"] = isActive == "true"
	}

	// 调用服务获取用户列表
	users, total, err := h.adminSvc.GetUsers(c.Request.Context(), params)
	if err != nil {
		middleware.ResponseError(c, 500, "GET_USERS_FAILED", "获取用户列表失败")
		return
	}

	middleware.ResponseOK(c, gin.H{
		"list":  users,
		"total": total,
	})
}

// CreateUser 创建老师/家长账号，设置初始密码
// POST /api/v1/admin/users
func (h *AdminHandler) CreateUser(c *gin.Context) {
	var req UserCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResponseValidationError(c, err.Error())
		return
	}

	// 构建用户对象
	user := &model.User{
		Username:  req.Username,
		Role:      model.UserRole(req.Role),
		AvatarURL: req.AvatarURL,
		IsActive:  true,
	}

	// 调用服务创建用户
	if err := h.adminSvc.CreateUser(c.Request.Context(), user, req.InitialPassword); err != nil {
		if err == service.ErrDuplicateResource {
			middleware.ResponseError(c, 400, "USER_EXISTS", "用户已存在")
			return
		}
		middleware.ResponseError(c, 500, "CREATE_USER_FAILED", "创建用户失败")
		return
	}

	middleware.ResponseOKWithMessage(c, "创建用户成功", user)
}

// UpdateUser 更新用户信息
// PUT /api/v1/admin/users/:id
func (h *AdminHandler) UpdateUser(c *gin.Context) {
	var req UserUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResponseValidationError(c, err.Error())
		return
	}

	// 获取用户ID
	id, err := parseUint64(c.Param("id"))
	if err != nil {
		middleware.ResponseValidationError(c, "无效的用户ID")
		return
	}

	// 构建用户对象
	user := &model.User{
		Base: model.Base{
			ID: id,
		},
		Username:  req.Username,
		Role:      model.UserRole(req.Role),
		AvatarURL: req.AvatarURL,
		IsActive:  req.IsActive,
	}

	// 调用服务更新用户
	if err := h.adminSvc.UpdateUser(c.Request.Context(), user); err != nil {
		if err == service.ErrResourceNotFound {
			middleware.ResponseError(c, 404, "USER_NOT_FOUND", "用户不存在")
			return
		}
		if err == service.ErrDuplicateResource {
			middleware.ResponseError(c, 400, "USER_EXISTS", "用户已存在")
			return
		}
		middleware.ResponseError(c, 500, "UPDATE_USER_FAILED", "更新用户失败")
		return
	}

	middleware.ResponseOKWithMessage(c, "更新用户成功", user)
}

// UpdateUserStatus 启用/停用账号
// PATCH /api/v1/admin/users/:id/status
func (h *AdminHandler) UpdateUserStatus(c *gin.Context) {
	var req UserStatusReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResponseValidationError(c, err.Error())
		return
	}

	// 获取用户ID
	id, err := parseUint64(c.Param("id"))
	if err != nil {
		middleware.ResponseValidationError(c, "无效的用户ID")
		return
	}

	// 调用服务更新用户状态
	if err := h.adminSvc.UpdateUserStatus(c.Request.Context(), id, req.IsActive); err != nil {
		if err == service.ErrResourceNotFound {
			middleware.ResponseError(c, 404, "USER_NOT_FOUND", "用户不存在")
			return
		}
		middleware.ResponseError(c, 500, "UPDATE_USER_STATUS_FAILED", "更新用户状态失败")
		return
	}

	middleware.ResponseOKWithMessage(c, "更新用户状态成功", nil)
}

// ResetUserPassword 重置用户密码（无需旧密码）
// PATCH /api/v1/admin/users/:id/reset-pwd
func (h *AdminHandler) ResetUserPassword(c *gin.Context) {
	var req UserResetPasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResponseValidationError(c, err.Error())
		return
	}

	// 获取用户ID
	id, err := parseUint64(c.Param("id"))
	if err != nil {
		middleware.ResponseValidationError(c, "无效的用户ID")
		return
	}

	// 调用服务重置用户密码
	if err := h.adminSvc.ResetUserPassword(c.Request.Context(), id, req.NewPassword); err != nil {
		if err == service.ErrResourceNotFound {
			middleware.ResponseError(c, 404, "USER_NOT_FOUND", "用户不存在")
			return
		}
		middleware.ResponseError(c, 500, "RESET_PASSWORD_FAILED", "重置用户密码失败")
		return
	}

	middleware.ResponseOKWithMessage(c, "重置用户密码成功", nil)
}

// ─── 学生账号管理 ──────────────────────────────────────────────

// StudentCreateReq 创建学生请求参数
type StudentCreateReq struct {
	UserID       uint64 `json:"user_id" binding:"required"`
	ClassID      uint64 `json:"class_id" binding:"required"`
	DisplayName  string `json:"display_name" binding:"required"`
	RealNameEnc  string `json:"real_name_enc" binding:"required"`
	StudentNoEnc string `json:"student_no_enc" binding:"required"`
	Gender       string `json:"gender" binding:"required,oneof=M F"`
	BirthYear    int    `json:"birth_year" binding:"required"`
	EnrollYear   int    `json:"enroll_year" binding:"required"`
	CurrentGrade int    `json:"current_grade" binding:"required"`
	IsActive     bool   `json:"is_active"`
}

// StudentUpdateReq 更新学生请求参数
type StudentUpdateReq struct {
	ClassID      uint64 `json:"class_id" binding:"required"`
	DisplayName  string `json:"display_name" binding:"required"`
	RealNameEnc  string `json:"real_name_enc" binding:"required"`
	StudentNoEnc string `json:"student_no_enc" binding:"required"`
	Gender       string `json:"gender" binding:"required,oneof=M F"`
	BirthYear    int    `json:"birth_year" binding:"required"`
	EnrollYear   int    `json:"enroll_year" binding:"required"`
	CurrentGrade int    `json:"current_grade" binding:"required"`
	IsActive     bool   `json:"is_active"`
}

// BatchImportStudents 批量导入学生（CSV上传）
// POST /api/v1/admin/students/batch-import
func (h *AdminHandler) BatchImportStudents(c *gin.Context) {
	// 这里需要处理CSV文件上传
	// 暂时实现一个简单的JSON批量导入
	var students []*model.Child
	if err := c.ShouldBindJSON(&students); err != nil {
		middleware.ResponseValidationError(c, err.Error())
		return
	}

	// 调用服务批量导入学生
	if err := h.adminSvc.BatchImportStudents(c.Request.Context(), students); err != nil {
		middleware.ResponseError(c, 500, "BATCH_IMPORT_FAILED", "批量导入学生失败")
		return
	}

	middleware.ResponseOKWithMessage(c, "批量导入学生成功", nil)
}

// GetStudents 学生列表（按班级/学年筛选，脱敏）
// GET /api/v1/admin/students
func (h *AdminHandler) GetStudents(c *gin.Context) {
	// 构建查询参数
	params := make(map[string]interface{})
	if classID := c.Query("class_id"); classID != "" {
		if id, err := parseUint64(classID); err == nil {
			params["class_id"] = id
		}
	}
	if isActive := c.Query("is_active"); isActive != "" {
		params["is_active"] = isActive == "true"
	}

	// 调用服务获取学生列表
	students, total, err := h.adminSvc.GetStudents(c.Request.Context(), params)
	if err != nil {
		middleware.ResponseError(c, 500, "GET_STUDENTS_FAILED", "获取学生列表失败")
		return
	}

	middleware.ResponseOK(c, gin.H{
		"list":  students,
		"total": total,
	})
}

// CreateStudent 单个创建学生账号
// POST /api/v1/admin/students
func (h *AdminHandler) CreateStudent(c *gin.Context) {
	var req StudentCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResponseValidationError(c, err.Error())
		return
	}

	// 构建学生对象
	child := &model.Child{
		UserID:       req.UserID,
		ClassID:      req.ClassID,
		DisplayName:  req.DisplayName,
		RealNameEnc:  req.RealNameEnc,
		StudentNoEnc: req.StudentNoEnc,
		Gender:       req.Gender,
		BirthYear:    req.BirthYear,
		EnrollYear:   req.EnrollYear,
		CurrentGrade: req.CurrentGrade,
		IsActive:     req.IsActive,
	}

	// 调用服务创建学生
	if err := h.adminSvc.CreateStudent(c.Request.Context(), child); err != nil {
		if err == service.ErrDuplicateResource {
			middleware.ResponseError(c, 400, "STUDENT_EXISTS", "学生已存在")
			return
		}
		middleware.ResponseError(c, 500, "CREATE_STUDENT_FAILED", "创建学生失败")
		return
	}

	middleware.ResponseOKWithMessage(c, "创建学生成功", child)
}

// UpdateStudent 更新学生信息
// PUT /api/v1/admin/students/:id
func (h *AdminHandler) UpdateStudent(c *gin.Context) {
	var req StudentUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResponseValidationError(c, err.Error())
		return
	}

	// 获取学生ID
	id, err := parseUint64(c.Param("id"))
	if err != nil {
		middleware.ResponseValidationError(c, "无效的学生ID")
		return
	}

	// 构建学生对象
	child := &model.Child{
		Base: model.Base{
			ID: id,
		},
		ClassID:      req.ClassID,
		DisplayName:  req.DisplayName,
		RealNameEnc:  req.RealNameEnc,
		StudentNoEnc: req.StudentNoEnc,
		Gender:       req.Gender,
		BirthYear:    req.BirthYear,
		EnrollYear:   req.EnrollYear,
		CurrentGrade: req.CurrentGrade,
		IsActive:     req.IsActive,
	}

	// 调用服务更新学生
	if err := h.adminSvc.UpdateStudent(c.Request.Context(), child); err != nil {
		if err == service.ErrResourceNotFound {
			middleware.ResponseError(c, 404, "STUDENT_NOT_FOUND", "学生不存在")
			return
		}
		if err == service.ErrDuplicateResource {
			middleware.ResponseError(c, 400, "STUDENT_EXISTS", "学生已存在")
			return
		}
		middleware.ResponseError(c, 500, "UPDATE_STUDENT_FAILED", "更新学生失败")
		return
	}

	middleware.ResponseOKWithMessage(c, "更新学生成功", child)
}

// ─── 老师班级权限分配 ──────────────────────────────────────────

// AssignmentCreateReq 创建权限分配请求参数
type AssignmentCreateReq struct {
	TeacherUserID  uint64 `json:"teacher_user_id" binding:"required"`
	ClassID        uint64 `json:"class_id" binding:"required"`
	PermissionType string `json:"permission_type" binding:"required,oneof=manage score view"`
	SchoolYear     string `json:"school_year" binding:"required"`
	IsActive       bool   `json:"is_active"`
}

// GetAssignments 查看所有老师-班级分配关系
// GET /api/v1/admin/assignments
func (h *AdminHandler) GetAssignments(c *gin.Context) {
	// 构建查询参数
	params := make(map[string]interface{})
	if teacherUserID := c.Query("teacher_user_id"); teacherUserID != "" {
		if id, err := parseUint64(teacherUserID); err == nil {
			params["teacher_user_id"] = id
		}
	}
	if classID := c.Query("class_id"); classID != "" {
		if id, err := parseUint64(classID); err == nil {
			params["class_id"] = id
		}
	}
	if schoolYear := c.Query("school_year"); schoolYear != "" {
		params["school_year"] = schoolYear
	}

	// 调用服务获取权限分配列表
	assignments, total, err := h.adminSvc.GetAssignments(c.Request.Context(), params)
	if err != nil {
		middleware.ResponseError(c, 500, "GET_ASSIGNMENTS_FAILED", "获取权限分配列表失败")
		return
	}

	middleware.ResponseOK(c, gin.H{
		"list":  assignments,
		"total": total,
	})
}

// CreateAssignment 为老师分配班级权限
// POST /api/v1/admin/assignments
func (h *AdminHandler) CreateAssignment(c *gin.Context) {
	var req AssignmentCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResponseValidationError(c, err.Error())
		return
	}

	// 构建权限分配对象
	assignment := &model.AdminPermission{
		TeacherUserID:  req.TeacherUserID,
		ClassID:        req.ClassID,
		PermissionType: model.PermissionType(req.PermissionType),
		SchoolYear:     req.SchoolYear,
		IsActive:       req.IsActive,
	}

	// 调用服务创建权限分配
	if err := h.adminSvc.CreateAssignment(c.Request.Context(), assignment); err != nil {
		if err == service.ErrResourceNotFound {
			middleware.ResponseError(c, 404, "RESOURCE_NOT_FOUND", "资源不存在")
			return
		}
		if err == service.ErrDuplicateResource {
			middleware.ResponseError(c, 400, "ASSIGNMENT_EXISTS", "权限分配已存在")
			return
		}
		middleware.ResponseError(c, 500, "CREATE_ASSIGNMENT_FAILED", "创建权限分配失败")
		return
	}

	middleware.ResponseOKWithMessage(c, "创建权限分配成功", assignment)
}

// DeleteAssignment 撤销老师班级权限（软删除）
// DELETE /api/v1/admin/assignments/:id
func (h *AdminHandler) DeleteAssignment(c *gin.Context) {
	// 获取权限分配ID
	id, err := parseUint64(c.Param("id"))
	if err != nil {
		middleware.ResponseValidationError(c, "无效的权限分配ID")
		return
	}

	// 调用服务删除权限分配
	if err := h.adminSvc.DeleteAssignment(c.Request.Context(), id); err != nil {
		if err == service.ErrResourceNotFound {
			middleware.ResponseError(c, 404, "ASSIGNMENT_NOT_FOUND", "权限分配不存在")
			return
		}
		middleware.ResponseError(c, 500, "DELETE_ASSIGNMENT_FAILED", "删除权限分配失败")
		return
	}

	middleware.ResponseOKWithMessage(c, "删除权限分配成功", nil)
}

// BatchCreateAssignments 批量为老师分配多个班级
// POST /api/v1/admin/assignments/batch
func (h *AdminHandler) BatchCreateAssignments(c *gin.Context) {
	var assignments []*model.AdminPermission
	if err := c.ShouldBindJSON(&assignments); err != nil {
		middleware.ResponseValidationError(c, err.Error())
		return
	}

	// 调用服务批量创建权限分配
	if err := h.adminSvc.BatchCreateAssignments(c.Request.Context(), assignments); err != nil {
		middleware.ResponseError(c, 500, "BATCH_CREATE_ASSIGNMENTS_FAILED", "批量创建权限分配失败")
		return
	}

	middleware.ResponseOKWithMessage(c, "批量创建权限分配成功", nil)
}

// ─── 家长-学生绑定管理 ────────────────────────────────────────

// ParentBindingCreateReq 创建家长绑定请求参数
type ParentBindingCreateReq struct {
	ParentUserID uint64 `json:"parent_user_id" binding:"required"`
	ChildID      uint64 `json:"child_id" binding:"required"`
	Relationship string `json:"relationship" binding:"required"`
	IsPrimary    bool   `json:"is_primary"`
	IsActive     bool   `json:"is_active"`
}

// GetParentBindings 查看家长绑定关系
// GET /api/v1/admin/parent-bindings
func (h *AdminHandler) GetParentBindings(c *gin.Context) {
	// 构建查询参数
	params := make(map[string]interface{})
	if parentUserID := c.Query("parent_user_id"); parentUserID != "" {
		if id, err := parseUint64(parentUserID); err == nil {
			params["parent_user_id"] = id
		}
	}
	if childID := c.Query("child_id"); childID != "" {
		if id, err := parseUint64(childID); err == nil {
			params["child_id"] = id
		}
	}
	if isPrimary := c.Query("is_primary"); isPrimary != "" {
		params["is_primary"] = isPrimary == "true"
	}

	// 调用服务获取家长绑定列表
	bindings, total, err := h.adminSvc.GetParentBindings(c.Request.Context(), params)
	if err != nil {
		middleware.ResponseError(c, 500, "GET_PARENT_BINDINGS_FAILED", "获取家长绑定列表失败")
		return
	}

	middleware.ResponseOK(c, gin.H{
		"list":  bindings,
		"total": total,
	})
}

// CreateParentBinding 建立家长-学生绑定
// POST /api/v1/admin/parent-bindings
func (h *AdminHandler) CreateParentBinding(c *gin.Context) {
	var req ParentBindingCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResponseValidationError(c, err.Error())
		return
	}

	// 构建家长绑定对象
	binding := &model.ParentChildRelation{
		ParentUserID: req.ParentUserID,
		ChildID:      req.ChildID,
		Relationship: req.Relationship,
		IsPrimary:    req.IsPrimary,
		IsActive:     req.IsActive,
	}

	// 调用服务创建家长绑定
	if err := h.adminSvc.CreateParentBinding(c.Request.Context(), binding); err != nil {
		if err == service.ErrResourceNotFound {
			middleware.ResponseError(c, 404, "RESOURCE_NOT_FOUND", "资源不存在")
			return
		}
		if err == service.ErrDuplicateResource {
			middleware.ResponseError(c, 400, "BINDING_EXISTS", "绑定已存在")
			return
		}
		middleware.ResponseError(c, 500, "CREATE_BINDING_FAILED", "创建家长绑定失败")
		return
	}

	middleware.ResponseOKWithMessage(c, "创建家长绑定成功", binding)
}

// DeleteParentBinding 解除家长-学生绑定
// DELETE /api/v1/admin/parent-bindings/:id
func (h *AdminHandler) DeleteParentBinding(c *gin.Context) {
	// 获取绑定ID
	id, err := parseUint64(c.Param("id"))
	if err != nil {
		middleware.ResponseValidationError(c, "无效的绑定ID")
		return
	}

	// 调用服务删除绑定
	if err := h.adminSvc.DeleteParentBinding(c.Request.Context(), id); err != nil {
		if err == service.ErrResourceNotFound {
			middleware.ResponseError(c, 404, "BINDING_NOT_FOUND", "绑定不存在")
			return
		}
		middleware.ResponseError(c, 500, "DELETE_BINDING_FAILED", "删除绑定失败")
		return
	}

	middleware.ResponseOKWithMessage(c, "删除绑定成功", nil)
}

// ─── 数据概览 ──────────────────────────────────────────────────

// GetDashboard 全局数据概览（学校/班级/学生数等）
// GET /api/v1/admin/dashboard
func (h *AdminHandler) GetDashboard(c *gin.Context) {
	// 调用服务获取数据概览
	dashboard, err := h.adminSvc.GetDashboard(c.Request.Context())
	if err != nil {
		middleware.ResponseError(c, 500, "GET_DASHBOARD_FAILED", "获取数据概览失败")
		return
	}

	middleware.ResponseOKWithMessage(c, "获取数据概览成功", dashboard)
}

// GetAuditLogs 查看管理员操作审计日志
// GET /api/v1/admin/audit-logs
func (h *AdminHandler) GetAuditLogs(c *gin.Context) {
	// 构建查询参数
	params := make(map[string]interface{})
	if action := c.Query("action"); action != "" {
		params["action"] = action
	}
	if adminID := c.Query("admin_id"); adminID != "" {
		if id, err := parseUint64(adminID); err == nil {
			params["admin_id"] = id
		}
	}
	if startDate := c.Query("start_date"); startDate != "" {
		params["start_date"] = startDate
	}
	if endDate := c.Query("end_date"); endDate != "" {
		params["end_date"] = endDate
	}

	// 调用服务获取审计日志列表
	logs, total, err := h.adminSvc.GetAuditLogs(c.Request.Context(), params)
	if err != nil {
		middleware.ResponseError(c, 500, "GET_AUDIT_LOGS_FAILED", "获取审计日志失败")
		return
	}

	middleware.ResponseOK(c, gin.H{
		"list":  logs,
		"total": total,
	})
}

// ─── 伙伴模板管理 ──────────────────────────────────────────────

// PartnerTemplateCreateReq 创建伙伴模板请求参数
type PartnerTemplateCreateReq struct {
	Code           string `json:"code" binding:"required"`
	Name           string `json:"name" binding:"required"`
	Type           string `json:"type" binding:"required,oneof=pet plant anime"`
	Description    string `json:"description"`
	Slogan         string `json:"slogan"`
	LowStageAsset  string `json:"low_stage_asset"`
	MidStageAsset  string `json:"mid_stage_asset"`
	HighStageAsset string `json:"high_stage_asset"`
	IsActive       bool   `json:"is_active"`
	SortOrder      int    `json:"sort_order"`
}

// PartnerTemplateUpdateReq 更新伙伴模板请求参数
type PartnerTemplateUpdateReq struct {
	Name           string `json:"name" binding:"required"`
	Type           string `json:"type" binding:"required,oneof=pet plant anime"`
	Description    string `json:"description"`
	Slogan         string `json:"slogan"`
	LowStageAsset  string `json:"low_stage_asset"`
	MidStageAsset  string `json:"mid_stage_asset"`
	HighStageAsset string `json:"high_stage_asset"`
	IsActive       bool   `json:"is_active"`
	SortOrder      int    `json:"sort_order"`
}

// CreatePartnerTemplate 新增伙伴模板
// POST /api/v1/admin/partner-templates
func (h *AdminHandler) CreatePartnerTemplate(c *gin.Context) {
	var req PartnerTemplateCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResponseValidationError(c, err.Error())
		return
	}

	// 构建伙伴模板对象
	template := &model.PartnerTemplate{
		Code:           req.Code,
		Name:           req.Name,
		Type:           model.PartnerType(req.Type),
		Description:    req.Description,
		Slogan:         req.Slogan,
		LowStageAsset:  req.LowStageAsset,
		MidStageAsset:  req.MidStageAsset,
		HighStageAsset: req.HighStageAsset,
		IsActive:       req.IsActive,
		SortOrder:      req.SortOrder,
	}

	// 调用服务创建伙伴模板
	if err := h.adminSvc.CreatePartnerTemplate(c.Request.Context(), template); err != nil {
		if err == service.ErrDuplicateResource {
			middleware.ResponseError(c, 400, "TEMPLATE_EXISTS", "模板已存在")
			return
		}
		middleware.ResponseError(c, 500, "CREATE_TEMPLATE_FAILED", "创建伙伴模板失败")
		return
	}

	middleware.ResponseOKWithMessage(c, "创建伙伴模板成功", template)
}

// UpdatePartnerTemplate 更新模板信息（资源/鼓励语/启停）
// PUT /api/v1/admin/partner-templates/:id
func (h *AdminHandler) UpdatePartnerTemplate(c *gin.Context) {
	var req PartnerTemplateUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResponseValidationError(c, err.Error())
		return
	}

	// 获取模板ID
	id, err := parseUint64(c.Param("id"))
	if err != nil {
		middleware.ResponseValidationError(c, "无效的模板ID")
		return
	}

	// 构建伙伴模板对象
	template := &model.PartnerTemplate{
		Base: model.Base{
			ID: id,
		},
		Name:           req.Name,
		Type:           model.PartnerType(req.Type),
		Description:    req.Description,
		Slogan:         req.Slogan,
		LowStageAsset:  req.LowStageAsset,
		MidStageAsset:  req.MidStageAsset,
		HighStageAsset: req.HighStageAsset,
		IsActive:       req.IsActive,
		SortOrder:      req.SortOrder,
	}

	// 调用服务更新伙伴模板
	if err := h.adminSvc.UpdatePartnerTemplate(c.Request.Context(), template); err != nil {
		if err == service.ErrResourceNotFound {
			middleware.ResponseError(c, 404, "TEMPLATE_NOT_FOUND", "模板不存在")
			return
		}
		middleware.ResponseError(c, 500, "UPDATE_TEMPLATE_FAILED", "更新伙伴模板失败")
		return
	}

	middleware.ResponseOKWithMessage(c, "更新伙伴模板成功", template)
}

// SeedPartnerTemplates 一键初始化30个预设模板（幂等）
// POST /api/v1/admin/partner-templates/seed
func (h *AdminHandler) SeedPartnerTemplates(c *gin.Context) {
	// 调用服务初始化预设模板
	if err := h.adminSvc.SeedPartnerTemplates(c.Request.Context()); err != nil {
		middleware.ResponseError(c, 500, "SEED_TEMPLATES_FAILED", "初始化预设模板失败")
		return
	}

	middleware.ResponseOKWithMessage(c, "初始化预设模板成功", nil)
}
