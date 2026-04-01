// growth-partner/backend/internal/websocket/message.go
// 消息协议定义：定义服务器与客户端之间交互的数据结构

package websocket

import "encoding/json"

// MessageType 消息类型枚举
type MessageType string

const (
	// TypePartnerMessage 伙伴对主人的私信（鼓励/进化通知）
	TypePartnerMessage MessageType = "partner_message"
	// TypeClassBroadcast 园长对全班的广播
	TypeClassBroadcast MessageType = "class_broadcast"
	// TypeSystemNotice 系统全局通知
	TypeSystemNotice MessageType = "system_notice"
	// TypeHeartbeat 心跳包（保持连接）
	TypeHeartbeat MessageType = "heartbeat"
)

// WSMessage 统一的 WebSocket 消息包装结构
type WSMessage struct {
	Type    MessageType     `json:"type"`              // 消息类型
	Payload json.RawMessage `json:"payload,omitempty"` // 具体的消息体内容
	SentAt  int64           `json:"sent_at"`           // 发送时间戳（毫秒）
}

// PartnerMessagePayload 伙伴私信的具体负载
type PartnerMessagePayload struct {
	Text      string `json:"text"`                   // 鼓励的话
	Animation string `json:"animation,omitempty"`    // 触发的前端 Lottie 动画 Key
	Growth    int    `json:"growth_added,omitempty"` // 本次增加的成长值
}

// ClassBroadcastPayload 班级广播的具体负载
type ClassBroadcastPayload struct {
	Content string `json:"content"`            // 广播文字
	Sender  string `json:"sender"`             // 发送者（如“成长乐园园长”）
	IsAlert bool   `json:"is_alert,omitempty"` // 是否需要弹窗强提示
}
