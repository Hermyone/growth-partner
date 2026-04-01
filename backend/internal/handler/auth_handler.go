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
