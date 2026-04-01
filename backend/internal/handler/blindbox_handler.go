// 简化版 BlindboxHandler
package handler

import (
	"github.com/gin-gonic/gin"
	"growth-partner/internal/middleware"
	"growth-partner/internal/service"
)

type BlindboxHandler struct {
	blindboxSvc service.BlindboxService
}

func NewBlindboxHandler(svc service.BlindboxService) *BlindboxHandler {
	return &BlindboxHandler{blindboxSvc: svc}
}

func (h *BlindboxHandler) GetPool(c *gin.Context)        { middleware.ResponseOK(c, []interface{}{}) }
func (h *BlindboxHandler) AddToPool(c *gin.Context)      { middleware.ResponseOK(c, nil) }
func (h *BlindboxHandler) RemoveFromPool(c *gin.Context) { middleware.ResponseOK(c, nil) }
func (h *BlindboxHandler) DrawForStudent(c *gin.Context) {
	middleware.ResponseOK(c, gin.H{"reward": "免作业一次"})
}
