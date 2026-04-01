// growth-partner/backend/internal/service/broadcast_service.go

package service

import (
	"context"
	"encoding/json"
	"fmt"
	"growth-partner/internal/model"

	"github.com/redis/go-redis/v9"
)

type BroadcastService interface {
	// 推送伙伴鼓励消息（针对个人）
	BroadcastPartnerMessage(ctx context.Context, childID uint64, message string) error
	// 推送进化通知（针对个人，带特效指令）
	BroadcastEvolution(ctx context.Context, childID, partnerID uint64, from, to model.EvolutionStage) error
	// 推送班级广播（针对全班）
	BroadcastToClass(ctx context.Context, classID uint64, content string) error
}

type broadcastServiceImpl struct {
	rdb *redis.Client
}

func NewBroadcastService(rdb *redis.Client) BroadcastService {
	return &broadcastServiceImpl{rdb: rdb}
}

// 消息协议结构
type WSMessage struct {
	Type    string      `json:"type"` // "partner_msg", "evolution", "class_broadcast"
	Payload interface{} `json:"payload"`
}

func (s *broadcastServiceImpl) BroadcastPartnerMessage(ctx context.Context, childID uint64, message string) error {
	channel := fmt.Sprintf("ws:child:%d", childID)
	msg := WSMessage{
		Type:    "partner_msg",
		Payload: map[string]string{"text": message},
	}
	data, _ := json.Marshal(msg)
	return s.rdb.Publish(ctx, channel, data).Err()
}

func (s *broadcastServiceImpl) BroadcastEvolution(ctx context.Context, childID, partnerID uint64, from, to model.EvolutionStage) error {
	channel := fmt.Sprintf("ws:child:%d", childID)
	msg := WSMessage{
		Type: "evolution",
		Payload: map[string]interface{}{
			"partner_id": partnerID,
			"from_stage": from,
			"to_stage":   to,
			"animation":  "evolution_effect_lottie", // 前端触发 Lottie 的 Key
		},
	}
	data, _ := json.Marshal(msg)
	return s.rdb.Publish(ctx, channel, data).Err()
}

func (s *broadcastServiceImpl) BroadcastToClass(ctx context.Context, classID uint64, content string) error {
	channel := fmt.Sprintf("ws:class:%d", classID)
	msg := WSMessage{
		Type:    "class_broadcast",
		Payload: map[string]string{"content": content},
	}
	data, _ := json.Marshal(msg)
	return s.rdb.Publish(ctx, channel, data).Err()
}
