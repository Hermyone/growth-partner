// growth-partner/backend/internal/handler/teacher_handler.go
// 老师端模块控制器

package handler

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"growth-partner/internal/middleware"
	"growth-partner/internal/model"
	"growth-partner/internal/service"

	"github.com/gin-gonic/gin"
)

// TeacherHandler 老师端控制器
type TeacherHandler struct {
	teacherSvc service.TeacherService
}

// NewTeacherHandler 创建老师端控制器实例
func NewTeacherHandler(teacherSvc service.TeacherService) *TeacherHandler {
	return &TeacherHandler{
		teacherSvc: teacherSvc,
	}
}

// ─── 3.1 我的班级管理 ──────────────────────────────────────

// GetMyClasses 获取当前老师被授权的所有班级（含权限类型）
// GET /api/v1/teacher/my-classes
func (h *TeacherHandler) GetMyClasses(c *gin.Context) {
	// 从上下文获取老师ID
	teacherID, exists := c.Get("user_id")
	if !exists {
		middleware.ResponseError(c, 401, "UNAUTHORIZED", "未授权访问")
		return
	}

	// 调用服务获取班级列表
	teacherIDUint, ok := teacherID.(uint64)
	if !ok {
		middleware.ResponseError(c, 401, "UNAUTHORIZED", "无效的用户ID")
		return
	}
	classes, err := h.teacherSvc.GetMyClasses(c.Request.Context(), teacherIDUint)
	if err != nil {
		middleware.ResponseError(c, 500, "GET_CLASSES_FAILED", "获取班级列表失败")
		return
	}

	middleware.ResponseOKWithMessage(c, "获取班级列表成功", classes)
}

// GetClassOverview 获取班级概览（学生数/行为数/成长值等）
// GET /api/v1/teacher/classes/:classId/overview
func (h *TeacherHandler) GetClassOverview(c *gin.Context) {
	// 获取班级ID
	classID, err := strconv.ParseUint(c.Param("classId"), 10, 64)
	if err != nil {
		middleware.ResponseValidationError(c, "无效的班级ID")
		return
	}

	// 调用服务获取班级概览
	overview, err := h.teacherSvc.GetClassOverview(c.Request.Context(), classID)
	if err != nil {
		if err == service.ErrInvalidClassID {
			middleware.ResponseError(c, 404, "CLASS_NOT_FOUND", "班级不存在")
			return
		}
		middleware.ResponseError(c, 500, "GET_CLASS_OVERVIEW_FAILED", "获取班级概览失败")
		return
	}

	middleware.ResponseOKWithMessage(c, "获取班级概览成功", overview)
}

// GetClassStudents 获取班级学生列表（含伙伴/成长值信息）
// GET /api/v1/teacher/classes/:classId/students
func (h *TeacherHandler) GetClassStudents(c *gin.Context) {
	// 获取班级ID
	classID, err := strconv.ParseUint(c.Param("classId"), 10, 64)
	if err != nil {
		middleware.ResponseValidationError(c, "无效的班级ID")
		return
	}

	// 构建查询参数
	params := make(map[string]interface{})
	if name := c.Query("name"); name != "" {
		params["name"] = name
	}

	// 调用服务获取学生列表
	students, total, err := h.teacherSvc.GetClassStudents(c.Request.Context(), classID, params)
	if err != nil {
		if err == service.ErrInvalidClassID {
			middleware.ResponseError(c, 404, "CLASS_NOT_FOUND", "班级不存在")
			return
		}
		middleware.ResponseError(c, 500, "GET_STUDENTS_FAILED", "获取学生列表失败")
		return
	}

	middleware.ResponseOK(c, gin.H{
		"list":  students,
		"total": total,
	})
}

// ─── 3.2 正向行为打分 ──────────────────────────────────────

// RecordBehavior 为学生添加正向行为记录，触发成长值/伙伴进化
// POST /api/v1/teacher/behaviors
// 注意：该接口已在 behavior_handler.go 中实现

// GetBehaviors 查看班级行为记录列表（多条件筛选，分页）
// GET /api/v1/teacher/behaviors
// 注意：该接口已在 behavior_handler.go 中实现

// GetBehavior 获取单条行为记录详情
// GET /api/v1/teacher/behaviors/:id
func (h *TeacherHandler) GetBehavior(c *gin.Context) {
	// 获取行为记录ID
	behaviorID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		middleware.ResponseValidationError(c, "无效的行为记录ID")
		return
	}

	// 调用服务获取行为记录
	behavior, err := h.teacherSvc.GetBehavior(c.Request.Context(), behaviorID)
	if err != nil {
		if err == service.ErrBehaviorNotFound {
			middleware.ResponseError(c, 404, "BEHAVIOR_NOT_FOUND", "行为记录不存在")
			return
		}
		middleware.ResponseError(c, 500, "GET_BEHAVIOR_FAILED", "获取行为记录失败")
		return
	}

	middleware.ResponseOKWithMessage(c, "获取行为记录成功", behavior)
}

// DeleteBehavior 撤销行为记录（24小时内，扣减成长值）
// DELETE /api/v1/teacher/behaviors/:id
func (h *TeacherHandler) DeleteBehavior(c *gin.Context) {
	// 获取行为记录ID
	behaviorID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		middleware.ResponseValidationError(c, "无效的行为记录ID")
		return
	}

	// 从上下文获取老师ID
	teacherID, exists := c.Get("user_id")
	if !exists {
		middleware.ResponseError(c, 401, "UNAUTHORIZED", "未授权访问")
		return
	}

	// 调用服务撤销行为记录
	teacherIDUint, ok := teacherID.(uint64)
	if !ok {
		middleware.ResponseError(c, 401, "UNAUTHORIZED", "无效的用户ID")
		return
	}
	err = h.teacherSvc.DeleteBehavior(c.Request.Context(), behaviorID, teacherIDUint)
	if err != nil {
		if err == service.ErrBehaviorNotFound {
			middleware.ResponseError(c, 404, "BEHAVIOR_NOT_FOUND", "行为记录不存在")
			return
		}
		if err == service.ErrBehaviorExpired {
			middleware.ResponseError(c, 400, "BEHAVIOR_EXPIRED", "行为记录已超过撤销时限")
			return
		}
		middleware.ResponseError(c, 500, "DELETE_BEHAVIOR_FAILED", "撤销行为记录失败")
		return
	}

	middleware.ResponseOKWithMessage(c, "撤销行为记录成功", nil)
}

// BatchRecordBehaviors 批量为多个学生打分
// POST /api/v1/teacher/behaviors/batch
func (h *TeacherHandler) BatchRecordBehaviors(c *gin.Context) {
	// 定义请求结构
	type BatchRecordBehaviorsReq struct {
		Requests []service.RecordBehaviorRequest `json:"requests" binding:"required"`
	}

	var req BatchRecordBehaviorsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResponseValidationError(c, err.Error())
		return
	}

	// 调用服务批量为多个学生打分
	behaviors, err := h.teacherSvc.BatchRecordBehaviors(c.Request.Context(), req.Requests)
	if err != nil {
		middleware.ResponseError(c, 500, "BATCH_RECORD_BEHAVIORS_FAILED", "批量为多个学生打分失败")
		return
	}

	middleware.ResponseOKWithMessage(c, "批量为多个学生打分成功", behaviors)
}

// ─── 3.3 广播发送 ──────────────────────────────────────────

// GetBroadcasts 查看自己发送的广播列表（已发/定时待发）
// GET /api/v1/teacher/broadcasts
func (h *TeacherHandler) GetBroadcasts(c *gin.Context) {
	// 从上下文获取老师ID
	teacherID, exists := c.Get("user_id")
	if !exists {
		middleware.ResponseError(c, 401, "UNAUTHORIZED", "未授权访问")
		return
	}

	// 构建查询参数
	params := make(map[string]interface{})
	if status := c.Query("status"); status != "" {
		params["status"] = status
	}
	if sendTimeStart := c.Query("send_time_start"); sendTimeStart != "" {
		params["send_time_start"] = sendTimeStart
	}
	if sendTimeEnd := c.Query("send_time_end"); sendTimeEnd != "" {
		params["send_time_end"] = sendTimeEnd
	}

	// 调用服务获取广播列表
	teacherIDUint, ok := teacherID.(uint64)
	if !ok {
		middleware.ResponseError(c, 401, "UNAUTHORIZED", "无效的用户ID")
		return
	}
	broadcasts, total, err := h.teacherSvc.GetBroadcasts(c.Request.Context(), teacherIDUint, params)
	if err != nil {
		middleware.ResponseError(c, 500, "GET_BROADCASTS_FAILED", "获取广播列表失败")
		return
	}

	middleware.ResponseOK(c, gin.H{
		"list":  broadcasts,
		"total": total,
	})
}

// SendBroadcast 发送广播（立即/定时）
// POST /api/v1/teacher/broadcasts
// 注意：该接口已在 broadcast_handler.go 中实现

// CancelBroadcast 取消定时广播（仅未发送）
// DELETE /api/v1/teacher/broadcasts/:id
func (h *TeacherHandler) CancelBroadcast(c *gin.Context) {
	// 获取广播ID
	broadcastID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		middleware.ResponseValidationError(c, "无效的广播ID")
		return
	}

	// 从上下文获取老师ID
	teacherID, exists := c.Get("user_id")
	if !exists {
		middleware.ResponseError(c, 401, "UNAUTHORIZED", "未授权访问")
		return
	}

	// 调用服务取消广播
	teacherIDUint, ok := teacherID.(uint64)
	if !ok {
		middleware.ResponseError(c, 401, "UNAUTHORIZED", "无效的用户ID")
		return
	}
	err = h.teacherSvc.CancelBroadcast(c.Request.Context(), broadcastID, teacherIDUint)
	if err != nil {
		if err == service.ErrTeacherUnauthorized {
			middleware.ResponseError(c, 403, "UNAUTHORIZED", "无权限操作该广播")
			return
		}
		middleware.ResponseError(c, 500, "CANCEL_BROADCAST_FAILED", "取消广播失败")
		return
	}

	middleware.ResponseOKWithMessage(c, "取消广播成功", nil)
}

// ─── 3.4 集体挑战管理 ──────────────────────────────────────

// GetChallenges 查看班级当前进行中的集体挑战
// GET /api/v1/teacher/challenges
func (h *TeacherHandler) GetChallenges(c *gin.Context) {
	// 获取班级ID
	classID, err := strconv.ParseUint(c.Query("class_id"), 10, 64)
	if err != nil || classID == 0 {
		middleware.ResponseValidationError(c, "无效的班级ID")
		return
	}

	// 调用服务获取挑战列表
	challenges, err := h.teacherSvc.GetChallenges(c.Request.Context(), classID)
	if err != nil {
		if err == service.ErrInvalidClassID {
			middleware.ResponseError(c, 404, "CLASS_NOT_FOUND", "班级不存在")
			return
		}
		middleware.ResponseError(c, 500, "GET_CHALLENGES_FAILED", "获取挑战列表失败")
		return
	}

	middleware.ResponseOKWithMessage(c, "获取挑战列表成功", challenges)
}

// CreateChallenge 创建集体挑战，配置条件+奖励
// POST /api/v1/teacher/challenges
func (h *TeacherHandler) CreateChallenge(c *gin.Context) {
	// 定义请求结构
	type CreateChallengeReq struct {
		ClassID             uint64 `json:"class_id" binding:"required"`
		Title               string `json:"title" binding:"required"`
		Description         string `json:"description"`
		TargetBehaviorCount int    `json:"target_behavior_count" binding:"required"`
		RewardGrowthPoints  int    `json:"reward_growth_points" binding:"required"`
		StartDate           string `json:"start_date" binding:"required"`
		EndDate             string `json:"end_date" binding:"required"`
	}

	var req CreateChallengeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResponseValidationError(c, err.Error())
		return
	}

	// 从上下文获取老师ID
	teacherID, exists := c.Get("user_id")
	if !exists {
		middleware.ResponseError(c, 401, "UNAUTHORIZED", "未授权访问")
		return
	}

	// 构建挑战对象
	// 解析开始时间和结束时间
	startAt, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		middleware.ResponseValidationError(c, "无效的开始时间格式")
		return
	}

	var endAt *time.Time
	if req.EndDate != "" {
		parsedEndAt, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			middleware.ResponseValidationError(c, "无效的结束时间格式")
			return
		}
		endAt = &parsedEndAt
	}

	teacherIDUint, ok := teacherID.(uint64)
	if !ok {
		middleware.ResponseError(c, 401, "UNAUTHORIZED", "无效的用户ID")
		return
	}
	challenge := &model.Challenge{
		ClassID:     req.ClassID,
		Title:       req.Title,
		Description: req.Description,
		TargetType:  "behavior_count",
		TargetValue: req.TargetBehaviorCount,
		RewardType:  "growth_points",
		RewardValue: req.RewardGrowthPoints,
		StartAt:     startAt,
		EndAt:       endAt,
		CreatedBy:   teacherIDUint,
	}

	// 调用服务创建挑战
	err = h.teacherSvc.CreateChallenge(c.Request.Context(), challenge)
	if err != nil {
		if err == service.ErrInvalidClassID {
			middleware.ResponseError(c, 404, "CLASS_NOT_FOUND", "班级不存在")
			return
		}
		middleware.ResponseError(c, 500, "CREATE_CHALLENGE_FAILED", "创建挑战失败")
		return
	}

	middleware.ResponseOKWithMessage(c, "创建挑战成功", challenge)
}

// CompleteChallenge 手动标记挑战完成，批量发放成长值
// PATCH /api/v1/teacher/challenges/:id/complete
func (h *TeacherHandler) CompleteChallenge(c *gin.Context) {
	// 获取挑战ID
	challengeID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		middleware.ResponseValidationError(c, "无效的挑战ID")
		return
	}

	// 从上下文获取老师ID
	teacherID, exists := c.Get("user_id")
	if !exists {
		middleware.ResponseError(c, 401, "UNAUTHORIZED", "未授权访问")
		return
	}

	// 调用服务完成挑战
	teacherIDUint, ok := teacherID.(uint64)
	if !ok {
		middleware.ResponseError(c, 401, "UNAUTHORIZED", "无效的用户ID")
		return
	}
	err = h.teacherSvc.CompleteChallenge(c.Request.Context(), challengeID, teacherIDUint)
	if err != nil {
		middleware.ResponseError(c, 500, "COMPLETE_CHALLENGE_FAILED", "完成挑战失败")
		return
	}

	middleware.ResponseOKWithMessage(c, "完成挑战成功", nil)
}

// ─── 3.5 题库管理 ──────────────────────────────────────────

// GetQuestions 查看班级题库（公共+专属）
// GET /api/v1/teacher/questions
func (h *TeacherHandler) GetQuestions(c *gin.Context) {
	// 获取班级ID
	classID, err := strconv.ParseUint(c.Query("class_id"), 10, 64)
	if err != nil || classID == 0 {
		middleware.ResponseValidationError(c, "无效的班级ID")
		return
	}

	// 构建查询参数
	params := make(map[string]interface{})
	if subject := c.Query("subject"); subject != "" {
		params["subject"] = subject
	}
	if difficulty := c.Query("difficulty"); difficulty != "" {
		params["difficulty"] = difficulty
	}

	// 调用服务获取题目列表
	questions, total, err := h.teacherSvc.GetQuestions(c.Request.Context(), classID, params)
	if err != nil {
		if err == service.ErrInvalidClassID {
			middleware.ResponseError(c, 404, "CLASS_NOT_FOUND", "班级不存在")
			return
		}
		middleware.ResponseError(c, 500, "GET_QUESTIONS_FAILED", "获取题目列表失败")
		return
	}

	middleware.ResponseOK(c, gin.H{
		"list":  questions,
		"total": total,
	})
}

// CreateQuestion 添加班级专属题目
// POST /api/v1/teacher/questions
func (h *TeacherHandler) CreateQuestion(c *gin.Context) {
	// 定义请求结构
	type CreateQuestionReq struct {
		ClassID       uint64   `json:"class_id" binding:"required"`
		Subject       string   `json:"subject" binding:"required"`
		Content       string   `json:"content" binding:"required"`
		Options       []string `json:"options" binding:"required"`
		CorrectAnswer int      `json:"correct_answer" binding:"required"`
		Difficulty    string   `json:"difficulty" binding:"required"`
		Explanation   string   `json:"explanation"`
	}

	var req CreateQuestionReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResponseValidationError(c, err.Error())
		return
	}

	// 从上下文获取老师ID
	teacherID, exists := c.Get("user_id")
	if !exists {
		middleware.ResponseError(c, 401, "UNAUTHORIZED", "未授权访问")
		return
	}

	// 构建题目对象
	// 将选项转换为JSON字符串
	optionsJSON, err := json.Marshal(req.Options)
	if err != nil {
		middleware.ResponseValidationError(c, "无效的选项格式")
		return
	}

	// 解析难度等级
	difficulty := 1
	switch req.Difficulty {
	case "easy":
		difficulty = 1
	case "medium":
		difficulty = 3
	case "hard":
		difficulty = 5
	}

	// 转换ClassID为指针类型
	classID := req.ClassID

	teacherIDUint, ok := teacherID.(uint64)
	if !ok {
		middleware.ResponseError(c, 401, "UNAUTHORIZED", "无效的用户ID")
		return
	}
	question := &model.Question{
		SubjectID:    1, // 暂时设置为1，后续需要根据实际科目ID设置
		ClassID:      &classID,
		Content:      req.Content,
		QuestionType: "single", // 暂时设置为单选题，后续可以根据需求调整
		Options:      string(optionsJSON),
		Answer:       fmt.Sprintf("%d", req.CorrectAnswer),
		Explanation:  req.Explanation,
		Difficulty:   difficulty,
		IsPublic:     false,
		CreatedBy:    teacherIDUint,
		IsActive:     true,
	}

	// 调用服务创建题目
	err = h.teacherSvc.CreateQuestion(c.Request.Context(), question)
	if err != nil {
		if err == service.ErrInvalidClassID {
			middleware.ResponseError(c, 404, "CLASS_NOT_FOUND", "班级不存在")
			return
		}
		middleware.ResponseError(c, 500, "CREATE_QUESTION_FAILED", "创建题目失败")
		return
	}

	middleware.ResponseOKWithMessage(c, "创建题目成功", question)
}

// UpdateQuestion 编辑题目
// PUT /api/v1/teacher/questions/:id
func (h *TeacherHandler) UpdateQuestion(c *gin.Context) {
	// 获取题目ID
	questionID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		middleware.ResponseValidationError(c, "无效的题目ID")
		return
	}

	// 从上下文获取老师ID
	teacherID, exists := c.Get("user_id")
	if !exists {
		middleware.ResponseError(c, 401, "UNAUTHORIZED", "未授权访问")
		return
	}

	// 定义请求结构
	type UpdateQuestionReq struct {
		ClassID       uint64   `json:"class_id" binding:"required"`
		Subject       string   `json:"subject" binding:"required"`
		Content       string   `json:"content" binding:"required"`
		Options       []string `json:"options" binding:"required"`
		CorrectAnswer int      `json:"correct_answer" binding:"required"`
		Difficulty    string   `json:"difficulty" binding:"required"`
		Explanation   string   `json:"explanation"`
		IsActive      bool     `json:"is_active"`
	}

	var req UpdateQuestionReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResponseValidationError(c, err.Error())
		return
	}

	// 构建题目对象
	// 将选项转换为JSON字符串
	optionsJSON, err := json.Marshal(req.Options)
	if err != nil {
		middleware.ResponseValidationError(c, "无效的选项格式")
		return
	}

	// 解析难度等级
	difficulty := 1
	switch req.Difficulty {
	case "easy":
		difficulty = 1
	case "medium":
		difficulty = 3
	case "hard":
		difficulty = 5
	}

	// 转换ClassID为指针类型
	classID := req.ClassID

	teacherIDUint, ok := teacherID.(uint64)
	if !ok {
		middleware.ResponseError(c, 401, "UNAUTHORIZED", "无效的用户ID")
		return
	}
	question := &model.Question{
		Base: model.Base{
			ID: questionID,
		},
		SubjectID:    1, // 暂时设置为1，后续需要根据实际科目ID设置
		ClassID:      &classID,
		Content:      req.Content,
		QuestionType: "single", // 暂时设置为单选题，后续可以根据需求调整
		Options:      string(optionsJSON),
		Answer:       fmt.Sprintf("%d", req.CorrectAnswer),
		Explanation:  req.Explanation,
		Difficulty:   difficulty,
		IsPublic:     false,
		CreatedBy:    teacherIDUint,
		IsActive:     req.IsActive,
	}

	// 调用服务更新题目
	err = h.teacherSvc.UpdateQuestion(c.Request.Context(), question, teacherIDUint)
	if err != nil {
		if err == service.ErrTeacherUnauthorized {
			middleware.ResponseError(c, 403, "UNAUTHORIZED", "无权限编辑该题目")
			return
		}
		middleware.ResponseError(c, 500, "UPDATE_QUESTION_FAILED", "更新题目失败")
		return
	}

	middleware.ResponseOKWithMessage(c, "更新题目成功", question)
}

// DeleteQuestion 删除题目（软删除）
// DELETE /api/v1/teacher/questions/:id
func (h *TeacherHandler) DeleteQuestion(c *gin.Context) {
	// 获取题目ID
	questionID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		middleware.ResponseValidationError(c, "无效的题目ID")
		return
	}

	// 从上下文获取老师ID
	teacherID, exists := c.Get("user_id")
	if !exists {
		middleware.ResponseError(c, 401, "UNAUTHORIZED", "未授权访问")
		return
	}

	// 调用服务删除题目
	teacherIDUint, ok := teacherID.(uint64)
	if !ok {
		middleware.ResponseError(c, 401, "UNAUTHORIZED", "无效的用户ID")
		return
	}
	err = h.teacherSvc.DeleteQuestion(c.Request.Context(), questionID, teacherIDUint)
	if err != nil {
		if err == service.ErrTeacherUnauthorized {
			middleware.ResponseError(c, 403, "UNAUTHORIZED", "无权限删除该题目")
			return
		}
		middleware.ResponseError(c, 500, "DELETE_QUESTION_FAILED", "删除题目失败")
		return
	}

	middleware.ResponseOKWithMessage(c, "删除题目成功", nil)
}

// BatchImportQuestions 批量导入题目（CSV）
// POST /api/v1/teacher/questions/batch-import
func (h *TeacherHandler) BatchImportQuestions(c *gin.Context) {
	// 从上下文获取老师ID
	teacherID, exists := c.Get("user_id")
	if !exists {
		middleware.ResponseError(c, 401, "UNAUTHORIZED", "未授权访问")
		return
	}

	// 定义请求结构
	type BatchImportQuestionsReq struct {
		Questions []*model.Question `json:"questions" binding:"required"`
	}

	var req BatchImportQuestionsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResponseValidationError(c, err.Error())
		return
	}

	// 调用服务批量导入题目
	teacherIDUint, ok := teacherID.(uint64)
	if !ok {
		middleware.ResponseError(c, 401, "UNAUTHORIZED", "无效的用户ID")
		return
	}
	err := h.teacherSvc.BatchImportQuestions(c.Request.Context(), req.Questions, teacherIDUint)
	if err != nil {
		if err == service.ErrTeacherUnauthorized {
			middleware.ResponseError(c, 403, "UNAUTHORIZED", "无权限导入题目到该班级")
			return
		}
		middleware.ResponseError(c, 500, "BATCH_IMPORT_QUESTIONS_FAILED", "批量导入题目失败")
		return
	}

	middleware.ResponseOKWithMessage(c, "批量导入题目成功", nil)
}

// ─── 3.6 盲盒奖励池管理 ──────────────────────────────────

// GetBlindboxPool 查看本班盲盒奖励池
// GET /api/v1/teacher/blindbox/pool
// 注意：该接口已在 blindbox_handler.go 中实现

// AddToBlindboxPool 向奖励池添加奖励
// POST /api/v1/teacher/blindbox/pool
// 注意：该接口已在 blindbox_handler.go 中实现

// UpdateBlindboxPoolItem 编辑奖励配置
// PUT /api/v1/teacher/blindbox/pool/:id
func (h *TeacherHandler) UpdateBlindboxPoolItem(c *gin.Context) {
	// 获取盲盒奖励ID
	blindboxID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		middleware.ResponseValidationError(c, "无效的盲盒奖励ID")
		return
	}

	// 从上下文获取老师ID
	teacherID, exists := c.Get("user_id")
	if !exists {
		middleware.ResponseError(c, 401, "UNAUTHORIZED", "未授权访问")
		return
	}

	// 定义请求结构
	type UpdateBlindboxPoolItemReq struct {
		ClassID     uint64  `json:"class_id" binding:"required"`
		Name        string  `json:"name" binding:"required"`
		Description string  `json:"description"`
		Rarity      string  `json:"rarity" binding:"required"`
		Probability float64 `json:"probability" binding:"required"`
		Stock       int     `json:"stock" binding:"required"`
		IsActive    bool    `json:"is_active"`
	}

	var req UpdateBlindboxPoolItemReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResponseValidationError(c, err.Error())
		return
	}

	// 构建盲盒奖励对象
	blindbox := &model.BlindBoxPool{
		Base: model.Base{
			ID: blindboxID,
		},
		ClassID:     req.ClassID,
		Title:       req.Name, // 使用Title字段代替Name
		Description: req.Description,
		Stock:       req.Stock,
		IsActive:    req.IsActive,
	}

	// 调用服务更新盲盒奖励
	teacherIDUint, ok := teacherID.(uint64)
	if !ok {
		middleware.ResponseError(c, 401, "UNAUTHORIZED", "无效的用户ID")
		return
	}
	err = h.teacherSvc.UpdateBlindboxPoolItem(c.Request.Context(), blindbox, teacherIDUint)
	if err != nil {
		if err == service.ErrTeacherUnauthorized {
			middleware.ResponseError(c, 403, "UNAUTHORIZED", "无权限编辑该盲盒奖励")
			return
		}
		middleware.ResponseError(c, 500, "UPDATE_BLINDbox_POOL_ITEM_FAILED", "更新盲盒奖励失败")
		return
	}

	middleware.ResponseOKWithMessage(c, "更新盲盒奖励成功", blindbox)
}

// RemoveFromBlindboxPool 下架奖励（软删除）
// DELETE /api/v1/teacher/blindbox/pool/:id
// 注意：该接口已在 blindbox_handler.go 中实现

// DrawBlindboxForStudent 为学生触发抽盲盒
// POST /api/v1/teacher/blindbox/draw/:childId
// 注意：该接口已在 blindbox_handler.go 中实现

// ConfirmBlindboxRedeem 确认兑换学生盲盒奖励
// PATCH /api/v1/teacher/blindbox/draws/:drawId/redeem
func (h *TeacherHandler) ConfirmBlindboxRedeem(c *gin.Context) {
	// 获取盲盒抽取记录ID
	drawID, err := strconv.ParseUint(c.Param("drawId"), 10, 64)
	if err != nil {
		middleware.ResponseValidationError(c, "无效的盲盒抽取记录ID")
		return
	}

	// 从上下文获取老师ID
	teacherID, exists := c.Get("user_id")
	if !exists {
		middleware.ResponseError(c, 401, "UNAUTHORIZED", "未授权访问")
		return
	}

	// 调用服务确认兑换
	teacherIDUint, ok := teacherID.(uint64)
	if !ok {
		middleware.ResponseError(c, 401, "UNAUTHORIZED", "无效的用户ID")
		return
	}
	err = h.teacherSvc.ConfirmBlindboxRedeem(c.Request.Context(), drawID, teacherIDUint)
	if err != nil {
		if err == service.ErrTeacherUnauthorized {
			middleware.ResponseError(c, 403, "UNAUTHORIZED", "无权限操作该盲盒奖励")
			return
		}
		middleware.ResponseError(c, 500, "CONFIRM_BLINDbox_REDEEM_FAILED", "确认兑换盲盒奖励失败")
		return
	}

	middleware.ResponseOKWithMessage(c, "确认兑换盲盒奖励成功", nil)
}

// ─── 3.7 周报PDF生成 ──────────────────────────────────────

// GenerateWeeklyReport 触发生成本班本周正能量周报PDF（异步）
// POST /api/v1/teacher/reports/weekly
func (h *TeacherHandler) GenerateWeeklyReport(c *gin.Context) {
	// 从上下文获取老师ID
	teacherID, exists := c.Get("user_id")
	if !exists {
		middleware.ResponseError(c, 401, "UNAUTHORIZED", "未授权访问")
		return
	}

	// 定义请求结构
	type GenerateWeeklyReportReq struct {
		ClassID uint64 `json:"class_id" binding:"required"`
	}

	var req GenerateWeeklyReportReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResponseValidationError(c, err.Error())
		return
	}

	// 调用服务生成周报
	teacherIDUint, ok := teacherID.(uint64)
	if !ok {
		middleware.ResponseError(c, 401, "UNAUTHORIZED", "无效的用户ID")
		return
	}
	err := h.teacherSvc.GenerateWeeklyReport(c.Request.Context(), req.ClassID, teacherIDUint)
	if err != nil {
		if err == service.ErrInvalidClassID {
			middleware.ResponseError(c, 404, "CLASS_NOT_FOUND", "班级不存在")
			return
		}
		if err == service.ErrTeacherUnauthorized {
			middleware.ResponseError(c, 403, "UNAUTHORIZED", "无权限操作该班级")
			return
		}
		middleware.ResponseError(c, 500, "GENERATE_WEEKLY_REPORT_FAILED", "生成周报失败")
		return
	}

	middleware.ResponseOKWithMessage(c, "生成周报请求已提交，正在处理中", nil)
}

// GetWeeklyReports 查看历史周报列表（含下载链接）
// GET /api/v1/teacher/reports/weekly
func (h *TeacherHandler) GetWeeklyReports(c *gin.Context) {
	// 获取班级ID
	classID, err := strconv.ParseUint(c.Query("class_id"), 10, 64)
	if err != nil || classID == 0 {
		middleware.ResponseValidationError(c, "无效的班级ID")
		return
	}

	// 从上下文获取老师ID
	teacherID, exists := c.Get("user_id")
	if !exists {
		middleware.ResponseError(c, 401, "UNAUTHORIZED", "未授权访问")
		return
	}

	// 构建查询参数
	params := make(map[string]interface{})
	if startDate := c.Query("start_date"); startDate != "" {
		params["start_date"] = startDate
	}
	if endDate := c.Query("end_date"); endDate != "" {
		params["end_date"] = endDate
	}

	// 调用服务获取周报列表
	teacherIDUint, ok := teacherID.(uint64)
	if !ok {
		middleware.ResponseError(c, 401, "UNAUTHORIZED", "无效的用户ID")
		return
	}
	reports, total, err := h.teacherSvc.GetWeeklyReports(c.Request.Context(), classID, teacherIDUint, params)
	if err != nil {
		if err == service.ErrInvalidClassID {
			middleware.ResponseError(c, 404, "CLASS_NOT_FOUND", "班级不存在")
			return
		}
		if err == service.ErrTeacherUnauthorized {
			middleware.ResponseError(c, 403, "UNAUTHORIZED", "无权限操作该班级")
			return
		}
		middleware.ResponseError(c, 500, "GET_WEEKLY_REPORTS_FAILED", "获取周报列表失败")
		return
	}

	middleware.ResponseOK(c, gin.H{
		"list":  reports,
		"total": total,
	})
}

// DownloadWeeklyReport 下载指定周报PDF
// GET /api/v1/teacher/reports/weekly/:id/download
func (h *TeacherHandler) DownloadWeeklyReport(c *gin.Context) {
	// 获取报告ID
	reportID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		middleware.ResponseValidationError(c, "无效的报告ID")
		return
	}

	// 调用服务下载周报
	pdfData, fileName, err := h.teacherSvc.DownloadWeeklyReport(c.Request.Context(), reportID)
	if err != nil {
		middleware.ResponseError(c, 500, "DOWNLOAD_WEEKLY_REPORT_FAILED", "下载周报失败")
		return
	}

	// 设置响应头
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Length", strconv.Itoa(len(pdfData)))

	// 写入响应体
	c.Data(200, "application/pdf", pdfData)
}
