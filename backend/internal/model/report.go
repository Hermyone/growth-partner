// growth-partner/backend/internal/model/report.go
// 周报模型

package model

import "time"

// ReportType 报告类型
type ReportType string

const (
	ReportTypeWeekly  ReportType = "weekly"  // 周报
	ReportTypeMonthly ReportType = "monthly" // 月报
	ReportTypeAnnual  ReportType = "annual"  // 年报
)

// ReportStatus 报告状态
type ReportStatus string

const (
	ReportStatusGenerating ReportStatus = "generating" // 生成中
	ReportStatusReady      ReportStatus = "ready"      // 已生成
	ReportStatusFailed     ReportStatus = "failed"     // 生成失败
)

// Report 班级报告表
type Report struct {
	Base

	ClassID        uint64        `gorm:"index;not null" json:"class_id"`          // 班级ID
	ReportType     ReportType    `gorm:"size:20;not null" json:"report_type"`     // 报告类型
	ReportPeriod   string        `gorm:"size:50;not null" json:"report_period"`   // 报告周期（如：2024-03-18 至 2024-03-24）
	StartDate      time.Time     `gorm:"not null" json:"start_date"`              // 报告开始日期
	EndDate        time.Time     `gorm:"not null" json:"end_date"`                // 报告结束日期
	Status         ReportStatus  `gorm:"size:20;default:generating" json:"status"` // 报告状态
	FileURL        string        `gorm:"size:512" json:"file_url,omitempty"`       // 报告文件URL
	FileSize       int64         `json:"file_size,omitempty"`                        // 报告文件大小（字节）
	GeneratedBy    uint64        `gorm:"index;not null" json:"generated_by"`        // 生成人ID（老师）
	GeneratedAt    *time.Time    `json:"generated_at,omitempty"`                  // 生成完成时间
	DataSummary    string        `gorm:"type:text" json:"data_summary,omitempty"`    // 数据摘要（JSON格式存储）
	ErrorMsg       string        `gorm:"size:512" json:"error_msg,omitempty"`       // 生成失败时的错误信息
}

func (Report) TableName() string { return "reports" }
