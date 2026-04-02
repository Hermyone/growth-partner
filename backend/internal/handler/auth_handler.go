// growth-partner/backend/internal/handler/auth_handler.go

package handler

import (
	"net/http"

	"growth-partner/internal/middleware"
	"growth-partner/internal/service"

	"github.com/gin-gonic/gin"
)

// AuthHandler 鉴权控制器
type AuthHandler struct {
	authSvc service.AuthService
}

// NewAuthHandler 创建鉴权控制器实例
func NewAuthHandler(authSvc service.AuthService) *AuthHandler {
	return &AuthHandler{
		authSvc: authSvc,
	}
}

// LoginReq 登录请求参数
type LoginReq struct {
	Username string `json:"username" binding:"required,min=3,max=64"`
	Password string `json:"password" binding:"required,min=6,max=32"`
	Role     string `json:"role" binding:"required,oneof=student teacher parent admin"`
}

// Login 用户登录
// POST /api/v1/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResponseValidationError(c, err.Error())
		return
	}

	// 调用鉴权服务进行登录校验
	resp, err := h.authSvc.Login(c.Request.Context(), req.Username, req.Password, req.Role)
	if err != nil {
		// 统一返回 401 错误，不暴露具体的账号或密码错误细节
		middleware.ResponseError(c, http.StatusUnauthorized, "LOGIN_FAILED", "账号或密码错误")
		return
	}

	middleware.ResponseOKWithMessage(c, "登录成功", resp)
}

// RefreshTokenReq 刷新 Token 请求参数
type RefreshTokenReq struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// RefreshToken 刷新访问令牌
// POST /api/v1/auth/refresh
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req RefreshTokenReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResponseValidationError(c, err.Error())
		return
	}

	resp, err := h.authSvc.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		middleware.ResponseError(c, http.StatusUnauthorized, "REFRESH_FAILED", "刷新令牌无效或已过期，请重新登录")
		return
	}

	middleware.ResponseOK(c, resp)
}

// LogoutReq 登出请求参数
type LogoutReq struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// Logout 用户登出
// POST /api/v1/auth/logout
func (h *AuthHandler) Logout(c *gin.Context) {
	var req LogoutReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResponseValidationError(c, err.Error())
		return
	}

	if err := h.authSvc.Logout(c.Request.Context(), req.RefreshToken); err != nil {
		middleware.ResponseError(c, http.StatusInternalServerError, "LOGOUT_FAILED", "登出失败")
		return
	}

	middleware.ResponseOKWithMessage(c, "登出成功", nil)
}

// Me 获取当前登录用户信息
// GET /api/v1/auth/me
func (h *AuthHandler) Me(c *gin.Context) {
	// 从上下文获取用户ID（由auth中间件设置）
	userID, exists := c.Get("user_id")
	if !exists {
		middleware.ResponseError(c, http.StatusUnauthorized, "UNAUTHORIZED", "未授权访问")
		return
	}

	user, err := h.authSvc.GetCurrentUser(c.Request.Context(), userID.(uint64))
	if err != nil {
		middleware.ResponseError(c, http.StatusInternalServerError, "GET_USER_FAILED", "获取用户信息失败")
		return
	}

	middleware.ResponseOK(c, user)
}

// ChangePasswordReq 修改密码请求参数
type ChangePasswordReq struct {
	OldPassword string `json:"old_password" binding:"required,min=6,max=32"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=32"`
}

// ChangePassword 修改密码
// PATCH /api/v1/auth/password
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req ChangePasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResponseValidationError(c, err.Error())
		return
	}

	// 从上下文获取用户ID（由auth中间件设置）
	userID, exists := c.Get("user_id")
	if !exists {
		middleware.ResponseError(c, http.StatusUnauthorized, "UNAUTHORIZED", "未授权访问")
		return
	}

	if err := h.authSvc.ChangePassword(c.Request.Context(), userID.(uint64), req.OldPassword, req.NewPassword); err != nil {
		if err == service.ErrInvalidPassword {
			middleware.ResponseError(c, http.StatusBadRequest, "INVALID_OLD_PASSWORD", "旧密码错误")
			return
		}
		middleware.ResponseError(c, http.StatusInternalServerError, "CHANGE_PASSWORD_FAILED", "修改密码失败")
		return
	}

	middleware.ResponseOKWithMessage(c, "密码修改成功", nil)
}
