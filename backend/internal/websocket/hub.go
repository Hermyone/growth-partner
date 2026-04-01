// growth-partner/backend/internal/websocket/hub.go
// 连接管理器：负责维护所有在线客户端，并根据 UserID/ClassID 进行消息路由

package websocket

import (
	"log"
	"sync"
)

// Hub 维护活跃连接并广播消息
type Hub struct {
	// 存储所有在线客户端，Key 为 UserID
	// 考虑到一个学生可能多端登录（如平板+手机），Value 为 Client 的 Slice
	clients map[uint64][]*Client

	// 注册通道
	register chan *Client

	// 注销通道
	unregister chan *Client

	// 互斥锁，保证并发安全
	mu sync.RWMutex
}

// NewHub 创建一个新的 Hub
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[uint64][]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Run 启动 Hub 的主循环（监听注册/注销事件）
func (h *Hub) Run() {
	log.Println("[WS] Hub 运行中...")
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.UserID] = append(h.clients[client.UserID], client)
			h.mu.Unlock()
			log.Printf("[WS] 用户 %d 已连接 (当前在线用户数: %d)", client.UserID, len(h.clients))

		case client := <-h.unregister:
			h.mu.Lock()
			if clients, ok := h.clients[client.UserID]; ok {
				// 从 Slice 中移除该 Client
				for i, c := range clients {
					if c == client {
						h.clients[client.UserID] = append(clients[:i], clients[i+1:]...)
						break
					}
				}
				// 如果该用户没有任何连接了，删除 Key
				if len(h.clients[client.UserID]) == 0 {
					delete(h.clients, client.UserID)
				}
				close(client.send) // 关闭客户端的发送通道
			}
			h.mu.Unlock()
			log.Printf("[WS] 用户 %d 已断开连接", client.UserID)
		}
	}
}

// SendToUser 将消息发送给特定用户
func (h *Hub) SendToUser(userID uint64, msg []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if clients, ok := h.clients[userID]; ok {
		for _, client := range clients {
			select {
			case client.send <- msg:
			default:
				// 如果发送通道阻塞，说明客户端由于网络问题接收慢，主动断开
				log.Printf("[WS] 用户 %d 连接拥塞，正在主动断开", userID)
				go func(c *Client) { h.unregister <- c }(client)
			}
		}
	}
}

// SendToClass 将消息广播给班级内的所有在线学生
// 注意：此逻辑通常配合 Redis 订阅使用，Hub 只需要遍历内存中的 client 匹配 ClassID
func (h *Hub) SendToClass(classID uint64, msg []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, userClients := range h.clients {
		for _, client := range userClients {
			if client.ClassID == classID {
				client.send <- msg
			}
		}
	}
}
