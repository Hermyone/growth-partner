// growth-partner/backend/internal/handler/websocket_handler.go
// WebSocket处理器

package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// WebSocketHandler WebSocket处理器
type WebSocketHandler struct {
	upgrader websocket.Upgrader
}

// NewWebSocketHandler 创建WebSocket处理器实例
func NewWebSocketHandler() *WebSocketHandler {
	return &WebSocketHandler{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// 允许所有来源的WebSocket连接，生产环境应该根据实际情况进行限制
				return true
			},
		},
	}
}

// HandleWebSocket 处理WebSocket连接
func (h *WebSocketHandler) HandleWebSocket(c *gin.Context) {
	// 从上下文获取用户信息
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		return
	}

	// 升级HTTP连接为WebSocket连接
	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("升级WebSocket连接失败: %v", err)
		return
	}
	defer conn.Close()

	// 处理WebSocket消息
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("读取WebSocket消息失败: %v", err)
			break
		}

		log.Printf("接收到用户 %v 的消息: %s", userID, message)

		// 回显消息
		if err := conn.WriteMessage(messageType, message); err != nil {
			log.Printf("发送WebSocket消息失败: %v", err)
			break
		}
	}
}

// HandleBattleWebSocket 处理对战WebSocket连接
func (h *WebSocketHandler) HandleBattleWebSocket(c *gin.Context) {
	// 从上下文获取用户信息
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		return
	}

	// 获取房间代码
	roomCode := c.Query("room_code")
	if roomCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少房间代码"})
		return
	}

	// 升级HTTP连接为WebSocket连接
	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("升级WebSocket连接失败: %v", err)
		return
	}
	defer conn.Close()

	// 处理对战WebSocket消息
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("读取WebSocket消息失败: %v", err)
			break
		}

		log.Printf("接收到用户 %v 在房间 %v 的消息: %s", userID, roomCode, message)

		// 回显消息
		if err := conn.WriteMessage(messageType, message); err != nil {
			log.Printf("发送WebSocket消息失败: %v", err)
			break
		}
	}
}
