// growth-partner/backend/internal/websocket/client.go
// 客户端处理：负责单个 WebSocket 连接的读写操作及心跳维护

package websocket

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// 写消息超时
	writeWait = 10 * time.Second
	// 读下一个 Pong 消息的超时时间
	pongWait = 60 * time.Second
	// 发送 Ping 的频率
	pingPeriod = (pongWait * 9) / 10
	// 允许的最大消息大小
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// 跨域检查（由于使用了 Token 认证，此处可根据环境放行）
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Client 是 WebSocket 连接与 Hub 之间的媒介
type Client struct {
	Hub *Hub

	// 实际的网络连接
	Conn *websocket.Conn

	// 业务属性
	UserID  uint64
	ClassID uint64

	// 待发送的消息通道
	send chan []byte
}

// NewClient 升级 HTTP 连接并创建 Client
func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request, userID, classID uint64) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[WS] 升级协议失败: %v", err)
		return
	}

	client := &Client{
		Hub:     hub,
		Conn:    conn,
		UserID:  userID,
		ClassID: classID,
		send:    make(chan []byte, 256),
	}

	client.Hub.register <- client

	// 启动读写处理
	go client.writePump()
	go client.readPump()
}

// readPump 从 WebSocket 连接读取消息并处理（主要处理心跳和异常断开）
func (c *Client) readPump() {
	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, _, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("[WS] 读取消息异常: %v", err)
			}
			break
		}
		// 目前项目主要由服务端推送，客户端上行消息暂不处理（预留心跳由 Pong 处理）
	}
}

// writePump 将消息从 Hub 泵送到 WebSocket 连接
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// Hub 关闭了通道
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.BinaryMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// 检查队列中是否还有积压消息，一并发送
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte("\n"))
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
