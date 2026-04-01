// growth-partner/backend/internal/handler/child_handler.go

package handler

import (
	"github.com/gin-gonic/gin"
	"growth-partner/internal/middleware"
	"growth-partner/internal/repository"
	"strconv"
)

type ChildHandler struct {
	childRepo repository.ChildRepository
}

func NewChildHandler(childRepo repository.ChildRepository) *ChildHandler {
	return &ChildHandler{childRepo: childRepo}
}

// GetClassStudents 获取班级内所有学生概览
func (h *ChildHandler) GetClassStudents(c *gin.Context) {
	classIDStr := c.Param("classId")
	classID, _ := strconv.ParseUint(classIDStr, 10, 64)

	// TODO: 获取学生列表并脱敏处理
	middleware.ResponseOK(c, gin.H{"class_id": classID, "students": []interface{}{}})
}

// GetMyChildren 获取家长绑定的孩子档案
func (h *ChildHandler) GetMyChildren(c *gin.Context) {
	parentUserID := middleware.GetUserID(c)
	// TODO: 查询 parent_child_relations 表
	middleware.ResponseOK(c, gin.H{"parent_id": parentUserID, "children": []interface{}{}})
}
