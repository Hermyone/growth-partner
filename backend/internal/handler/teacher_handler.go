// growth-partner/backend/internal/handler/teacher_handler.go
// 老师端模块控制器

package handler

import (
	"growth-partner/internal/middleware"
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
	// TODO: 实现获取老师班级列表功能
	middleware.ResponseOKWithMessage(c, "获取老师班级列表接口待实现", nil)
}

// GetClassOverview 获取班级概览（学生数/行为数/成长值等）
// GET /api/v1/teacher/classes/:classId/overview
func (h *TeacherHandler) GetClassOverview(c *gin.Context) {
	// TODO: 实现获取班级概览功能
	middleware.ResponseOKWithMessage(c, "获取班级概览接口待实现", nil)
}

// GetClassStudents 获取班级学生列表（含伙伴/成长值信息）
// GET /api/v1/teacher/classes/:classId/students
func (h *TeacherHandler) GetClassStudents(c *gin.Context) {
	// TODO: 实现获取班级学生列表功能
	middleware.ResponseOKWithMessage(c, "获取班级学生列表接口待实现", nil)
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
	// TODO: 实现获取单条行为记录详情功能
	middleware.ResponseOKWithMessage(c, "获取单条行为记录详情接口待实现", nil)
}

// DeleteBehavior 撤销行为记录（24小时内，扣减成长值）
// DELETE /api/v1/teacher/behaviors/:id
func (h *TeacherHandler) DeleteBehavior(c *gin.Context) {
	// TODO: 实现撤销行为记录功能
	middleware.ResponseOKWithMessage(c, "撤销行为记录接口待实现", nil)
}

// BatchRecordBehaviors 批量为多个学生打分
// POST /api/v1/teacher/behaviors/batch
func (h *TeacherHandler) BatchRecordBehaviors(c *gin.Context) {
	// TODO: 实现批量为多个学生打分功能
	middleware.ResponseOKWithMessage(c, "批量为多个学生打分接口待实现", nil)
}

// ─── 3.3 广播发送 ──────────────────────────────────────────

// GetBroadcasts 查看自己发送的广播列表（已发/定时待发）
// GET /api/v1/teacher/broadcasts
func (h *TeacherHandler) GetBroadcasts(c *gin.Context) {
	// TODO: 实现查看自己发送的广播列表功能
	middleware.ResponseOKWithMessage(c, "查看自己发送的广播列表接口待实现", nil)
}

// SendBroadcast 发送广播（立即/定时）
// POST /api/v1/teacher/broadcasts
// 注意：该接口已在 broadcast_handler.go 中实现

// CancelBroadcast 取消定时广播（仅未发送）
// DELETE /api/v1/teacher/broadcasts/:id
func (h *TeacherHandler) CancelBroadcast(c *gin.Context) {
	// TODO: 实现取消定时广播功能
	middleware.ResponseOKWithMessage(c, "取消定时广播接口待实现", nil)
}

// ─── 3.4 集体挑战管理 ──────────────────────────────────────

// GetChallenges 查看班级当前进行中的集体挑战
// GET /api/v1/teacher/challenges
func (h *TeacherHandler) GetChallenges(c *gin.Context) {
	// TODO: 实现查看班级当前进行中的集体挑战功能
	middleware.ResponseOKWithMessage(c, "查看班级当前进行中的集体挑战接口待实现", nil)
}

// CreateChallenge 创建集体挑战，配置条件+奖励
// POST /api/v1/teacher/challenges
func (h *TeacherHandler) CreateChallenge(c *gin.Context) {
	// TODO: 实现创建集体挑战功能
	middleware.ResponseOKWithMessage(c, "创建集体挑战接口待实现", nil)
}

// CompleteChallenge 手动标记挑战完成，批量发放成长值
// PATCH /api/v1/teacher/challenges/:id/complete
func (h *TeacherHandler) CompleteChallenge(c *gin.Context) {
	// TODO: 实现手动标记挑战完成功能
	middleware.ResponseOKWithMessage(c, "手动标记挑战完成接口待实现", nil)
}

// ─── 3.5 题库管理 ──────────────────────────────────────────

// GetQuestions 查看班级题库（公共+专属）
// GET /api/v1/teacher/questions
func (h *TeacherHandler) GetQuestions(c *gin.Context) {
	// TODO: 实现查看班级题库功能
	middleware.ResponseOKWithMessage(c, "查看班级题库接口待实现", nil)
}

// CreateQuestion 添加班级专属题目
// POST /api/v1/teacher/questions
func (h *TeacherHandler) CreateQuestion(c *gin.Context) {
	// TODO: 实现添加班级专属题目功能
	middleware.ResponseOKWithMessage(c, "添加班级专属题目接口待实现", nil)
}

// UpdateQuestion 编辑题目
// PUT /api/v1/teacher/questions/:id
func (h *TeacherHandler) UpdateQuestion(c *gin.Context) {
	// TODO: 实现编辑题目功能
	middleware.ResponseOKWithMessage(c, "编辑题目接口待实现", nil)
}

// DeleteQuestion 删除题目（软删除）
// DELETE /api/v1/teacher/questions/:id
func (h *TeacherHandler) DeleteQuestion(c *gin.Context) {
	// TODO: 实现删除题目功能
	middleware.ResponseOKWithMessage(c, "删除题目接口待实现", nil)
}

// BatchImportQuestions 批量导入题目（CSV）
// POST /api/v1/teacher/questions/batch-import
func (h *TeacherHandler) BatchImportQuestions(c *gin.Context) {
	// TODO: 实现批量导入题目功能
	middleware.ResponseOKWithMessage(c, "批量导入题目接口待实现", nil)
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
	// TODO: 实现编辑奖励配置功能
	middleware.ResponseOKWithMessage(c, "编辑奖励配置接口待实现", nil)
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
	// TODO: 实现确认兑换学生盲盒奖励功能
	middleware.ResponseOKWithMessage(c, "确认兑换学生盲盒奖励接口待实现", nil)
}

// ─── 3.7 周报PDF生成 ──────────────────────────────────────

// GenerateWeeklyReport 触发生成本班本周正能量周报PDF（异步）
// POST /api/v1/teacher/reports/weekly
func (h *TeacherHandler) GenerateWeeklyReport(c *gin.Context) {
	// TODO: 实现触发生成本班本周正能量周报PDF功能
	middleware.ResponseOKWithMessage(c, "触发生成本班本周正能量周报PDF接口待实现", nil)
}

// GetWeeklyReports 查看历史周报列表（含下载链接）
// GET /api/v1/teacher/reports/weekly
func (h *TeacherHandler) GetWeeklyReports(c *gin.Context) {
	// TODO: 实现查看历史周报列表功能
	middleware.ResponseOKWithMessage(c, "查看历史周报列表接口待实现", nil)
}

// DownloadWeeklyReport 下载指定周报PDF
// GET /api/v1/teacher/reports/weekly/:id/download
func (h *TeacherHandler) DownloadWeeklyReport(c *gin.Context) {
	// TODO: 实现下载指定周报PDF功能
	middleware.ResponseOKWithMessage(c, "下载指定周报PDF接口待实现", nil)
}
