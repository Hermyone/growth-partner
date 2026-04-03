package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// WebSocketMessage 定义 WebSocket 消息结构
type WebSocketMessage struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

func main() {
	// 教师 Token
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOjEsInVzciI6InRlYWNoZXIwMDEiLCJyb2wiOiJ0ZWFjaGVyIiwidHlwIjoiYWNjZXNzIiwiaXNzIjoiZ3Jvd3RoLXBhcnRuZXIiLCJzdWIiOiJ0ZWFjaGVyMDAxIiwiZXhwIjoxNzc1MjcxMzk1LCJpYXQiOjE3NzUxODQ5OTUsImp0aSI6IjIwMjYwNDAzMTA1NjM1S3U4NEtvMksifQ.BfcGR4IVfW5ERnAmLeVI4ljRMa9lQIP1vRVBZfthjHk"

	// 测试 1: 基本连接测试
	log.Println("=== 测试 1: 基本连接测试 ===")
	testBasicConnection(token)

	// 测试 2: 心跳测试
	log.Println("\n=== 测试 2: 心跳测试 ===")
	testHeartbeat(token)

	// 测试 3: 并发连接测试
	log.Println("\n=== 测试 3: 并发连接测试 ===")
	testConcurrentConnections(token, 100)

	log.Println("\n=== WebSocket 测试完成 ===")
}

// 测试基本连接
func testBasicConnection(token string) {
	dialer := websocket.DefaultDialer
	headers := http.Header{}
	headers.Set("Authorization", "Bearer "+token)

	conn, _, err := dialer.Dial("ws://localhost:8080/api/v1/ws", headers)
	if err != nil {
		log.Fatalf("连接失败: %v", err)
	}
	defer conn.Close()

	// 发送测试消息
	testMessage := WebSocketMessage{
		Type:    "test",
		Payload: json.RawMessage(`{"message": "Hello WebSocket!"}`),
	}

	if err := conn.WriteJSON(testMessage); err != nil {
		log.Fatalf("发送消息失败: %v", err)
	}

	// 接收响应
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	var response WebSocketMessage
	if err := conn.ReadJSON(&response); err != nil {
		log.Fatalf("接收消息失败: %v", err)
	}

	log.Printf("✅ 基本连接测试通过，收到响应: %+v", response)
}

// 测试心跳
func testHeartbeat(token string) {
	dialer := websocket.DefaultDialer
	headers := http.Header{}
	headers.Set("Authorization", "Bearer "+token)

	conn, _, err := dialer.Dial("ws://localhost:8080/api/v1/ws", headers)
	if err != nil {
		log.Fatalf("连接失败: %v", err)
	}
	defer conn.Close()

	log.Println("等待 45 秒，测试心跳机制...")
	time.Sleep(45 * time.Second)

	// 尝试发送消息，验证连接是否仍然活跃
	testMessage := WebSocketMessage{
		Type:    "ping",
		Payload: json.RawMessage(`{}`),
	}

	if err := conn.WriteJSON(testMessage); err != nil {
		log.Printf("❌ 心跳测试失败: %v", err)
		return
	}

	// 接收响应
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	var response WebSocketMessage
	if err := conn.ReadJSON(&response); err != nil {
		log.Printf("❌ 心跳测试失败: %v", err)
		return
	}

	log.Println("✅ 心跳测试通过，连接保持活跃")
}

// 测试并发连接
func testConcurrentConnections(token string, count int) {
	var wg sync.WaitGroup
	errorChan := make(chan error, count)

	for i := 0; i < count; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			dialer := websocket.DefaultDialer
			headers := http.Header{}
			headers.Set("Authorization", "Bearer "+token)

			conn, _, err := dialer.Dial("ws://localhost:8080/api/v1/ws", headers)
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

			if id%10 == 0 {
				log.Printf("✅ 连接 %d 测试通过", id)
			}
		}(i)
	}

	wg.Wait()
	close(errorChan)

	var errors []error
	for err := range errorChan {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		log.Printf("❌ 并发测试失败，%d 个连接失败: %v", len(errors), errors)
	} else {
		log.Printf("✅ 并发测试通过，%d 个连接全部成功", count)
	}
}
