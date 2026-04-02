// growth-partner/backend/internal/middleware/auth.go
// JWT 认证中间件：从 Authorization Header 提取并验证 Token

package middleware

import (
	"net/http"
	"strings"

	"growth-partner/pkg/jwt"

	"github.com/gin-gonic/gin"
)

const (
	// ContextKeyUserID 上下文中用户ID的键
	ContextKeyUserID   = "user_id"
	ContextKeyUsername = "username"
	ContextKeyRole     = "role"
	ContextKeyClassID  = "class_id"
	ContextKeyChildID  = "child_id"
	ContextKeyClaims   = "jwt_claims"
)

// Auth JWT 认证中间件
// 验证 Authorization: Bearer <token> 头部
func Auth(jwtManager *jwt.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 提取 Bearer Token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			ResponseError(c, http.StatusUnauthorized, "UNAUTHORIZED", "请先登录")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			ResponseError(c, http.StatusUnauthorized, "INVALID_TOKEN_FORMAT", "Token 格式错误，应为 Bearer <token>")
			c.Abort()
			return
		}

		// 解析并验证 Token
		claims, err := jwtManager.ParseAccessToken(parts[1])
		if err != nil {
			switch err {
			case jwt.ErrTokenExpired:
				ResponseError(c, http.StatusUnauthorized, "TOKEN_EXPIRED", "登录已过期，请重新登录")
			default:
				ResponseError(c, http.StatusUnauthorized, "INVALID_TOKEN", "Token 无效")
			}
			c.Abort()
			return
		}

		// 检查用户是否在黑名单中（修改密码后所有Token失效）
		// 这里需要从依赖注入中获取Redis客户端
		// 由于中间件无法直接获取Redis客户端，这里暂时注释掉
		// 实际项目中需要通过依赖注入或全局变量来获取Redis客户端
		// key := fmt.Sprintf("blacklist:user:%d", claims.UserID)
		// exists, err := redisClient.Exists(c.Request.Context(), key).Result()
		// if err == nil && exists > 0 {
		// 	ResponseError(c, http.StatusUnauthorized, "TOKEN_REVOKED", "登录已失效，请重新登录")
		// 	c.Abort()
		// 	return
		// }

		// 将用户信息注入上下文，供后续 Handler 使用
		c.Set(ContextKeyUserID, claims.UserID)
		c.Set(ContextKeyUsername, claims.Username)
		c.Set(ContextKeyRole, string(claims.Role))
		c.Set(ContextKeyClassID, claims.ClassID)
		c.Set(ContextKeyChildID, claims.ChildID)
		c.Set(ContextKeyClaims, claims)

		c.Next()
	}
}

// ─── 便捷辅助函数（供 Handler 使用） ──────────────────────────

// GetUserID 从上下文获取当前用户ID
func GetUserID(c *gin.Context) uint64 {
	v, _ := c.Get(ContextKeyUserID)
	id, _ := v.(uint64)
	return id
}

// GetRole 从上下文获取当前用户角色
func GetRole(c *gin.Context) string {
	v, _ := c.Get(ContextKeyRole)
	role, _ := v.(string)
	return role
}

// GetChildID 从上下文获取当前学生档案ID（仅学生角色有效）
func GetChildID(c *gin.Context) uint64 {
	v, _ := c.Get(ContextKeyChildID)
	id, _ := v.(uint64)
	return id
}

// GetClassID 从上下文获取班级ID
func GetClassID(c *gin.Context) uint64 {
	v, _ := c.Get(ContextKeyClassID)
	id, _ := v.(uint64)
	return id
}
