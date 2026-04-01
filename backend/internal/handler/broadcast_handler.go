// growth-partner/backend/internal/handler/broadcast_handler.go

package handler

import (
	"github.com/gin-gonic/gin"
	"growth-partner/internal/middleware"
	"growth-partner/internal/service"
)

type BroadcastHandler struct {
	broadcastSvc service.BroadcastService
}

func NewBroadcastHandler(broadcastSvc service.BroadcastService) *BroadcastHandler {
	return &BroadcastHandler{broadcastSvc: broadcastSvc}
}

// WebSocketUpgrade 升级到 WebSocket 连接
func (h *BroadcastHandler) WebSocketUpgrade(c *gin.Context) {
	// 实际开发中会在这里调用 websocket 库进行 Upgrade
	// 并根据 userID 注册到 Hub 中
	middleware.ResponseOK(c, gin.H{"message": "WebSocket 握手成功（模拟）"})
}

// Send 园长发送广播
func (h *BroadcastHandler) Send(c *gin.Context) {
	// TODO: 解析广播内容并调用 Redis Pub/Sub
	middleware.ResponseOK(c, gin.H{"message": "广播已发出"})
}

// GetBroadcastList 获取历史通知列表
func (h *BroadcastHandler) GetBroadcastList(c *gin.Context) {
	page, size := getPaginationParams(c)
	middleware.ResponseOK(c, gin.H{"page": page, "size": size, "list": []interface{}{}})
}

// MarkAsRead 标记广播为已读
func (h *BroadcastHandler) MarkAsRead(c *gin.Context) {
	middleware.ResponseOK(c, nil)
}
