package websocket

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

// WebSocketMessage 定义 WebSocket 消息结构
type WebSocketMessage struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

// TestWebSocketConnection 测试 WebSocket 连接
func TestWebSocketConnection(t *testing.T) {
	// 获取 Token
	token, err := getTestToken(t, "teacher001", "123456", "teacher")
	if err != nil {
		t.Fatalf("获取 Token 失败: %v", err)
	}

	// WebSocket 连接
	conn, err := connectWebSocket(token)
	if err != nil {
		t.Fatalf("WebSocket 连接失败: %v", err)
	}
	defer conn.Close()

	// 发送测试消息
	testMessage := WebSocketMessage{
		Type:    "test",
		Payload: json.RawMessage(`{"message": "Hello WebSocket!"}`),
	}

	if err := conn.WriteJSON(testMessage); err != nil {
		t.Fatalf("发送消息失败: %v", err)
	}

	// 接收响应
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	var response WebSocketMessage
	if err := conn.ReadJSON(&response); err != nil {
		t.Fatalf("接收消息失败: %v", err)
	}

	if response.Type != "test" {
		t.Errorf("期望消息类型 'test'，实际收到 '%s'", response.Type)
	}

	t.Log("WebSocket 连接测试通过")
}

// TestWebSocketConcurrentConnections 测试并发 WebSocket 连接
func TestWebSocketConcurrentConnections(t *testing.T) {
	const concurrentCount = 100

	// 获取 Token
	token, err := getTestToken(t, "teacher001", "123456", "teacher")
	if err != nil {
		t.Fatalf("获取 Token 失败: %v", err)
	}

	var wg sync.WaitGroup
	errorChan := make(chan error, concurrentCount)

	for i := 0; i < concurrentCount; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			conn, err := connectWebSocket(token)
			if err != nil {
				errorChan <- err
				return
			}
			defer conn.Close()

			// 发送测试消息
			testMessage := WebSocketMessage{
				Type:    "test",
				Payload: json.RawMessage(`{"id": ` + string(rune(id)) + `}`),
			}

			if err := conn.WriteJSON(testMessage); err != nil {
				errorChan <- err
				return
			}

			// 接收响应
			conn.SetReadDeadline(time.Now().Add(5 * time.Second))
			var response WebSocketMessage
			if err := conn.ReadJSON(&response); err != nil {
				errorChan <- err
				return
			}

			t.Logf("连接 %d 测试通过", id)
		}(i)
	}

	wg.Wait()
	close(errorChan)

	var errors []error
	for err := range errorChan {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		t.Errorf("有 %d 个连接失败: %v", len(errors), errors)
	} else {
		t.Logf("%d 个并发 WebSocket 连接测试通过", concurrentCount)
	}
}

// TestWebSocketHeartbeat 测试 WebSocket 心跳
func TestWebSocketHeartbeat(t *testing.T) {
	// 获取 Token
	token, err := getTestToken(t, "teacher001", "123456", "teacher")
	if err != nil {
		t.Fatalf("获取 Token 失败: %v", err)
	}

	// WebSocket 连接
	conn, err := connectWebSocket(token)
	if err != nil {
		t.Fatalf("WebSocket 连接失败: %v", err)
	}
	defer conn.Close()

	// 等待 45 秒，测试心跳机制
	t.Log("测试 WebSocket 心跳，等待 45 秒...")
	time.Sleep(45 * time.Second)

	// 尝试发送消息，验证连接是否仍然活跃
	testMessage := WebSocketMessage{
		Type:    "ping",
		Payload: json.RawMessage(`{}`),
	}

	if err := conn.WriteJSON(testMessage); err != nil {
		t.Fatalf("发送心跳消息失败: %v", err)
	}

	// 接收响应
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	var response WebSocketMessage
	if err := conn.ReadJSON(&response); err != nil {
		t.Fatalf("接收心跳响应失败: %v", err)
	}

	t.Log("WebSocket 心跳测试通过")
}

// 辅助函数：获取测试 Token
func getTestToken(t *testing.T, username, password, role string) (string, error) {
	// 实现获取 Token 的逻辑
	// 这里简化处理，实际应该调用登录接口
	return "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOjEsInVzciI6InRlYWNoZXIwMDEiLCJyb2wiOiJ0ZWFjaGVyIiwidHlwIjoiYWNjZXNzIiwiaXNzIjoiZ3Jvd3RoLXBhcnRuZXIiLCJzdWIiOiJ0ZWFjaGVyMDAxIiwiZXhwIjoxNzc1MjY3MzA0LCJpYXQiOjE3NzUxODA5MDQsImp0aSI6IjIwMjYwNDAzMDk0ODI0eVk1UEJtN1AifQ.HlHgUqAFTYhcaUy3P7s7ggom9Wd44gi2MWM4N0Hg9mM", nil
}

// 辅助函数：连接 WebSocket
func connectWebSocket(token string) (*websocket.Conn, error) {
	dialer := websocket.DefaultDialer
	headers := http.Header{}
	headers.Set("Authorization", "Bearer "+token)

	conn, _, err := dialer.Dial("ws://localhost:8080/api/v1/ws", headers)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
