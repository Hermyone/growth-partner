// 简化版 BattleHandler
package handler

import (
	"github.com/gin-gonic/gin"
	"growth-partner/internal/middleware"
	"growth-partner/internal/service"
)

type BattleHandler struct {
	battleSvc service.BattleService
}

func NewBattleHandler(svc service.BattleService) *BattleHandler {
	return &BattleHandler{battleSvc: svc}
}

func (h *BattleHandler) CreateRoom(c *gin.Context)       { middleware.ResponseOK(c, gin.H{"room_id": "123"}) }
func (h *BattleHandler) JoinRoom(c *gin.Context)         { middleware.ResponseOK(c, nil) }
func (h *BattleHandler) WebSocketUpgrade(c *gin.Context) { /* 对战实时通信 */ }
func (h *BattleHandler) GetMyHistory(c *gin.Context)     { middleware.ResponseOK(c, []interface{}{}) }
