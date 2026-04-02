// growth-partner/backend/internal/middleware/response.go
// 统一 API 响应格式：所有接口均使用此格式返回，保证前端体验一致

package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// ─── 统一响应结构体 ────────────────────────────────────────────

// Response 统一 API 响应格式
type Response struct {
	Code      int         `json:"code"`                 // 业务状态码（0=成功，非0=业务错误）
	Message   string      `json:"message"`              // 人类可读的消息
	Data      interface{} `json:"data,omitempty"`       // 响应数据（成功时有值）
	Error     *ErrorInfo  `json:"error,omitempty"`      // 错误详情（失败时有值）
	Timestamp int64       `json:"timestamp"`            // 服务器时间戳（毫秒）
	RequestID string      `json:"request_id,omitempty"` // 请求追踪ID
}

// ErrorInfo 错误详情结构体
type ErrorInfo struct {
	ErrorCode string      `json:"error_code"`        // 机器可读的错误码，如 "INVALID_TOKEN"
	Details   interface{} `json:"details,omitempty"` // 校验错误的字段详情
}

// ─── 便捷响应函数 ──────────────────────────────────────────────

// ResponseOK 返回成功响应（HTTP 200）
func ResponseOK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:      0,
		Message:   "success",
		Data:      data,
		Timestamp: time.Now().UnixMilli(),
		RequestID: c.GetString("request_id"),
	})
}

// ResponseSuccess 返回成功响应（HTTP 200）
func ResponseSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:      0,
		Message:   "success",
		Data:      data,
		Timestamp: time.Now().UnixMilli(),
		RequestID: c.GetString("request_id"),
	})
}

// ResponseOKWithMessage 返回带自定义消息的成功响应
func ResponseOKWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:      0,
		Message:   message,
		Data:      data,
		Timestamp: time.Now().UnixMilli(),
		RequestID: c.GetString("request_id"),
	})
}

// ResponseError 返回业务错误响应
func ResponseError(c *gin.Context, httpStatus int, errorCode, message string) {
	c.JSON(httpStatus, Response{
		Code:    httpStatus,
		Message: message,
		Error: &ErrorInfo{
			ErrorCode: errorCode,
		},
		Timestamp: time.Now().UnixMilli(),
		RequestID: c.GetString("request_id"),
	})
}

// ResponseValidationError 返回参数校验失败响应（HTTP 422）
func ResponseValidationError(c *gin.Context, details interface{}) {
	c.JSON(http.StatusUnprocessableEntity, Response{
		Code:    422,
		Message: "请求参数校验失败",
		Error: &ErrorInfo{
			ErrorCode: "VALIDATION_ERROR",
			Details:   details,
		},
		Timestamp: time.Now().UnixMilli(),
		RequestID: c.GetString("request_id"),
	})
}

// ResponseInternalError 返回服务器内部错误（HTTP 500）
// 生产环境不暴露内部错误详情
func ResponseInternalError(c *gin.Context, isDev bool, err error) {
	resp := Response{
		Code:      http.StatusInternalServerError,
		Message:   "服务器内部错误，请稍后再试",
		Timestamp: time.Now().UnixMilli(),
		RequestID: c.GetString("request_id"),
	}
	// 开发环境暴露错误详情，方便调试
	if isDev && err != nil {
		resp.Error = &ErrorInfo{
			ErrorCode: "INTERNAL_ERROR",
			Details:   err.Error(),
		}
	}
	c.JSON(http.StatusInternalServerError, resp)
}
