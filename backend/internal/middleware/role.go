// growth-partner/backend/internal/middleware/role.go
// 角色权限中间件：基于 RBAC 的接口访问控制
// 三个主角色：学生（student）、老师/园长（teacher）、家长（parent）

package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RequireRoles 要求请求者拥有指定角色之一
// 用法：router.GET("/teacher/xxx", middleware.Auth(jwtMgr), middleware.RequireRoles("teacher","admin"), handler)
func RequireRoles(roles ...string) gin.HandlerFunc {
	// 将允许的角色存入 map，O(1) 查找
	allowedRoles := make(map[string]struct{}, len(roles))
	for _, r := range roles {
		allowedRoles[r] = struct{}{}
	}

	return func(c *gin.Context) {
		role := GetRole(c)
		if role == "" {
			// 没有角色信息说明未认证（Auth 中间件应先运行）
			ResponseError(c, http.StatusUnauthorized, "UNAUTHORIZED", "请先登录")
			c.Abort()
			return
		}

		if _, ok := allowedRoles[role]; !ok {
			ResponseError(c, http.StatusForbidden, "FORBIDDEN",
				"您没有权限执行此操作（需要角色: "+joinRoles(roles)+"）")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireTeacher 快捷方法：仅教师和管理员可访问
func RequireTeacher() gin.HandlerFunc {
	return RequireRoles("teacher", "admin")
}

// RequireParent 快捷方法：仅家长可访问
func RequireParent() gin.HandlerFunc {
	return RequireRoles("parent", "admin")
}

// RequireStudent 快捷方法：仅学生可访问
func RequireStudent() gin.HandlerFunc {
	return RequireRoles("student", "admin")
}

// RequireTeacherOrParent 快捷方法：教师或家长
func RequireTeacherOrParent() gin.HandlerFunc {
	return RequireRoles("teacher", "parent", "admin")
}

// joinRoles 将角色列表拼接成字符串用于错误提示
func joinRoles(roles []string) string {
	result := ""
	for i, r := range roles {
		if i > 0 {
			result += " 或 "
		}
		result += r
	}
	return result
}
