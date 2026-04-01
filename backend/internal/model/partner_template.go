// growth-partner/backend/internal/model/partner_template.go
// 伙伴模板：30 个预设形象，由管理员初始化导入

package model

import "database/sql/driver"
import "encoding/json"
import "fmt"

// PartnerType 伙伴大类
type PartnerType string

const (
	PartnerTypePet   PartnerType = "pet"   // 宠物
	PartnerTypeAnime PartnerType = "anime" // 二次元纸片人
	PartnerTypePlant PartnerType = "plant" // 植物
)

// EvolutionStage 进化阶段
type EvolutionStage int

const (
	StageLow  EvolutionStage = 1 // 低阶（0-99 成长值）
	StageMid  EvolutionStage = 2 // 中阶（100-299）
	StageHigh EvolutionStage = 3 // 高阶（300+）
)

// StageThresholds 进化阈值配置
var StageThresholds = map[EvolutionStage]int{
	StageLow:  0,
	StageMid:  100,
	StageHigh: 300,
}

// StringSlice 用于存储 JSON 数组到数据库
type StringSlice []string

func (s StringSlice) Value() (driver.Value, error) {
	b, err := json.Marshal(s)
	return string(b), err
}

func (s *StringSlice) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan StringSlice: %v", value)
	}
	return json.Unmarshal(bytes, s)
}

// PartnerTemplate 伙伴模板表（30 个预设形象）
type PartnerTemplate struct {
	Base

	// ─── 基础信息 ───
	Code        string      `gorm:"uniqueIndex;size:32;not null" json:"code"`        // 唯一标识，如 "pet_001"
	Name        string      `gorm:"size:32;not null"             json:"name"`        // 中文名，如 "星辉小奶狗"
	Type        PartnerType `gorm:"size:16;index;not null"       json:"type"`        // 大类
	Description string      `gorm:"size:128;not null"            json:"description"` // 一句话描述
	Slogan      string      `gorm:"size:128"                     json:"slogan"`      // 伙伴口号

	// ─── 形象资源（各阶段） ───
	// 实际使用时存放 Lottie JSON 文件路径或 CDN URL
	LowStageAsset  string `gorm:"size:512" json:"low_stage_asset"`  // 低阶形象资源路径
	MidStageAsset  string `gorm:"size:512" json:"mid_stage_asset"`  // 中阶形象资源路径
	HighStageAsset string `gorm:"size:512" json:"high_stage_asset"` // 高阶形象资源路径

	// ─── 预置鼓励语（按维度，JSONB存储） ───
	// 格式：{"virtue":["...","..."], "study":["..."], ...}
	EncourageMessages JSON `gorm:"type:jsonb;default:'{}'" json:"encourage_messages"`

	// ─── 推荐年龄段 ───
	RecommendGrades StringSlice `gorm:"type:jsonb" json:"recommend_grades"` // 如 ["1","2","3"]

	SortOrder int  `gorm:"default:0;index" json:"sort_order"` // 展示排序
	IsActive  bool `gorm:"default:true"    json:"is_active"`
}

// JSON 类型用于 JSONB 字段
type JSON map[string]interface{}

func (j JSON) Value() (driver.Value, error) {
	b, err := json.Marshal(j)
	return string(b), err
}

func (j *JSON) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan JSON: %v", value)
	}
	return json.Unmarshal(bytes, j)
}

func (PartnerTemplate) TableName() string { return "partner_templates" }
