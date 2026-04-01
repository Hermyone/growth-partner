// growth-partner/backend/internal/middleware/cors.go
// 跨域中间件：支持开发环境全开放，生产环境需配置白名单

package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CORS 跨域配置中间件
func CORS(isDev bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			// 如果是开发环境，允许所有来源
			if isDev {
				c.Header("Access-Control-Allow-Origin", origin)
			} else {
				// 生产环境建议将此处改为具体的域名白名单
				c.Header("Access-Control-Allow-Origin", origin)
			}

			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE, PATCH")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization, X-Request-Id")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type, X-Request-Id")
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		// 处理预检请求 (OPTIONS)
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
