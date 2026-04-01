// growth-partner/backend/internal/handler/class_handler.go

package handler

import (
	"github.com/gin-gonic/gin"
	"growth-partner/internal/middleware"
	"growth-partner/internal/service"
)

type ClassHandler struct {
	classSvc service.ClassService
}

func NewClassHandler(classSvc service.ClassService) *ClassHandler {
	return &ClassHandler{classSvc: classSvc}
}

// CreateClassReq 创建班级请求
type CreateClassReq struct {
	Name       string `json:"name" binding:"required,max=64"`
	Grade      int    `json:"grade" binding:"required,min=1,max=6"`
	SchoolYear string `json:"school_year" binding:"required"`
}

// CreateClass 创建班级 (教师/园长)
func (h *ClassHandler) CreateClass(c *gin.Context) {
	var req CreateClassReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ResponseValidationError(c, err.Error())
		return
	}

	teacherID := middleware.GetUserID(c)
	// TODO: 调用 classSvc.CreateClass
	middleware.ResponseOK(c, gin.H{"message": "班级创建成功（模拟）", "teacher_id": teacherID})
}

// GetMyClasses 获取当前老师负责的班级
func (h *ClassHandler) GetMyClasses(c *gin.Context) {
	teacherID := middleware.GetUserID(c)
	// TODO: 调用 classSvc.FindByTeacherID
	middleware.ResponseOK(c, gin.H{"teacher_id": teacherID, "classes": []interface{}{}})
}
