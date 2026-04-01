// growth-partner/backend/internal/middleware/logger.go
// 日志与追踪中间件：实现请求链路追踪 (RequestID) 与耗时监控

package middleware

import (
	"crypto/rand"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// RequestID 为每个请求注入唯一追踪码
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查 Header 中是否已有 Request-ID
		reqID := c.GetHeader("X-Request-Id")
		if reqID == "" {
			// 如果没有，生成一个简单的随机字符串作为 ID
			// 实际项目建议使用 uuid: github.com/google/uuid
			reqID = generateSimpleID()
		}

		// 存入上下文，供后续 Handler 或日志取出
		c.Set("request_id", reqID)

		// 写入响应 Header，方便前端在控制台查看或上报
		c.Header("X-Request-Id", reqID)

		c.Next()
	}
}

// RequestLogger 记录每个 API 的请求流水、状态码及耗时
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// 继续执行后续逻辑
		c.Next()

		// 记录结束信息
		endTime := time.Now()
		latency := endTime.Sub(startTime)

		if raw != "" {
			path = path + "?" + raw
		}

		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		reqID := c.GetString("request_id")

		// 耗时格式化
		latencyDisplay := fmt.Sprintf("%v", latency)
		if latency > time.Second {
			latencyDisplay = fmt.Sprintf("%.2fs", float64(latency)/float64(time.Second))
		}

		// 标准日志输出
		log.Printf("[API] %s | %3d | %10s | %15s | %-7s %s",
			reqID,
			statusCode,
			latencyDisplay,
			clientIP,
			method,
			path,
		)
	}
}

// generateSimpleID 生成一个简单的 8 位随机 ID
func generateSimpleID() string {
	b := make([]byte, 4)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
