// growth-partner/backend/internal/handler/battle_handler.go
// 知识对战模块控制器

package handler

import (
	"growth-partner/internal/middleware"
	"growth-partner/internal/service"

	"github.com/gin-gonic/gin"
)

// BattleHandler 对战控制器
type BattleHandler struct {
	battleSvc service.BattleService
}

// NewBattleHandler 创建对战控制器实例
func NewBattleHandler(battleSvc service.BattleService) *BattleHandler {
	return &BattleHandler{
		battleSvc: battleSvc,
	}
}

// ─── 对战模块核心功能 ────────────────────────────────────────

// GET /api/v1/battle/subjects
func (h *BattleHandler) GetBattleSubjects(c *gin.Context) {
	// 调用服务获取可用对战科目列表
	subjects, err := h.battleSvc.GetBattleSubjects(c.Request.Context())
	if err != nil {
		middleware.ResponseError(c, 500, "INTERNAL_ERROR", err.Error())
		return
	}

	middleware.ResponseOK(c, subjects)
}

// POST /api/v1/battle/rooms
func (h *BattleHandler) CreateBattleRoom(c *gin.Context) {
	// 从上下文获取学生ID
	studentID, exists := c.Get("user_id")
	if !exists {
		middleware.ResponseError(c, 401, "UNAUTHORIZED", "未授权访问")
		return
	}

	// 定义请求结构
	type CreateBattleRoomReq struct {
		Subject string `json:"subject" binding:"required"`
	}

	var req CreateBattleRoomReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResponseValidationError(c, err.Error())
		return
	}

	// 调用服务创建对战房间
	studentIDUint, ok := studentID.(uint64)
	if !ok {
		middleware.ResponseError(c, 401, "UNAUTHORIZED", "无效的用户ID")
		return
	}
	room, err := h.battleSvc.CreateBattleRoom(c.Request.Context(), studentIDUint, req.Subject)
	if err != nil {
		middleware.ResponseError(c, 500, "INTERNAL_ERROR", err.Error())
		return
	}

	middleware.ResponseOK(c, room)
}

// POST /api/v1/battle/rooms/:roomCode/join
func (h *BattleHandler) JoinBattleRoom(c *gin.Context) {
	// 从上下文获取学生ID
	studentID, exists := c.Get("user_id")
	if !exists {
		middleware.ResponseError(c, 401, "UNAUTHORIZED", "未授权访问")
		return
	}

	// 获取房间代码
	roomCode := c.Param("roomCode")
	if roomCode == "" {
		middleware.ResponseValidationError(c, "无效的房间代码")
		return
	}

	// 调用服务加入对战房间
	studentIDUint, ok := studentID.(uint64)
	if !ok {
		middleware.ResponseError(c, 401, "UNAUTHORIZED", "无效的用户ID")
		return
	}
	room, err := h.battleSvc.JoinBattleRoom(c.Request.Context(), studentIDUint, roomCode)
	if err != nil {
		middleware.ResponseError(c, 500, "INTERNAL_ERROR", err.Error())
		return
	}

	middleware.ResponseOK(c, room)
}

// GET /api/v1/battle/rooms/:roomCode
func (h *BattleHandler) GetBattleRoom(c *gin.Context) {
	// 从上下文获取学生ID
	_, exists := c.Get("user_id")
	if !exists {
		middleware.ResponseError(c, 401, "UNAUTHORIZED", "未授权访问")
		return
	}

	// 获取房间代码
	roomCode := c.Param("roomCode")
	if roomCode == "" {
		middleware.ResponseValidationError(c, "无效的房间代码")
		return
	}

	// 调用服务查询房间状态
	room, err := h.battleSvc.GetBattleRoom(c.Request.Context(), roomCode)
	if err != nil {
		middleware.ResponseError(c, 500, "INTERNAL_ERROR", err.Error())
		return
	}

	middleware.ResponseOK(c, room)
}

// GET /api/v1/battle/history
func (h *BattleHandler) GetBattleHistory(c *gin.Context) {
	// 从上下文获取学生ID
	studentID, exists := c.Get("user_id")
	if !exists {
		middleware.ResponseError(c, 401, "UNAUTHORIZED", "未授权访问")
		return
	}

	// 调用服务查看对战历史
	studentIDUint, ok := studentID.(uint64)
	if !ok {
		middleware.ResponseError(c, 401, "UNAUTHORIZED", "无效的用户ID")
		return
	}
	history, err := h.battleSvc.GetBattleHistory(c.Request.Context(), studentIDUint)
	if err != nil {
		middleware.ResponseError(c, 500, "INTERNAL_ERROR", err.Error())
		return
	}

	middleware.ResponseOK(c, history)
}

// GET /api/v1/battle/history/:roomId/review
func (h *BattleHandler) GetBattleReview(c *gin.Context) {
	// 从上下文获取学生ID
	studentID, exists := c.Get("user_id")
	if !exists {
		middleware.ResponseError(c, 401, "UNAUTHORIZED", "未授权访问")
		return
	}

	// 获取房间ID
	roomID := c.Param("roomId")
	if roomID == "" {
		middleware.ResponseValidationError(c, "无效的房间ID")
		return
	}

	// 调用服务获取对战复盘
	studentIDUint, ok := studentID.(uint64)
	if !ok {
		middleware.ResponseError(c, 401, "UNAUTHORIZED", "无效的用户ID")
		return
	}
	review, err := h.battleSvc.GetBattleReview(c.Request.Context(), studentIDUint, roomID)
	if err != nil {
		middleware.ResponseError(c, 500, "INTERNAL_ERROR", err.Error())
		return
	}

	middleware.ResponseOK(c, review)
}
