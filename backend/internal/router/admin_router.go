// growth-partner/backend/internal/router/admin_router.go
// 管理员模块路由配置

package router

import (
	"growth-partner/config"
	"growth-partner/internal/handler"
	"growth-partner/internal/middleware"
	"growth-partner/pkg/jwt"

	"github.com/gin-gonic/gin"
)

// SetupAdminRoutes 配置管理员相关路由
func SetupAdminRoutes(
	v1 *gin.RouterGroup,
	cfg *config.Config,
	jwtManager *jwt.Manager,
	adminHandler *handler.AdminHandler,
) {
	// ─── 管理员接口（需要 admin 角色）───────────────────────────
	adminAPI := v1.Group("/admin",
		middleware.Auth(jwtManager),
		middleware.RequireRoles("admin"),
	)
	{
		// 2.1 学校管理
		schoolGroup := adminAPI.Group("/schools")
		{
			schoolGroup.GET("", adminHandler.GetSchools)
			schoolGroup.POST("", adminHandler.CreateSchool)
			schoolGroup.PUT("/:id", adminHandler.UpdateSchool)
			schoolGroup.PATCH("/:id/status", adminHandler.UpdateSchoolStatus)
		}

		// 2.2 班级管理
		classGroup := adminAPI.Group("/classes")
		{
			classGroup.GET("", adminHandler.GetClasses)
			classGroup.POST("", adminHandler.CreateClass)
			classGroup.PUT("/:id", adminHandler.UpdateClass)
			classGroup.POST("/:id/promote", adminHandler.PromoteClass)
			classGroup.PATCH("/:id/status", adminHandler.UpdateClassStatus)
		}

		// 2.3 老师/家长用户管理
		userGroup := adminAPI.Group("/users")
		{
			userGroup.GET("", adminHandler.GetUsers)
			userGroup.POST("", adminHandler.CreateUser)
			userGroup.PUT("/:id", adminHandler.UpdateUser)
			userGroup.PATCH("/:id/status", adminHandler.UpdateUserStatus)
			userGroup.PATCH("/:id/reset-pwd", adminHandler.ResetUserPassword)
		}

		// 2.4 学生账号管理
		studentGroup := adminAPI.Group("/students")
		{
			studentGroup.POST("/batch-import", adminHandler.BatchImportStudents)
			studentGroup.GET("", adminHandler.GetStudents)
			studentGroup.POST("", adminHandler.CreateStudent)
			studentGroup.PUT("/:id", adminHandler.UpdateStudent)
		}

		// 2.5 老师班级权限分配
		assignmentGroup := adminAPI.Group("/assignments")
		{
			assignmentGroup.GET("", adminHandler.GetAssignments)
			assignmentGroup.POST("", adminHandler.CreateAssignment)
			assignmentGroup.DELETE("/:id", adminHandler.DeleteAssignment)
			assignmentGroup.POST("/batch", adminHandler.BatchCreateAssignments)
		}

		// 2.6 家长-学生绑定管理
		parentBindingGroup := adminAPI.Group("/parent-bindings")
		{
			parentBindingGroup.GET("", adminHandler.GetParentBindings)
			parentBindingGroup.POST("", adminHandler.CreateParentBinding)
			parentBindingGroup.DELETE("/:id", adminHandler.DeleteParentBinding)
		}

		// 2.7 数据概览
		adminAPI.GET("/dashboard", adminHandler.GetDashboard)
		adminAPI.GET("/audit-logs", adminHandler.GetAuditLogs)

		// 2.8 伙伴模板管理
		templateGroup := adminAPI.Group("/partner-templates")
		{
			templateGroup.POST("", adminHandler.CreatePartnerTemplate)
			templateGroup.PUT("/:id", adminHandler.UpdatePartnerTemplate)
			templateGroup.POST("/seed", adminHandler.SeedPartnerTemplates)
		}
	}
}
