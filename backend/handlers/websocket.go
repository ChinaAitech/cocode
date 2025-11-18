package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"cocode/backend/models"
	"cocode/backend/services"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有来源（生产环境需要限制）
	},
}

var hub *services.CollaborationHub

// InitWebSocketHub 初始化WebSocket中心
func InitWebSocketHub() {
	hub = services.NewCollaborationHub()
	go hub.Run()
}

// HandleWebSocket 处理WebSocket连接
func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// 验证会话
	sessionID := r.URL.Query().Get("session")
	session, err := services.ValidateSession(sessionID)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// 升级为WebSocket连接
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket升级失败: %v", err)
		return
	}

	client := &services.Client{
		Username: session.Username,
		Conn:     conn,
		Send:     make(chan []byte, 256),
		Hub:      hub,
	}

	hub.RegisterClient(client)

	// 发送当前代码状态
	codeState := hub.GetCodeState()
	initMessage := models.WebSocketMessage{
		Type:      "init",
		Username:  "system",
		Timestamp: time.Now().Unix(),
		Data:      codeState,
	}
	initData, _ := json.Marshal(initMessage)
	client.Send <- initData

	// 广播用户加入
	joinMessage := models.WebSocketMessage{
		Type:      "user_join",
		Username:  session.Username,
		Timestamp: time.Now().Unix(),
		Data: map[string]interface{}{
			"username": session.Username,
			"users":    hub.GetOnlineUsers(),
		},
	}
	joinData, _ := json.Marshal(joinMessage)
	hub.BroadcastMessage(joinData)

	// 启动读写协程
	go writePump(client)
	go readPump(client, hub)
}

// readPump 读取客户端消息
func readPump(client *services.Client, hub *services.CollaborationHub) {
	defer func() {
		hub.UnregisterClient(client)
		client.Conn.Close()

		// 广播用户离开
		leaveMessage := models.WebSocketMessage{
			Type:      "user_leave",
			Username:  client.Username,
			Timestamp: time.Now().Unix(),
			Data: map[string]interface{}{
				"username": client.Username,
				"users":    hub.GetOnlineUsers(),
			},
		}
		leaveData, _ := json.Marshal(leaveMessage)
		hub.BroadcastMessage(leaveData)
	}()

	client.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	client.Conn.SetPongHandler(func(string) error {
		client.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket错误: %v", err)
			}
			break
		}

		// 解析消息
		var wsMsg models.WebSocketMessage
		if err := json.Unmarshal(message, &wsMsg); err != nil {
			log.Printf("消息解析失败: %v", err)
			continue
		}

		wsMsg.Username = client.Username
		wsMsg.Timestamp = time.Now().Unix()

		// 处理不同类型的消息
		switch wsMsg.Type {
		case "edit":
			// 更新代码状态
			if data, ok := wsMsg.Data.(map[string]interface{}); ok {
				if code, ok := data["code"].(string); ok {
					hub.UpdateCodeState(code)
				}
			}
		case "compile":
			// 编译请求（异步处理）
			go handleCompileRequest(client, wsMsg)
			continue // 不广播编译请求
		}

		// 广播消息给所有客户端
		broadcastData, _ := json.Marshal(wsMsg)
		hub.BroadcastMessage(broadcastData)
	}
}

// writePump 向客户端写入消息
func writePump(client *services.Client) {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		client.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-client.Send:
			client.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := client.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}

		case <-ticker.C:
			client.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// handleCompileRequest 处理编译请求
func handleCompileRequest(client *services.Client, msg models.WebSocketMessage) {
	data, ok := msg.Data.(map[string]interface{})
	if !ok {
		return
	}

	code, _ := data["code"].(string)
	input, _ := data["input"].(string)

	// 执行编译
	result := services.CompileAndRun(code, input)

	// 发送结果给请求者
	responseMsg := models.WebSocketMessage{
		Type:      "compile_result",
		Username:  "system",
		Timestamp: time.Now().Unix(),
		Data:      result,
	}

	responseData, _ := json.Marshal(responseMsg)
	client.Send <- responseData
}
