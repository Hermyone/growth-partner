// growth-partner/backend/internal/router/parent_router.go
// 家长端路由

package router

import (
	"growth-partner/internal/handler"
	"growth-partner/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterParentRoutes 注册家长端路由
func RegisterParentRoutes(router *gin.RouterGroup, parentHandler *handler.ParentHandler) {
	// 家长端路由组，需要家长角色权限
	parentRoutes := router.Group("/parent")
	parentRoutes.Use(middleware.RequireParent())
	{
		// 获取绑定的孩子列表
		parentRoutes.GET("/children", parentHandler.GetMyChildren)
		
		// 孩子相关接口
		childRoutes := parentRoutes.Group("/children/:childId")
		{
			// 查看孩子当前伙伴状态
			childRoutes.GET("/partner", parentHandler.GetChildPartner)
			
			// 查看孩子的历史伙伴列表
			childRoutes.GET("/partners", parentHandler.GetChildPartners)
			
			// 查看孩子的正向行为记录
			childRoutes.GET("/behaviors", parentHandler.GetChildBehaviors)
			
			// 查看孩子收到的伙伴鼓励广播
			childRoutes.GET("/broadcasts", parentHandler.GetChildBroadcasts)
			
			// 查看孩子的里程碑贴纸
			childRoutes.GET("/milestones", parentHandler.GetChildMilestones)
			
			// 查看孩子本月/历史月度成长卡
			childRoutes.GET("/monthly-card", parentHandler.GetChildMonthlyCard)
			
			// 查看孩子年度成长画卷
			childRoutes.GET("/annual-report", parentHandler.GetChildAnnualReport)
			
			// 查看孩子的对战参与记录
			childRoutes.GET("/battles", parentHandler.GetChildBattles)
		}
	}
}
