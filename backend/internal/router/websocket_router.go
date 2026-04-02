// growth-partner/backend/internal/router/websocket_router.go
// WebSocket路由

package router

import (
	"growth-partner/internal/handler"
	"growth-partner/internal/middleware"
	"growth-partner/pkg/jwt"

	"github.com/gin-gonic/gin"
)

// RegisterWebSocketRoutes 注册WebSocket路由
func RegisterWebSocketRoutes(router *gin.RouterGroup, wsHandler *handler.WebSocketHandler, jwtManager *jwt.Manager) {
	// 主WebSocket路由，需要JWT认证
	router.GET("/ws", middleware.Auth(jwtManager), wsHandler.HandleWebSocket)

	// 对战WebSocket路由，需要学生角色权限
	router.GET("/battle/ws", middleware.Auth(jwtManager), middleware.RequireStudent(), wsHandler.HandleBattleWebSocket)
}
