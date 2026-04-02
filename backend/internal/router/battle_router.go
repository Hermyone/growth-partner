// growth-partner/backend/internal/router/battle_router.go
// 知识对战路由

package router

import (
	"growth-partner/internal/handler"
	"growth-partner/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterBattleRoutes 注册对战路由
func RegisterBattleRoutes(router *gin.RouterGroup, battleHandler *handler.BattleHandler) {
	// 对战路由组，需要学生角色权限
	battleRoutes := router.Group("/battle")
	battleRoutes.Use(middleware.RequireStudent())
	{
		// 获取可用对战科目列表
		battleRoutes.GET("/subjects", battleHandler.GetBattleSubjects)
		
		// 房间管理
		battleRoutes.POST("/rooms", battleHandler.CreateBattleRoom)
		battleRoutes.POST("/rooms/:roomCode/join", battleHandler.JoinBattleRoom)
		battleRoutes.GET("/rooms/:roomCode", battleHandler.GetBattleRoom)
		
		// 对战历史
		battleRoutes.GET("/history", battleHandler.GetBattleHistory)
		battleRoutes.GET("/history/:roomId/review", battleHandler.GetBattleReview)
		
		// WebSocket接口
		// 注意：WebSocket接口需要单独处理，这里只注册HTTP路由
	}
}
