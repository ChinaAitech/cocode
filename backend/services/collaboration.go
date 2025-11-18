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
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	codeState  *models.CodeState
	mu         sync.RWMutex
}

// NewCollaborationHub 创建新的协同中心
func NewCollaborationHub() *CollaborationHub {
	return &CollaborationHub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		codeState: &models.CodeState{
			Code:    "// 欢迎使用协同编程平台\n// 在此编写你的C++代码\n\n#include <iostream>\n\nint main() {\n    std::cout << \"Hello, World!\" << std::endl;\n    return 0;\n}\n",
			Version: 0,
			Updated: time.Now(),
		},
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
