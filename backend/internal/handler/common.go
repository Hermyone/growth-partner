// growth-partner/backend/internal/handler/common.go
// Handler 层通用的辅助函数

package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// getPaginationParams 统一获取并校验分页参数
func getPaginationParams(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	if page < 1 {
		page = 1
	}
	// 限制单页最大数量，保护数据库性能
	if size < 1 || size > 100 {
		size = 20
	}
	return page, size
}
