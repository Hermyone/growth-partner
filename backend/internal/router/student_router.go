// growth-partner/backend/internal/router/student_router.go
// 学生端路由

package router

import (
	"growth-partner/internal/handler"
	"growth-partner/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterStudentRoutes 注册学生端路由
func RegisterStudentRoutes(router *gin.RouterGroup, studentHandler *handler.StudentHandler) {
	// 学生端路由组，需要学生角色权限
	studentRoutes := router.Group("/student")
	studentRoutes.Use(middleware.RequireStudent())
	{
		// ─── 4.1 伙伴系统 ──────────────────────────────────────────
		studentRoutes.GET("/partner", studentHandler.GetCurrentPartner)
		studentRoutes.GET("/partners", studentHandler.GetAllPartners)
		studentRoutes.POST("/partner", studentHandler.SelectNewPartner)
		studentRoutes.PATCH("/partner/nickname", studentHandler.UpdatePartnerNickname)
		studentRoutes.GET("/partner/growth-history", studentHandler.GetPartnerGrowthHistory)
		studentRoutes.GET("/partner/templates", studentHandler.GetPartnerTemplates)

		// ─── 4.2 行为记录查看 ───────────────────────────────────────
		studentRoutes.GET("/behaviors", studentHandler.GetMyBehaviors)
		studentRoutes.GET("/behaviors/stats", studentHandler.GetBehaviorStats)

		// ─── 4.3 广播中心 ──────────────────────────────────────────
		studentRoutes.GET("/broadcasts", studentHandler.GetMyBroadcasts)
		studentRoutes.PATCH("/broadcasts/:id/read", studentHandler.MarkBroadcastAsRead)
		studentRoutes.POST("/broadcasts/read-all", studentHandler.MarkAllBroadcastsAsRead)

		// ─── 4.4 成长年历 ──────────────────────────────────────────
		studentRoutes.GET("/growth-calendar/months", studentHandler.GetGrowthCalendarMonths)
		studentRoutes.GET("/growth-calendar/months/:month", studentHandler.GetGrowthCalendarMonth)
		studentRoutes.GET("/growth-calendar/annual/:year", studentHandler.GetGrowthCalendarAnnual)
		studentRoutes.GET("/milestones", studentHandler.GetMilestones)

		// ─── 4.5 盲盒查看 ──────────────────────────────────────────
		studentRoutes.GET("/blindbox/my-draws", studentHandler.GetMyBlindboxDraws)
	}
}
