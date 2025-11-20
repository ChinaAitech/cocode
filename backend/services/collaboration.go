package services

import (
	"sync"
	"time"

	"cocode/backend/models"

	"github.com/gorilla/websocket"
)

// Client WebSocket客户端
type Client struct {
	Username string
	Conn     *websocket.Conn
	Send     chan []byte
	Hub      *CollaborationHub
}

// CollaborationHub 协同编辑中心
type CollaborationHub struct {
	clients        map[*Client]bool
	broadcast      chan []byte
	register       chan *Client
	unregister     chan *Client
	codeState      *models.CodeState
	sharedState    *models.SharedState
	compileRecords []models.CompileRecord
	mu             sync.RWMutex
}

// NewCollaborationHub 创建新的协同中心
func NewCollaborationHub() *CollaborationHub {
	return &CollaborationHub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		codeState: &models.CodeState{
			Code:    "// 欢迎使用协同编程平台\n// 在此编写你的C++代码\n\n#include <bits/stdc++.h>\nusing namespace std;\n\nint main() {\n    int n;\n    cin >> n;\n    cout << n*n << '\\n';\n    return 0;\n}\n",
			Version: 0,
			Updated: time.Now(),
		},
		sharedState: &models.SharedState{
			InputData:  "",
			OutputData: "",
			CompileLog: "等待编译...\n",
			Answer:     "",
			Updated:    time.Now(),
		},
		compileRecords: make([]models.CompileRecord, 0),
	}
}

// Run 运行协同中心
func (h *CollaborationHub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.Send)
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.clients, client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// GetCodeState 获取当前代码状态
func (h *CollaborationHub) GetCodeState() *models.CodeState {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.codeState
}

// UpdateCodeState 更新代码状态
func (h *CollaborationHub) UpdateCodeState(code string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.codeState.Code = code
	h.codeState.Version++
	h.codeState.Updated = time.Now()
}

// BroadcastMessage 广播消息
func (h *CollaborationHub) BroadcastMessage(message []byte) {
	h.broadcast <- message
}

// GetOnlineUsers 获取在线用户列表
func (h *CollaborationHub) GetOnlineUsers() []string {
	h.mu.RLock()
	defer h.mu.RUnlock()

	users := make([]string, 0, len(h.clients))
	for client := range h.clients {
		users = append(users, client.Username)
	}
	return users
}

// RegisterClient 注册客户端
func (h *CollaborationHub) RegisterClient(client *Client) {
	h.register <- client
}

// UnregisterClient 注销客户端
func (h *CollaborationHub) UnregisterClient(client *Client) {
	h.unregister <- client
}

// GetSharedState 获取共享状态
func (h *CollaborationHub) GetSharedState() *models.SharedState {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.sharedState
}

// UpdateInputData 更新输入数据
func (h *CollaborationHub) UpdateInputData(input string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.sharedState.InputData = input
	h.sharedState.Updated = time.Now()
}

// UpdateOutputData 更新输出数据
func (h *CollaborationHub) UpdateOutputData(output string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.sharedState.OutputData = output
	h.sharedState.Updated = time.Now()
}

// UpdateCompileLog 更新编译日志
func (h *CollaborationHub) UpdateCompileLog(log string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.sharedState.CompileLog = log
	h.sharedState.Updated = time.Now()
}

// UpdateAnswer 更新标准答案
func (h *CollaborationHub) UpdateAnswer(answer string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.sharedState.Answer = answer
	h.sharedState.Updated = time.Now()
}

// AddCompileRecord 添加编译记录
func (h *CollaborationHub) AddCompileRecord(username string, success bool) {
	h.mu.Lock()
	defer h.mu.Unlock()
	record := models.CompileRecord{
		Username:  username,
		Timestamp: time.Now(),
		Success:   success,
	}
	h.compileRecords = append(h.compileRecords, record)

	// 只保留最近10条记录
	if len(h.compileRecords) > 10 {
		h.compileRecords = h.compileRecords[len(h.compileRecords)-10:]
	}
}

// GetCompileRecords 获取编译记录
func (h *CollaborationHub) GetCompileRecords() []models.CompileRecord {
	h.mu.RLock()
	defer h.mu.RUnlock()
	records := make([]models.CompileRecord, len(h.compileRecords))
	copy(records, h.compileRecords)
	return records
}

// KickUser 踢出用户
func (h *CollaborationHub) KickUser(username string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for client := range h.clients {
		if client.Username == username {
			close(client.Send)
			delete(h.clients, client)
			client.Conn.Close()
		}
	}
}
